package optical

import (
	"math/rand"
	"testing"
)

func TestNovoSimulador(t *testing.T) {
	params := ParametrosPadrao()
	sim := NovoSimulador(params)

	if sim.T1NominalPS <= 0 || sim.T0NominalPS <= 0 {
		t.Fatalf("Tempos nominais de propagação inválidos: t1=%f, t0=%f", sim.T1NominalPS, sim.T0NominalPS)
	}

	if sim.DeltaTNominalPS <= 0 {
		t.Fatalf("Separação temporal delta_t deve ser positiva, obtido: %f", sim.DeltaTNominalPS)
	}

	if sim.SigmaTotalPS <= 0 {
		t.Fatalf("Sigma total de jitter deve ser positivo, obtido: %f", sim.SigmaTotalPS)
	}
}

func TestPortaNOT(t *testing.T) {
	sim := NovoSimulador(ParametrosPadrao())
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		ent := i % 2
		saida, _, ok := sim.TestarPortaNOT(ent, r)
		if ok && saida == ent {
			t.Errorf("Falha lógica na porta NOT: entrada=%d, saída=%d", ent, saida)
		}
	}
}

func BenchmarkMonteCarloConcorrente(b *testing.B) {
	sim := NovoSimulador(ParametrosPadrao())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sim.SimularMonteCarloConcorrente(100000)
	}
}
