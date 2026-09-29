package main

import (
	"flag"
	"fmt"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// Compares the photonic race (simulated RL-16) with a synchronous CMOS race
// logic that solves the same map with one clock cycle per unit of delay.
func main() {
	side := flag.Int("side", 16, "grid side (side x side map)")
	flag.Parse()

	p := optical.DefaultRaceLogicParams()
	g := optical.NewGridGraph(*side, *side, p.MaxWeight(), 42)
	photonic := optical.SimulateRaceLogic(g, p, 200, 7)
	hw := optical.BuildRaceHardware(g, p, 42)
	budget := optical.ComputeEnergyBudget(hw, photonic.TotalQueryTimeNS, optical.DefaultEnergyParams())

	fmt.Printf("Map %dx%d, weights 1..%d\n\n", *side, *side, p.MaxWeight())
	fmt.Printf("Photonic race (simulated): unit %.0f ps, race %.1f ns + readout %.1f ns = %.1f ns, %.0f mm2 of spirals, %.0f nJ per query, %.1f W\n\n",
		p.UnitDelayPS, photonic.RaceSolveTimeNS, photonic.ReadoutTimeNS, photonic.TotalQueryTimeNS, photonic.DelayAreaMM2,
		budget.EnergyPerQueryNJ, budget.TotalW)

	fmt.Println("Synchronous CMOS race logic (exact, one cycle per unit; area and energy are order-of-magnitude ASSUMPTIONS):")
	fmt.Printf("%-7s %8s %10s %14s %15s %12s %12s %12s\n", "clock", "cycles", "race", "+on-die read", "+same link", "flip-flops", "area", "energy")
	for _, f := range []float64{1, 2, 3, 4} {
		r := optical.EstimateCMOSRace(g, p.MaxWeight(), f, p.ReadoutLinkGbps)
		fmt.Printf("%4.0f GHz %8.1f %7.1f ns %11.1f ns %12.1f ns %12d %9.3f mm2 %9.2f nJ\n",
			f, r.MeanCycles, r.RaceNS, r.RaceNS+r.OnChipReadNS, r.RaceNS+r.OffChipReadNS, r.FlipFlops, r.AreaMM2, r.EnergyPerQuery)
	}
	r := optical.EstimateCMOSRace(g, p.MaxWeight(), 3, p.ReadoutLinkGbps)
	fmt.Printf("\nAt 3 GHz: photonic total %.1f ns vs CMOS %.1f ns (on-die read) or %.1f ns (same off-chip link, %d-bit distances)\n",
		photonic.TotalQueryTimeNS, r.RaceNS+r.OnChipReadNS, r.RaceNS+r.OffChipReadNS, r.DistanceBits)
	fmt.Printf("Area ratio photonic/CMOS: ~%.0fx; energy ratio: ~%.0fx\n",
		photonic.DelayAreaMM2/r.AreaMM2, budget.EnergyPerQueryNJ/r.EnergyPerQuery)
}
