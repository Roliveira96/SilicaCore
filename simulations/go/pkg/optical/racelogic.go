package optical

import (
	"math"
	"math/rand"
	"time"
)

// ============================================================================
// PHOTONIC RACE LOGIC: SHORTEST PATH BY TIME OF FLIGHT
// A graph is mapped onto hardware: every edge is a programmable delay line
// (binary-weighted Si3N4 spirals selected by Sb2Se3 switches) and every node
// detects the first arriving pulse and re-fires it to all outgoing edges
// (UTC photodiode + comparator + TFLN modulator). The first-arrival time at a
// node is its shortest distance (Madhavan, Sherwood & Strukov, ISCA 2014).
// ============================================================================

const (
	// DefaultRaceUnitDelayPS is the delay that represents one unit of edge weight. 100 ps gives 0 errors in
	// 2.55e7 decoded distances (16x16, 10 chips x 1e4 queries); 50 ps gives ~1.2e-4 per distance.
	DefaultRaceUnitDelayPS = 100.0

	// DefaultRaceWeightBits is the number of binary-weighted delay stages per edge (weights 1..2^bits-1).
	DefaultRaceWeightBits = 4

	// DefaultNodeRegenLatencyPS is the detect-and-refire latency of a node (photodiode + comparator + TFLN driver).
	// ASSUMPTION: no published integrated figure; compensated in each edge's programmed delay.
	DefaultNodeRegenLatencyPS = 20.0

	// DefaultNodeTimingJitterPS is the rms timing jitter added at each node firing (UTC photodiode + comparator).
	// ASSUMPTION: sub-2 ps rms is typical of high-speed electro-optic receivers; measure on bench.
	DefaultNodeTimingJitterPS = 1.5

	// DefaultEdgeDelayErrorPS is the rms static fabrication error of each programmed delay after calibration (ASSUMPTION).
	DefaultEdgeDelayErrorPS = 0.5

	// DefaultPcmSwitchLossDb is the insertion loss of one Sb2Se3 switch; ~0.1-0.5 dB per pi reported [arXiv:2604.11649].
	DefaultPcmSwitchLossDb = 0.25

	// DefaultPcmProgramTimeNS is the time to program all delay switches in parallel (ITO micro-heater pulses, ASSUMPTION).
	DefaultPcmProgramTimeNS = 1000.0

	// DefaultTdcBitsPerNode is the resolution of the time-to-digital converter that captures each node's firing time.
	DefaultTdcBitsPerNode = 12

	// DefaultReadoutLinkGbps is the bandwidth of the link that reads all node times back to electronics (ASSUMPTION).
	DefaultReadoutLinkGbps = 100.0

	// DefaultSpiralPitchUM is the centre-to-centre pitch of Si3N4 delay spirals, set by crosstalk (ASSUMPTION).
	DefaultSpiralPitchUM = 3.0

	// DefaultPerEdgeModulators places one TFLN modulator on every outgoing edge, so fan-out splits the CW laser feed
	// (separate power budget) instead of the data pulse. A shared modulator per node pays 10*log10(fan-out) on data.
	DefaultPerEdgeModulators = true
)

// RaceLogicParams holds the physical parameters of the photonic race-logic accelerator.
type RaceLogicParams struct {
	UnitDelayPS        float64
	WeightBits         int
	NodeRegenLatencyPS float64
	NodeJitterPS       float64
	EdgeDelayErrorPS   float64
	PcmSwitchLossDb    float64
	PcmProgramTimeNS   float64
	PerEdgeModulators  bool
	TdcBitsPerNode     int
	ReadoutLinkGbps    float64
	Platform           RoutingPlatform
}

// DefaultRaceLogicParams returns the nominal race-logic configuration on the Si3N4 + TFLN platform.
func DefaultRaceLogicParams() RaceLogicParams {
	return RaceLogicParams{
		UnitDelayPS:        DefaultRaceUnitDelayPS,
		WeightBits:         DefaultRaceWeightBits,
		NodeRegenLatencyPS: DefaultNodeRegenLatencyPS,
		NodeJitterPS:       DefaultNodeTimingJitterPS,
		EdgeDelayErrorPS:   DefaultEdgeDelayErrorPS,
		PcmSwitchLossDb:    DefaultPcmSwitchLossDb,
		PcmProgramTimeNS:   DefaultPcmProgramTimeNS,
		PerEdgeModulators:  DefaultPerEdgeModulators,
		TdcBitsPerNode:     DefaultTdcBitsPerNode,
		ReadoutLinkGbps:    DefaultReadoutLinkGbps,
		Platform:           SiliconNitrideTflnPlatform(),
	}
}

// MaxWeight returns the largest edge weight the binary delay stages can encode.
func (p RaceLogicParams) MaxWeight() int {
	return (1 << p.WeightBits) - 1
}

// Edge is a directed weighted edge.
type Edge struct {
	To     int
	Weight int
}

// Graph is a directed graph in adjacency-list form with positive integer weights.
type Graph struct {
	Adj [][]Edge
}

// NumEdges returns the total number of directed edges.
func (g Graph) NumEdges() int {
	n := 0
	for _, a := range g.Adj {
		n += len(a)
	}
	return n
}

// MaxOutDegree returns the largest fan-out of any node.
func (g Graph) MaxOutDegree() int {
	m := 0
	for _, a := range g.Adj {
		if len(a) > m {
			m = len(a)
		}
	}
	return m
}

// NewGridGraph builds a W x H 4-neighbour grid (a game map) with random terrain costs in [1, maxWeight].
// Each undirected neighbour pair gets the same cost in both directions.
func NewGridGraph(width, height, maxWeight int, seed int64) Graph {
	r := rand.New(rand.NewSource(seed))
	n := width * height
	adj := make([][]Edge, n)
	id := func(x, y int) int { return y*width + x }

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x+1 < width {
				w := 1 + r.Intn(maxWeight)
				adj[id(x, y)] = append(adj[id(x, y)], Edge{To: id(x+1, y), Weight: w})
				adj[id(x+1, y)] = append(adj[id(x+1, y)], Edge{To: id(x, y), Weight: w})
			}
			if y+1 < height {
				w := 1 + r.Intn(maxWeight)
				adj[id(x, y)] = append(adj[id(x, y)], Edge{To: id(x, y+1), Weight: w})
				adj[id(x, y+1)] = append(adj[id(x, y+1)], Edge{To: id(x, y), Weight: w})
			}
		}
	}
	return Graph{Adj: adj}
}

type pqItem struct {
	node int
	key  float64
}

// priorityQueue is a typed binary min-heap; it avoids container/heap interface boxing so the
// classical baseline is not slowed down by per-push allocations.
type priorityQueue []pqItem

func (q *priorityQueue) push(it pqItem) {
	*q = append(*q, it)
	h := *q
	i := len(h) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if h[parent].key <= h[i].key {
			break
		}
		h[parent], h[i] = h[i], h[parent]
		i = parent
	}
}

func (q *priorityQueue) pop() pqItem {
	h := *q
	top := h[0]
	last := len(h) - 1
	h[0] = h[last]
	h = h[:last]
	i := 0
	for {
		l, r, m := 2*i+1, 2*i+2, i
		if l < len(h) && h[l].key < h[m].key {
			m = l
		}
		if r < len(h) && h[r].key < h[m].key {
			m = r
		}
		if m == i {
			break
		}
		h[i], h[m] = h[m], h[i]
		i = m
	}
	*q = h
	return top
}

// Dijkstra returns exact integer shortest distances from src (classical electronic baseline). Unreachable = -1.
func Dijkstra(g Graph, src int) []int {
	n := len(g.Adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	best := make([]float64, n)
	for i := range best {
		best[i] = math.Inf(1)
	}
	best[src] = 0
	q := make(priorityQueue, 0, n)
	q.push(pqItem{node: src, key: 0})

	for len(q) > 0 {
		it := q.pop()
		if dist[it.node] >= 0 {
			continue
		}
		dist[it.node] = int(it.key)
		for _, e := range g.Adj[it.node] {
			k := it.key + float64(e.Weight)
			if dist[e.To] < 0 && k < best[e.To] {
				best[e.To] = k
				q.push(pqItem{node: e.To, key: k})
			}
		}
	}
	return dist
}

// RaceHardware is the physical realisation of a graph on the race-logic accelerator.
type RaceHardware struct {
	Params           RaceLogicParams
	Graph            Graph
	EdgeDelayPS      [][]float64 // Programmed delay per edge, node latency already subtracted
	EdgeErrorPS      [][]float64 // Static fabrication/calibration error per edge
	MaxEdgeLengthMM  float64
	TotalDelayM      float64 // Total delay-line waveguide: every edge carries all binary stages
	DelayAreaMM2     float64 // Footprint of all delay spirals at DefaultSpiralPitchUM
	FitsOnReticle    bool
	WorstEdgeLossDb  float64
	PowerMarginDb    float64
	LinkBudgetOK     bool
	Detectors        int // One photodiode per incoming edge (electronic OR, no lossy combiner)
	Tdcs             int // One time-to-digital converter per node to capture its firing time
	Modulators       int // TFLN re-fire modulators: one per node, or one per edge with PerEdgeModulators
	PcmSwitches      int // Two 2x2 switches per binary delay stage per edge
	WeightsSupported bool
}

// BuildRaceHardware programs the graph onto the accelerator and evaluates its optical link budget.
func BuildRaceHardware(g Graph, p RaceLogicParams, seed int64) RaceHardware {
	r := rand.New(rand.NewSource(seed))
	delays := make([][]float64, len(g.Adj))
	errs := make([][]float64, len(g.Adj))
	supported := p.UnitDelayPS >= p.NodeRegenLatencyPS
	maxW := 0

	for u, edges := range g.Adj {
		delays[u] = make([]float64, len(edges))
		errs[u] = make([]float64, len(edges))
		for i, e := range edges {
			if e.Weight < 1 || e.Weight > p.MaxWeight() {
				supported = false
			}
			if e.Weight > maxW {
				maxW = e.Weight
			}
			delays[u][i] = float64(e.Weight)*p.UnitDelayPS - p.NodeRegenLatencyPS
			errs[u][i] = r.NormFloat64() * p.EdgeDelayErrorPS
		}
	}

	maxDelay := float64(maxW)*p.UnitDelayPS - p.NodeRegenLatencyPS
	length := PathLengthForDelayMM(maxDelay, p.Platform.GroupIndex)
	loss := p.Platform.PropagationDbPerCm*length/10.0 +
		float64(2*p.WeightBits)*p.PcmSwitchLossDb +
		p.Platform.SwitchLossDb // TFLN re-fire modulator

	modulators := g.NumEdges()
	if !p.PerEdgeModulators {
		loss += 10 * math.Log10(float64(g.MaxOutDegree()))
		modulators = len(g.Adj)
	}
	margin := DefaultChannelLaserPowerDbm - DefaultPhotodiodeSensitivityDbm

	// Binary-weighted stages mean every edge holds the full maximum delay, whatever its weight.
	stagesDelay := float64(p.MaxWeight()) * p.UnitDelayPS
	totalM := PathLengthForDelayMM(stagesDelay, p.Platform.GroupIndex) * float64(g.NumEdges()) / 1000.0
	areaMM2 := totalM * 1000.0 * DefaultSpiralPitchUM * 1e-3

	return RaceHardware{
		Params:           p,
		Graph:            g,
		EdgeDelayPS:      delays,
		EdgeErrorPS:      errs,
		MaxEdgeLengthMM:  length,
		TotalDelayM:      totalM,
		DelayAreaMM2:     areaMM2,
		FitsOnReticle:    areaMM2 <= ReticleAreaCM2*100,
		WorstEdgeLossDb:  loss,
		PowerMarginDb:    margin,
		LinkBudgetOK:     loss <= margin,
		Detectors:        g.NumEdges(),
		Tdcs:             len(g.Adj),
		Modulators:       modulators,
		PcmSwitches:      g.NumEdges() * 2 * p.WeightBits,
		WeightsSupported: supported,
	}
}

// Race fires a pulse from src and returns each node's firing time in ps (Inf if unreachable).
// Physically every node fires once, at its first arrival plus regeneration latency and jitter;
// processing events in time order is exactly an event-driven simulation of the race.
func (h RaceHardware) Race(src int, r *rand.Rand) []float64 {
	n := len(h.Graph.Adj)
	fire := make([]float64, n)
	done := make([]bool, n)
	for i := range fire {
		fire[i] = math.Inf(1)
	}
	arrival := make([]float64, n)
	for i := range arrival {
		arrival[i] = math.Inf(1)
	}
	arrival[src] = -h.Params.NodeRegenLatencyPS // source fires at t = 0
	q := make(priorityQueue, 0, n)
	q.push(pqItem{node: src, key: arrival[src]})

	for len(q) > 0 {
		it := q.pop()
		u := it.node
		if done[u] {
			continue
		}
		done[u] = true

		jitter := 0.0
		if u != src {
			jitter = r.NormFloat64() * h.Params.NodeJitterPS
		}
		fire[u] = it.key + h.Params.NodeRegenLatencyPS + jitter

		for i, e := range h.Graph.Adj[u] {
			if done[e.To] {
				continue
			}
			t := fire[u] + h.EdgeDelayPS[u][i] + h.EdgeErrorPS[u][i]
			if t < arrival[e.To] {
				arrival[e.To] = t
				q.push(pqItem{node: e.To, key: t})
			}
		}
	}
	return fire
}

// RaceLogicResult summarises correctness and speed of the photonic race against classical Dijkstra.
type RaceLogicResult struct {
	Nodes               int     `json:"nodes"`
	Edges               int     `json:"edges"`
	Trials              int     `json:"trials"`
	UnitDelayPS         float64 `json:"unit_delay_ps"`
	NodeErrorRate       float64 `json:"node_error_rate"`
	QueryErrorRate      float64 `json:"query_error_rate"`
	MaxHops             int     `json:"max_hops"`
	WorstPathSigmaPS    float64 `json:"worst_path_sigma_ps"`
	DecisionQ           float64 `json:"decision_q"`
	RaceSolveTimeNS     float64 `json:"race_solve_time_ns"`
	ReadoutTimeNS       float64 `json:"readout_time_ns"`
	TotalQueryTimeNS    float64 `json:"total_query_time_ns"`
	DijkstraSolveTimeNS float64 `json:"dijkstra_solve_time_ns"`
	SpeedupPerQuery     float64 `json:"speedup_per_query"`
	ProgramTimeNS       float64 `json:"program_time_ns"`
	BreakEvenQueries    float64 `json:"break_even_queries"`
	LinkBudgetOK        bool    `json:"link_budget_ok"`
	WorstEdgeLossDb     float64 `json:"worst_edge_loss_db"`
	MaxEdgeLengthMM     float64 `json:"max_edge_length_mm"`
	DelayAreaMM2        float64 `json:"delay_area_mm2"`
	FitsOnReticle       bool    `json:"fits_on_reticle"`
	WeightsSupported    bool    `json:"weights_supported"`
	Detectors           int     `json:"detectors"`
	Tdcs                int     `json:"tdcs"`
	Modulators          int     `json:"modulators"`
	PcmSwitches         int     `json:"pcm_switches"`
}

// hopCounts returns the number of edges on one shortest path to every node (BFS over the shortest-path DAG).
func hopCounts(g Graph, dist []int, src int) []int {
	hops := make([]int, len(g.Adj))
	for i := range hops {
		hops[i] = -1
	}
	hops[src] = 0
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, e := range g.Adj[u] {
			if hops[e.To] < 0 && dist[e.To] == dist[u]+e.Weight {
				hops[e.To] = hops[u] + 1
				queue = append(queue, e.To)
			}
		}
	}
	return hops
}

// SimulateRaceLogic runs Monte Carlo races from random sources and checks every decoded distance against Dijkstra.
func SimulateRaceLogic(g Graph, p RaceLogicParams, trials int, seed int64) RaceLogicResult {
	hw := BuildRaceHardware(g, p, seed)
	r := rand.New(rand.NewSource(seed + 1))
	n := len(g.Adj)

	nodeErrors, queryErrors, maxHops := 0, 0, 0
	var raceTimeSum, dijkstraNsSum float64

	for t := 0; t < trials; t++ {
		src := r.Intn(n)

		start := time.Now()
		dist := Dijkstra(g, src)
		dijkstraNsSum += float64(time.Since(start).Nanoseconds())

		fire := hw.Race(src, r)
		wrong := false
		maxFire := 0.0
		for v := 0; v < n; v++ {
			if dist[v] < 0 {
				continue
			}
			if fire[v] > maxFire {
				maxFire = fire[v]
			}
			if int(math.Round(fire[v]/p.UnitDelayPS)) != dist[v] {
				nodeErrors++
				wrong = true
			}
		}
		if wrong {
			queryErrors++
		}
		raceTimeSum += maxFire

		for _, h := range hopCounts(g, dist, src) {
			if h > maxHops {
				maxHops = h
			}
		}
	}

	// Worst-case path noise: node jitter and static edge error both accumulate once per hop.
	sigmaPath := math.Sqrt(float64(maxHops)) * math.Hypot(p.NodeJitterPS, p.EdgeDelayErrorPS)
	raceNS := raceTimeSum / float64(trials) / 1000.0
	readoutNS := float64(n*p.TdcBitsPerNode) / p.ReadoutLinkGbps
	totalNS := raceNS + readoutNS
	dijkstraNS := dijkstraNsSum / float64(trials)
	speedup := dijkstraNS / totalNS

	breakEven := math.Inf(1)
	if dijkstraNS > totalNS {
		breakEven = p.PcmProgramTimeNS / (dijkstraNS - totalNS)
	}

	return RaceLogicResult{
		Nodes:               n,
		Edges:               g.NumEdges(),
		Trials:              trials,
		UnitDelayPS:         p.UnitDelayPS,
		NodeErrorRate:       float64(nodeErrors) / float64(trials*n),
		QueryErrorRate:      float64(queryErrors) / float64(trials),
		MaxHops:             maxHops,
		WorstPathSigmaPS:    sigmaPath,
		DecisionQ:           p.UnitDelayPS / (2 * sigmaPath),
		RaceSolveTimeNS:     raceNS,
		ReadoutTimeNS:       readoutNS,
		TotalQueryTimeNS:    totalNS,
		DijkstraSolveTimeNS: dijkstraNS,
		SpeedupPerQuery:     speedup,
		ProgramTimeNS:       p.PcmProgramTimeNS,
		BreakEvenQueries:    breakEven,
		LinkBudgetOK:        hw.LinkBudgetOK,
		WorstEdgeLossDb:     hw.WorstEdgeLossDb,
		MaxEdgeLengthMM:     hw.MaxEdgeLengthMM,
		DelayAreaMM2:        hw.DelayAreaMM2,
		FitsOnReticle:       hw.FitsOnReticle,
		WeightsSupported:    hw.WeightsSupported,
		Detectors:           hw.Detectors,
		Tdcs:                hw.Tdcs,
		Modulators:          hw.Modulators,
		PcmSwitches:         hw.PcmSwitches,
	}
}

// HopErrorBin compares measured and predicted decoding error for nodes at a given shortest-path hop count.
type HopErrorBin struct {
	Hops          int     `json:"hops"`
	Samples       int     `json:"samples"`
	Errors        int     `json:"errors"`
	MeasuredRate  float64 `json:"measured_rate"`
	PredictedRate float64 `json:"predicted_rate"`
}

// RaceStatisticsResult aggregates a large multi-chip Monte Carlo campaign of race-logic queries.
type RaceStatisticsResult struct {
	UnitDelayPS       float64       `json:"unit_delay_ps"`
	Chips             int           `json:"chips"`
	Queries           int           `json:"queries"`
	NodeSamples       int           `json:"node_samples"`
	NodeErrors        int           `json:"node_errors"`
	QueryErrors       int           `json:"query_errors"`
	NodeErrorRate     float64       `json:"node_error_rate"`
	NodeErrorUpper95  float64       `json:"node_error_upper95"`
	QueryErrorRate    float64       `json:"query_error_rate"`
	QueryErrorUpper95 float64       `json:"query_error_upper95"`
	WorstPathQ        float64       `json:"worst_path_q"`
	PredictedNodeRate float64       `json:"predicted_node_rate"`
	ByHops            []HopErrorBin `json:"by_hops"`
}

// upper95 returns a 95% upper confidence bound on a binomial rate (rule of three when there are no errors,
// Wilson score interval otherwise).
func upper95(errors, n int) float64 {
	if n == 0 {
		return 1
	}
	if errors == 0 {
		return 3.0 / float64(n)
	}
	z := 1.959964
	p := float64(errors) / float64(n)
	nf := float64(n)
	den := 1 + z*z/nf
	centre := p + z*z/(2*nf)
	margin := z * math.Sqrt(p*(1-p)/nf+z*z/(4*nf*nf))
	return (centre + margin) / den
}

// predictedNodeError is the single-path Gaussian prediction: the decoded distance is wrong when the
// accumulated timing noise over h hops exceeds half a unit delay.
func predictedNodeError(p RaceLogicParams, hops int) float64 {
	if hops == 0 {
		return 0
	}
	sigma := math.Sqrt(float64(hops)) * math.Hypot(p.NodeJitterPS, p.EdgeDelayErrorPS)
	return math.Erfc((p.UnitDelayPS / 2) / (sigma * math.Sqrt2))
}

// RaceLogicStatistics runs queriesPerChip random-source races on each of `chips` independently fabricated
// chips (fresh static edge errors per chip) in parallel, and bins errors by shortest-path hop count.
func RaceLogicStatistics(g Graph, p RaceLogicParams, chips, queriesPerChip int, seed int64) RaceStatisticsResult {
	type chipRes struct {
		nodeSamples, nodeErrors, queryErrors int
		binSamples, binErrors                map[int]int
	}
	results := make(chan chipRes, chips)
	n := len(g.Adj)

	for c := 0; c < chips; c++ {
		go func(chip int) {
			hw := BuildRaceHardware(g, p, seed+int64(chip)*7919)
			r := rand.New(rand.NewSource(seed + int64(chip)*104729 + 1))
			res := chipRes{binSamples: map[int]int{}, binErrors: map[int]int{}}
			for q := 0; q < queriesPerChip; q++ {
				src := r.Intn(n)
				dist := Dijkstra(g, src)
				hops := hopCounts(g, dist, src)
				fire := hw.Race(src, r)
				wrong := false
				for v := 0; v < n; v++ {
					if dist[v] < 0 || v == src {
						continue
					}
					res.nodeSamples++
					res.binSamples[hops[v]]++
					if int(math.Round(fire[v]/p.UnitDelayPS)) != dist[v] {
						res.nodeErrors++
						res.binErrors[hops[v]]++
						wrong = true
					}
				}
				if wrong {
					res.queryErrors++
				}
			}
			results <- res
		}(c)
	}

	total := chipRes{binSamples: map[int]int{}, binErrors: map[int]int{}}
	for c := 0; c < chips; c++ {
		res := <-results
		total.nodeSamples += res.nodeSamples
		total.nodeErrors += res.nodeErrors
		total.queryErrors += res.queryErrors
		for h, v := range res.binSamples {
			total.binSamples[h] += v
		}
		for h, v := range res.binErrors {
			total.binErrors[h] += v
		}
	}

	maxHops := 0
	predictedSum := 0.0
	for h, cnt := range total.binSamples {
		if h > maxHops {
			maxHops = h
		}
		predictedSum += predictedNodeError(p, h) * float64(cnt)
	}
	bins := make([]HopErrorBin, 0, maxHops)
	for h := 1; h <= maxHops; h++ {
		s := total.binSamples[h]
		if s == 0 {
			continue
		}
		bins = append(bins, HopErrorBin{
			Hops:          h,
			Samples:       s,
			Errors:        total.binErrors[h],
			MeasuredRate:  float64(total.binErrors[h]) / float64(s),
			PredictedRate: predictedNodeError(p, h),
		})
	}

	queries := chips * queriesPerChip
	sigmaWorst := math.Sqrt(float64(maxHops)) * math.Hypot(p.NodeJitterPS, p.EdgeDelayErrorPS)
	return RaceStatisticsResult{
		UnitDelayPS:       p.UnitDelayPS,
		Chips:             chips,
		Queries:           queries,
		NodeSamples:       total.nodeSamples,
		NodeErrors:        total.nodeErrors,
		QueryErrors:       total.queryErrors,
		NodeErrorRate:     float64(total.nodeErrors) / float64(total.nodeSamples),
		NodeErrorUpper95:  upper95(total.nodeErrors, total.nodeSamples),
		QueryErrorRate:    float64(total.queryErrors) / float64(queries),
		QueryErrorUpper95: upper95(total.queryErrors, queries),
		WorstPathQ:        p.UnitDelayPS / (2 * sigmaWorst),
		PredictedNodeRate: predictedSum / float64(total.nodeSamples),
		ByHops:            bins,
	}
}
