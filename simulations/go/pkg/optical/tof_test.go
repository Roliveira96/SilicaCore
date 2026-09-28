package optical

import (
	"math/rand"
	"testing"
)

func TestNewSimulator(t *testing.T) {
	params := DefaultParams()
	sim := NewSimulator(params)

	if sim.T1NominalPS <= 0 || sim.T0NominalPS <= 0 {
		t.Fatalf("Invalid nominal propagation times: t1=%f, t0=%f", sim.T1NominalPS, sim.T0NominalPS)
	}

	if sim.DeltaTNominalPS <= 0 {
		t.Fatalf("Temporal separation delta_t must be positive, got: %f", sim.DeltaTNominalPS)
	}

	if sim.TotalSigmaPS <= 0 {
		t.Fatalf("Total jitter sigma must be positive, got: %f", sim.TotalSigmaPS)
	}
}

func TestNOTGate(t *testing.T) {
	sim := NewSimulator(DefaultParams())
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		input := i % 2
		output, _, ok := sim.TestNOTGate(input, r)
		if ok && output == input {
			t.Errorf("Logic failure in NOT gate: input=%d, output=%d", input, output)
		}
	}
}

func BenchmarkMonteCarloConcurrent(b *testing.B) {
	sim := NewSimulator(DefaultParams())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sim.SimulateMonteCarloConcurrent(100000)
	}
}
