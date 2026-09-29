package optical

import "math"

// Synchronous CMOS race logic, the electronic counterpart of the photonic race.
//
// As in the original race-logic work (Madhavan, Sherwood and Strukov, 2014),
// one unit of delay is one clock cycle: every edge is a small down-counter
// loaded with its weight, every node is a latch that fires once and captures
// the global cycle counter. The race is digital and exact, so there are no
// timing errors; a query takes (eccentricity of the source + 2) cycles.

const (
	// CMOSFlipFlopAreaUM2 is the area of one flip-flop plus its share of logic
	// in a 7 nm-class process (ASSUMPTION, order of magnitude).
	CMOSFlipFlopAreaUM2 = 1.0
	// CMOSFlipFlopEnergyFJ is the energy of one clocked flip-flop per cycle,
	// including clock tree and logic share (ASSUMPTION, order of magnitude).
	CMOSFlipFlopEnergyFJ = 1.0
	// CMOSOnChipBusBits is the width of an on-die bus that reads the distances.
	CMOSOnChipBusBits = 512
)

// CMOSRaceResult summarises the synchronous CMOS race for one graph and clock.
type CMOSRaceResult struct {
	ClockGHz        float64
	MeanCycles      float64
	MaxCycles       int
	RaceNS          float64 // mean race time per query
	OnChipReadNS    float64 // distances read over an on-die bus
	OffChipReadNS   float64 // distances read over the same link as the photonic chip
	FlipFlops       int
	AreaMM2         float64
	EnergyPerQuery  float64 // nJ, every flip-flop clocked every cycle (upper bound)
	DistanceBits    int
	ReadoutLinkGbps float64
}

// EstimateCMOSRace computes latency, size and energy of a synchronous CMOS race
// on graph g for the given clock, averaging over all sources.
func EstimateCMOSRace(g Graph, maxWeight int, clockGHz, readoutLinkGbps float64) CMOSRaceResult {
	n := len(g.Adj)
	st := NewDialState(n, maxWeight)
	sum, maxEcc, maxDist := 0, 0, 0
	for s := 0; s < n; s++ {
		d := DijkstraDial(g, s, st)
		ecc := 0
		for _, v := range d {
			if v > ecc {
				ecc = v
			}
		}
		sum += ecc
		if ecc > maxEcc {
			maxEcc = ecc
		}
		if ecc > maxDist {
			maxDist = ecc
		}
	}
	meanCycles := float64(sum)/float64(n) + 2
	distBits := int(math.Ceil(math.Log2(float64(maxDist + 1))))
	counterBits := int(math.Ceil(math.Log2(float64(maxWeight + 1))))

	// Per edge: down-counter plus an armed bit. Per node: fired latch plus captured distance.
	flipFlops := g.NumEdges()*(counterBits+1) + n*(1+distBits) + distBits
	cycleNS := 1 / clockGHz
	busCycles := math.Ceil(float64(n*distBits) / CMOSOnChipBusBits)

	return CMOSRaceResult{
		ClockGHz:        clockGHz,
		MeanCycles:      meanCycles,
		MaxCycles:       maxEcc + 2,
		RaceNS:          meanCycles * cycleNS,
		OnChipReadNS:    busCycles * cycleNS,
		OffChipReadNS:   float64(n*distBits) / readoutLinkGbps,
		FlipFlops:       flipFlops,
		AreaMM2:         float64(flipFlops) * CMOSFlipFlopAreaUM2 / 1e6,
		EnergyPerQuery:  float64(flipFlops) * CMOSFlipFlopEnergyFJ * meanCycles / 1e6,
		DistanceBits:    distBits,
		ReadoutLinkGbps: readoutLinkGbps,
	}
}
