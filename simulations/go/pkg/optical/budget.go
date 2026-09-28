package optical

import (
	"math"
)

// ============================================================================
// PHYSICAL LINK BUDGET: ROUTING PLATFORMS, SWITCHING & DETECTION LIMITS
// Values are taken from the literature cited next to each constant (see
// docs/architecture/11-roteamento-e-comutacao-optica.md). Values marked
// ASSUMPTION have no direct published figure yet and must be measured.
// ============================================================================

const (
	// PlanckConstant is the Planck constant in J*s.
	PlanckConstant = 6.62607015e-34

	// DefaultTargetBER is the design bit error rate for logic-level decisions.
	DefaultTargetBER = 1e-12

	// DefaultTelecomWavelengthNM is the C-band operating wavelength (Kerr microcomb / DWDM grid).
	DefaultTelecomWavelengthNM = 1550.0

	// DefaultChannelLaserPowerDbm is the optical power available per DWDM channel (comb line) in dBm.
	DefaultChannelLaserPowerDbm = 0.0

	// DefaultPhotodiodeSensitivityDbm is the high-speed Ge/InGaAs photodiode sensitivity at the target BER in dBm.
	DefaultPhotodiodeSensitivityDbm = -10.0

	// DefaultPhotodiodeBandwidthGHz is the 3 dB bandwidth of a uni-travelling-carrier (UTC) photodiode in GHz.
	DefaultPhotodiodeBandwidthGHz = 100.0

	// DefaultSpadDeadTimeNS is the SPAD recovery (dead) time after each avalanche in ns.
	// Best reported free-running InGaAs/InP: ~1-2 ns (500 Mcount/s) [Sensors 2023, PMC9961215].
	DefaultSpadDeadTimeNS = 2.0

	// DefaultSpadDetectionEfficiency is the SPAD photon detection efficiency (0..1).
	DefaultSpadDetectionEfficiency = 0.5

	// DefaultBeamWaistUM is the launched Gaussian beam waist radius in um.
	DefaultBeamWaistUM = 5.0

	// DefaultDetectorApertureUM is the capture radius of the receiving detector or waveguide in um.
	DefaultDetectorApertureUM = 10.0

	// DefaultTurnsPerGate is the number of mirror reflections or waveguide bends in one ToF gate path.
	DefaultTurnsPerGate = 4

	// FusedSilicaDnDT is the thermo-optic coefficient of fused silica (1/K).
	FusedSilicaDnDT = 1.0e-5

	// FusedSilicaExpansion is the linear thermal expansion coefficient of fused silica (1/K).
	FusedSilicaExpansion = 0.55e-6

	// FusedSilicaAcousticVelocity is the longitudinal sound speed in fused silica (m/s); it bounds AOM switching.
	FusedSilicaAcousticVelocity = 5960.0

	// FusedSilicaMaxFieldDeltaN is an optimistic upper bound of field-induced index change in unpoled silica (DC Kerr).
	FusedSilicaMaxFieldDeltaN = 1e-9

	// SiliconNitrideDnDT is the thermo-optic coefficient of LPCVD silicon nitride (1/K).
	SiliconNitrideDnDT = 2.45e-5

	// SiliconNitrideExpansion is the linear thermal expansion coefficient of silicon nitride (1/K).
	SiliconNitrideExpansion = 3.3e-6

	// SiliconNitrideGroupIndex is the typical group index of a Si3N4 strip waveguide at 1550 nm.
	SiliconNitrideGroupIndex = 2.0
)

// RoutingPlatform describes how light is steered between gates: the "mirror" problem.
type RoutingPlatform struct {
	Name               string  `json:"name"`
	PropagationDbPerCm float64 `json:"propagation_db_per_cm"` // Propagation loss in the medium
	TurnLossDb         float64 `json:"turn_loss_db"`          // Loss per mirror reflection or waveguide bend
	SwitchTechnology   string  `json:"switch_technology"`     // Fast (per-bit) switching element
	SwitchLossDb       float64 `json:"switch_loss_db"`        // Insertion loss of the fast switch
	SwitchRiseTimePS   float64 `json:"switch_rise_time_ps"`   // 10-90% switching time of the fast switch
	GroupIndex         float64 `json:"group_index"`           // Group index, sets delay per mm
	ThermoOpticDnDT    float64 `json:"thermo_optic_dndt"`     // dn/dT (1/K)
	ThermalExpansion   float64 `json:"thermal_expansion"`     // Linear expansion (1/K)
	Guided             bool    `json:"guided"`                // false = free-space beam in bulk, diffraction applies
	FastSwitchFeasible bool    `json:"fast_switch_feasible"`  // Whether a ps-scale switch exists on this platform
}

// BulkFreeSpaceMirrorsPlatform is the current SilicaCore concept: free beams in a fused silica cube, steered by internal mirrors.
func BulkFreeSpaceMirrorsPlatform() RoutingPlatform {
	return RoutingPlatform{
		Name:               "Bulk SiO2 free-space + internal mirrors (current)",
		PropagationDbPerCm: 0.001,
		TurnLossDb:         0.2,
		SwitchTechnology:   "EOM/Kerr in SiO2 (no Pockels) or AOM (acoustic)",
		SwitchLossDb:       0.0,
		SwitchRiseTimePS:   AomAccessTimePS(2 * DefaultBeamWaistUM * 10), // 100 um beam through an AOM
		GroupIndex:         DefaultRefractiveIndex,
		ThermoOpticDnDT:    FusedSilicaDnDT,
		ThermalExpansion:   FusedSilicaExpansion,
		Guided:             false,
		FastSwitchFeasible: false,
	}
}

// FsWaveguideTirPlatform uses femtosecond-written 3D waveguides in glass with TIR micro-mirrors etched as air trenches (FLICE).
func FsWaveguideTirPlatform() RoutingPlatform {
	return RoutingPlatform{
		Name: "Fs-written 3D waveguides in SiO2 + air-trench TIR micro-mirrors",
		// Typical fs-written loss 0.3 dB/cm; benchmark 0.05 dB/cm [Sci. Rep. 8, 10.1038/s41598-018-28631-3].
		PropagationDbPerCm: 0.3,
		// ASSUMPTION: no published per-mirror loss for buried air-trench TIR mirrors in fused silica.
		TurnLossDb:       0.5,
		SwitchTechnology: "Thin-film lithium niobate (TFLN) EO switch, hybrid-coupled to glass",
		// Two glass<->TFLN mode transitions (~1.25 dB each, fiber-to-chip class) + 0.5 dB device.
		SwitchLossDb:       3.0,
		SwitchRiseTimePS:   5.0,
		GroupIndex:         DefaultRefractiveIndex,
		ThermoOpticDnDT:    FusedSilicaDnDT,
		ThermalExpansion:   FusedSilicaExpansion,
		Guided:             true,
		FastSwitchFeasible: true,
	}
}

// SiliconNitrideTflnPlatform uses stacked Si3N4-on-SiO2 layers with tight bends (no mirrors) and heterogeneous TFLN switches.
func SiliconNitrideTflnPlatform() RoutingPlatform {
	return RoutingPlatform{
		Name: "Si3N4-on-SiO2 multilayer (bends, no mirrors) + heterogeneous TFLN",
		// Si3N4 <0.1 dB/cm and LN<->Si3N4 adiabatic transitions <0.1 dB [Churaev et al., Nat. Commun. 14, 3499 (2023)].
		PropagationDbPerCm: 0.1,
		// 50 um bends, 0.01 dB interlayer couplers [Shang et al., Opt. Express 23, 21334 (2015)].
		TurnLossDb:       0.01,
		SwitchTechnology: "Thin-film lithium niobate (TFLN) EO switch, heterogeneous on Si3N4",
		// 2 x 0.1 dB transitions + ~0.8 dB MZI switch.
		SwitchLossDb:       1.0,
		SwitchRiseTimePS:   5.0,
		GroupIndex:         SiliconNitrideGroupIndex,
		ThermoOpticDnDT:    SiliconNitrideDnDT,
		ThermalExpansion:   SiliconNitrideExpansion,
		Guided:             true,
		FastSwitchFeasible: true,
	}
}

// AllRoutingPlatforms returns the platforms compared by the link budget, current concept first.
func AllRoutingPlatforms() []RoutingPlatform {
	return []RoutingPlatform{
		BulkFreeSpaceMirrorsPlatform(),
		FsWaveguideTirPlatform(),
		SiliconNitrideTflnPlatform(),
	}
}

// GateBudgetResult holds the optical power budget of one ToF gate and how many gates can be cascaded.
type GateBudgetResult struct {
	Platform           string  `json:"platform"`
	PathLengthMM       float64 `json:"path_length_mm"`
	PropagationLossDb  float64 `json:"propagation_loss_db"`
	TurnsLossDb        float64 `json:"turns_loss_db"`
	DiffractionLossDb  float64 `json:"diffraction_loss_db"`
	SwitchLossDb       float64 `json:"switch_loss_db"`
	TotalLossPerGateDb float64 `json:"total_loss_per_gate_db"`
	PowerMarginDb      float64 `json:"power_margin_db"`
	MaxCascadedGates   int     `json:"max_cascaded_gates"`
	SwitchRiseTimePS   float64 `json:"switch_rise_time_ps"`
	PhaseDriftRadPerK  float64 `json:"phase_drift_rad_per_k"`
	Feasible           bool    `json:"feasible"`
}

// PathLengthForDelayMM returns the physical path length that produces the given propagation delay.
func PathLengthForDelayMM(delayPS, groupIndex float64) float64 {
	return delayPS * (SpeedOfLightVacuo / groupIndex) * 1e-9
}

// GaussianBeamRadiusUM returns the 1/e^2 radius of a Gaussian beam after propagating zMM in a medium of index n.
func GaussianBeamRadiusUM(waistUM, wavelengthNM, n, zMM float64) float64 {
	w0 := waistUM * 1e-6
	zR := math.Pi * w0 * w0 * n / (wavelengthNM * 1e-9)
	z := zMM * 1e-3
	return w0 * math.Sqrt(1+(z/zR)*(z/zR)) * 1e6
}

// CaptureFraction returns the fraction of a centred Gaussian beam of radius wUM collected by a circular aperture of radius aUM.
func CaptureFraction(aUM, wUM float64) float64 {
	return 1 - math.Exp(-2*aUM*aUM/(wUM*wUM))
}

// LinearToDb converts a power ratio (0..1) into a positive loss in dB.
func LinearToDb(ratio float64) float64 {
	return -10 * math.Log10(ratio)
}

// ThermalPhaseDriftRadPerK returns the optical phase drift per kelvin over a path of lengthMM.
func ThermalPhaseDriftRadPerK(lengthMM, n, dnDT, expansion, wavelengthNM float64) float64 {
	return 2 * math.Pi * (lengthMM * 1e-3) * (dnDT + n*expansion) / (wavelengthNM * 1e-9)
}

// TirGlancingAngleDeg returns the maximum grazing angle (from the interface) at which an index step deltaN still gives total internal reflection.
func TirGlancingAngleDeg(n, deltaN float64) float64 {
	return math.Acos((n-deltaN)/n) * 180 / math.Pi
}

// AomAccessTimePS returns the acoustic transit time across a beam of the given diameter in fused silica, the AOM switching floor.
func AomAccessTimePS(beamDiameterUM float64) float64 {
	return beamDiameterUM * 1e-6 / FusedSilicaAcousticVelocity * 1e12
}

// ComputeGateBudget evaluates the worst-case (delayed line) optical path of one ToF gate on a routing platform.
func (s *ToFSimulator) ComputeGateBudget(p RoutingPlatform) GateBudgetResult {
	length := PathLengthForDelayMM(s.T0NominalPS, p.GroupIndex)

	propLoss := p.PropagationDbPerCm * length / 10.0
	turnsLoss := float64(DefaultTurnsPerGate) * p.TurnLossDb

	diffLoss := 0.0
	if !p.Guided {
		w := GaussianBeamRadiusUM(DefaultBeamWaistUM, DefaultTelecomWavelengthNM, p.GroupIndex, length)
		diffLoss = LinearToDb(CaptureFraction(DefaultDetectorApertureUM, w))
	}

	total := propLoss + turnsLoss + diffLoss + p.SwitchLossDb
	margin := DefaultChannelLaserPowerDbm - DefaultPhotodiodeSensitivityDbm

	maxGates := 0
	if p.FastSwitchFeasible && total > 0 {
		maxGates = int(math.Floor(margin / total))
	}

	return GateBudgetResult{
		Platform:           p.Name,
		PathLengthMM:       length,
		PropagationLossDb:  propLoss,
		TurnsLossDb:        turnsLoss,
		DiffractionLossDb:  diffLoss,
		SwitchLossDb:       p.SwitchLossDb,
		TotalLossPerGateDb: total,
		PowerMarginDb:      margin,
		MaxCascadedGates:   maxGates,
		SwitchRiseTimePS:   p.SwitchRiseTimePS,
		PhaseDriftRadPerK:  ThermalPhaseDriftRadPerK(length, p.GroupIndex, p.ThermoOpticDnDT, p.ThermalExpansion, DefaultTelecomWavelengthNM),
		Feasible:           p.FastSwitchFeasible && maxGates > 0,
	}
}

// TimingBudgetResult holds the corrected BER and realistic per-channel symbol rate of the ToF scheme.
type TimingBudgetResult struct {
	QFactor                float64 `json:"q_factor"`
	TheoreticalBER         float64 `json:"theoretical_ber"`
	RequiredQForTarget     float64 `json:"required_q_for_target"`
	RequiredSigmaPS        float64 `json:"required_sigma_ps"`
	SymbolSlotPS           float64 `json:"symbol_slot_ps"`
	ToFSymbolRateGHz       float64 `json:"tof_symbol_rate_ghz"`
	SpadMaxRateGHz         float64 `json:"spad_max_rate_ghz"`
	SpadPhotonsPerBit      float64 `json:"spad_photons_per_bit"`
	PhotodiodeMaxRateGbaud float64 `json:"photodiode_max_rate_gbaud"`
	DetectorLimitedRateGHz float64 `json:"detector_limited_rate_ghz"`
}

// RequiredQFactor returns the Q factor needed for a Gaussian decision to reach the target BER.
func RequiredQFactor(targetBER float64) float64 {
	return math.Sqrt2 * math.Erfcinv(2*targetBER)
}

// ComputeTimingBudget corrects the separation-margin claim and derives the achievable ToF symbol rate per channel.
func (s *ToFSimulator) ComputeTimingBudget(targetBER float64) TimingBudgetResult {
	q := s.DeltaTNominalPS / (2.0 * s.TotalSigmaPS)
	qReq := RequiredQFactor(targetBER)

	// A symbol slot must hold both gating windows: from t1 - W/2 to t0 + W/2.
	slot := s.DeltaTNominalPS + s.Params.WindowWidthPS
	tofRate := 1000.0 / slot

	spadRate := 1.0 / DefaultSpadDeadTimeNS
	pdRate := 1.5 * DefaultPhotodiodeBandwidthGHz

	return TimingBudgetResult{
		QFactor:                q,
		TheoreticalBER:         0.5 * math.Erfc(q/math.Sqrt2),
		RequiredQForTarget:     qReq,
		RequiredSigmaPS:        s.DeltaTNominalPS / (2.0 * qReq),
		SymbolSlotPS:           slot,
		ToFSymbolRateGHz:       tofRate,
		SpadMaxRateGHz:         spadRate,
		SpadPhotonsPerBit:      -math.Log(targetBER) / DefaultSpadDetectionEfficiency,
		PhotodiodeMaxRateGbaud: pdRate,
		DetectorLimitedRateGHz: math.Min(tofRate, pdRate),
	}
}
