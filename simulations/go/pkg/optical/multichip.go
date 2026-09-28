package optical

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// ============================================================================
// MULTI-CHIP RACE LOGIC: COMPOSING 16x16 TILES FOR LARGER GAME MAPS
// Exact mode: one global race; edges that cross a tile border go chip-to-chip
// (facet coupling loss, link delay, extra jitter and static error).
// Hierarchical mode (HPA*-style): each tile precomputes, with its own optical
// race, distances between a few border transition points; a query runs two
// local races (source and target tiles) plus a small electronic Dijkstra over
// the transitions. No optical chip-to-chip links, but paths may be suboptimal.
// ============================================================================

const (
	// DefaultTileSide is the side of one race-logic chip; 16x16 is the largest grid that fits one reticle at 100 ps.
	DefaultTileSide = 16

	// DefaultInterChipCouplingLossDb is the loss per facet crossing between chips (fiber/interposer coupling;
	// <2.5 dB per facet reported for Si3N4 [Churaev et al., Nat. Commun. 2023]). Two facets per crossing edge.
	DefaultInterChipCouplingLossDb = 1.5

	// DefaultInterChipLinkDelayPS is the flight time of the chip-to-chip link (~5 mm interposer), compensated in
	// the programmed delay like the node latency (ASSUMPTION).
	DefaultInterChipLinkDelayPS = 25.0

	// DefaultInterChipJitterPS is the extra rms timing jitter added on every chip-to-chip crossing (ASSUMPTION).
	DefaultInterChipJitterPS = 1.0

	// DefaultInterChipStaticErrorPS is the rms static delay error of a chip-to-chip link after calibration (ASSUMPTION).
	DefaultInterChipStaticErrorPS = 1.0

	// DefaultTransitionStride is the spacing of HPA*-style transition points along each tile border
	// (one transition every 4 border cells: 4 per 16-cell border).
	DefaultTransitionStride = 4
)

// MultiChipParams configures how 16x16 tiles are composed.
type MultiChipParams struct {
	TileSide               int
	CouplingLossDbPerFacet float64
	LinkDelayPS            float64
	LinkJitterPS           float64
	LinkStaticErrorPS      float64
	TransitionStride       int
}

// DefaultMultiChipParams returns the nominal inter-chip configuration.
func DefaultMultiChipParams() MultiChipParams {
	return MultiChipParams{
		TileSide:               DefaultTileSide,
		CouplingLossDbPerFacet: DefaultInterChipCouplingLossDb,
		LinkDelayPS:            DefaultInterChipLinkDelayPS,
		LinkJitterPS:           DefaultInterChipJitterPS,
		LinkStaticErrorPS:      DefaultInterChipStaticErrorPS,
		TransitionStride:       DefaultTransitionStride,
	}
}

// TiledGrid is a side x side grid map split into tiles of TileSide x TileSide nodes.
type TiledGrid struct {
	Side     int
	TileSide int
	Graph    Graph
}

// tileOf returns the tile index of a node.
func (t TiledGrid) tileOf(node int) int {
	x, y := node%t.Side, node/t.Side
	tilesPerRow := t.Side / t.TileSide
	return (y/t.TileSide)*tilesPerRow + x/t.TileSide
}

// NumTiles returns the number of tiles (chips).
func (t TiledGrid) NumTiles() int {
	n := t.Side / t.TileSide
	return n * n
}

// IsCrossing reports whether the edge u->v crosses a tile border.
func (t TiledGrid) IsCrossing(u, v int) bool {
	return t.tileOf(u) != t.tileOf(v)
}

// NewTiledGrid builds a side x side grid map (side must be a multiple of tileSide).
func NewTiledGrid(side, tileSide, maxWeight int, seed int64) TiledGrid {
	return TiledGrid{Side: side, TileSide: tileSide, Graph: NewGridGraph(side, side, maxWeight, seed)}
}

// ExactMultiChipHardware is the global race hardware with chip-to-chip crossings annotated.
type ExactMultiChipHardware struct {
	RaceHardware
	CrossingEdges    int
	CrossingLossDb   float64 // Worst crossing-edge loss: on-chip edge loss + two facets
	CrossingBudgetOK bool
	LinkDelayOK      bool // Unit delay must cover node latency plus link flight time
}

// BuildExactMultiChip programs the global graph and adds inter-chip link effects on crossing edges.
func BuildExactMultiChip(tg TiledGrid, p RaceLogicParams, mc MultiChipParams, seed int64) ExactMultiChipHardware {
	hw := BuildRaceHardware(tg.Graph, p, seed)
	r := rand.New(rand.NewSource(seed + 31))
	jitter := make([][]float64, len(tg.Graph.Adj))
	crossings := 0

	for u, edges := range tg.Graph.Adj {
		jitter[u] = make([]float64, len(edges))
		for i, e := range edges {
			if !tg.IsCrossing(u, e.To) {
				continue
			}
			crossings++
			// Link flight time is compensated in the programmed delay, so the net nominal delay is unchanged;
			// only the link's static error and dynamic jitter remain.
			hw.EdgeErrorPS[u][i] += r.NormFloat64() * mc.LinkStaticErrorPS
			jitter[u][i] = mc.LinkJitterPS
		}
	}
	hw.EdgeJitterPS = jitter

	crossingLoss := hw.WorstEdgeLossDb + 2*mc.CouplingLossDbPerFacet
	return ExactMultiChipHardware{
		RaceHardware:     hw,
		CrossingEdges:    crossings,
		CrossingLossDb:   crossingLoss,
		CrossingBudgetOK: crossingLoss <= hw.PowerMarginDb,
		LinkDelayOK:      p.UnitDelayPS >= p.NodeRegenLatencyPS+mc.LinkDelayPS,
	}
}

// tileSubgraph extracts one tile as a standalone graph with local indices.
func tileSubgraph(tg TiledGrid, tile int) (Graph, map[int]int, []int) {
	toLocal := map[int]int{}
	var toGlobal []int
	for v := range tg.Graph.Adj {
		if tg.tileOf(v) == tile {
			toLocal[v] = len(toGlobal)
			toGlobal = append(toGlobal, v)
		}
	}
	adj := make([][]Edge, len(toGlobal))
	for li, v := range toGlobal {
		for _, e := range tg.Graph.Adj[v] {
			if lj, ok := toLocal[e.To]; ok {
				adj[li] = append(adj[li], Edge{To: lj, Weight: e.Weight})
			}
		}
	}
	return Graph{Adj: adj}, toLocal, toGlobal
}

// hierTile is one chip in hierarchical mode.
type hierTile struct {
	hw        RaceHardware
	toLocal   map[int]int
	toGlobal  []int
	endpoints []int // Global ids of transition endpoints inside this tile
}

// HierarchicalMultiChip holds per-tile hardware and the precomputed abstract graph over transition endpoints.
type HierarchicalMultiChip struct {
	Grid            TiledGrid
	Params          RaceLogicParams
	tiles           []hierTile
	endpointIdx     map[int]int // global node -> abstract index
	endpoints       []int
	abstractAdj     [][]Edge // intra-tile (optically precomputed) + crossing edges
	PrecomputeRaces int
}

// transitionPairs lists the crossing edges chosen as HPA* transitions (every stride-th cell of each border).
func transitionPairs(tg TiledGrid, stride int) [][2]int {
	var pairs [][2]int
	n := tg.Side / tg.TileSide
	offset := stride / 2
	for ty := 0; ty < n; ty++ {
		for tx := 0; tx < n; tx++ {
			for k := offset; k < tg.TileSide; k += stride {
				if tx+1 < n { // vertical border between (tx,ty) and (tx+1,ty)
					x := (tx+1)*tg.TileSide - 1
					y := ty*tg.TileSide + k
					pairs = append(pairs, [2]int{y*tg.Side + x, y*tg.Side + x + 1})
				}
				if ty+1 < n { // horizontal border between (tx,ty) and (tx,ty+1)
					x := tx*tg.TileSide + k
					y := (ty+1)*tg.TileSide - 1
					pairs = append(pairs, [2]int{y*tg.Side + x, (y+1)*tg.Side + x})
				}
			}
		}
	}
	return pairs
}

func edgeWeight(g Graph, u, v int) int {
	for _, e := range g.Adj[u] {
		if e.To == v {
			return e.Weight
		}
	}
	return -1
}

// BuildHierarchical fabricates one chip per tile and precomputes, with noisy optical races, the intra-tile
// distances between transition endpoints (one race per endpoint, done once per map).
func BuildHierarchical(tg TiledGrid, p RaceLogicParams, mc MultiChipParams, seed int64) HierarchicalMultiChip {
	h := HierarchicalMultiChip{Grid: tg, Params: p, endpointIdx: map[int]int{}}
	h.tiles = make([]hierTile, tg.NumTiles())
	for t := range h.tiles {
		sub, toLocal, toGlobal := tileSubgraph(tg, t)
		h.tiles[t] = hierTile{hw: BuildRaceHardware(sub, p, seed+int64(t)*6151), toLocal: toLocal, toGlobal: toGlobal}
	}

	addEndpoint := func(v int) int {
		if i, ok := h.endpointIdx[v]; ok {
			return i
		}
		i := len(h.endpoints)
		h.endpointIdx[v] = i
		h.endpoints = append(h.endpoints, v)
		tile := tg.tileOf(v)
		h.tiles[tile].endpoints = append(h.tiles[tile].endpoints, v)
		return i
	}

	pairs := transitionPairs(tg, mc.TransitionStride)
	type crossing struct{ a, b, w int }
	var crossings []crossing
	for _, pr := range pairs {
		a, b := addEndpoint(pr[0]), addEndpoint(pr[1])
		crossings = append(crossings, crossing{a, b, edgeWeight(tg.Graph, pr[0], pr[1])})
	}

	h.abstractAdj = make([][]Edge, len(h.endpoints))
	for _, c := range crossings {
		h.abstractAdj[c.a] = append(h.abstractAdj[c.a], Edge{To: c.b, Weight: c.w})
		h.abstractAdj[c.b] = append(h.abstractAdj[c.b], Edge{To: c.a, Weight: c.w})
	}

	r := rand.New(rand.NewSource(seed + 97))
	for _, tile := range h.tiles {
		for _, src := range tile.endpoints {
			fire := tile.hw.Race(tile.toLocal[src], r)
			h.PrecomputeRaces++
			for _, dst := range tile.endpoints {
				if dst == src {
					continue
				}
				d := int(math.Round(fire[tile.toLocal[dst]] / p.UnitDelayPS))
				h.abstractAdj[h.endpointIdx[src]] = append(h.abstractAdj[h.endpointIdx[src]], Edge{To: h.endpointIdx[dst], Weight: d})
			}
		}
	}
	return h
}

// Query returns the hierarchical estimate of dist(s, t): two local optical races plus an electronic Dijkstra
// over the abstract graph. It also returns the electronic search time.
func (h HierarchicalMultiChip) Query(s, t int, r *rand.Rand) (int, time.Duration) {
	p := h.Params
	ts, tt := h.Grid.tileOf(s), h.Grid.tileOf(t)
	tileS, tileT := h.tiles[ts], h.tiles[tt]
	fireS := tileS.hw.Race(tileS.toLocal[s], r)
	fireT := tileT.hw.Race(tileT.toLocal[t], r) // grid costs are symmetric: t's race gives dist to t
	decode := func(ps float64) int { return int(math.Round(ps / p.UnitDelayPS)) }

	start := time.Now()
	n := len(h.endpoints)
	srcIdx, dstIdx := n, n+1
	adj := make([][]Edge, n+2)
	copy(adj, h.abstractAdj)
	for _, e := range tileS.endpoints {
		adj[srcIdx] = append(adj[srcIdx], Edge{To: h.endpointIdx[e], Weight: decode(fireS[tileS.toLocal[e]])})
	}
	for _, e := range tileT.endpoints {
		i := h.endpointIdx[e]
		adj[i] = append(append([]Edge(nil), adj[i]...), Edge{To: dstIdx, Weight: decode(fireT[tileT.toLocal[e]])})
	}
	if ts == tt {
		adj[srcIdx] = append(adj[srcIdx], Edge{To: dstIdx, Weight: decode(fireS[tileS.toLocal[t]])})
	}
	d := Dijkstra(Graph{Adj: adj}, srcIdx)[dstIdx]
	return d, time.Since(start)
}

// MultiChipResult compares exact (cascade) and hierarchical composition on one map size.
type MultiChipResult struct {
	Side             int     `json:"side"`
	Tiles            int     `json:"tiles"`
	Queries          int     `json:"queries"`
	CrossingEdges    int     `json:"crossing_edges"`
	CrossingLossDb   float64 `json:"crossing_loss_db"`
	CrossingBudgetOK bool    `json:"crossing_budget_ok"`
	LinkDelayOK      bool    `json:"link_delay_ok"`

	ExactNodeErrorRate float64 `json:"exact_node_error_rate"` // all-destination decoding errors per distance
	ExactNodeUpper95   float64 `json:"exact_node_upper95"`
	MaxHops            int     `json:"max_hops"`
	ExactErrorRate     float64 `json:"exact_error_rate"` // point-to-point distance decoding errors
	ExactErrorUpper95  float64 `json:"exact_error_upper95"`
	ExactMeanCrossings float64 `json:"exact_mean_crossings"` // chip borders crossed by the optimal path
	ExactLatencyNS     float64 `json:"exact_latency_ns"`     // mean time until the target fires + single TDC read
	ExactAllDestNS     float64 `json:"exact_all_dest_ns"`    // full race + parallel per-chip readout

	HierOptimalFrac     float64 `json:"hier_optimal_frac"`
	HierMeanExcessPct   float64 `json:"hier_mean_excess_pct"`
	HierP99ExcessPct    float64 `json:"hier_p99_excess_pct"`
	HierMaxExcessPct    float64 `json:"hier_max_excess_pct"`
	HierBelowOptimal    int     `json:"hier_below_optimal"` // noise-induced underestimates
	HierLatencyNS       float64 `json:"hier_latency_ns"`    // two local races + transition readout + electronic search
	HierPrecomputeRaces int     `json:"hier_precompute_races"`
	Transitions         int     `json:"transitions"`

	DijkstraNS float64 `json:"dijkstra_ns"` // classical point-to-point baseline (full single-source run)
}

// crossingsOnPath counts tile borders crossed by one shortest path from s to t (walked back over the DAG).
func crossingsOnPath(tg TiledGrid, dist []int, t int) int {
	count := 0
	v := t
	for dist[v] > 0 {
		for _, e := range tg.Graph.Adj[v] {
			if dist[e.To] >= 0 && dist[e.To]+e.Weight == dist[v] {
				if tg.IsCrossing(v, e.To) {
					count++
				}
				v = e.To
				break
			}
		}
	}
	return count
}

// SimulateMultiChip runs point-to-point queries on a tiled map in both composition modes.
func SimulateMultiChip(side int, p RaceLogicParams, mc MultiChipParams, queries int, seed int64) MultiChipResult {
	tg := NewTiledGrid(side, mc.TileSide, p.MaxWeight(), seed)
	exact := BuildExactMultiChip(tg, p, mc, seed)
	hier := BuildHierarchical(tg, p, mc, seed)
	r := rand.New(rand.NewSource(seed + 5))
	n := len(tg.Graph.Adj)

	exactErrors, nodeErrors, nodeSamples, maxHops := 0, 0, 0, 0
	var crossSum, exactLatSum, exactAllSum, hierLatSum, dijkstraSum float64
	var excess []float64
	optimal, below := 0, 0
	readOnePS := float64(p.TdcBitsPerNode) / p.ReadoutLinkGbps * 1000
	tileReadoutNS := float64(tg.TileSide*tg.TileSide*p.TdcBitsPerNode) / p.ReadoutLinkGbps
	maxEndpoints := 0
	for _, tl := range hier.tiles {
		if len(tl.endpoints) > maxEndpoints {
			maxEndpoints = len(tl.endpoints)
		}
	}
	endpointReadNS := float64(maxEndpoints*p.TdcBitsPerNode) / p.ReadoutLinkGbps

	for q := 0; q < queries; q++ {
		s, t := r.Intn(n), r.Intn(n)
		for t == s {
			t = r.Intn(n)
		}
		start := time.Now()
		dist := Dijkstra(tg.Graph, s)
		dijkstraSum += float64(time.Since(start).Nanoseconds())
		opt := dist[t]

		fire := exact.Race(s, r)
		if int(math.Round(fire[t]/p.UnitDelayPS)) != opt {
			exactErrors++
		}
		maxFire := 0.0
		for v, f := range fire {
			if !math.IsInf(f, 1) && f > maxFire {
				maxFire = f
			}
			if v != s && dist[v] >= 0 {
				nodeSamples++
				if int(math.Round(f/p.UnitDelayPS)) != dist[v] {
					nodeErrors++
				}
			}
		}
		if q < 50 { // hop depth of the map, sampled from the first sources
			for _, hc := range hopCounts(tg.Graph, dist, s) {
				if hc > maxHops {
					maxHops = hc
				}
			}
		}
		crossSum += float64(crossingsOnPath(tg, dist, t))
		exactLatSum += (fire[t] + readOnePS) / 1000
		exactAllSum += maxFire/1000 + tileReadoutNS

		est, searchTime := hier.Query(s, t, r)
		hierLatSum += float64(searchTime.Nanoseconds()) + endpointReadNS
		switch {
		case est == opt:
			optimal++
			excess = append(excess, 0)
		case est < opt:
			below++
			excess = append(excess, 100*float64(est-opt)/float64(opt))
		default:
			excess = append(excess, 100*float64(est-opt)/float64(opt))
		}
	}

	// Local races run in parallel on the source and target chips; add the slowest tile race span.
	localRaceNS := 0.0
	for _, tl := range hier.tiles {
		fire := tl.hw.Race(0, rand.New(rand.NewSource(seed)))
		for _, f := range fire {
			if !math.IsInf(f, 1) && f/1000 > localRaceNS {
				localRaceNS = f / 1000
			}
		}
	}

	sort.Float64s(excess)
	meanExcess := 0.0
	for _, e := range excess {
		meanExcess += e
	}
	meanExcess /= float64(len(excess))

	return MultiChipResult{
		Side:                side,
		Tiles:               tg.NumTiles(),
		Queries:             queries,
		CrossingEdges:       exact.CrossingEdges,
		CrossingLossDb:      exact.CrossingLossDb,
		CrossingBudgetOK:    exact.CrossingBudgetOK,
		LinkDelayOK:         exact.LinkDelayOK,
		ExactNodeErrorRate:  float64(nodeErrors) / float64(nodeSamples),
		ExactNodeUpper95:    upper95(nodeErrors, nodeSamples),
		MaxHops:             maxHops,
		ExactErrorRate:      float64(exactErrors) / float64(queries),
		ExactErrorUpper95:   upper95(exactErrors, queries),
		ExactMeanCrossings:  crossSum / float64(queries),
		ExactLatencyNS:      exactLatSum / float64(queries),
		ExactAllDestNS:      exactAllSum / float64(queries),
		HierOptimalFrac:     float64(optimal) / float64(queries),
		HierMeanExcessPct:   meanExcess,
		HierP99ExcessPct:    excess[int(0.99*float64(len(excess)-1))],
		HierMaxExcessPct:    excess[len(excess)-1],
		HierBelowOptimal:    below,
		HierLatencyNS:       hierLatSum/float64(queries) + localRaceNS,
		HierPrecomputeRaces: hier.PrecomputeRaces,
		Transitions:         len(hier.endpoints) / 2,
		DijkstraNS:          dijkstraSum / float64(queries),
	}
}
