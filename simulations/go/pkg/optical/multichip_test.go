package optical

import (
	"math"
	"math/rand"
	"testing"
)

func noiselessRaceParams() RaceLogicParams {
	p := DefaultRaceLogicParams()
	p.NodeJitterPS = 0
	p.EdgeDelayErrorPS = 0
	return p
}

func noiselessLinks() MultiChipParams {
	mc := DefaultMultiChipParams()
	mc.LinkJitterPS = 0
	mc.LinkStaticErrorPS = 0
	return mc
}

func TestExactMultiChipNoiselessEqualsDijkstra(t *testing.T) {
	tg := NewTiledGrid(32, 16, 15, 3)
	p := noiselessRaceParams()
	hw := BuildExactMultiChip(tg, p, noiselessLinks(), 3)
	r := rand.New(rand.NewSource(1))
	for src := 0; src < len(tg.Graph.Adj); src += 97 {
		dist := Dijkstra(tg.Graph, src)
		fire := hw.Race(src, r)
		for v := range dist {
			if math.Abs(fire[v]-float64(dist[v])*p.UnitDelayPS) > 1e-6 {
				t.Fatalf("src %d node %d: fire %.3f ps, want %.3f ps", src, v, fire[v], float64(dist[v])*p.UnitDelayPS)
			}
		}
	}
	if hw.CrossingEdges != 2*2*32 {
		t.Fatalf("Expected %d directed crossing edges on a 2x2-tile 32x32 grid, got %d", 2*2*32, hw.CrossingEdges)
	}
}

func TestHierarchicalWithAllTransitionsIsExact(t *testing.T) {
	// With every border cell as a transition and no noise, the abstract graph preserves all shortest paths.
	tg := NewTiledGrid(32, 16, 15, 5)
	mc := noiselessLinks()
	mc.TransitionStride = 1
	h := BuildHierarchical(tg, noiselessRaceParams(), mc, 5)
	r := rand.New(rand.NewSource(2))
	for q := 0; q < 200; q++ {
		s, d := r.Intn(len(tg.Graph.Adj)), r.Intn(len(tg.Graph.Adj))
		want := Dijkstra(tg.Graph, s)[d]
		if got, _ := h.Query(s, d, r); got != want {
			t.Fatalf("query %d->%d: hierarchical %d, exact %d", s, d, got, want)
		}
	}
}

func TestHierarchicalSparseTransitionsNeverBeatOptimal(t *testing.T) {
	tg := NewTiledGrid(32, 16, 15, 7)
	h := BuildHierarchical(tg, noiselessRaceParams(), noiselessLinks(), 7)
	r := rand.New(rand.NewSource(3))
	for q := 0; q < 200; q++ {
		s, d := r.Intn(len(tg.Graph.Adj)), r.Intn(len(tg.Graph.Adj))
		want := Dijkstra(tg.Graph, s)[d]
		if got, _ := h.Query(s, d, r); got < want {
			t.Fatalf("noiseless hierarchical estimate %d below optimum %d", got, want)
		}
	}
}

func TestCrossingLinkBudget(t *testing.T) {
	tg := NewTiledGrid(32, 16, 15, 1)
	p := DefaultRaceLogicParams()
	mc := DefaultMultiChipParams()
	mc.CouplingLossDbPerFacet = 2.0
	if hw := BuildExactMultiChip(tg, p, mc, 1); !hw.CrossingBudgetOK {
		t.Fatalf("2 dB per facet should still close the 10 dB margin, got %.2f dB", hw.CrossingLossDb)
	}
	mc.CouplingLossDbPerFacet = 2.5
	if hw := BuildExactMultiChip(tg, p, mc, 1); hw.CrossingBudgetOK {
		t.Fatalf("2.5 dB per facet should exceed the 10 dB margin, got %.2f dB", hw.CrossingLossDb)
	}
}
