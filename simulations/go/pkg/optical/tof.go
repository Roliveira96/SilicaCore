package optical

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// EstatisticasMemoria armazena os resultados de simulação da hierarquia de memória óptica.
type EstatisticasMemoria struct {
	HitsCacheL1          int     `json:"hits_cache_l1"`
	MissesCacheL1        int     `json:"misses_cache_l1"`
	AcessosRAMLoop       int     `json:"acessos_ram_loop"`
	AcessosROMKernel     int     `json:"acessos_rom_kernel"`
	LatenciaMediaGlobalPS float64 `json:"latencia_media_global_ps"`
}

// ResultadoMonteCarlo armazena as estatísticas de desempenho da simulação de tempo de voo e memória.
type ResultadoMonteCarlo struct {
	AmostrasTotal    int                 `json:"amostras_total"`
	ErrosDetectados  int                 `json:"erros_detectados"`
	TaxaAcerto       float64             `json:"taxa_acerto_pct"`
	FatorQ           float64             `json:"fator_q"`
	BerEmpirico      float64             `json:"ber_empirico"`
	BerTeorico       float64             `json:"ber_teorico"`
	DuracaoExecucao  string              `json:"duracao_execucao"`
	EstatisticasMem  EstatisticasMemoria `json:"estatisticas_memoria"`
}

// SimularMonteCarloConcorrente executa N disparos de pulso e acessos de memória em goroutines paralelas.
func (s *SimuladorToF) SimularMonteCarloConcorrente(numAmostras int) ResultadoMonteCarlo {
	inicio := time.Now()
	numWorkers := runtime.NumCPU()
	amostrasPorWorker := numAmostras / numWorkers
	sobra := numAmostras % numWorkers

	type workerRes struct {
		erros    int
		hitsL1   int
		missesL1 int
		acessosRAM int
		acessosROM int
		somaLatenciasPS float64
	}

	var wg sync.WaitGroup
	resChan := make(chan workerRes, numWorkers)

	for w := 0; w < numWorkers; w++ {
		totalWorker := amostrasPorWorker
		if w == 0 {
			totalWorker += sobra
		}
		wg.Add(1)

		go func(workerID int, total int) {
			defer wg.Done()
			src := rand.NewSource(time.Now().UnixNano() + int64(workerID*10007))
			r := rand.New(src)

			wRes := workerRes{}
			halfWindow := s.Params.LarguraJanelaPS / 2.0

			for i := 0; i < total; i++ {
				// Simulação de pulso ToF (Linha Rápida vs Linha Atrasada)
				bitEsperado := i % 2
				tempoNominal := s.T1NominalPS
				if bitEsperado == 1 {
					tempoNominal = s.T0NominalPS
				}

				ruido := r.NormFloat64() * s.SigmaTotalPS
				tempoMedido := tempoNominal + ruido
				lsb := s.Params.ResolucaoTDCps
				tempoQuantizado := math.Round(tempoMedido/lsb) * lsb

				if math.Abs(tempoQuantizado-tempoNominal) > halfWindow {
					wRes.erros++
				}

				// Simulação da Hierarquia de Memória Óptica
				// 70% acessos de dados dinâmicos, 30% acessos a instruções do Kernel (ROM)
				isKernelAccess := r.Float64() < 0.30

				if isKernelAccess {
					wRes.acessosROM++
					wRes.somaLatenciasPS += s.T1NominalPS // Leitura direta em SiO2 na velocidade c/n
				} else {
					// Acesso a Dados: Tenta Cache L1 (Micro-anéis <= 5ps)
					if r.Float64() < s.Params.TaxaAcertoCacheL1 {
						wRes.hitsL1++
						wRes.somaLatenciasPS += s.Params.LatenciaCacheL1PS
					} else {
						// Miss na Cache L1 -> Acessa RAM Óptica Recirculante (Delay Loop ~96.73ps)
						wRes.missesL1++
						wRes.acessosRAM++
						wRes.somaLatenciasPS += s.Params.LatenciaRamLoopPS
					}
				}
			}

			resChan <- wRes
		}(w, totalWorker)
	}

	wg.Wait()
	close(resChan)

	totalErros := 0
	totHitsL1 := 0
	totMissesL1 := 0
	totRAM := 0
	totROM := 0
	var somaLatenciasTotal float64

	for r := range resChan {
		totalErros += r.erros
		totHitsL1 += r.hitsL1
		totMissesL1 += r.missesL1
		totRAM += r.acessosRAM
		totROM += r.acessosROM
		somaLatenciasTotal += r.somaLatenciasPS
	}

	duracao := time.Since(inicio)

	berEmpirico := float64(totalErros) / float64(numAmostras)
	fatorQ := s.DeltaTNominalPS / (2.0 * s.SigmaTotalPS)
	berTeorico := 0.5 * math.Erfc(fatorQ/math.Sqrt(2.0))
	taxaAcerto := (1.0 - berEmpirico) * 100.0
	latenciaMedia := somaLatenciasTotal / float64(numAmostras)

	return ResultadoMonteCarlo{
		AmostrasTotal:   numAmostras,
		ErrosDetectados: totalErros,
		TaxaAcerto:      taxaAcerto,
		FatorQ:          fatorQ,
		BerEmpirico:     berEmpirico,
		BerTeorico:      berTeorico,
		DuracaoExecucao: duracao.String(),
		EstatisticasMem: EstatisticasMemoria{
			HitsCacheL1:           totHitsL1,
			MissesCacheL1:         totMissesL1,
			AcessosRAMLoop:        totRAM,
			AcessosROMKernel:      totROM,
			LatenciaMediaGlobalPS: latenciaMedia,
		},
	}
}

// TestarPortaNOT simula uma porta inversora ToF por modulação de trajetória.
func (s *SimuladorToF) TestarPortaNOT(bitEntrada int, r *rand.Rand) (bitSaida int, tempoChegadaPS float64, ok bool) {
	tempoNominal := s.T1NominalPS
	if bitEntrada == 1 {
		tempoNominal = s.T0NominalPS
	}

	ruido := r.NormFloat64() * s.SigmaTotalPS
	tempoChegadaPS = tempoNominal + ruido
	halfWindow := s.Params.LarguraJanelaPS / 2.0

	if math.Abs(tempoChegadaPS-s.T1NominalPS) <= halfWindow {
		bitSaida = 1
		ok = (bitEntrada == 0)
		return
	}

	if math.Abs(tempoChegadaPS-s.T0NominalPS) <= halfWindow {
		bitSaida = 0
		ok = (bitEntrada == 1)
		return
	}

	bitSaida = -1
	ok = false
	return
}
