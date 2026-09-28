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
	CacheL2Hits        int     `json:"cache_l2_hits"`
	UnifiedRamAccesses int     `json:"unified_ram_accesses"`
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
		hitsL2       int
		ramAccesses  int
		sumLatencyPS float64
	}

	// Unified hierarchy latencies: light-speed transport + storage-cell time (see memory.go).
	tiers := UnifiedMemoryHierarchy()
	l1LatencyPS := TransportTimePS(tiers[0].DistanceMM, SiliconNitrideGroupIndex) + s.Params.CacheL1LatencyPS
	l2LatencyPS := tiers[1].TotalReadPS
	ramLatencyPS := tiers[3].TotalReadPS

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

				// Unified memory access: L1 photonic SRAM -> L2/L3 3D SRAM -> unified HBM over optical I/O
				if r.Float64() < s.Params.CacheL1HitRate {
					wRes.hitsL1++
					wRes.sumLatencyPS += l1LatencyPS
				} else {
					wRes.missesL1++
					if r.Float64() < s.Params.CacheL2HitRate {
						wRes.hitsL2++
						wRes.sumLatencyPS += l2LatencyPS
					} else {
						wRes.ramAccesses++
						wRes.sumLatencyPS += ramLatencyPS
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
	totHitsL2 := 0
	totRAM := 0
	var sumLatencyTotal float64

	for r := range resChan {
		totalErrors += r.errors
		totHitsL1 += r.hitsL1
		totMissesL1 += r.missesL1
		totHitsL2 += r.hitsL2
		totRAM += r.ramAccesses
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
			CacheL2Hits:        totHitsL2,
			UnifiedRamAccesses: totRAM,
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

// CalibrateOptimalWindow performs a startup calibration pass over N pulses,
// measuring empirical jitter and setting the optimal time-gating window width W_gate.
func (s *ToFSimulator) CalibrateOptimalWindow(numSamples int) float64 {
	r := rand.New(rand.NewSource(999))
	maxDeviation := 0.0

	for i := 0; i < numSamples; i++ {
		noise := r.NormFloat64() * s.TotalSigmaPS
		dev := math.Abs(noise)
		if dev > maxDeviation {
			maxDeviation = dev
		}
	}

	// Optimal half-window: max observed deviation + 3.0ps margin, safely bounded by deltaT / 2
	halfWindow := maxDeviation + 3.0
	maxHalfWindow := (s.DeltaTNominalPS / 2.0) - 1.0
	if halfWindow > maxHalfWindow {
		halfWindow = maxHalfWindow
	}

	s.Params.WindowWidthPS = halfWindow * 2.0
	return s.Params.WindowWidthPS
}

// TestNOTGateBatch executes N NOT gate inverter tests, returning success/failure counts and overall accuracy.
func (s *ToFSimulator) TestNOTGateBatch(numRuns int, r *rand.Rand) (successes int, failures int, accuracyPct float64) {
	for i := 0; i < numRuns; i++ {
		input := i % 2
		_, _, ok := s.TestNOTGate(input, r)
		if ok {
			successes++
		} else {
			failures++
		}
	}
	accuracyPct = (float64(successes) / float64(numRuns)) * 100.0
	return
}

// MArySymbolResult stores statistics for M-ary dense Byte/Hex symbol transmission over CW lasers.
type MArySymbolResult struct {
	TotalSymbols     int     `json:"total_symbols"`
	SymbolErrors     int     `json:"symbol_errors"`
	SymbolErrorRate  float64 `json:"symbol_error_rate"`
	ThroughputBoostX float64 `json:"throughput_boost_x"`
}

// QuantumLOQCResult stores quantum photon interference dip and gate fidelity simulation stats.
type QuantumLOQCResult struct {
	QubitsTested     int     `json:"qubits_tested"`
	HomVisibilityPct float64 `json:"hom_visibility_pct"`
	CnotFidelityPct  float64 `json:"cnot_fidelity_pct"`
	QuantumBER       float64 `json:"quantum_ber"`
}

// OpticalTensorResult stores photonic AI tensor matrix multiplication simulation stats.
type OpticalTensorResult struct {
	MatrixDimension  int     `json:"matrix_dimension"`
	EffectiveTOPS    float64 `json:"effective_tops"`
	EnergyEffTOPSW   float64 `json:"energy_eff_topsw"`
	MziPhaseErrorRad float64 `json:"mzi_phase_error_rad"`
	MvmAccuracyPct   float64 `json:"mvm_accuracy_pct"`
}

// SimulateMAryEncoding tests dense M-ary 8-bit symbol (Byte) transmission under CW laser RIN noise and phase noise.
func (s *ToFSimulator) SimulateMAryEncoding(numSymbols int) MArySymbolResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	errors := 0
	numStates := s.Params.MultiLevelStatesCount // 256 states for 8-bit Byte

	// Phase state separation angle (2*pi / 256)
	deltaPhase := (2.0 * math.Pi) / float64(numStates)

	for i := 0; i < numSymbols; i++ {
		targetByte := i % numStates
		nominalPhase := float64(targetByte) * deltaPhase

		// Add electro-optic phase noise & RIN fluctuation
		phaseNoise := r.NormFloat64() * s.Params.PhaseNoiseSigmaRad
		measuredPhase := nominalPhase + phaseNoise

		// Normalize phase to [0, 2*pi)
		for measuredPhase < 0 {
			measuredPhase += 2.0 * math.Pi
		}
		measuredPhase = math.Mod(measuredPhase, 2.0*math.Pi)

		// Demodulate to discrete byte symbol
		demodulatedByte := int(math.Round(measuredPhase/deltaPhase)) % numStates

		if demodulatedByte != targetByte {
			errors++
		}
	}

	ser := float64(errors) / float64(numSymbols)
	return MArySymbolResult{
		TotalSymbols:     numSymbols,
		SymbolErrors:     errors,
		SymbolErrorRate:  ser,
		ThroughputBoostX: float64(s.Params.BitsPerSymbol), // 8x for 8-bit Byte
	}
}

// SimulateQuantumLOQC computes two-photon Hong-Ou-Mandel (HOM) quantum interference visibility and CNOT fidelity.
// MODEL ASSUMPTIONS: published HOM visibility and CNOT fidelity scaled by waveguide transmissivity. In real LOQC,
// loss mainly lowers the heralded success rate rather than fidelity; source indistinguishability, multi-photon
// emission and detector efficiency are ignored. Detection requires cryogenic SNSPDs (~1-4 K).
func (s *ToFSimulator) SimulateQuantumLOQC(numQubits int) QuantumLOQCResult {
	// Attenuation loss along 20mm waveguide: Loss = 0.2 dB/cm * 2 cm = 0.4 dB
	lossDb := s.Params.GlassLossDbPerCm * (s.Params.FastDistanceMM / 10.0)
	transmissivity := math.Pow(10.0, -lossDb/10.0) // ~0.912

	// HOM interference dip depth V = V_0 * T (where T is transmissivity)
	actualHomVis := s.Params.QuantumHomVisibilityPct * transmissivity
	actualCnotFid := s.Params.QuantumCnotFidelityPct * transmissivity

	quantumBER := (100.0 - actualCnotFid) / 100.0

	return QuantumLOQCResult{
		QubitsTested:     numQubits,
		HomVisibilityPct: actualHomVis,
		CnotFidelityPct:  actualCnotFid,
		QuantumBER:       quantumBER,
	}
}

// SimulateOpticalTensorEngine simulates Matrix-Vector Multiplication (MVM) on a photonic MZI mesh with phase drift.
func (s *ToFSimulator) SimulateOpticalTensorEngine(matrixDim int) OpticalTensorResult {
	r := rand.New(rand.NewSource(12345))
	totalElements := matrixDim * matrixDim
	errorSum := 0.0

	for i := 0; i < totalElements; i++ {
		// Ideal MZI transmission T = cos^2(theta/2)
		thetaIdeal := r.Float64() * math.Pi
		tIdeal := math.Pow(math.Cos(thetaIdeal/2.0), 2)

		// MZI phase noise error
		phaseError := r.NormFloat64() * s.Params.MziPhaseErrorRad
		tMeasured := math.Pow(math.Cos((thetaIdeal+phaseError)/2.0), 2)

		errorSum += math.Abs(tMeasured - tIdeal)
	}

	avgError := errorSum / float64(totalElements)
	accuracyPct := (1.0 - avgError) * 100.0

	return OpticalTensorResult{
		MatrixDimension:  matrixDim,
		EffectiveTOPS:    s.Params.AiTensorDensityTOPS,
		EnergyEffTOPSW:   s.Params.AiTensorEfficiency,
		MziPhaseErrorRad: s.Params.MziPhaseErrorRad,
		MvmAccuracyPct:   accuracyPct,
	}
}

// PowerEfficiencyResult stores power metrics and comparative silicon benchmark ratios.
type PowerEfficiencyResult struct {
	TotalOperations          int     `json:"total_operations"`
	SilicaCoreTdpWatts       float64 `json:"silica_core_tdp_watts"`
	AggregateBitRateGbps     float64 `json:"aggregate_bit_rate_gbps"`
	EnergyPerBitFj           float64 `json:"energy_per_bit_fj"`
	EnergyEfficiencyTOPSW    float64 `json:"energy_efficiency_topsw"`
	IntelI9TdpWatts          float64 `json:"intel_i9_tdp_watts"`
	IntelI9EfficiencyMult    float64 `json:"intel_i9_efficiency_mult"`
	NvidiaH100TdpWatts       float64 `json:"nvidia_h100_tdp_watts"`
	NvidiaH100EfficiencyMult float64 `json:"nvidia_h100_efficiency_mult"`
	AmdEpycTdpWatts          float64 `json:"amd_epyc_tdp_watts"`
	AmdEpycEfficiencyMult    float64 `json:"amd_epyc_efficiency_mult"`
}

// SimulatePowerEfficiency calculates energy consumption metrics and compares them against silicon baselines.
func (s *ToFSimulator) SimulatePowerEfficiency(numOps int) PowerEfficiencyResult {
	// Energy per bit: E_EOM + E_SPAD + static CW laser power shared over the real aggregate bit rate
	// (ToF symbol rate per channel x bits per symbol x DWDM channels), not over the inverse flight time.
	aggregateGbps := s.ComputeTimingBudget(DefaultTargetBER).ToFSymbolRateGHz * float64(s.Params.BitsPerSymbol) * float64(s.Params.DwdmChannelsCount)
	staticLaserFjPerBit := (s.Params.CwLaserPowerWatts * 1e15) / (aggregateGbps * 1e9)
	totalEnergyFjPerBit := s.Params.EomEnergyFjPerBit + s.Params.SpadEnergyFjPerPhoton + staticLaserFjPerBit

	// Efficiency multipliers vs Silicon CPUs and GPUs
	intelMult := s.Params.IntelI9TdpWatts / s.Params.SilicaCoreTdpWatts
	h100Mult := s.Params.NvidiaH100TdpWatts / s.Params.SilicaCoreTdpWatts
	epycMult := s.Params.AmdEpycTdpWatts / s.Params.SilicaCoreTdpWatts

	return PowerEfficiencyResult{
		TotalOperations:          numOps,
		SilicaCoreTdpWatts:       s.Params.SilicaCoreTdpWatts,
		AggregateBitRateGbps:     aggregateGbps,
		EnergyPerBitFj:           totalEnergyFjPerBit,
		EnergyEfficiencyTOPSW:    s.Params.AiTensorEfficiency,
		IntelI9TdpWatts:          s.Params.IntelI9TdpWatts,
		IntelI9EfficiencyMult:    intelMult,
		NvidiaH100TdpWatts:       s.Params.NvidiaH100TdpWatts,
		NvidiaH100EfficiencyMult: h100Mult,
		AmdEpycTdpWatts:          s.Params.AmdEpycTdpWatts,
		AmdEpycEfficiencyMult:    epycMult,
	}
}

