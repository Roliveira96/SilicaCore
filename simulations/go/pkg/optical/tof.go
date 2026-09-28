package optical

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// MemoryStats holds simulation metrics for the photonic memory hierarchy access pattern.
type MemoryStats struct {
	CacheL1Hits        int     `json:"cache_l1_hits"`
	CacheL1Misses      int     `json:"cache_l1_misses"`
	RamLoopAccesses    int     `json:"ram_loop_accesses"`
	RomKernelAccesses  int     `json:"rom_kernel_accesses"`
	GlobalAvgLatencyPS float64 `json:"global_avg_latency_ps"`
}

// MonteCarloResult stores performance statistics for the ToF simulation and memory hierarchy.
type MonteCarloResult struct {
	TotalSamples    int         `json:"total_samples"`
	DetectedErrors  int         `json:"detected_errors"`
	AccuracyRatePct float64     `json:"accuracy_rate_pct"`
	QFactor         float64     `json:"q_factor"`
	EmpiricalBER    float64     `json:"empirical_ber"`
	TheoreticalBER  float64     `json:"theoretical_ber"`
	ExecutionTime   string      `json:"execution_time"`
	MemStats        MemoryStats `json:"memory_stats"`
}

// SimulateMonteCarloConcurrent executes N pulse firings and memory accesses concurrently using goroutines.
func (s *ToFSimulator) SimulateMonteCarloConcurrent(numSamples int) MonteCarloResult {
	startTime := time.Now()
	numWorkers := runtime.NumCPU()
	samplesPerWorker := numSamples / numWorkers
	remainder := numSamples % numWorkers

	type workerRes struct {
		errors       int
		hitsL1       int
		missesL1     int
		ramAccesses  int
		romAccesses  int
		sumLatencyPS float64
	}

	var wg sync.WaitGroup
	resChan := make(chan workerRes, numWorkers)

	for w := 0; w < numWorkers; w++ {
		totalWorker := samplesPerWorker
		if w == 0 {
			totalWorker += remainder
		}
		wg.Add(1)

		go func(workerID int, total int) {
			defer wg.Done()
			// Per-worker pseudo-random generator to avoid mutex lock contention
			src := rand.NewSource(time.Now().UnixNano() + int64(workerID*10007))
			r := rand.New(src)

			wRes := workerRes{}
			halfWindow := s.Params.WindowWidthPS / 2.0

			for i := 0; i < total; i++ {
				// ToF pulse firing simulation (Fast Line t1 vs Delayed Line t0)
				expectedBit := i % 2
				nominalTime := s.T1NominalPS
				if expectedBit == 1 {
					nominalTime = s.T0NominalPS
				}

				// Box-Muller Gaussian jitter noise
				noise := r.NormFloat64() * s.TotalSigmaPS
				measuredTime := nominalTime + noise
				lsb := s.Params.TdcResolutionPS
				quantizedTime := math.Round(measuredTime/lsb) * lsb

				// Time-gating window check
				if math.Abs(quantizedTime-nominalTime) > halfWindow {
					wRes.errors++
				}

				// Photonic Memory Hierarchy Access Simulation
				// 70% dynamic data accesses, 30% static OS kernel instruction accesses (ROM)
				isKernelAccess := r.Float64() < 0.30

				if isKernelAccess {
					wRes.romAccesses++
					wRes.sumLatencyPS += s.T1NominalPS // Direct speed-of-light readout in SiO2 (c/n)
				} else {
					// Data Access: Attempt L1 Cache (Micro-ring resonators <= 5ps)
					if r.Float64() < s.Params.CacheL1HitRate {
						wRes.hitsL1++
						wRes.sumLatencyPS += s.Params.CacheL1LatencyPS
					} else {
						// L1 Cache Miss -> Access Recirculating Photonic RAM Loop (~96.73ps)
						wRes.missesL1++
						wRes.ramAccesses++
						wRes.sumLatencyPS += s.Params.RamLoopLatencyPS
					}
				}
			}

			resChan <- wRes
		}(w, totalWorker)
	}

	wg.Wait()
	close(resChan)

	totalErrors := 0
	totHitsL1 := 0
	totMissesL1 := 0
	totRAM := 0
	totROM := 0
	var sumLatencyTotal float64

	for r := range resChan {
		totalErrors += r.errors
		totHitsL1 += r.hitsL1
		totMissesL1 += r.missesL1
		totRAM += r.ramAccesses
		totROM += r.romAccesses
		sumLatencyTotal += r.sumLatencyPS
	}

	elapsed := time.Since(startTime)

	empiricalBER := float64(totalErrors) / float64(numSamples)
	qFactor := s.DeltaTNominalPS / (2.0 * s.TotalSigmaPS)
	theoreticalBER := 0.5 * math.Erfc(qFactor/math.Sqrt(2.0))
	accuracyRate := (1.0 - empiricalBER) * 100.0
	globalAvgLatency := sumLatencyTotal / float64(numSamples)

	return MonteCarloResult{
		TotalSamples:    numSamples,
		DetectedErrors:  totalErrors,
		AccuracyRatePct: accuracyRate,
		QFactor:         qFactor,
		EmpiricalBER:    empiricalBER,
		TheoreticalBER:  theoreticalBER,
		ExecutionTime:   elapsed.String(),
		MemStats: MemoryStats{
			CacheL1Hits:        totHitsL1,
			CacheL1Misses:      totMissesL1,
			RamLoopAccesses:    totRAM,
			RomKernelAccesses:  totROM,
			GlobalAvgLatencyPS: globalAvgLatency,
		},
	}
}

// TestNOTGate simulates a ToF trajectory modulation inverter gate.
// If inputBit == 0, signal pulse follows Fast Line (d1 -> t1). Output = 1.
// If inputBit == 1, signal pulse is deflected to Delayed Line (d0 -> t0). Output = 0.
func (s *ToFSimulator) TestNOTGate(inputBit int, r *rand.Rand) (outputBit int, arrivalTimePS float64, ok bool) {
	nominalTime := s.T1NominalPS
	if inputBit == 1 {
		nominalTime = s.T0NominalPS
	}

	noise := r.NormFloat64() * s.TotalSigmaPS
	arrivalTimePS = nominalTime + noise
	halfWindow := s.Params.WindowWidthPS / 2.0

	// Fast Window (t1) -> Bit 1 registered
	if math.Abs(arrivalTimePS-s.T1NominalPS) <= halfWindow {
		outputBit = 1
		ok = (inputBit == 0) // NOT(0) -> 1 (OK)
		return
	}

	// Delayed Window (t0) -> Bit 0 registered
	if math.Abs(arrivalTimePS-s.T0NominalPS) <= halfWindow {
		outputBit = 0
		ok = (inputBit == 1) // NOT(1) -> 0 (OK)
		return
	}

	// Out of time-gating window (Error)
	outputBit = -1
	ok = false
	return
}
