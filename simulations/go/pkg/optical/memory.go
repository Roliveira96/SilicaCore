package optical

import (
	"math"
)

// ============================================================================
// PHYSICAL MEMORY MODEL: DELAY-LINE CAPACITY, UNIFIED HIERARCHY & LOCAL AI
// Separates the light-speed transport time (distance * n_g / c) from the
// storage-cell mechanism time, which is what actually dominates latency.
// ============================================================================

const (
	// DefaultOpticalLineRateGbps is the per-channel line rate of the optical fabric (TFLN modulators, UTC photodiodes).
	DefaultOpticalLineRateGbps = 100.0

	// DefaultPcmBitsPerCell is the multi-level depth of an Sb2Se3 cell; >6 bits reported [Yu et al., arXiv:2604.11649 (2026)].
	DefaultPcmBitsPerCell = 6

	// DefaultPcmCellPitchUM is the crossbar pitch of one in-waveguide PCM weight cell in um (optimistic; ASSUMPTION).
	DefaultPcmCellPitchUM = 5.0

	// DefaultUnifiedMemoryBandwidthGBps is the bandwidth of an HBM3-class unified memory pool in GB/s.
	DefaultUnifiedMemoryBandwidthGBps = 3350.0

	// DefaultPackageDistanceMM is the worst-case distance between compute and unified memory inside the package.
	DefaultPackageDistanceMM = 20.0
)

// MemoryTier is one level of the proposed unified photonic-electronic hierarchy.
type MemoryTier struct {
	Name        string  `json:"name"`
	Technology  string  `json:"technology"`
	Capacity    string  `json:"capacity"`
	DistanceMM  float64 `json:"distance_mm"`
	CellReadPS  float64 `json:"cell_read_ps"`
	TransportPS float64 `json:"transport_ps"`
	TotalReadPS float64 `json:"total_read_ps"`
	Volatile    bool    `json:"volatile"`
	Reference   string  `json:"reference"`
}

// TransportTimePS returns the light-speed flight time over a distance at the given group index.
func TransportTimePS(distanceMM, groupIndex float64) float64 {
	return distanceMM * groupIndex / (SpeedOfLightVacuo * 1e-9)
}

func newTier(name, tech, capacity string, distanceMM, cellPS float64, volatile bool, ref string) MemoryTier {
	transport := TransportTimePS(distanceMM, SiliconNitrideGroupIndex)
	return MemoryTier{
		Name:        name,
		Technology:  tech,
		Capacity:    capacity,
		DistanceMM:  distanceMM,
		CellReadPS:  cellPS,
		TransportPS: transport,
		TotalReadPS: transport + cellPS,
		Volatile:    volatile,
		Reference:   ref,
	}
}

// UnifiedMemoryHierarchy returns the proposed hierarchy for gaming + local AI: every tier is reached over
// light-speed optical links, but each keeps the storage cell that physics allows at its capacity.
func UnifiedMemoryHierarchy() []MemoryTier {
	return []MemoryTier{
		newTier("L1 / registers", "Photonic SRAM (cross-coupled microrings, 40 GHz)", "KB (area-limited)",
			1, 25, true, "arXiv:2503.19544 (2025)"),
		newTier("L2 / L3", "Electronic SRAM hybrid-bonded under the photonic die", "64-256 MB",
			2, 2000, true, "3D-stacked SRAM (AMD V-Cache class)"),
		newTier("AI weights (read-mostly)", "Sb2Se3 PCM in Si3N4 waveguides, >6 bits/cell", "limited by area (see footprint)",
			5, 100, false, "arXiv:2604.11649 (2026)"),
		newTier("Unified RAM", "HBM/LPDDR over co-packaged optical I/O", "64-192 GB",
			DefaultPackageDistanceMM, 30000, true, "DRAM tRCD+CL ~30 ns"),
		newTier("Storage", "NVMe flash over optical link", "1-8 TB",
			50, 50e6, false, "NAND read ~50 us"),
		newTier("Archive", "Fs-written 5D voxels in glass (write-once)", "TB per platter",
			50, 1e9, false, "Project Silica, SOSP 2023"),
	}
}

// DelayLineCapacityBits returns how many bits a recirculating delay loop holds at once.
func DelayLineCapacityBits(lineRateGbps, loopDelayPS float64, channels int) float64 {
	return lineRateGbps * loopDelayPS * 1e-3 * float64(channels)
}

// DelayLineLengthForCapacityM returns the waveguide length a delay-line RAM needs to hold the given number of bits.
func DelayLineLengthForCapacityM(bits, lineRateGbps float64, channels int, groupIndex float64) float64 {
	delayS := bits / (lineRateGbps * 1e9 * float64(channels))
	return delayS * SpeedOfLightVacuo / groupIndex
}

// LocalAIResult sizes an on-device language model against photonic in-memory weights and unified RAM streaming.
type LocalAIResult struct {
	ParamsBillions     float64 `json:"params_billions"`
	BitsPerWeight      int     `json:"bits_per_weight"`
	ModelGB            float64 `json:"model_gb"`
	PcmCells           float64 `json:"pcm_cells"`
	PcmAreaCM2         float64 `json:"pcm_area_cm2"`
	DecodeTokensPerSec float64 `json:"decode_tokens_per_sec"`
	FitsOnReticle      bool    `json:"fits_on_reticle"`
}

// ReticleAreaCM2 is the maximum single-exposure lithography field (26 x 33 mm).
const ReticleAreaCM2 = 8.58

// SizeLocalAI computes the PCM area needed to hold all weights on-chip, and the decode speed bound when weights
// are instead streamed from unified RAM (autoregressive decoding reads every weight once per token).
func SizeLocalAI(paramsBillions float64, bitsPerWeight int) LocalAIResult {
	bits := paramsBillions * 1e9 * float64(bitsPerWeight)
	cellsPerWeight := math.Ceil(float64(bitsPerWeight) / float64(DefaultPcmBitsPerCell))
	cells := paramsBillions * 1e9 * cellsPerWeight
	pitchCM := DefaultPcmCellPitchUM * 1e-4
	area := cells * pitchCM * pitchCM
	modelGB := bits / 8 / 1e9

	return LocalAIResult{
		ParamsBillions:     paramsBillions,
		BitsPerWeight:      bitsPerWeight,
		ModelGB:            modelGB,
		PcmCells:           cells,
		PcmAreaCM2:         area,
		DecodeTokensPerSec: DefaultUnifiedMemoryBandwidthGBps / modelGB,
		FitsOnReticle:      area <= ReticleAreaCM2,
	}
}
