package optical

import "testing"

func TestDialMatchesDijkstra(t *testing.T) {
	for _, side := range []int{4, 9, 16} {
		g := NewGridGraph(side, side, 15, int64(side))
		st := NewDialState(len(g.Adj), 15)
		for src := 0; src < len(g.Adj); src++ {
			want := Dijkstra(g, src)
			got := DijkstraDial(g, src, st)
			for v := range want {
				if got[v] != want[v] {
					t.Fatalf("side %d src %d node %d: Dial %d, Dijkstra %d", side, src, v, got[v], want[v])
				}
			}
		}
	}
}

func TestAllPairsTableMatchesDijkstra(t *testing.T) {
	g := NewGridGraph(8, 8, 15, 7)
	table := AllPairsTable(g, 15)
	for src := range table {
		want := Dijkstra(g, src)
		for v := range want {
			if int(table[src][v]) != want[v] {
				t.Fatalf("src %d node %d: table %d, Dijkstra %d", src, v, table[src][v], want[v])
			}
		}
	}
}

func TestCMOSRaceCyclesFollowDistances(t *testing.T) {
	g := NewGridGraph(6, 6, 15, 3)
	r := EstimateCMOSRace(g, 15, 2, 100)
	maxEcc := 0
	for s := range g.Adj {
		for _, d := range Dijkstra(g, s) {
			if d > maxEcc {
				maxEcc = d
			}
		}
	}
	if r.MaxCycles != maxEcc+2 {
		t.Fatalf("MaxCycles %d, want %d", r.MaxCycles, maxEcc+2)
	}
	if r.RaceNS <= 0 || r.FlipFlops <= g.NumEdges() {
		t.Fatalf("implausible estimate: %+v", r)
	}
}
