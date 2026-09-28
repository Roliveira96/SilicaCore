package optical

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// ResultadoMonteCarlo armazena as estatísticas de desempenho da simulação de tempo de voo.
type ResultadoMonteCarlo struct {
	AmostrasTotal    int     `json:"amostras_total"`
	ErrosDetectados  int     `json:"erros_detectados"`
	TaxaAcerto       float64 `json:"taxa_acerto_pct"`
	FatorQ           float64 `json:"fator_q"`
	BerEmpirico      float64 `json:"ber_empirico"`
	BerTeorico       float64 `json:"ber_teorico"`
	DuracaoExecucao  string  `json:"duracao_execucao"`
}

// SimularMonteCarloConcorrente executa N disparos de pulso distribuídos paralelamente em goroutines.
func (s *SimuladorToF) SimularMonteCarloConcorrente(numAmostras int) ResultadoMonteCarlo {
	inicio := time.Now()
	numWorkers := runtime.NumCPU()
	amostrasPorWorker := numAmostras / numWorkers
	sobra := numAmostras % numWorkers

	var wg sync.WaitGroup
	errosChan := make(chan int, numWorkers)

	for w := 0; w < numWorkers; w++ {
		totalWorker := amostrasPorWorker
		if w == 0 {
			totalWorker += sobra
		}
		wg.Add(1)

		go func(workerID int, total int) {
			defer wg.Done()
			// Gerador aleatório local por worker para evitar lock contention
			src := rand.NewSource(time.Now().UnixNano() + int64(workerID*10007))
			r := rand.New(src)

			errosLocais := 0
			halfWindow := s.Params.LarguraJanelaPS / 2.0

			for i := 0; i < total; i++ {
				// Alterna simulação de bit 0 (t1 - Linha Rápida) e bit 1 (t0 - Linha Atrasada)
				bitEsperado := i % 2
				tempoNominal := s.T1NominalPS
				if bitEsperado == 1 {
					tempoNominal = s.T0NominalPS
				}

				// Ruído Gaussiano Box-Muller
				ruido := r.NormFloat64() * s.SigmaTotalPS
				tempoMedido := tempoNominal + ruido

				// Aplica quantização TDC LSB
				lsb := s.Params.ResolucaoTDCps
				tempoQuantizado := math.Round(tempoMedido/lsb) * lsb

				// Verificação por janelamento temporal (Time-Gating)
				if math.Abs(tempoQuantizado-tempoNominal) > halfWindow {
					errosLocais++
				}
			}

			errosChan <- errosLocais
		}(w, totalWorker)
	}

	wg.Wait()
	close(errosChan)

	totalErros := 0
	for e := range errosChan {
		totalErros += e
	}

	duracao := time.Since(inicio)

	berEmpirico := float64(totalErros) / float64(numAmostras)
	fatorQ := s.DeltaTNominalPS / (2.0 * s.SigmaTotalPS)

	// BER Teórico para modulação de 2 estados em canal Gaussiano: 0.5 * erfc(Q / sqrt(2))
	berTeorico := 0.5 * math.Erfc(fatorQ/math.Sqrt(2.0))

	taxaAcerto := (1.0 - berEmpirico) * 100.0

	return ResultadoMonteCarlo{
		AmostrasTotal:   numAmostras,
		ErrosDetectados: totalErros,
		TaxaAcerto:      taxaAcerto,
		FatorQ:          fatorQ,
		BerEmpirico:     berEmpirico,
		BerTeorico:      berTeorico,
		DuracaoExecucao: duracao.String(),
	}
}

// TestarPortaNOT simula uma porta inversora ToF por modulação de trajetória.
// Se bitEntrada == 0, o feixe de sinal segue pela Linha Rápida (d1 -> t1). Saída = 1.
// Se bitEntrada == 1, o feixe é defletido para a Linha Atrasada (d0 -> t0). Saída = 0.
func (s *SimuladorToF) TestarPortaNOT(bitEntrada int, r *rand.Rand) (bitSaida int, tempoChegadaPS float64, ok bool) {
	tempoNominal := s.T1NominalPS
	if bitEntrada == 1 {
		tempoNominal = s.T0NominalPS
	}

	ruido := r.NormFloat64() * s.SigmaTotalPS
	tempoChegadaPS = tempoNominal + ruido

	halfWindow := s.Params.LarguraJanelaPS / 2.0

	// Janela Rápida (t1) -> Leitura de bit 1
	if math.Abs(tempoChegadaPS-s.T1NominalPS) <= halfWindow {
		bitSaida = 1
		ok = (bitEntrada == 0) // NOT(0) -> 1 (OK)
		return
	}

	// Janela Atrasada (t0) -> Leitura de bit 0
	if math.Abs(tempoChegadaPS-s.T0NominalPS) <= halfWindow {
		bitSaida = 0
		ok = (bitEntrada == 1) // NOT(1) -> 0 (OK)
		return
	}

	// Fora da janela de janelamento (Erro)
	bitSaida = -1
	ok = false
	return
}
