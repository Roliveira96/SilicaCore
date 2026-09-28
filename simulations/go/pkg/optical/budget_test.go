package optical

import (
	"math"
	"testing"
)

func TestBulkFreeSpaceMirrorsNotFeasible(t *testing.T) {
	sim := NewSimulator(DefaultParams())
	res := sim.ComputeGateBudget(BulkFreeSpaceMirrorsPlatform())

	if res.Feasible {
		t.Fatalf("Free-space mirrors in bulk SiO2 must not be feasible (no ps switch), got feasible")
	}

	if res.DiffractionLossDb < 20.0 {
		t.Fatalf("Expected >20 dB diffraction loss over the delayed line in bulk, got %.2f dB", res.DiffractionLossDb)
	}
}

func TestGuidedPlatformsCascade(t *testing.T) {
	sim := NewSimulator(DefaultParams())
	fs := sim.ComputeGateBudget(FsWaveguideTirPlatform())
	sin := sim.ComputeGateBudget(SiliconNitrideTflnPlatform())

	if !fs.Feasible || !sin.Feasible {
		t.Fatalf("Guided platforms with TFLN switches must cascade at least one gate: fs=%v sin=%v", fs.Feasible, sin.Feasible)
	}

	if sin.MaxCascadedGates <= fs.MaxCascadedGates {
		t.Fatalf("Expected Si3N4 to cascade more gates than fs-glass: sin=%d fs=%d", sin.MaxCascadedGates, fs.MaxCascadedGates)
	}
}

func TestSilicaHasNoUsableFieldInducedTIR(t *testing.T) {
	angle := TirGlancingAngleDeg(DefaultRefractiveIndex, FusedSilicaMaxFieldDeltaN)
	if angle > 0.01 {
		t.Fatalf("Expected TIR grazing angle < 0.01 deg for field-induced dn in silica, got %.4f", angle)
	}
}

func TestAomIsNanosecondClass(t *testing.T) {
	if AomAccessTimePS(100) < 10000 {
		t.Fatalf("Expected AOM access time > 10 ns for a 100 um beam, got %.1f ps", AomAccessTimePS(100))
	}
}

func TestTimingBudgetCorrectsBER(t *testing.T) {
	sim := NewSimulator(DefaultParams())
	res := sim.ComputeTimingBudget(DefaultTargetBER)

	if res.TheoreticalBER < 1e-7 {
		t.Fatalf("Default geometry must not reach 1e-12 (Q=4.45); got BER %.2e", res.TheoreticalBER)
	}

	if math.Abs(res.RequiredQForTarget-7.03) > 0.05 {
		t.Fatalf("Expected Q ~7.03 for BER 1e-12, got %.3f", res.RequiredQForTarget)
	}

	if res.ToFSymbolRateGHz > 10 {
		t.Fatalf("Default ToF slot must limit rate to single-digit GHz, got %.2f GHz", res.ToFSymbolRateGHz)
	}
}

func TestDelayLineRamCannotHoldGigabytes(t *testing.T) {
	bits := 16.0 * 8e9 // 16 GB
	length := DelayLineLengthForCapacityM(bits, DefaultOpticalLineRateGbps, DefaultDwdmChannelsCount, SiliconNitrideGroupIndex)
	if length < 1e6 {
		t.Fatalf("Expected >1000 km of waveguide for 16 GB delay-line RAM, got %.0f m", length)
	}
}

func TestLocalAIDoesNotFitOnChipPCM(t *testing.T) {
	res := SizeLocalAI(8, 4)
	if res.FitsOnReticle {
		t.Fatalf("An 8B-parameter model should not fit on one reticle of PCM cells, got %.2f cm^2", res.PcmAreaCM2)
	}

	if res.DecodeTokensPerSec < 100 {
		t.Fatalf("Expected >100 tokens/s when streaming 4 GB over HBM3-class bandwidth, got %.1f", res.DecodeTokensPerSec)
	}
}
