package optical

import (
	"math"
)

// SpeedOfLightVacuo is the speed of light in vacuum in meters per second (m/s).
const SpeedOfLightVacuo = 299792458.0

// OpticalParams encapsulates physical, geometric, hardware, memory, GPU, AI, Quantum, and Noritsu CW parameters.
type OpticalParams struct {
	RefractiveIndex   float64 // Effective refractive index of substrate (e.g., 1.4500 for SiO2)
	FastDistanceMM    float64 // d1: Straight-line path distance (mm)
	DelayedDistanceMM float64 // d0: Deflected/extended reflection path distance (mm)
	LaserJitterFwhmPS float64 // Laser pulse temporal jitter in ps (FWHM)
	SpadJitterFwhmPS  float64 // Detector SPAD jitter in ps (FWHM)
	TdcResolutionPS   float64 // Time-to-Digital Converter LSB resolution in ps
	WindowWidthPS     float64 // Sampling time-gating window width in ps

	// Noritsu-Style Continuous Wave (CW) Laser Engine & M-ary Multi-Level Encoding Parameters
	LaserEngineMode       string  // Laser Mode: "Continuous Wave (Noritsu CW Style - Always-ON)"
	EncodingFormat        string  // Data Encoding: "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)"
	BitsPerSymbol         int     // Bits per optical symbol (8 bits = 1 Byte / symbol)
	MultiLevelStatesCount int     // Number of discrete spectral/phase states (256 states for 8-bit Byte)

	// Photonic Memory & Photonic SSD Parameters
	CacheL1LatencyPS          float64 // L1 Photonic Cache latency (Micro-ring resonators): <= 5ps (Alexoudi et al., 2020)
	RamLoopLatencyPS          float64 // Dynamic Photonic RAM latency (Recirculating Delay Loops): ~96.73ps (Yao, 1993)
	CacheL1HitRate            float64 // Nominal empirical hit rate probability for L1 Cache (e.g., 0.92)
	PhotonicSsdThroughputTBps float64 // Photonic Glass SSD parallel read throughput in Terabytes per second (e.g., 1.2 TB/s)
	PhotonicSsdCapacityTB     float64 // Photonic Glass SSD volumetric storage capacity per cube (e.g., 100 TB)

	// Optical GPU & Photonic AI Tensor Engine Parameters
	GpuWdmChannelsCount int     // Number of WDM spectral channels (Red 635nm, Green 532nm, Blue 450nm -> 3)
	AiTensorDensityTOPS float64 // Photonic AI Tensor Core density (e.g., 11 TOPS/mm^2) [Xu et al., Nature 2021]
	AiTensorEfficiency  float64 // Photonic AI Tensor Core energy efficiency (e.g., 100 TOPS/W) [Shen et al., 2017]

	// Photonic Quantum LOQC Core Parameters
	QuantumQubitsCount      int     // Number of room-temperature dual-rail photonic qubits (e.g., 16)
	QuantumHomVisibilityPct float64 // Hong-Ou-Mandel (HOM) 2-photon interference visibility in % (e.g., 99.4%) [Crespi et al., 2013]
	QuantumCnotFidelityPct  float64 // Conditional photonic CNOT gate fidelity in % (e.g., 98.7%) [Carolan et al., Science 2015]

	// Advanced Physical Noise & Attenuation Parameters
	LaserRinDbHz       float64 // Laser Relative Intensity Noise (RIN) in dB/Hz (e.g. -155.0 dB/Hz for stabilized CW)
	PhaseNoiseSigmaRad float64 // Electro-optic phase noise standard deviation in radians (e.g. 0.012 rad)
	GlassLossDbPerCm   float64 // Propagation attenuation in fused silica (e.g. 0.2 dB/cm)
	MziPhaseErrorRad   float64 // MZI mesh phase drift standard deviation in radians (e.g. 0.010 rad)
}

// DefaultParams returns the nominal configuration for ToF CPU, Photonic SSD, Optical GPU, AI Tensor, LOQC, and Noritsu CW.
func DefaultParams() OpticalParams {
	return OpticalParams{
		RefractiveIndex:           1.4500,
		FastDistanceMM:            20.0,
		DelayedDistanceMM:         40.675,
		LaserJitterFwhmPS:         8.0,
		SpadJitterFwhmPS:          25.0,
		TdcResolutionPS:           5.0,
		WindowWidthPS:             35.0,
		LaserEngineMode:           "Continuous Wave (Noritsu CW Style - Always-ON)",
		EncodingFormat:            "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)",
		BitsPerSymbol:             8,
		MultiLevelStatesCount:     256,
		CacheL1LatencyPS:          5.0,
		RamLoopLatencyPS:          96.73,
		CacheL1HitRate:            0.92,
		PhotonicSsdThroughputTBps: 1.2,
		PhotonicSsdCapacityTB:     100.0,
		GpuWdmChannelsCount:       3,     // RGB Channels
		AiTensorDensityTOPS:       11.0,  // 11 TOPS/mm^2
		AiTensorEfficiency:        100.0, // 100 TOPS/W
		QuantumQubitsCount:        16,    // 16 Dual-Rail Photonic Qubits
		QuantumHomVisibilityPct:   99.4,  // 99.4% HOM Visibility
		QuantumCnotFidelityPct:    98.7,  // 98.7% CNOT Fidelity
		LaserRinDbHz:              -155.0,// Stabilized CW laser RIN (-155 dB/Hz)
		PhaseNoiseSigmaRad:        0.012, // Phase jitter (12 mrad)
		GlassLossDbPerCm:          0.2,   // Waveguide loss (0.2 dB/cm)
		MziPhaseErrorRad:          0.010, // MZI Mesh phase error (10 mrad)
	}
}

// FwhmToSigma converts Full Width at Half Maximum (FWHM) of a Gaussian distribution to standard deviation (sigma).
func FwhmToSigma(fwhm float64) float64 {
	return fwhm / (2.0 * math.Sqrt(2.0*math.Log(2.0)))
}

// ToFSimulator represents the core simulation engine holding physical parameters and noise model.
type ToFSimulator struct {
	Params            OpticalParams
	MediumSpeedMMps   float64 // Speed of light in the substrate medium in mm/ps
	SpecificDelayPSmm float64 // Propagation delay rate in ps/mm
	T1NominalPS       float64 // Fast Line nominal propagation time (ps)
	T0NominalPS       float64 // Delayed Line nominal propagation time (ps)
	DeltaTNominalPS   float64 // Temporal difference t0 - t1 (ps)
	TotalSigmaPS      float64 // Total convoluted Gaussian jitter (ps)
}

// NewSimulator initializes the simulation engine with the provided optical parameters.
func NewSimulator(params OpticalParams) *ToFSimulator {
	vMedium := (SpeedOfLightVacuo / params.RefractiveIndex) * 1e-9
	delaySpec := 1.0 / vMedium

	t1 := params.FastDistanceMM * delaySpec
	t0 := params.DelayedDistanceMM * delaySpec
	deltaT := t0 - t1

	sigmaLaser := FwhmToSigma(params.LaserJitterFwhmPS)
	sigmaSpad := FwhmToSigma(params.SpadJitterFwhmPS)
	sigmaTDC := params.TdcResolutionPS / math.Sqrt(12.0)

	sigmaTotal := math.Sqrt(sigmaLaser*sigmaLaser + sigmaSpad*sigmaSpad + sigmaTDC*sigmaTDC)

	return &ToFSimulator{
		Params:            params,
		MediumSpeedMMps:   vMedium,
		SpecificDelayPSmm: delaySpec,
		T1NominalPS:       t1,
		T0NominalPS:       t0,
		DeltaTNominalPS:   deltaT,
		TotalSigmaPS:      sigmaTotal,
	}
}

// SeparationMarginSigmas calculates the temporal separation ratio in multiples of total sigma.
func (s *ToFSimulator) SeparationMarginSigmas() float64 {
	return s.DeltaTNominalPS / s.TotalSigmaPS
}
