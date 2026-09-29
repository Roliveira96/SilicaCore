package optical

import "testing"

func TestEnergyBudgetBreakdownAddsUp(t *testing.T) {
	hw := BuildRaceHardware(NewGridGraph(16, 16, 15, 42), DefaultRaceLogicParams(), 42)
	e := ComputeEnergyBudget(hw, 42.2, DefaultEnergyParams())
	sum := e.LaserElectricalW + e.ReceiversW + e.TdcsW + e.ModulatorsW + e.ReadoutW
	if d := sum - e.TotalW; d > 1e-9 || d < -1e-9 {
		t.Fatalf("Blocks sum to %.6f W, total says %.6f W", sum, e.TotalW)
	}
	if e.EnergyPerQueryNJ <= 0 || e.PowerDensityWPerCm2 <= 0 {
		t.Fatalf("Expected positive energy per query and power density, got %.3f nJ, %.3f W/cm^2", e.EnergyPerQueryNJ, e.PowerDensityWPerCm2)
	}
}

func TestAlwaysOnReceiversDominate(t *testing.T) {
	hw := BuildRaceHardware(NewGridGraph(16, 16, 15, 42), DefaultRaceLogicParams(), 42)
	e := ComputeEnergyBudget(hw, 42.2, DefaultEnergyParams())
	if e.LargestBlock != "receivers" {
		t.Fatalf("Expected the 960 always-on receivers to dominate, got %s", e.LargestBlock)
	}
}

func TestCPUQueryEnergy(t *testing.T) {
	// 17 W package for 69.6 us is 1.18 mJ = 1,183,200 nJ.
	if got := CPUQueryEnergy(17, 69.6); got < 1.18e6 || got > 1.19e6 {
		t.Fatalf("Expected ~1.18e6 nJ, got %.0f", got)
	}
}
