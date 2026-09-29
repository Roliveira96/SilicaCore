package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// Recomputes the time-of-flight gate on the adopted platform and compares it
// with the legacy numbers of the original concept.
//
// Legacy (DefaultParams): fused-silica paths (n = 1.45), 8 ps FWHM laser,
// 25 ps FWHM SPAD, 5 ps TDC -> sigma 11.24 ps, Q = 4.45.
// Platform: Si3N4 spirals (n_g = 2.0), mode-locked/comb laser (1 ps FWHM),
// InGaAs photodiode + comparator with the race-logic node jitter (1.5 ps rms,
// ASSUMPTION) and the same 5 ps TDC. The arrival times t1 and t0 are kept;
// only the path lengths and the noise change.
func main() {
	samples := flag.Int("samples", 1000000, "Monte Carlo pulses per configuration")
	flag.Parse()

	legacy := optical.DefaultParams()

	platform := optical.DefaultParams()
	legacySim := optical.NewSimulator(legacy)
	platform.RefractiveIndex = optical.SiliconNitrideGroupIndex
	scale := legacy.RefractiveIndex / optical.SiliconNitrideGroupIndex
	platform.FastDistanceMM = legacy.FastDistanceMM * scale
	platform.DelayedDistanceMM = legacy.DelayedDistanceMM * scale
	platform.LaserJitterFwhmPS = optical.MicroLaserJitterFwhmPS
	nodeJitterRMS := optical.DefaultRaceLogicParams().NodeJitterPS
	platform.SpadJitterFwhmPS = nodeJitterRMS * 2 * math.Sqrt(2*math.Ln2) // detector field holds the photodiode chain

	fmt.Println("Time-of-flight gate: legacy concept vs adopted platform")
	fmt.Println()
	report("Legacy: fused silica (n = 1.45), SPAD 25 ps FWHM", legacySim, *samples)
	report(fmt.Sprintf("Platform: Si3N4 (n_g = %.1f), InGaAs photodiode %.1f ps rms", optical.SiliconNitrideGroupIndex, nodeJitterRMS),
		optical.NewSimulator(platform), *samples)
}

func report(title string, sim *optical.ToFSimulator, samples int) {
	p := sim.Params
	sigma := sim.TotalSigmaPS
	// Window half-width for a two-tailed 1e-12 error (7.13 sigma), capped at dt/2 - 1 ps.
	halfWindow := math.Min(7.13*sigma, sim.DeltaTNominalPS/2-1)
	sim.Params.WindowWidthPS = 2 * halfWindow
	q := sim.DeltaTNominalPS / (2 * sigma)
	thresholdBER := 0.5 * math.Erfc(q/math.Sqrt2)
	windowBER := math.Erfc(halfWindow / (sigma * math.Sqrt2)) // both tails leave the window

	r := rand.New(rand.NewSource(2026))
	_, failures, _ := sim.TestNOTGateBatch(samples, r)

	minDeltaT := 2 * 7.034 * sigma // Q = 7.03 gives BER 1e-12 with a midpoint threshold
	mmPerPS := optical.SpeedOfLightVacuo / p.RefractiveIndex * 1e-9

	fmt.Println(title)
	fmt.Printf("  paths:                     d1 = %.2f mm, d0 = %.2f mm (t1 = %.2f ps, t0 = %.2f ps)\n",
		p.FastDistanceMM, p.DelayedDistanceMM, sim.T1NominalPS, sim.T0NominalPS)
	fmt.Printf("  sigma total:               %.2f ps\n", sigma)
	fmt.Printf("  Q = dt / 2 sigma:          %.2f\n", q)
	fmt.Printf("  BER, midpoint threshold:   %.2e\n", thresholdBER)
	fmt.Printf("  BER, window +/- %.2f ps:   %.2e (both tails)\n", halfWindow, windowBER)
	fmt.Printf("  Monte Carlo (window):      %d errors in %d pulses (%.2e)\n", failures, samples, float64(failures)/float64(samples))
	fmt.Printf("  slot dt + W:               %.1f ps -> %.2f GHz per channel\n", sim.DeltaTNominalPS+2*halfWindow, 1000/(sim.DeltaTNominalPS+2*halfWindow))
	fmt.Printf("  dt for BER 1e-12:          %.1f ps (%.2f mm of extra path), slot %.1f ps -> %.1f GHz\n",
		minDeltaT, minDeltaT*mmPerPS, minDeltaT+2*halfWindowFor(sigma, minDeltaT), 1000/(minDeltaT+2*halfWindowFor(sigma, minDeltaT)))
	fmt.Println("  (timing-only limits; modulator, detector and TDC bandwidth are not included)")
	fmt.Println()
}

func halfWindowFor(sigma, deltaT float64) float64 {
	return math.Min(7.13*sigma, deltaT/2-1)
}
