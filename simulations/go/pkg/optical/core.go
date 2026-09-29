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
	// 1.1 MICRO-CUBE GEOMETRIC DOWNSCALING & HIGH-DENSITY INTEGRATION PARAMETERS
	// ------------------------------------------------------------------------

	// MicroFastDistanceMM is the scaled-down straight-line path distance d1 (2.0 mm micro-cube).
	MicroFastDistanceMM = 2.0

	// MicroDelayedDistanceMM is the scaled-down deflected path distance d0 (4.0675 mm micro-cube).
	MicroDelayedDistanceMM = 4.0675

	// MicroLaserJitterFwhmPS is the ultra-fast integrated laser pulse temporal jitter in ps (FWHM).
	MicroLaserJitterFwhmPS = 1.0

	// MicroSpadJitterFwhmPS is the integrated SPAD detector jitter in ps (FWHM).
	MicroSpadJitterFwhmPS = 3.0

	// MicroTdcResolutionPS is the high-resolution TDC LSB resolution in ps.
	MicroTdcResolutionPS = 1.0

	// MicroRamLoopLatencyPS is the micro-waveguide Photonic RAM recirculating delay loop latency in ps.
	MicroRamLoopLatencyPS = 9.67

	// DefaultDwdmChannelsCount is the number of dense wavelength division multiplexing (DWDM) channels.
	DefaultDwdmChannelsCount = 64

	// ------------------------------------------------------------------------
	// 2. CONTINUOUS WAVE LASER ENGINE (SOLID-STATE CW) & M-ARY DENSE ENCODING
	// ------------------------------------------------------------------------

	// DefaultLaserEngineMode is the operational mode of the continuous wave lasers.
	DefaultLaserEngineMode = "Continuous Wave (Solid-State CW Engine - Always-ON)"

	// DefaultEncodingFormat is the spectral/phase data encoding scheme.
	DefaultEncodingFormat = "M-ary Multi-Level WDM (Hexadecimal 4-bit / 16-QAM/PAM-4 Equivalent)"

	// DefaultBitsPerSymbol is the number of bits encoded per optical symbol (4 bits = 1 Nibble).
	DefaultBitsPerSymbol = 4

	// DefaultMultiLevelStatesCount is the number of discrete optical states per symbol (2^4 = 16).
	DefaultMultiLevelStatesCount = 16

	// ------------------------------------------------------------------------
	// 3. PHOTONIC MEMORY HIERARCHY & PHOTONIC GLASS SSD
	// ------------------------------------------------------------------------

	// DefaultCacheL1LatencyPS is the L1 photonic SRAM cell latency in ps: cross-coupled microring pSRAM at 40 GHz
	// on GlobalFoundries 45 nm [arXiv:2503.19544, 2025]. Capacity is KB-class (area-limited).
	DefaultCacheL1LatencyPS = 25.0

	// DefaultRamLoopLatencyPS is the recirculating delay-loop period in ps [Yao, 1993]. It is a buffer/register
	// (619 bits per loop at 64 ch x 100 Gb/s), not main RAM; main RAM is unified HBM (see memory.go).
	DefaultRamLoopLatencyPS = 96.73

	// DefaultCacheL1HitRate is the nominal empirical hit rate probability for L1 Photonic Cache.
	DefaultCacheL1HitRate = 0.92

	// DefaultCacheL2HitRate is the hit rate of the 3D-stacked electronic L2/L3 SRAM for L1 misses (ASSUMPTION).
	DefaultCacheL2HitRate = 0.90

	// DefaultPhotonicSsdCapacityTB is the glass storage capacity in 15.6 cm^3 in TB (ASSUMPTION: 6.4 TB/cm^3;
	// published fs-written glass storage is write-once archival media [Project Silica, SOSP 2023]).
	DefaultPhotonicSsdCapacityTB = 100.0

	// DefaultPhotonicSsdThroughputTBps is the parallel WDM read throughput in TB/s (ASSUMPTION: not demonstrated;
	// Project Silica reads by polarization microscopy).
	DefaultPhotonicSsdThroughputTBps = 1.2

	// ------------------------------------------------------------------------
	// 4. OPTICAL GPU (WDM RGB) & PHOTONIC AI TENSOR ENGINE (MZI MESH)
	// ------------------------------------------------------------------------

	// DefaultGpuWdmChannelsCount is the number of WDM spectral channels (Red 635nm, Green 532nm, Blue 450nm).
	DefaultGpuWdmChannelsCount = 3

	// DefaultAiTensorDensityTOPS is the Photonic AI Tensor Core compute density in TOPS/mm^2 [Xu et al., Nature 2021].
	DefaultAiTensorDensityTOPS = 11.0

	// DefaultAiTensorEfficiency is the optical-core-only efficiency in TOPS/W [Shen et al., Nature Phot. 2017].
	// Full-system state of the art is ~0.84 TOPS/W [Ahmed et al., Nature 2025].
	DefaultAiTensorEfficiency = 100.0

	// ------------------------------------------------------------------------
	// 5. PHOTONIC QUANTUM PROCESSOR (LOQC): circuit at room temperature, SNSPDs and sources cryogenic (~1-4 K)
	// ------------------------------------------------------------------------

	// DefaultQuantumQubitsCount is the number of dual-rail photonic qubits (cryogenic detection subsystem).
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

	// ------------------------------------------------------------------------
	// 7. POWER CONSUMPTION & COMPARATIVE SILICON BASELINES
	// ------------------------------------------------------------------------

	// DefaultCwLaserPowerWatts is the electrical power consumed by continuous wave RGB lasers (Watts).
	DefaultCwLaserPowerWatts = 2.5

	// DefaultEomEnergyFjPerBit is the energy consumed per electro-optic modulation bit (fJ/bit).
	DefaultEomEnergyFjPerBit = 0.8

	// DefaultSpadEnergyFjPerPhoton is the energy consumed per SPAD single-photon detection event (fJ/photon).
	DefaultSpadEnergyFjPerPhoton = 12.0

	// DefaultSilicaCoreTdpWatts is the board TDP in Watts (ASSUMPTION: not derived from the component model).
	DefaultSilicaCoreTdpWatts = 18.5

	// DefaultIntelI9TdpWatts is the TDP of Intel Core i9-14900KS CPU in Watts.
	DefaultIntelI9TdpWatts = 253.0

	// DefaultNvidiaH100TdpWatts is the TDP of NVIDIA H100 Tensor Core GPU in Watts.
	DefaultNvidiaH100TdpWatts = 700.0

	// DefaultAmdEpycTdpWatts is the TDP of AMD EPYC 9654 96-core server CPU in Watts.
	DefaultAmdEpycTdpWatts = 360.0
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

	// Solid-State Continuous Wave (CW) Laser Engine & M-ary Multi-Level Encoding
	LaserEngineMode       string // Laser Mode: "Continuous Wave (Solid-State CW Engine - Always-ON)"
	EncodingFormat        string // Data Encoding: "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)"
	BitsPerSymbol         int    // Bits per optical symbol (8 bits = 1 Byte / symbol)
	MultiLevelStatesCount int    // Number of discrete spectral/phase states (256 states for 8-bit Byte)

	// Photonic Memory & Photonic SSD
	CacheL1LatencyPS          float64 // L1 photonic SRAM cell latency: ~25 ps (arXiv:2503.19544)
	RamLoopLatencyPS          float64 // Delay-loop buffer period: ~96.73 ps (Yao, 1993); buffer, not main RAM
	CacheL1HitRate            float64 // Nominal empirical hit rate probability for L1 Cache (0.92)
	CacheL2HitRate            float64 // Hit rate of 3D-stacked L2/L3 SRAM for L1 misses (0.90, ASSUMPTION)
	PhotonicSsdThroughputTBps float64 // Photonic Glass SSD parallel read throughput in TB/s (1.2 TB/s)
	PhotonicSsdCapacityTB     float64 // Photonic Glass SSD volumetric storage capacity per cube (100 TB)

	// Dense Wavelength Division Multiplexing (DWDM) & Optical GPU
	GpuWdmChannelsCount int     // WDM spectral channels for Optical GPU (RGB: 3 channels)
	DwdmChannelsCount   int     // DWDM spectral channels (64 channels)
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

	// Power Consumption & Comparative Silicon Baselines
	CwLaserPowerWatts     float64 // Electrical power of CW RGB lasers in Watts (2.5W)
	EomEnergyFjPerBit     float64 // Energy per EOM modulation bit in fJ/bit (0.8 fJ)
	SpadEnergyFjPerPhoton float64 // Energy per SPAD detection in fJ/photon (12 fJ)
	SilicaCoreTdpWatts    float64 // Total SilicaCore TDP in Watts (18.5W)
	IntelI9TdpWatts       float64 // Baseline Intel Core i9-14900KS TDP in Watts (253W)
	NvidiaH100TdpWatts    float64 // Baseline NVIDIA H100 GPU TDP in Watts (700W)
	AmdEpycTdpWatts       float64 // Baseline AMD EPYC 9654 CPU TDP in Watts (360W)
}

// DwdmBitsPerPulse returns the total number of parallel bits transmitted per optical pulse over DWDM channels.
func (p OpticalParams) DwdmBitsPerPulse() int {
	return p.DwdmChannelsCount * p.BitsPerSymbol
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
		CacheL2HitRate:            DefaultCacheL2HitRate,
		PhotonicSsdThroughputTBps: DefaultPhotonicSsdThroughputTBps,
		PhotonicSsdCapacityTB:     DefaultPhotonicSsdCapacityTB,
		GpuWdmChannelsCount:       DefaultGpuWdmChannelsCount,
		DwdmChannelsCount:         DefaultDwdmChannelsCount,
		AiTensorDensityTOPS:       DefaultAiTensorDensityTOPS,
		AiTensorEfficiency:        DefaultAiTensorEfficiency,
		QuantumQubitsCount:        DefaultQuantumQubitsCount,
		QuantumHomVisibilityPct:   DefaultQuantumHomVisibilityPct,
		QuantumCnotFidelityPct:    DefaultQuantumCnotFidelityPct,
		LaserRinDbHz:              DefaultLaserRinDbHz,
		PhaseNoiseSigmaRad:        DefaultPhaseNoiseSigmaRad,
		GlassLossDbPerCm:          DefaultGlassLossDbPerCm,
		MziPhaseErrorRad:          DefaultMziPhaseErrorRad,
		CwLaserPowerWatts:        DefaultCwLaserPowerWatts,
		EomEnergyFjPerBit:        DefaultEomEnergyFjPerBit,
		SpadEnergyFjPerPhoton:    DefaultSpadEnergyFjPerPhoton,
		SilicaCoreTdpWatts:       DefaultSilicaCoreTdpWatts,
		IntelI9TdpWatts:          DefaultIntelI9TdpWatts,
		NvidiaH100TdpWatts:       DefaultNvidiaH100TdpWatts,
		AmdEpycTdpWatts:          DefaultAmdEpycTdpWatts,
	}
}

// MicroCubeParams returns the ultra-high throughput configuration for a 2.0 mm micro-cube with 64-channel dense DWDM.
func MicroCubeParams() OpticalParams {
	p := DefaultParams()
	p.FastDistanceMM = MicroFastDistanceMM
	p.DelayedDistanceMM = MicroDelayedDistanceMM
	p.LaserJitterFwhmPS = MicroLaserJitterFwhmPS
	p.SpadJitterFwhmPS = MicroSpadJitterFwhmPS
	p.TdcResolutionPS = MicroTdcResolutionPS
	p.RamLoopLatencyPS = MicroRamLoopLatencyPS
	p.DwdmChannelsCount = DefaultDwdmChannelsCount
	p.EncodingFormat = "Massive 64-Channel DWDM + M-ary Multi-Level (Hexadecimal 4-bit / 16-QAM)"
	return p
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
