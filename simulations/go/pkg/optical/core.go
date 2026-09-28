package optical

import (
	"math"
)

// ============================================================================
// PHYSICAL CONSTANTS & DEFAULT SIMULATION CONFIGURATIONS (SILICA CORE)
// Modify the constant values below to adjust physical processor parameters.
// ============================================================================

const (
	// SpeedOfLightVacuo is the speed of light in vacuum in meters per second (m/s).
	SpeedOfLightVacuo = 299792458.0

	// ------------------------------------------------------------------------
	// 1. SUBSTRATE & PROPAGATION PHYSICAL PARAMETERS (ToF LOGIC)
	// ------------------------------------------------------------------------

	// DefaultRefractiveIndex is the effective refractive index of fused silica (SiO2).
	DefaultRefractiveIndex = 1.4500

	// DefaultFastDistanceMM is the straight-line propagation path distance d1 (mm).
	DefaultFastDistanceMM = 20.0

	// DefaultDelayedDistanceMM is the deflected/extended propagation path distance d0 (mm).
	DefaultDelayedDistanceMM = 40.675

	// DefaultLaserJitterFwhmPS is the laser pulse temporal jitter in picoseconds (FWHM).
	DefaultLaserJitterFwhmPS = 8.0

	// DefaultSpadJitterFwhmPS is the single-photon avalanche diode (SPAD) detector jitter in picoseconds (FWHM).
	DefaultSpadJitterFwhmPS = 25.0

	// DefaultTdcResolutionPS is the Time-to-Digital Converter (TDC) LSB resolution in picoseconds.
	DefaultTdcResolutionPS = 5.0

	// DefaultWindowWidthPS is the sampling time-gating window width in picoseconds.
	DefaultWindowWidthPS = 75.0

	// ------------------------------------------------------------------------
	// 2. CONTINUOUS WAVE LASER ENGINE (NORITSU CW) & M-ARY DENSE ENCODING
	// ------------------------------------------------------------------------

	// DefaultLaserEngineMode is the operational mode of the continuous wave lasers.
	DefaultLaserEngineMode = "Continuous Wave (Noritsu CW Style - Always-ON)"

	// DefaultEncodingFormat is the spectral/phase data encoding scheme.
	DefaultEncodingFormat = "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)"

	// DefaultBitsPerSymbol is the number of bits encoded per optical symbol (8 bits = 1 Byte).
	DefaultBitsPerSymbol = 8

	// DefaultMultiLevelStatesCount is the number of discrete optical states per symbol (2^8 = 256).
	DefaultMultiLevelStatesCount = 256

	// ------------------------------------------------------------------------
	// 3. PHOTONIC MEMORY HIERARCHY & PHOTONIC GLASS SSD
	// ------------------------------------------------------------------------

	// DefaultCacheL1LatencyPS is the L1 Photonic Cache latency (micro-ring resonators) in ps [Alexoudi et al., 2020].
	DefaultCacheL1LatencyPS = 5.0

	// DefaultRamLoopLatencyPS is the Photonic RAM recirculating delay loop latency in ps [Yao, 1993].
	DefaultRamLoopLatencyPS = 96.73

	// DefaultCacheL1HitRate is the nominal empirical hit rate probability for L1 Photonic Cache.
	DefaultCacheL1HitRate = 0.92

	// DefaultPhotonicSsdCapacityTB is the storage capacity per fused silica cube in Terabytes [Project Silica / Microsoft].
	DefaultPhotonicSsdCapacityTB = 100.0

	// DefaultPhotonicSsdThroughputTBps is the parallel WDM read throughput in Terabytes per second (TB/s).
	DefaultPhotonicSsdThroughputTBps = 1.2

	// ------------------------------------------------------------------------
	// 4. OPTICAL GPU (WDM RGB) & PHOTONIC AI TENSOR ENGINE (MZI MESH)
	// ------------------------------------------------------------------------

	// DefaultGpuWdmChannelsCount is the number of WDM spectral channels (Red 635nm, Green 532nm, Blue 450nm).
	DefaultGpuWdmChannelsCount = 3

	// DefaultAiTensorDensityTOPS is the Photonic AI Tensor Core compute density in TOPS/mm^2 [Xu et al., Nature 2021].
	DefaultAiTensorDensityTOPS = 11.0

	// DefaultAiTensorEfficiency is the Photonic AI Tensor Core energy efficiency in TOPS/W [Shen et al., Nature Phot. 2017].
	DefaultAiTensorEfficiency = 100.0

	// ------------------------------------------------------------------------
	// 5. PHOTONIC QUANTUM PROCESSOR (ROOM-TEMPERATURE LOQC - 298K)
	// ------------------------------------------------------------------------

	// DefaultQuantumQubitsCount is the number of room-temperature dual-rail photonic qubits (298K).
	DefaultQuantumQubitsCount = 16

	// DefaultQuantumHomVisibilityPct is the Hong-Ou-Mandel 2-photon interference dip visibility in % [Crespi et al., 2013].
	DefaultQuantumHomVisibilityPct = 99.4

	// DefaultQuantumCnotFidelityPct is the photonic CNOT logic gate fidelity in % [Carolan et al., Science 2015].
	DefaultQuantumCnotFidelityPct = 98.7

	// ------------------------------------------------------------------------
	// 6. PHYSICAL NOISE MODEL & WAVEGUIDE PROPAGATION ATTENUATION
	// ------------------------------------------------------------------------

	// DefaultLaserRinDbHz is the continuous wave laser Relative Intensity Noise (RIN) in dB/Hz.
	DefaultLaserRinDbHz = -155.0

	// DefaultPhaseNoiseSigmaRad is the electro-optic phase noise standard deviation in radians.
	DefaultPhaseNoiseSigmaRad = 0.012

	// DefaultGlassLossDbPerCm is the optical propagation loss in fused silica waveguide in dB/cm.
	DefaultGlassLossDbPerCm = 0.2

	// DefaultMziPhaseErrorRad is the MZI mesh phase drift standard deviation in radians.
	DefaultMziPhaseErrorRad = 0.010
)

// OpticalParams encapsulates all physical, geometric, and functional parameters of SilicaCore.
type OpticalParams struct {
	RefractiveIndex   float64 // Effective refractive index of substrate (e.g., 1.4500 for SiO2)
	FastDistanceMM    float64 // d1: Straight-line path distance (mm)
	DelayedDistanceMM float64 // d0: Deflected/extended reflection path distance (mm)
	LaserJitterFwhmPS float64 // Laser pulse temporal jitter in ps (FWHM)
	SpadJitterFwhmPS  float64 // Detector SPAD jitter in ps (FWHM)
	TdcResolutionPS   float64 // Time-to-Digital Converter LSB resolution in ps
	WindowWidthPS     float64 // Sampling time-gating window width in ps

	// Noritsu-Style Continuous Wave (CW) Laser Engine & M-ary Multi-Level Encoding
	LaserEngineMode       string // Laser Mode: "Continuous Wave (Noritsu CW Style - Always-ON)"
	EncodingFormat        string // Data Encoding: "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)"
	BitsPerSymbol         int    // Bits per optical symbol (8 bits = 1 Byte / symbol)
	MultiLevelStatesCount int    // Number of discrete spectral/phase states (256 states for 8-bit Byte)

	// Photonic Memory & Photonic SSD
	CacheL1LatencyPS          float64 // L1 Photonic Cache latency (Micro-ring resonators): <= 5ps (Alexoudi et al., 2020)
	RamLoopLatencyPS          float64 // Dynamic Photonic RAM latency (Recirculating Delay Loops): ~96.73ps (Yao, 1993)
	CacheL1HitRate            float64 // Nominal empirical hit rate probability for L1 Cache (0.92)
	PhotonicSsdThroughputTBps float64 // Photonic Glass SSD parallel read throughput in TB/s (1.2 TB/s)
	PhotonicSsdCapacityTB     float64 // Photonic Glass SSD volumetric storage capacity per cube (100 TB)

	// Optical GPU & Photonic AI Tensor Engine
	GpuWdmChannelsCount int     // WDM spectral channels for Optical GPU (RGB: 3 channels)
	AiTensorDensityTOPS float64 // Photonic AI Tensor Core compute density (11 TOPS/mm^2)
	AiTensorEfficiency  float64 // Photonic AI Tensor Core energy efficiency (100 TOPS/W)

	// Room-Temperature Photonic Quantum Processor (LOQC - 298K)
	QuantumQubitsCount      int     // Room-temperature dual-rail photonic qubits (16 Qubits)
	QuantumHomVisibilityPct float64 // HOM 2-photon interference visibility in % (99.4%) [Crespi et al., 2013]
	QuantumCnotFidelityPct  float64 // CNOT logic gate fidelity in % (98.7%) [Carolan et al., 2015]

	// Advanced Physical Noise & Attenuation Parameters
	LaserRinDbHz       float64 // Laser Relative Intensity Noise (RIN) in dB/Hz
	PhaseNoiseSigmaRad float64 // Electro-optic phase noise standard deviation in radians
	GlassLossDbPerCm   float64 // Waveguide attenuation loss in fused silica (0.2 dB/cm)
	MziPhaseErrorRad   float64 // MZI mesh phase drift standard deviation in radians
}

// DefaultParams returns the nominal configuration for SilicaCore using defined constants.
func DefaultParams() OpticalParams {
	return OpticalParams{
		RefractiveIndex:           DefaultRefractiveIndex,
		FastDistanceMM:            DefaultFastDistanceMM,
		DelayedDistanceMM:         DefaultDelayedDistanceMM,
		LaserJitterFwhmPS:         DefaultLaserJitterFwhmPS,
		SpadJitterFwhmPS:          DefaultSpadJitterFwhmPS,
		TdcResolutionPS:           DefaultTdcResolutionPS,
		WindowWidthPS:             DefaultWindowWidthPS,
		LaserEngineMode:           DefaultLaserEngineMode,
		EncodingFormat:            DefaultEncodingFormat,
		BitsPerSymbol:             DefaultBitsPerSymbol,
		MultiLevelStatesCount:     DefaultMultiLevelStatesCount,
		CacheL1LatencyPS:          DefaultCacheL1LatencyPS,
		RamLoopLatencyPS:          DefaultRamLoopLatencyPS,
		CacheL1HitRate:            DefaultCacheL1HitRate,
		PhotonicSsdThroughputTBps: DefaultPhotonicSsdThroughputTBps,
		PhotonicSsdCapacityTB:     DefaultPhotonicSsdCapacityTB,
		GpuWdmChannelsCount:       DefaultGpuWdmChannelsCount,
		AiTensorDensityTOPS:       DefaultAiTensorDensityTOPS,
		AiTensorEfficiency:        DefaultAiTensorEfficiency,
		QuantumQubitsCount:        DefaultQuantumQubitsCount,
		QuantumHomVisibilityPct:   DefaultQuantumHomVisibilityPct,
		QuantumCnotFidelityPct:    DefaultQuantumCnotFidelityPct,
		LaserRinDbHz:              DefaultLaserRinDbHz,
		PhaseNoiseSigmaRad:        DefaultPhaseNoiseSigmaRad,
		GlassLossDbPerCm:          DefaultGlassLossDbPerCm,
		MziPhaseErrorRad:          DefaultMziPhaseErrorRad,
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
