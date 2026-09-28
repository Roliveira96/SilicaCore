package optical

import (
	"math"
	"math/rand"
	"testing"
)

func TestDijkstraOnSmallGraph(t *testing.T) {
	g := Graph{Adj: [][]Edge{
		{{To: 1, Weight: 4}, {To: 2, Weight: 1}},
		{{To: 3, Weight: 1}},
		{{To: 1, Weight: 2}, {To: 3, Weight: 5}},
		{},
	}}
	want := []int{0, 3, 1, 4}
	got := Dijkstra(g, 0)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Dijkstra dist[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestNoiselessRaceEqualsDijkstra(t *testing.T) {
	g := NewGridGraph(12, 12, 15, 7)
	p := DefaultRaceLogicParams()
	p.NodeJitterPS = 0
	p.EdgeDelayErrorPS = 0
	hw := BuildRaceHardware(g, p, 1)
	r := rand.New(rand.NewSource(1))

	for src := 0; src < len(g.Adj); src += 17 {
		dist := Dijkstra(g, src)
		fire := hw.Race(src, r)
		for v := range dist {
			if math.Abs(fire[v]-float64(dist[v])*p.UnitDelayPS) > 1e-6 {
				t.Fatalf("src %d node %d: fire %.3f ps, want %.3f ps", src, v, fire[v], float64(dist[v])*p.UnitDelayPS)
			}
		}
	}
}

func TestNodeLatencyDoesNotPenalizeHops(t *testing.T) {
	// Path 0->1->2 (weight 1+1) must beat direct 0->2 (weight 3) despite one extra regeneration.
	g := Graph{Adj: [][]Edge{
		{{To: 1, Weight: 1}, {To: 2, Weight: 3}},
		{{To: 2, Weight: 1}},
		{},
	}}
	p := DefaultRaceLogicParams()
	p.NodeJitterPS = 0
	p.EdgeDelayErrorPS = 0
	fire := BuildRaceHardware(g, p, 1).Race(0, rand.New(rand.NewSource(1)))
	if math.Abs(fire[2]-2*p.UnitDelayPS) > 1e-6 {
		t.Fatalf("Expected node 2 to fire at 2 units (%.1f ps), got %.1f ps", 2*p.UnitDelayPS, fire[2])
	}
}

func TestDefaultRaceLogicIsAccurate(t *testing.T) {
	g := NewGridGraph(16, 16, 15, 42)
	res := SimulateRaceLogic(g, DefaultRaceLogicParams(), 50, 42)

	if !res.WeightsSupported || !res.LinkBudgetOK {
		t.Fatalf("Default race hardware must support weights and close the link budget: weights=%v link=%v (%.2f dB)",
			res.WeightsSupported, res.LinkBudgetOK, res.WorstEdgeLossDb)
	}

	if res.NodeErrorRate > 1e-3 {
		t.Fatalf("Expected node error rate < 1e-3 at the default unit delay, got %.2e", res.NodeErrorRate)
	}
}

func TestShortUnitDelayCausesErrors(t *testing.T) {
	g := NewGridGraph(16, 16, 15, 42)
	p := DefaultRaceLogicParams()
	p.UnitDelayPS = 21 // barely above node latency, far below the accumulated path jitter margin
	res := SimulateRaceLogic(g, p, 50, 42)

	if res.NodeErrorRate < 1e-3 {
		t.Fatalf("Expected visible decoding errors at 21 ps unit delay, got %.2e", res.NodeErrorRate)
	}
}

func TestSharedModulatorFanoutBreaksLinkBudget(t *testing.T) {
	p := DefaultRaceLogicParams()
	p.PerEdgeModulators = false
	hw := BuildRaceHardware(NewGridGraph(16, 16, 15, 42), p, 1)
	if hw.LinkBudgetOK {
		t.Fatalf("Shared modulator with fan-out 4 should exceed the 10 dB margin, got %.2f dB", hw.WorstEdgeLossDb)
	}
}

func TestLargeGridExceedsReticleArea(t *testing.T) {
	hw := BuildRaceHardware(NewGridGraph(64, 64, 15, 1), DefaultRaceLogicParams(), 1)
	if hw.FitsOnReticle {
		t.Fatalf("A 64x64 grid of full binary delay spirals should exceed one reticle, got %.0f mm^2", hw.DelayAreaMM2)
	}
}

func TestUnitDelayBelowNodeLatencyIsUnsupported(t *testing.T) {
	p := DefaultRaceLogicParams()
	p.UnitDelayPS = p.NodeRegenLatencyPS - 1
	hw := BuildRaceHardware(NewGridGraph(4, 4, 15, 1), p, 1)
	if hw.WeightsSupported {
		t.Fatalf("Unit delay below node regeneration latency must be flagged as unsupported")
	}
}
