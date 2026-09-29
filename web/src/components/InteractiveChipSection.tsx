import React, { Suspense, useMemo, useState } from 'react';
import { Cpu, Info, Layers, Play, Radio, Sliders, Timer, Zap } from 'lucide-react';
import { ChipModel3D } from './ChipModel3D';
import {
  createGridGraph,
  dijkstra,
  measureDijkstraUs,
  simulatePhysicalRace,
  DEFAULT_PARAMS,
  RaceHardwareResult,
} from '../lib/racelogic';

const fmtNs = (ns: number) => `${ns.toFixed(1)} ns`;
const fmtUs = (us: number) => (us >= 1000 ? `${(us / 1000).toFixed(2)} ms` : `${us.toFixed(1)} µs`);

export const InteractiveChipSection: React.FC = () => {
  const [tab, setTab] = useState<'3d' | 'race'>('3d');

  // 3D chip state
  const [explode, setExplode] = useState(0);
  const [fireToken, setFireToken] = useState(0);
  const [racing3D, setRacing3D] = useState(false);
  const [race3D, setRace3D] = useState<RaceHardwareResult | null>(null);

  // 4x4 race vs browser Dijkstra
  const graph4 = useMemo(() => createGridGraph(4, 4, 15, 1234), []);
  const [sourceNode, setSourceNode] = useState(0);
  const [targetNode, setTargetNode] = useState(15);
  const [unitPS, setUnitPS] = useState(100);
  const [runId, setRunId] = useState(0);
  const [cpu, setCpu] = useState<{ small: number; large: number } | null>(null);

  const race4 = useMemo(
    () => simulatePhysicalRace(graph4, sourceNode, { ...DEFAULT_PARAMS, unitDelayPS: unitPS }),
    // runId re-samples the noise on every run
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [graph4, sourceNode, unitPS, runId],
  );
  const exact4 = useMemo(() => dijkstra(graph4, sourceNode), [graph4, sourceNode]);
  const decoded = Math.round(race4.fireTimesPS[targetNode] / unitPS);
  const correct = decoded === exact4.dist[targetNode];

  const path = useMemo(() => {
    const p: number[] = [];
    let v = targetNode;
    while (v !== -1) {
      p.unshift(v);
      if (v === sourceNode) break;
      v = race4.predecessors[v];
    }
    return p;
  }, [race4, sourceNode, targetNode]);

  const runRace = () => {
    setRunId((r) => r + 1);
    // Measure after the click so the page load stays light.
    setTimeout(() => setCpu({ small: measureDijkstraUs(4, 4000), large: measureDijkstraUs(16, 400) }), 30);
  };

  return (
    <section id="section-chip" className="relative overflow-hidden border-t border-white/10 bg-[#040814] py-28">
      <div className="pointer-events-none absolute left-1/2 top-1/2 -z-10 h-[550px] w-[850px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-gradient-to-tr from-cyan-600/15 via-blue-600/10 to-transparent blur-[200px]" />

      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto mb-12 max-w-3xl text-center">
          <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-cyan-500/40 bg-cyan-950/70 px-3.5 py-1.5 font-mono text-xs text-cyan-300">
            <Radio className="h-3.5 w-3.5 text-cyan-400" />
            <span>INSIDE THE CHIP · SIMULATED, NOT YET FABRICATED</span>
          </div>
          <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">See the race on the die.</h2>
          <p className="mt-4 text-base font-normal leading-relaxed text-slate-300 sm:text-lg">
            A physically scaled model of the 16×16 race-logic chip: a 25.5 mm die with 960 programmable delay spirals, wire
            bonded on its package and fed by a fiber array. Every glowing node below fires at the time computed by the
            project&apos;s physics model.
          </p>

          <div className="mt-8 inline-flex items-center rounded-full border border-white/15 bg-white/10 p-1.5 backdrop-blur-xl">
            {(
              [
                ['3d', 'The chip in 3D', Layers],
                ['race', 'Race vs. Dijkstra', Zap],
              ] as const
            ).map(([id, label, Icon]) => (
              <button
                key={id}
                onClick={() => setTab(id)}
                className={`flex items-center space-x-2 rounded-full px-6 py-2.5 text-xs font-bold transition-all ${
                  tab === id
                    ? 'bg-gradient-to-r from-cyan-400 to-blue-500 text-slate-950 shadow-[0_0_20px_rgba(6,182,212,0.5)]'
                    : 'text-slate-200 hover:text-white'
                }`}
              >
                <Icon className="h-4 w-4" />
                <span>{label}</span>
              </button>
            ))}
          </div>
        </div>

        {tab === '3d' && (
          <div className="relative overflow-hidden rounded-3xl border border-white/15 bg-black shadow-2xl">
            <div className="h-[560px] w-full sm:h-[640px]">
              <Suspense fallback={<div className="flex h-full items-center justify-center text-sm text-slate-400">Loading 3D model…</div>}>
                <ChipModel3D
                  explode={explode}
                  fireToken={fireToken}
                  onRaceComplete={(r) => {
                    setRace3D(r);
                    setRacing3D(false);
                  }}
                />
              </Suspense>
            </div>

            {/* Telemetry */}
            <div className="pointer-events-none absolute left-4 top-4 max-w-xs rounded-2xl border border-white/15 bg-black/70 p-4 font-mono text-[11px] text-slate-300 backdrop-blur-xl">
              <div className="mb-2 text-xs font-bold text-white">16×16 race · unit 100 ps</div>
              {race3D ? (
                <div className="space-y-1">
                  <div className="flex justify-between gap-6">
                    <span>Light race</span>
                    <span className="text-cyan-300">{fmtNs(race3D.raceSolveTimeNS)}</span>
                  </div>
                  <div className="flex justify-between gap-6">
                    <span>TDC readout</span>
                    <span className="text-cyan-300">{fmtNs(race3D.readoutTimeNS)}</span>
                  </div>
                  <div className="flex justify-between gap-6 border-t border-white/10 pt-1">
                    <span>Total per query</span>
                    <span className="font-bold text-white">{fmtNs(race3D.totalLatencyNS)}</span>
                  </div>
                </div>
              ) : (
                <div className="text-slate-400">Run the race to see the simulated timings.</div>
              )}
              <div className="mt-3 flex items-start gap-1.5 text-[10px] leading-snug text-slate-400">
                <Info className="mt-px h-3 w-3 flex-shrink-0" />
                <span>False color: 1550 nm light is infrared. The ~11 ns race is slowed down to 4 s.</span>
              </div>
            </div>

            {/* Controls */}
            <div className="pointer-events-none absolute bottom-5 left-0 right-0 flex flex-wrap items-center justify-between gap-4 px-5">
              <button
                onClick={() => {
                  setRacing3D(true);
                  setFireToken((t) => t + 1);
                }}
                disabled={racing3D}
                className="pointer-events-auto flex items-center space-x-2 rounded-full bg-gradient-to-r from-cyan-400 via-teal-300 to-blue-400 px-7 py-3 text-sm font-bold text-slate-950 shadow-[0_0_25px_rgba(6,182,212,0.6)] transition-all hover:from-cyan-300 hover:to-blue-300 disabled:opacity-50"
              >
                <Play className="h-4 w-4 fill-slate-950" />
                <span>{racing3D ? 'Light racing…' : 'Inject a pulse'}</span>
              </button>

              <div className="pointer-events-auto flex items-center space-x-3 rounded-full border border-white/20 bg-black/80 px-5 py-2.5 backdrop-blur-xl">
                <Sliders className="h-4 w-4 text-cyan-400" />
                <span className="font-mono text-xs text-slate-200">Explode layers</span>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  value={explode}
                  onChange={(e) => setExplode(parseFloat(e.target.value))}
                  className="w-32 cursor-pointer accent-cyan-400"
                />
              </div>
            </div>
          </div>
        )}

        {tab === 'race' && (
          <div className="space-y-8">
            <div className="flex flex-wrap items-center justify-between gap-6 rounded-3xl border border-white/15 bg-white/[0.04] p-6">
              <div className="flex flex-wrap items-center gap-3 font-mono text-xs text-slate-200">
                {(
                  [
                    ['Source', sourceNode, setSourceNode, 'text-cyan-300'],
                    ['Target', targetNode, setTargetNode, 'text-emerald-300'],
                  ] as const
                ).map(([label, value, set, color]) => (
                  <label key={label} className="flex items-center space-x-2 rounded-full border border-white/15 bg-white/10 px-3.5 py-2">
                    <span className="text-slate-400">{label}:</span>
                    <select
                      value={value}
                      onChange={(e) => set(Number(e.target.value))}
                      className={`cursor-pointer bg-transparent font-bold focus:outline-none ${color}`}
                    >
                      {Array.from({ length: 16 }).map((_, i) => (
                        <option key={i} value={i} className="bg-slate-900 text-white">
                          Node {i}
                        </option>
                      ))}
                    </select>
                  </label>
                ))}
                <button
                  onClick={() => setUnitPS(unitPS === 100 ? 50 : 100)}
                  className="rounded-full border border-white/15 bg-white/10 px-3.5 py-2 text-amber-300 hover:underline"
                >
                  Delay unit: {unitPS} ps
                </button>
              </div>
              <button
                onClick={runRace}
                className="flex items-center space-x-2 rounded-full bg-gradient-to-r from-cyan-400 via-teal-300 to-blue-400 px-8 py-3.5 text-sm font-bold text-slate-950 shadow-[0_0_30px_rgba(6,182,212,0.6)]"
              >
                <Play className="h-4 w-4 fill-slate-950" />
                <span>Run race + measure CPU</span>
              </button>
            </div>

            <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
              {/* Photonic race */}
              <div className="rounded-3xl border border-cyan-500/40 bg-gradient-to-b from-cyan-950/30 to-black p-8">
                <div className="mb-6 flex items-center justify-between border-b border-white/15 pb-5">
                  <div className="flex items-center space-x-3">
                    <Zap className="h-6 w-6 text-cyan-400" />
                    <div>
                      <h3 className="text-xl font-bold text-white">Photonic race (simulated)</h3>
                      <p className="font-mono text-xs text-slate-400">4×4 map · Si₃N₄ / TFLN model</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-mono text-2xl font-black text-cyan-400">{fmtNs(race4.totalLatencyNS)}</div>
                    <span className="font-mono text-xs text-slate-400">race + readout</span>
                  </div>
                </div>

                <div className="mb-6 grid grid-cols-4 gap-3">
                  {Array.from({ length: 16 }).map((_, id) => {
                    const onPath = path.includes(id);
                    return (
                      <div
                        key={id}
                        className={`flex h-14 flex-col items-center justify-center rounded-xl border font-mono text-xs ${
                          id === sourceNode
                            ? 'border-cyan-300 bg-cyan-500 font-extrabold text-slate-950'
                            : id === targetNode
                              ? 'border-emerald-300 bg-emerald-500 font-extrabold text-slate-950'
                              : onPath
                                ? 'border-cyan-400 bg-cyan-950 text-cyan-200'
                                : 'border-white/10 bg-white/5 text-slate-300'
                        }`}
                      >
                        <span className="font-bold">{id}</span>
                        <span className="text-[10px] opacity-70">{(race4.fireTimesPS[id] / 1000).toFixed(2)} ns</span>
                      </div>
                    );
                  })}
                </div>

                <div className="space-y-2 font-mono text-xs">
                  <Row label="Light race" value={fmtNs(race4.raceSolveTimeNS)} />
                  <Row label="TDC readout (16 nodes, 100 Gb/s)" value={fmtNs(race4.readoutTimeNS)} />
                  <Row label="Decoded distance to target" value={`${decoded} (exact: ${exact4.dist[targetNode]})`} tone={correct ? 'ok' : 'bad'} />
                  <Row label="Worst-path decision factor Q" value={race4.qFactor.toFixed(2)} />
                </div>
              </div>

              {/* Classical baseline */}
              <div className="rounded-3xl border border-white/15 bg-gradient-to-b from-white/[0.04] to-black p-8">
                <div className="mb-6 flex items-center space-x-3 border-b border-white/15 pb-5">
                  <Cpu className="h-6 w-6 text-slate-300" />
                  <div>
                    <h3 className="text-xl font-bold text-white">Classical Dijkstra (measured)</h3>
                    <p className="font-mono text-xs text-slate-400">Same graphs, timed live in your browser</p>
                  </div>
                </div>

                <div className="space-y-2 font-mono text-xs">
                  <Row label="4×4 map, per query" value={cpu ? fmtUs(cpu.small) : 'press Run'} />
                  <Row label="16×16 map, per query" value={cpu ? fmtUs(cpu.large) : 'press Run'} />
                  <Row label="Best native baseline, 16×16 (Dial, Go, i3-3217U, 2012)" value="25.4 µs" />
                  <Row label="Photonic model, 16×16 (race + readout)" value="42.2 ns" />
                </div>

                <div className="mt-6 flex items-start gap-2 rounded-2xl border border-white/10 bg-white/[0.03] p-4 text-xs leading-relaxed text-slate-300">
                  <Timer className="mt-0.5 h-4 w-4 flex-shrink-0 text-amber-300" />
                  <span>
                    Your browser runs JavaScript, so its times are slower than native code. Against the best native algorithm the model is
                    ~600× faster than a measured 2012 laptop CPU and ~42× faster than an Apple M5 Max (estimated), when the
                    map changes between queries; on a fixed map a precomputed table is faster. The
                    photonic numbers come from simulation: no chip has been fabricated yet.
                  </span>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </section>
  );
};

const Row: React.FC<{ label: string; value: string; tone?: 'ok' | 'bad' }> = ({ label, value, tone }) => (
  <div className="flex items-center justify-between rounded-xl border border-white/10 bg-white/[0.04] p-3">
    <span className="text-slate-300">{label}</span>
    <span className={`font-bold ${tone === 'ok' ? 'text-emerald-400' : tone === 'bad' ? 'text-rose-400' : 'text-cyan-300'}`}>
      {value}
    </span>
  </div>
);
