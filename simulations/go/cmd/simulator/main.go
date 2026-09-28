package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// ============================================================================
// CLI SIMULATOR EXECUTION PARAMETERS
// Modify the constants below to change simulation sample volumes and iterations.
// ============================================================================
const (
	// CalibrationPulseCount is the number of reference pulses used in initial time-gating window calibration.
	CalibrationPulseCount = 10000

	// MAryTestSymbolCount is the number of dense M-ary optical symbols (0x00..0xFF) tested during transmission.
	MAryTestSymbolCount = 100000

	// AiTensorMeshDimension is the N x N matrix dimension of the Mach-Zehnder Interferometer (MZI) mesh in the Photonic AI Tensor Core.
	AiTensorMeshDimension = 64

	// MonteCarloOperationsCount is the total number of parallel simulation operations executed across Goroutines.
	MonteCarloOperationsCount = 1000000

	// NotGateBatchTestCount is the number of batch test executions performed on the ToF NOT inverter gate.
	NotGateBatchTestCount = 100
)

func main() {
	fmt.Println("======================================================================")
	fmt.Println("   VOLUMETRIC OPTICAL PROCESSOR SIMULATION - SILICA CORE (GOLANG)    ")
	fmt.Println("   Substrate: Fused Silica (SiO2) | Detection: SPAD + TDC            ")
	fmt.Println("   Conceptual Inspiration: Continuous Wave Solid-State Laser Scanning")
	fmt.Println("   Special Acknowledgment: Ricardo Oliveira & Valmor Moreira (Grafis)")
	fmt.Printf("   Parallel Execution on %d CPU Cores (Goroutines)\n", runtime.NumCPU())
	fmt.Println("======================================================================")

	params := optical.DefaultParams()
	sim := optical.NewSimulator(params)

	fmt.Println("\n--- 1. PHYSICAL PARAMETERS & NOMINAL TIMINGS ---")
	fmt.Printf("Glass Refractive Index (n):                %.4f\n", sim.Params.RefractiveIndex)
	fmt.Printf("Speed of Light in Medium (v):              %.5f mm/ps\n", sim.MediumSpeedMMps)
	fmt.Printf("Specific Propagation Delay Rate:           %.4f ps/mm\n", sim.SpecificDelayPSmm)
	fmt.Printf("Direct Path Distance (d1 - Fast Line):     %.3f mm\n", sim.Params.FastDistanceMM)
	fmt.Printf("Deflected Path Distance (d0 - Delay Line): %.3f mm\n", sim.Params.DelayedDistanceMM)
	fmt.Printf("Geometric Path Difference (Delta d):       %.3f mm\n", sim.Params.DelayedDistanceMM-sim.Params.FastDistanceMM)
	fmt.Printf("Nominal Fast Line Time (t1):               %.2f ps\n", sim.T1NominalPS)
	fmt.Printf("Nominal Delay Line Time (t0):              %.2f ps\n", sim.T0NominalPS)
	fmt.Printf("Temporal Separation Difference (Delta t):   %.2f ps\n", sim.DeltaTNominalPS)
	fmt.Printf("Total Convoluted System Jitter (sigma):    %.2f ps\n", sim.TotalSigmaPS)
	fmt.Printf("Temporal Separation (Delta t/sigma):       %.2f sigmas (decision Q = Delta t/2 sigma = %.2f)\n", sim.SeparationMarginSigmas(), sim.SeparationMarginSigmas()/2.0)
	fmt.Printf("Laser Relative Intensity Noise (RIN):      %.1f dB/Hz\n", sim.Params.LaserRinDbHz)
	fmt.Printf("Electro-Optic Phase Noise (sigma):         %.3f rad\n", sim.Params.PhaseNoiseSigmaRad)
	fmt.Printf("Fused Silica Waveguide Attenuation:        %.2f dB/cm\n", sim.Params.GlassLossDbPerCm)

	calibratedWindow := sim.CalibrateOptimalWindow(CalibrationPulseCount)
	fmt.Printf("Auto-Calibrated Time Window (%d pulses): %.2f ps (Half-Window: +/-%.2f ps)\n", CalibrationPulseCount, calibratedWindow, calibratedWindow/2.0)

	fmt.Println("\n--- 2. CONTINUOUS WAVE LASERS, M-ARY ENCODING & SYSTEM HARDWARE ---")
	fmt.Printf("Laser Engine Mode:                         %s\n", sim.Params.LaserEngineMode)
	fmt.Printf("Data Encoding Scheme:                      %s\n", sim.Params.EncodingFormat)
	fmt.Printf("Bits Per Optical Symbol:                   %d bits (%d Discrete States / Symbol)\n", sim.Params.BitsPerSymbol, sim.Params.MultiLevelStatesCount)
	fmt.Printf("L1 Cache Latency (Micro-ring Resonators):  <= %.2f ps  [Alexoudi et al., 2020]\n", sim.Params.CacheL1LatencyPS)
	fmt.Printf("Photonic RAM Latency (Delay-Line Loop):   ~%.2f ps   [Yao, 1993]\n", sim.Params.RamLoopLatencyPS)
	fmt.Printf("Kernel ROM Access (Femtosecond SiO2):      ~%.2f ps (Direct c/n) [Zhang et al., 2014]\n", sim.T1NominalPS)
	fmt.Printf("Photonic Glass SSD Capacity:               %.0f TB / cube [Project Silica / Microsoft]\n", sim.Params.PhotonicSsdCapacityTB)
	fmt.Printf("Photonic Glass SSD Read Throughput:        %.1f TB/s (Parallel WDM)\n", sim.Params.PhotonicSsdThroughputTBps)
	fmt.Printf("Optical GPU WDM Spectral Channels:         %d (Red 635nm, Green 532nm, Blue 450nm) [Weng et al., 2020]\n", sim.Params.GpuWdmChannelsCount)

	fmt.Println("\n--- 3. M-ARY DENSE SYMBOL TRANSMISSION SIMULATION ---")
	mAryRes := sim.SimulateMAryEncoding(MAryTestSymbolCount)
	fmt.Printf("M-ary Byte Symbols Tested:                 %d (0x00..0xFF)\n", mAryRes.TotalSymbols)
	fmt.Printf("Symbol Error Rate (SER):                   %.2e\n", mAryRes.SymbolErrorRate)
	fmt.Printf("Throughput Boost Multiplier:               %.0fx (Byte-level direct delivery)\n", mAryRes.ThroughputBoostX)

	fmt.Println("\n--- 4. ROOM-TEMPERATURE PHOTONIC QUANTUM LOQC SIMULATION (298K) ---")
	quantumRes := sim.SimulateQuantumLOQC(sim.Params.QuantumQubitsCount)
	fmt.Printf("Dual-Rail Photonic Qubits Tested:          %d Qubits\n", quantumRes.QubitsTested)
	fmt.Printf("Waveguide Attenuation Loss:                %.2f dB over 20mm\n", sim.Params.GlassLossDbPerCm*0.2)
	fmt.Printf("Effective HOM 2-Photon Dip Visibility:    %.2f%% [Crespi et al., Nature Phot. 2013]\n", quantumRes.HomVisibilityPct)
	fmt.Printf("Effective CNOT Photonic Gate Fidelity:     %.2f%% [Carolan et al., Science 2015]\n", quantumRes.CnotFidelityPct)
	fmt.Printf("Quantum Gate Bit Error Rate (QBER):        %.2e\n", quantumRes.QuantumBER)

	fmt.Println("\n--- 5. PHOTONIC AI TENSOR CORE SIMULATION (MZI MESH MVM) ---")
	tensorRes := sim.SimulateOpticalTensorEngine(AiTensorMeshDimension)
	fmt.Printf("Matrix-Vector Multiplication Dimension:    %dx%d MZI Mesh\n", tensorRes.MatrixDimension, tensorRes.MatrixDimension)
	fmt.Printf("MZI Phase Drift Error (sigma):            %.3f rad\n", tensorRes.MziPhaseErrorRad)
	fmt.Printf("Photonic AI Tensor Compute Density:        %.0f TOPS/mm^2 [Xu et al., Nature 2021]\n", tensorRes.EffectiveTOPS)
	fmt.Printf("Photonic AI Energy Efficiency:             > %.0f TOPS/W [Shen et al., 2017]\n", tensorRes.EnergyEffTOPSW)
	fmt.Printf("Simulated MVM Precision Accuracy:          %.4f%%\n", tensorRes.MvmAccuracyPct)

	fmt.Println("\n--- 6. MONTE CARLO SIMULATION (1,000,000 CPU & MEMORY OPERATIONS) ---")
	res := sim.SimulateMonteCarloConcurrent(MonteCarloOperationsCount)

	fmt.Printf("Total Samples Tested:                      %d\n", res.TotalSamples)
	fmt.Printf("Detected ToF Logic Errors:                 %d\n", res.DetectedErrors)
	fmt.Printf("ToF Accuracy Rate:                         %.4f%%\n", res.AccuracyRatePct)
	fmt.Printf("Quality Factor Q:                          %.2f\n", res.QFactor)
	fmt.Printf("Empirical Bit Error Rate (BER):            %.2e\n", res.EmpiricalBER)
	fmt.Printf("Theoretical Bit Error Rate (BER):          %.2e\n", res.TheoreticalBER)
	fmt.Printf("L1 Cache Hits:                             %d\n", res.MemStats.CacheL1Hits)
	fmt.Printf("L1 Cache Misses (RAM Loop Accesses):       %d\n", res.MemStats.CacheL1Misses)
	fmt.Printf("Direct Kernel ROM Accesses in SiO2:        %d\n", res.MemStats.RomKernelAccesses)
	fmt.Printf("Global Average Data Latency:               %.2f ps\n", res.MemStats.GlobalAvgLatencyPS)
	fmt.Printf("Go Concurrent Execution Duration:          %s\n", res.ExecutionTime)

	fmt.Println("\n--- 6.1 POWER CONSUMPTION & SILICON COMPARISON BENCHMARK ---")
	pRes := sim.SimulatePowerEfficiency(MonteCarloOperationsCount)
	fmt.Printf("SilicaCore Board Thermal Design Power (TDP): %.1f Watts\n", pRes.SilicaCoreTdpWatts)
	fmt.Printf("Energy Consumed Per Bit Transmitted:        %.2f fJ/bit (EOM + SPAD + CW laser share at 206.75 GHz x 8)\n", pRes.EnergyPerBitFj)
	fmt.Printf("Photonic AI Compute Energy Efficiency:      > %.0f TOPS/W [Shen et al., 2017]\n", pRes.EnergyEfficiencyTOPSW)
	fmt.Printf("Baseline Intel Core i9-14900KS TDP:         %.0f Watts (PL2: 320W)\n", pRes.IntelI9TdpWatts)
	fmt.Printf("SilicaCore Energy Efficiency Ratio vs i9:   %.1fx (TDP ratio; 18.5 W is an assumption, not measured)\n", pRes.IntelI9EfficiencyMult)
	fmt.Printf("Baseline AMD EPYC 9654 Server CPU TDP:      %.0f Watts\n", pRes.AmdEpycTdpWatts)
	fmt.Printf("SilicaCore Energy Efficiency Ratio vs EPYC: %.1fx (TDP ratio; 18.5 W is an assumption, not measured)\n", pRes.AmdEpycEfficiencyMult)
	fmt.Printf("Baseline NVIDIA H100 Tensor GPU TDP:        %.0f Watts\n", pRes.NvidiaH100TdpWatts)
	fmt.Printf("SilicaCore Energy Efficiency Ratio vs H100: %.1fx (TDP ratio; 18.5 W is an assumption, not measured)\n", pRes.NvidiaH100EfficiencyMult)

	fmt.Printf("\n--- 7. CALIBRATED ToF NOT LOGIC GATE TEST (%d ITERATIONS BATCH) ---\n", NotGateBatchTestCount)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	halfWindow := sim.Params.WindowWidthPS / 2.0
	failuresCount := 0

	for i := 0; i < NotGateBatchTestCount; i++ {
		input := i % 2
		output, arrivalTime, ok := sim.TestNOTGate(input, r)

		if !ok || output == -1 {
			failuresCount++
			expectedNominal := sim.T1NominalPS
			if input == 1 {
				expectedNominal = sim.T0NominalPS
			}
			deviation := arrivalTime - expectedNominal
			fmt.Printf("[DEBUG FALHA] Iteração #%03d | Entrada: %d | Saída: %2d | Chegada: %6.2f ps | Esperado: %6.2f ps | Desvio: %+6.2f ps | Janela Máx Permitida: +/-%.2f ps\n",
				i+1, input, output, arrivalTime, expectedNominal, deviation, halfWindow)
		}
	}

	successes := NotGateBatchTestCount - failuresCount
	accuracyPct := (float64(successes) / float64(NotGateBatchTestCount)) * 100.0

	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Total Calibrated Test Runs:                %d Iterations\n", NotGateBatchTestCount)
	fmt.Printf("Successful NOT Inverter Operations:         %d / %d\n", successes, NotGateBatchTestCount)
	fmt.Printf("Failed NOT Inverter Operations:             %d / %d\n", failuresCount, NotGateBatchTestCount)
	fmt.Printf("Calibrated NOT Gate Logic Accuracy:         %.2f%%\n", accuracyPct)

	fmt.Println("\n--- 8. COMPARATIVO DE PERFORMANCE: SILICA CORE VS SILÍCIO CONVENCIONAL ---")
	fmt.Println("Métrica                      | Silício Comercial (6.0 GHz) | SilicaCore (Óptico Fused SiO2)")
	fmt.Println("-----------------------------------------------------------------------------------------")
	fmt.Printf("Tempo de Ciclo / Latência    | ~166.67 ps                  | %.2f ps (~%.1fx mais veloz)\n", res.MemStats.GlobalAvgLatencyPS, 166.67/res.MemStats.GlobalAvgLatencyPS)
	fmt.Printf("Acesso L1 Cache              | ~666.00 ps (4 ciclos)       | %.2f ps (~%.0fx mais veloz)\n", sim.Params.CacheL1LatencyPS, 666.00/sim.Params.CacheL1LatencyPS)
	timing := sim.ComputeTimingBudget(optical.DefaultTargetBER)
	fmt.Printf("Taxa por Canal               | SerDes 112 Gb/s (PAM4)      | %.2f GHz x %d bits = %.1f Gb/s (slot Delta t + W)\n", timing.ToFSymbolRateGHz, sim.Params.BitsPerSymbol, timing.ToFSymbolRateGHz*float64(sim.Params.BitsPerSymbol))
	fmt.Printf("BER da Decisão ToF           | < 1e-15                     | %.2e (Q = %.2f; 1e-12 exige Q = %.2f)\n", timing.TheoreticalBER, timing.QFactor, timing.RequiredQForTarget)
	fmt.Println("Geração Térmica / Fricção    | Altíssima (Efeito Joule)    | Próxima de zero no substrato óptico")
	fmt.Println("Estrutura de Interconexão    | Barramento elétrico de cobre| Guias de onda 3D na velocidade c/n")

	fmt.Println("\n--- 9. ESCALONAMENTO DE DESEMPENHO MASSIVO: MICRO-CUBO (2.0mm) + DWDM MASSIVO (64 CANAIS) ---")
	microParams := optical.MicroCubeParams()
	microSim := optical.NewSimulator(microParams)
	microCalibratedWindow := microSim.CalibrateOptimalWindow(CalibrationPulseCount)
	microMonteCarloRes := microSim.SimulateMonteCarloConcurrent(MonteCarloOperationsCount)

	dwdmBitsPerPulse := microParams.DwdmBitsPerPulse() // 64 channels * 4 bits = 256 bits/pulse
	microTiming := microSim.ComputeTimingBudget(optical.DefaultTargetBER)
	aggregateTbps := microTiming.ToFSymbolRateGHz * float64(dwdmBitsPerPulse) / 1000.0

	fmt.Printf("Micro-Cube Dimension (d1):                 %.1f mm (Redução de 10x na escala física)\n", microSim.Params.FastDistanceMM)
	fmt.Printf("Micro-Cube Direct Nominal Time (t1):       %.2f ps (flight time, not a clock period)\n", microSim.T1NominalPS)
	fmt.Printf("Micro-Cube Global Avg Data Latency:        %.2f ps\n", microMonteCarloRes.MemStats.GlobalAvgLatencyPS)
	fmt.Printf("Auto-Calibrated Micro Window (%d pulses): %.2f ps (Half-Window: +/-%.2f ps)\n", CalibrationPulseCount, microCalibratedWindow, microCalibratedWindow/2.0)
	fmt.Printf("Dense DWDM Spectral Grid:                 %d Wavelength Channels\n", microParams.DwdmChannelsCount)
	fmt.Printf("Parallel Data Density per Optical Pulse:   %d Bits / Pulse (%d channels x %d bits/symbol)\n", dwdmBitsPerPulse, microParams.DwdmChannelsCount, microParams.BitsPerSymbol)
	fmt.Printf("Micro-Cube Decision Q / BER:               %.2f / %.2e (1 ps laser + 3 ps detector jitter: SNSPD-class, cryogenic)\n", microTiming.QFactor, microTiming.TheoreticalBER)
	fmt.Printf("Micro-Cube ToF Symbol Rate per Channel:    %.2f GHz (slot %.1f ps = Delta t + W)\n", microTiming.ToFSymbolRateGHz, microTiming.SymbolSlotPS)
	fmt.Printf("Aggregate Raw Optical Line Rate:           %.2f Tb/s (%d bits x %.2f GHz; transport capacity, not compute)\n", aggregateTbps, dwdmBitsPerPulse, microTiming.ToFSymbolRateGHz)

	printPhysicalBudget(sim)
	printUnifiedMemoryAndLocalAI(sim)

	fmt.Println("\nSimulation Conclusion:")
	fmt.Printf("Default geometry: Q = %.2f (BER %.2e), %.2f GHz per channel. Micro-cube: Q = %.2f (BER %.2e), %.2f GHz per channel,\n", timing.QFactor, timing.TheoreticalBER, timing.ToFSymbolRateGHz, microTiming.QFactor, microTiming.TheoreticalBER, microTiming.ToFSymbolRateGHz)
	fmt.Printf("%.2f Tb/s raw over %d DWDM bits. BER 1e-12 requires total jitter sigma <= %.2f ps (default) / %.2f ps (micro).\n", aggregateTbps, dwdmBitsPerPulse, timing.RequiredSigmaPS, microTiming.RequiredSigmaPS)
	fmt.Println("======================================================================")
}


// printPhysicalBudget reports the routing/mirror link budget, switching limits and the corrected timing budget.
func printPhysicalBudget(sim *optical.ToFSimulator) {
	fmt.Println("\n--- 10. PHYSICAL LINK BUDGET: INTERNAL MIRRORS vs GUIDED PLATFORMS (1550 nm) ---")
	fmt.Printf("Beam radius after delayed line in bulk:    %.0f um (waist %.0f um, no waveguide)\n",
		optical.GaussianBeamRadiusUM(optical.DefaultBeamWaistUM, optical.DefaultTelecomWavelengthNM, sim.Params.RefractiveIndex, sim.Params.DelayedDistanceMM),
		optical.DefaultBeamWaistUM)
	fmt.Printf("Max TIR grazing angle, field-induced dn:   %.4f deg (dn = %.0e, unpoled SiO2 has no Pockels)\n",
		optical.TirGlancingAngleDeg(sim.Params.RefractiveIndex, optical.FusedSilicaMaxFieldDeltaN), optical.FusedSilicaMaxFieldDeltaN)
	fmt.Printf("AOM switching floor (100 um beam):         %.1f ns\n", optical.AomAccessTimePS(100)/1000)

	for _, p := range optical.AllRoutingPlatforms() {
		b := sim.ComputeGateBudget(p)
		fmt.Printf("\n[%s]\n", b.Platform)
		fmt.Printf("  Fast switch:                             %s (rise %.1f ps)\n", p.SwitchTechnology, b.SwitchRiseTimePS)
		fmt.Printf("  Delayed path length:                     %.2f mm\n", b.PathLengthMM)
		fmt.Printf("  Loss per gate:                           %.2f dB (prop %.2f + turns %.2f + diffraction %.2f + switch %.2f)\n",
			b.TotalLossPerGateDb, b.PropagationLossDb, b.TurnsLossDb, b.DiffractionLossDb, b.SwitchLossDb)
		fmt.Printf("  Gates in cascade before regeneration:    %d (power margin %.1f dB)\n", b.MaxCascadedGates, b.PowerMarginDb)
		fmt.Printf("  Thermal phase drift:                     %.2f rad/K\n", b.PhaseDriftRadPerK)
		fmt.Printf("  Feasible as ps logic:                    %v\n", b.Feasible)
	}

	tb := sim.ComputeTimingBudget(optical.DefaultTargetBER)
	fmt.Println("\n--- 10.1 CORRECTED TIMING BUDGET ---")
	fmt.Printf("Q factor (delta_t / 2 sigma):              %.2f -> BER %.2e\n", tb.QFactor, tb.TheoreticalBER)
	fmt.Printf("Q required for BER %.0e:                  %.2f -> max sigma %.2f ps\n", optical.DefaultTargetBER, tb.RequiredQForTarget, tb.RequiredSigmaPS)
	fmt.Printf("ToF symbol slot (delta_t + window):        %.1f ps -> %.2f GHz per channel\n", tb.SymbolSlotPS, tb.ToFSymbolRateGHz)
	fmt.Printf("SPAD max rate (dead time):                 %.2f GHz, %.0f photons/bit at BER target\n", tb.SpadMaxRateGHz, tb.SpadPhotonsPerBit)
	fmt.Printf("UTC photodiode max rate:                   %.0f Gbaud\n", tb.PhotodiodeMaxRateGbaud)
}

// printUnifiedMemoryAndLocalAI reports the unified hierarchy (transport vs cell time) and on-device LLM sizing.
func printUnifiedMemoryAndLocalAI(sim *optical.ToFSimulator) {
	fmt.Println("\n--- 11. UNIFIED MEMORY: LIGHT-SPEED TRANSPORT vs STORAGE-CELL TIME ---")
	fmt.Printf("%-26s %-12s %-14s %-14s %s\n", "Tier", "Transport", "Cell read", "Total", "Capacity")
	for _, m := range optical.UnifiedMemoryHierarchy() {
		fmt.Printf("%-26s %8.1f ps  %11.0f ps  %11.0f ps  %s\n", m.Name, m.TransportPS, m.CellReadPS, m.TotalReadPS, m.Capacity)
	}

	loopBits := optical.DelayLineCapacityBits(optical.DefaultOpticalLineRateGbps, sim.Params.RamLoopLatencyPS, optical.DefaultDwdmChannelsCount)
	length16GB := optical.DelayLineLengthForCapacityM(16*8e9, optical.DefaultOpticalLineRateGbps, optical.DefaultDwdmChannelsCount, optical.SiliconNitrideGroupIndex)
	fmt.Printf("\nDelay-line RAM loop (%.2f ps, %d ch @ %.0f Gbps): %.0f bits held\n",
		sim.Params.RamLoopLatencyPS, optical.DefaultDwdmChannelsCount, optical.DefaultOpticalLineRateGbps, loopBits)
	fmt.Printf("Waveguide needed for 16 GB delay-line RAM: %.0f km\n", length16GB/1000)

	fmt.Println("\n--- 11.1 LOCAL AI: ON-CHIP PCM WEIGHTS vs STREAMING FROM UNIFIED RAM ---")
	for _, m := range []struct {
		params float64
		bits   int
	}{{1, 4}, {8, 4}, {70, 4}} {
		r := optical.SizeLocalAI(m.params, m.bits)
		fmt.Printf("%4.0fB params @ %d-bit: %6.1f GB | PCM area %8.1f cm^2 (fits reticle: %v) | decode bound %6.1f tok/s\n",
			r.ParamsBillions, r.BitsPerWeight, r.ModelGB, r.PcmAreaCM2, r.FitsOnReticle, r.DecodeTokensPerSec)
	}
}
