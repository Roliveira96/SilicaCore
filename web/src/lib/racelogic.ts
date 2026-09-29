// ============================================================================
// PHOTONIC RACE LOGIC SIMULATION & PHYSICS ENGINE (TypeScript Port)
// Faithfully implements simulations/go/pkg/optical/racelogic.go
// Time-of-Flight shortest path calculation with physical noise, losses & jitter.
// ============================================================================

export interface RaceLogicParams {
  unitDelayPS: number;
  weightBits: number;
  nodeRegenLatencyPS: number;
  nodeJitterPS: number;
  edgeDelayErrorPS: number;
  pcmSwitchLossDb: number;
  pcmProgramTimeNS: number;
  perEdgeModulators: boolean;
  tdcBitsPerNode: number;
  readoutLinkGbps: number;
  groupIndex: number;
  propagationDbPerCm: number;
  switchLossDb: number;
}

export const DEFAULT_PARAMS: RaceLogicParams = {
  unitDelayPS: 100.0,
  weightBits: 4,
  nodeRegenLatencyPS: 20.0,
  nodeJitterPS: 1.5,
  edgeDelayErrorPS: 0.5,
  pcmSwitchLossDb: 0.25,
  pcmProgramTimeNS: 1000.0,
  perEdgeModulators: true,
  tdcBitsPerNode: 12,
  readoutLinkGbps: 100.0,
  groupIndex: 1.8836,
  propagationDbPerCm: 0.05,
  switchLossDb: 0.2,
};

export interface Edge {
  from: number;
  to: number;
  weight: number;
}

export interface Graph {
  width: number;
  height: number;
  numNodes: number;
  adj: Edge[][];
}

// Pseudo-random Gaussian using Box-Muller transform
export function gaussianRandom(mean = 0, stdev = 1): number {
  const u = 1 - Math.random();
  const v = Math.random();
  const z = Math.sqrt(-2.0 * Math.log(u)) * Math.cos(2.0 * Math.PI * v);
  return z * stdev + mean;
}

// Complementary error function approximation for Q-factor error prediction
export function erfc(x: number): number {
  // Chebyshev fitting approximation
  const z = Math.abs(x);
  const t = 1.0 / (1.0 + 0.5 * z);
  const r =
    t *
    Math.exp(
      -z * z -
        1.26551223 +
        t *
          (1.00002368 +
            t *
              (0.37409196 +
                t *
                  (0.09678418 +
                    t *
                      (-0.18628806 +
                        t *
                          (0.27886807 +
                            t *
                              (-1.13520398 +
                                t * (1.48851587 + t * (-0.82215223 + t * 0.17087277))))))))
    );
  return x >= 0 ? r : 2.0 - r;
}

// Creates an NxM grid graph with pseudo-random positive integer terrain weights
export function createGridGraph(width: number, height: number, maxWeight = 15, seed = 42): Graph {
  const numNodes = width * height;
  const adj: Edge[][] = Array.from({ length: numNodes }, () => []);
  const id = (x: number, y: number) => y * width + x;

  // Simple LCG PRNG for reproducible graphs
  let s = seed;
  const nextRand = () => {
    s = (s * 1664525 + 1013904223) % 4294967296;
    return s / 4294967296;
  };

  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const u = id(x, y);
      if (x + 1 < width) {
        const v = id(x + 1, y);
        const w = 1 + Math.floor(nextRand() * maxWeight);
        adj[u].push({ from: u, to: v, weight: w });
        adj[v].push({ from: v, to: u, weight: w });
      }
      if (y + 1 < height) {
        const v = id(x, y + 1);
        const w = 1 + Math.floor(nextRand() * maxWeight);
        adj[u].push({ from: u, to: v, weight: w });
        adj[v].push({ from: v, to: u, weight: w });
      }
    }
  }

  return { width, height, numNodes, adj };
}

// Priority Queue for classical electronic Dijkstra
interface PQItem {
  node: number;
  cost: number;
}

export function dijkstra(
  graph: Graph,
  src: number
): { dist: number[]; prev: number[]; hops: number[]; solveTimeUs: number } {
  const startTime = performance.now();
  const n = graph.numNodes;
  const dist = new Array<number>(n).fill(Infinity);
  const prev = new Array<number>(n).fill(-1);
  const hops = new Array<number>(n).fill(-1);

  dist[src] = 0;
  hops[src] = 0;

  const pq: PQItem[] = [{ node: src, cost: 0 }];

  while (pq.length > 0) {
    // extract min
    let minIdx = 0;
    for (let i = 1; i < pq.length; i++) {
      if (pq[i].cost < pq[minIdx].cost) minIdx = i;
    }
    const { node: u, cost: d } = pq.splice(minIdx, 1)[0];

    if (d > dist[u]) continue;

    for (const edge of graph.adj[u]) {
      const v = edge.to;
      const alt = d + edge.weight;
      if (alt < dist[v]) {
        dist[v] = alt;
        prev[v] = u;
        hops[v] = hops[u] + 1;
        pq.push({ node: v, cost: alt });
      }
    }
  }

  const solveTimeUs = (performance.now() - startTime) * 1000;
  return { dist, prev, hops, solveTimeUs };
}

// Physical Path length in mm for delay in ps: L = c * t / n_g
export function pathLengthForDelayMM(delayPS: number, groupIndex = 1.8836): number {
  const cMmPerPS = 0.299792458; // mm / ps in vacuum
  return (cMmPerPS * delayPS) / groupIndex;
}

// Physical Race Logic Waveguide Simulation
export interface RaceNodeEvent {
  node: number;
  arrivalTimePS: number;
  firingTimePS: number;
  predecessor: number;
  hopCount: number;
  accumulatedJitterPS: number;
}

export interface RacePulseState {
  edge: Edge;
  startTimePS: number;
  durationPS: number;
  progress: number; // 0 to 1
  active: boolean;
  isWinningPath: boolean;
}

export interface RaceHardwareResult {
  fireTimesPS: number[];
  predecessors: number[];
  hopCounts: number[];
  events: RaceNodeEvent[];
  maxToF_PS: number;
  raceSolveTimeNS: number;
  readoutTimeNS: number;
  totalLatencyNS: number;
  qFactor: number;
  worstPathSigmaPS: number;
  predictedErrorRate: number;
  lossDb: number;
  dieFootprintMM2: number;
  fitsOnReticle: boolean;
}

export function simulatePhysicalRace(
  graph: Graph,
  src: number,
  params: RaceLogicParams = DEFAULT_PARAMS
): RaceHardwareResult {
  const n = graph.numNodes;
  const fireTimes = new Array<number>(n).fill(Infinity);
  const arrivalTimes = new Array<number>(n).fill(Infinity);
  const predecessors = new Array<number>(n).fill(-1);
  const hopCounts = new Array<number>(n).fill(0);
  const accumulatedJitter = new Array<number>(n).fill(0);
  const done = new Array<boolean>(n).fill(false);
  const events: RaceNodeEvent[] = [];

  // Source fires at t = 0
  arrivalTimes[src] = -params.nodeRegenLatencyPS;
  const pq: { node: number; time: number }[] = [{ node: src, time: arrivalTimes[src] }];

  while (pq.length > 0) {
    let minIdx = 0;
    for (let i = 1; i < pq.length; i++) {
      if (pq[i].time < pq[minIdx].time) minIdx = i;
    }
    const it = pq.splice(minIdx, 1)[0];
    const u = it.node;

    if (done[u]) continue;
    done[u] = true;

    const jitter = u === src ? 0 : gaussianRandom(0, params.nodeJitterPS);
    const fireT = it.time + params.nodeRegenLatencyPS + jitter;
    fireTimes[u] = fireT;
    accumulatedJitter[u] = (u === src ? 0 : accumulatedJitter[predecessors[u]]) + Math.abs(jitter);

    events.push({
      node: u,
      arrivalTimePS: it.time,
      firingTimePS: fireT,
      predecessor: predecessors[u],
      hopCount: hopCounts[u],
      accumulatedJitterPS: accumulatedJitter[u],
    });

    for (const e of graph.adj[u]) {
      const v = e.to;
      if (done[v]) continue;

      const edgeDelay = e.weight * params.unitDelayPS - params.nodeRegenLatencyPS;
      const fabError = gaussianRandom(0, params.edgeDelayErrorPS);
      const arrivalAtV = fireT + edgeDelay + fabError;

      if (arrivalAtV < arrivalTimes[v]) {
        arrivalTimes[v] = arrivalAtV;
        predecessors[v] = u;
        hopCounts[v] = hopCounts[u] + 1;
        pq.push({ node: v, time: arrivalAtV });
      }
    }
  }

  let maxToF = 0;
  let maxHops = 0;
  for (let i = 0; i < n; i++) {
    if (fireTimes[i] < Infinity && fireTimes[i] > maxToF) {
      maxToF = fireTimes[i];
    }
    if (hopCounts[i] > maxHops) {
      maxHops = hopCounts[i];
    }
  }

  // Physical timing noise sigma over worst-case path
  const sigmaPath =
    Math.sqrt(Math.max(1, maxHops)) *
    Math.hypot(params.nodeJitterPS, params.edgeDelayErrorPS);
  const qFactor = params.unitDelayPS / (2 * Math.max(0.001, sigmaPath));
  const predictedErrorRate = erfc(qFactor / Math.SQRT2);

  const raceSolveTimeNS = maxToF / 1000.0;
  const readoutTimeNS = (n * params.tdcBitsPerNode) / params.readoutLinkGbps;
  const totalLatencyNS = raceSolveTimeNS + readoutTimeNS;

  // Die Area footprint: binary delay spirals at 3.0 um pitch
  const stagesDelayPS = ((1 << params.weightBits) - 1) * params.unitDelayPS;
  let totalEdges = 0;
  for (const list of graph.adj) totalEdges += list.length;
  const totalWaveguideM =
    (pathLengthForDelayMM(stagesDelayPS, params.groupIndex) * totalEdges) / 1000.0;
  const dieFootprintMM2 = totalWaveguideM * 1000.0 * 3.0 * 1e-3; // Spiral pitch 3um

  // Optical loss
  const worstEdgeLengthMM = pathLengthForDelayMM(
    15 * params.unitDelayPS - params.nodeRegenLatencyPS,
    params.groupIndex
  );
  const lossDb =
    (params.propagationDbPerCm * worstEdgeLengthMM) / 10.0 +
    2 * params.weightBits * params.pcmSwitchLossDb +
    params.switchLossDb;

  return {
    fireTimesPS: fireTimes,
    predecessors,
    hopCounts,
    events,
    maxToF_PS: maxToF,
    raceSolveTimeNS,
    readoutTimeNS,
    totalLatencyNS,
    qFactor,
    worstPathSigmaPS: sigmaPath,
    predictedErrorRate,
    lossDb,
    dieFootprintMM2,
    fitsOnReticle: dieFootprintMM2 <= 858.0, // standard 858 mm2 reticle
  };
}

// Backtracks the winning path from target to source
export function getShortestPath(predecessors: number[], src: number, target: number): number[] {
  const path: number[] = [];
  let curr = target;
  while (curr !== -1) {
    path.unshift(curr);
    if (curr === src) break;
    curr = predecessors[curr];
  }
  return path[0] === src ? path : [];
}
