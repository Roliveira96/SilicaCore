package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// Race-logic statistical campaign: many fabricated chips x many random-source queries per unit delay.
func main() {
	side := flag.Int("side", 16, "grid side (side x side game map)")
	chips := flag.Int("chips", 10, "independently fabricated chips (fresh static edge errors)")
	queries := flag.Int("queries", 10000, "random-source queries per chip")
	csvPath := flag.String("csv", "", "optional CSV output with per-hop measured vs predicted error")
	flag.Parse()

	p := optical.DefaultRaceLogicParams()
	g := optical.NewGridGraph(*side, *side, p.MaxWeight(), 42)
	fmt.Printf("Race-logic statistics: %dx%d grid, %d chips x %d queries = %d queries per unit delay\n",
		*side, *side, *chips, *queries, *chips**queries)
	fmt.Printf("Noise model: node jitter %.1f ps rms, static edge error %.1f ps rms per chip\n\n", p.NodeJitterPS, p.EdgeDelayErrorPS)

	fmt.Printf("%-9s %-7s %-12s %-10s %-12s %-12s %-12s %-10s %-12s %s\n",
		"Unit(ps)", "Q", "Distances", "NodeErr", "NodeRate", "Upper95", "Predicted", "QueryErr", "QueryUp95", "Time")

	var w *csv.Writer
	if *csvPath != "" {
		f, err := os.Create(*csvPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		w = csv.NewWriter(f)
		defer w.Flush()
		_ = w.Write([]string{"unit_ps", "hops", "samples", "errors", "measured_rate", "predicted_rate"})
	}

	for _, unit := range []float64{35, 50, 75, 100} {
		pu := p
		pu.UnitDelayPS = unit
		start := time.Now()
		r := optical.RaceLogicStatistics(g, pu, *chips, *queries, 42)
		fmt.Printf("%-9.0f %-7.2f %-12d %-10d %-12.2e %-12.2e %-12.2e %-10d %-12.2e %s\n",
			unit, r.WorstPathQ, r.NodeSamples, r.NodeErrors, r.NodeErrorRate, r.NodeErrorUpper95,
			r.PredictedNodeRate, r.QueryErrors, r.QueryErrorUpper95, time.Since(start).Round(time.Millisecond))

		if w != nil {
			for _, b := range r.ByHops {
				_ = w.Write([]string{strconv.FormatFloat(unit, 'f', 0, 64), strconv.Itoa(b.Hops), strconv.Itoa(b.Samples),
					strconv.Itoa(b.Errors), strconv.FormatFloat(b.MeasuredRate, 'e', 3, 64), strconv.FormatFloat(b.PredictedRate, 'e', 3, 64)})
			}
		}
	}
}
