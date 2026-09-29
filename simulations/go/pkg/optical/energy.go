package optical

import (
	"math"
)

// ============================================================================
// BOTTOM-UP ENERGY AND THERMAL MODEL OF THE RACE-LOGIC CHIP
// Adds up what every block spends while the chip answers back-to-back queries
// (one query per race + readout period). Values come from published devices;
// ASSUMPTION marks the ones without a direct figure.
// ============================================================================

const (
	// DefaultLaserWallPlugEfficiency is the electrical-to-optical efficiency of the 1550 nm source.
	// DFB lasers reach ~18% at 10 mW and ~35% at 250 mW optical output (published range); 20% is conservative.
	DefaultLaserWallPlugEfficiency = 0.20

	// DefaultLaserDistributionExcessDb is the excess loss of the 1:N splitter tree feeding every
	// re-fire modulator (about 0.3 dB per 1x2 stage over ~10 stages, ASSUMPTION).
	DefaultLaserDistributionExcessDb = 3.0

	// DefaultModulatorEnergyFJ is the switching energy per TFLN re-fire event including its driver.
	// Driverless TFLN modulators reach ~4.5 fJ/bit; x10 covers the CMOS driver (ASSUMPTION).
	DefaultModulatorEnergyFJ = 50.0

	// DefaultReceiverPowerMW is the static power of one photodiode + TIA + comparator, always biased.
	// Integrated receivers reach ~0.17 pJ/bit at 25 Gb/s, i.e. ~4-5 mW per channel.
	DefaultReceiverPowerMW = 5.0

	// DefaultTdcPowerMW is the power of one picosecond time-to-digital converter channel
	// (~4.1 mW per channel for 0.9 ps rms in 130 nm CMOS, shared oscillator).
	DefaultTdcPowerMW = 4.1

	// DefaultReadoutEnergyPJPerBit is the serializer / link energy to ship node times off chip (ASSUMPTION).
	DefaultReadoutEnergyPJPerBit = 2.0

	// DefaultPcmSwitchEnergyNJ is the energy to reprogram one Sb2Se3 switch with its micro-heater (ASSUMPTION);
	// it is spent only when the map changes.
	DefaultPcmSwitchEnergyNJ = 10.0

	// DefaultThermalResistanceKPerW is the junction-to-ambient thermal resistance of the package with a
	// small fan-cooled heatsink (ASSUMPTION); passive cooling is ~2-5 K/W.
	DefaultThermalResistanceKPerW = 0.5
)

// EnergyParams holds the per-block energy figures of the chip.
type EnergyParams struct {
	LaserWallPlugEfficiency   float64
	LaserDistributionExcessDb float64
	ModulatorEnergyFJ         float64
	ReceiverPowerMW           float64
	TdcPowerMW                float64
	ReadoutEnergyPJPerBit     float64
	PcmSwitchEnergyNJ         float64
	ThermalResistanceKPerW    float64
}

// DefaultEnergyParams returns the nominal energy figures.
func DefaultEnergyParams() EnergyParams {
	return EnergyParams{
		LaserWallPlugEfficiency:   DefaultLaserWallPlugEfficiency,
		LaserDistributionExcessDb: DefaultLaserDistributionExcessDb,
		ModulatorEnergyFJ:         DefaultModulatorEnergyFJ,
		ReceiverPowerMW:           DefaultReceiverPowerMW,
		TdcPowerMW:                DefaultTdcPowerMW,
		ReadoutEnergyPJPerBit:     DefaultReadoutEnergyPJPerBit,
		PcmSwitchEnergyNJ:         DefaultPcmSwitchEnergyNJ,
		ThermalResistanceKPerW:    DefaultThermalResistanceKPerW,
	}
}

// EnergyBudgetResult is the power and energy breakdown while answering back-to-back queries.
type EnergyBudgetResult struct {
	QueryPeriodNS       float64 `json:"query_period_ns"`
	QueriesPerSecond    float64 `json:"queries_per_second"`
	LaserOpticalW       float64 `json:"laser_optical_w"`
	LaserElectricalW    float64 `json:"laser_electrical_w"`
	ReceiversW          float64 `json:"receivers_w"`
	TdcsW               float64 `json:"tdcs_w"`
	ModulatorsW         float64 `json:"modulators_w"`
	ReadoutW            float64 `json:"readout_w"`
	TotalW              float64 `json:"total_w"`
	EnergyPerQueryNJ    float64 `json:"energy_per_query_nj"`
	PowerDensityWPerCm2 float64 `json:"power_density_w_per_cm2"`
	TemperatureRiseK    float64 `json:"temperature_rise_k"`
	ReprogramEnergyUJ   float64 `json:"reprogram_energy_uj"`
	LargestBlock        string  `json:"largest_block"`
}

func dbmToWatts(dbm float64) float64 { return math.Pow(10, dbm/10) / 1000 }

// ComputeEnergyBudget evaluates a programmed race-logic chip running queries back to back.
// queryPeriodNS is the time per query (race + readout) from the simulator.
func ComputeEnergyBudget(hw RaceHardware, queryPeriodNS float64, e EnergyParams) EnergyBudgetResult {
	p := hw.Params
	edges := float64(hw.Graph.NumEdges())
	nodes := float64(len(hw.Graph.Adj))

	// Each re-fire modulator must launch enough light to reach its receiver at the sensitivity limit after
	// the worst edge loss; the splitter tree feeding the modulators adds its excess loss.
	perEdgeDbm := DefaultPhotodiodeSensitivityDbm + hw.WorstEdgeLossDb + e.LaserDistributionExcessDb
	laserOptical := edges * dbmToWatts(perEdgeDbm)
	laserElectrical := laserOptical / e.LaserWallPlugEfficiency

	receivers := float64(hw.Detectors) * e.ReceiverPowerMW / 1000
	tdcs := float64(hw.Tdcs) * e.TdcPowerMW / 1000
	queriesPerSecond := 1e9 / queryPeriodNS
	modulators := edges * e.ModulatorEnergyFJ * 1e-15 * queriesPerSecond
	readoutBits := nodes * float64(p.TdcBitsPerNode)
	readout := readoutBits * e.ReadoutEnergyPJPerBit * 1e-12 * queriesPerSecond

	total := laserElectrical + receivers + tdcs + modulators + readout
	areaCm2 := math.Max(hw.DelayAreaMM2, 1) / 100

	blocks := map[string]float64{
		"laser": laserElectrical, "receivers": receivers, "tdcs": tdcs, "modulators": modulators, "readout": readout,
	}
	largest, largestW := "", -1.0
	for name, w := range blocks {
		if w > largestW {
			largest, largestW = name, w
		}
	}

	return EnergyBudgetResult{
		QueryPeriodNS:       queryPeriodNS,
		QueriesPerSecond:    queriesPerSecond,
		LaserOpticalW:       laserOptical,
		LaserElectricalW:    laserElectrical,
		ReceiversW:          receivers,
		TdcsW:               tdcs,
		ModulatorsW:         modulators,
		ReadoutW:            readout,
		TotalW:              total,
		EnergyPerQueryNJ:    total * queryPeriodNS,
		PowerDensityWPerCm2: total / areaCm2,
		TemperatureRiseK:    total * e.ThermalResistanceKPerW,
		ReprogramEnergyUJ:   float64(hw.PcmSwitches) * e.PcmSwitchEnergyNJ / 1000,
		LargestBlock:        largest,
	}
}

// CPUQueryEnergy is the energy of one classical query given a measured time and a package power.
// Using the full package power gives an upper bound for a single-threaded run.
func CPUQueryEnergy(packagePowerW, queryTimeUS float64) (energyNJ float64) {
	return packagePowerW * queryTimeUS * 1000
}
