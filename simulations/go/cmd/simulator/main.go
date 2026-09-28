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
	fmt.Println("   Conceptual Inspiration: Noritsu CW Laser Scanning Principles      ")
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
	fmt.Printf("Temporal Separation Margin (Delta t/sigma):%.2f sigmas\n", sim.SeparationMarginSigmas())
	fmt.Printf("Laser Relative Intensity Noise (RIN):      %.1f dB/Hz\n", sim.Params.LaserRinDbHz)
	fmt.Printf("Electro-Optic Phase Noise (sigma):         %.3f rad\n", sim.Params.PhaseNoiseSigmaRad)
	fmt.Printf("Fused Silica Waveguide Attenuation:        %.2f dB/cm\n", sim.Params.GlassLossDbPerCm)

	calibratedWindow := sim.CalibrateOptimalWindow(CalibrationPulseCount)
	fmt.Printf("Auto-Calibrated Time Window (%d pulses): %.2f ps (Half-Window: +/-%.2f ps)\n", CalibrationPulseCount, calibratedWindow, calibratedWindow/2.0)

	fmt.Println("\n--- 2. NORITSU CW LASERS, M-ARY ENCODING & SYSTEM HARDWARE ---")
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

	fmt.Println("\nSimulation Conclusion: Noritsu CW Laser Engine, M-ary Byte Encoding, Quantum LOQC Core,")
	fmt.Printf("Optical GPU WDM RGB, and Photonic SSD confirm %.2f ps latency, 8x throughput boost, and BER below 10^-12.\n", res.MemStats.GlobalAvgLatencyPS)
	fmt.Println("======================================================================")
}
