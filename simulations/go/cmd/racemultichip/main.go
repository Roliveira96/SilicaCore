package main

import (
	"flag"
	"fmt"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// Multi-chip race logic: exact cascade across chips vs hierarchical (HPA*-style) composition of 16x16 tiles.
func main() {
	queries := flag.Int("queries", 3000, "random point-to-point queries per configuration")
	mode := flag.String("mode", "all", "all | exact | hier")
	flag.Parse()

	p := optical.DefaultRaceLogicParams()
	base := optical.DefaultMultiChipParams()
	fmt.Printf("Multi-chip race logic: %dx%d tiles, unit %.0f ps, %d point-to-point queries per configuration\n",
		base.TileSide, base.TileSide, p.UnitDelayPS, *queries)
	fmt.Printf("Inter-chip link: %.1f dB/facet, %.0f ps flight (compensated), %.1f ps jitter, %.1f ps static error (rms)\n\n",
		base.CouplingLossDbPerFacet, base.LinkDelayPS, base.LinkJitterPS, base.LinkStaticErrorPS)

	if *mode != "hier" {
		fmt.Println("EXACT MODE (one optical race across chips; errors counted over all destinations of every race)")
		fmt.Printf("%-7s %-6s %-6s %-8s %-9s %-24s %-11s %-10s %-12s %s\n",
			"Map", "Unit", "Tiles", "MaxHops", "Crossing", "DistErrRate (up95)", "BordersHit", "s->t(ns)", "AllDest(ns)", "Dijkstra(ns)")
		for _, side := range []int{32, 64} {
			for _, unit := range []float64{100, 150} {
				pu := p
				pu.UnitDelayPS = unit
				r := optical.SimulateMultiChip(side, pu, base, *queries, 42)
				fmt.Printf("%-7s %-6.0f %-6d %-8d %-9s %-24s %-11.2f %-10.2f %-12.2f %.0f\n",
					fmt.Sprintf("%dx%d", side, side), unit, r.Tiles, r.MaxHops, fmt.Sprintf("%.2fdB", r.CrossingLossDb),
					fmt.Sprintf("%.2e (%.2e)", r.ExactNodeErrorRate, r.ExactNodeUpper95),
					r.ExactMeanCrossings, r.ExactLatencyNS, r.ExactAllDestNS, r.DijkstraNS)
			}
		}
		for _, facet := range []float64{1.0, 2.0, 2.5} {
			mc := base
			mc.CouplingLossDbPerFacet = facet
			hw := optical.BuildExactMultiChip(optical.NewTiledGrid(32, 16, p.MaxWeight(), 42), p, mc, 42)
			fmt.Printf("  coupling %.1f dB/facet -> crossing edge loss %.2f dB (closes 10 dB margin: %v)\n", facet, hw.CrossingLossDb, hw.CrossingBudgetOK)
		}

	}
	if *mode == "exact" {
		return
	}

	fmt.Println("\nHIERARCHICAL MODE (per-chip optical races + electronic search over border transitions, no chip-to-chip optics)")
	fmt.Printf("%-7s %-7s %-12s %-10s %-11s %-11s %-11s %-8s %-11s %s\n",
		"Map", "Stride", "Transitions", "Optimal", "MeanExc%", "P99Exc%", "MaxExc%", "Below", "Query(ns)", "Precompute")
	for _, side := range []int{32, 64} {
		for _, stride := range []int{8, 4, 2} {
			mc := base
			mc.TransitionStride = stride
			r := optical.SimulateMultiChip(side, p, mc, *queries, 42)
			fmt.Printf("%-7s %-7d %-12d %-10s %-11.2f %-11.2f %-11.2f %-8d %-11.0f %d races\n",
				fmt.Sprintf("%dx%d", side, side), stride, r.Transitions, fmt.Sprintf("%.1f%%", 100*r.HierOptimalFrac),
				r.HierMeanExcessPct, r.HierP99ExcessPct, r.HierMaxExcessPct, r.HierBelowOptimal, r.HierLatencyNS, r.HierPrecomputeRaces)
		}
	}
	fmt.Println("\nQuery(ns) in hierarchical mode includes the electronic abstract-graph search measured in Go on this machine.")
}
