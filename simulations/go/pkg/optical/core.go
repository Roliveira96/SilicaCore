package optical

import (
	"math"
)

// ============================================================================
// CONSTANTES FÍSICAS E CONFIGURAÇÕES PADRÃO DO SIMULADOR (SILICA CORE)
// Altere os valores abaixo para ajustar os parâmetros físicos do processador.
// ============================================================================

const (
	// SpeedOfLightVacuo: Velocidade da luz no vácuo em metros por segundo (m/s).
	SpeedOfLightVacuo = 299792458.0

	// ------------------------------------------------------------------------
	// 1. PARÂMETROS FÍSICOS DO SUBSTRATO E PROPAGAÇÃO (ToF)
	// ------------------------------------------------------------------------

	// DefaultRefractiveIndex: Índice de refração efetivo da sílica fundida (SiO2).
	DefaultRefractiveIndex = 1.4500

	// DefaultFastDistanceMM: Distância em linha reta da linha rápida d1 (mm).
	DefaultFastDistanceMM = 20.0

	// DefaultDelayedDistanceMM: Distância percorrida na linha de retardo d0 (mm).
	DefaultDelayedDistanceMM = 40.675

	// DefaultLaserJitterFwhmPS: Jitter temporal do diodo/canhão laser em picossegundos (FWHM).
	DefaultLaserJitterFwhmPS = 8.0

	// DefaultSpadJitterFwhmPS: Jitter do detector de fóton único SPAD em picossegundos (FWHM).
	DefaultSpadJitterFwhmPS = 25.0

	// DefaultTdcResolutionPS: Resolução LSB do conversor Tempo-para-Digital (TDC) em picossegundos.
	DefaultTdcResolutionPS = 5.0

	// DefaultWindowWidthPS: Largura da janela de amostragem (Time-Gating Window) em picossegundos.
	DefaultWindowWidthPS = 75.0

	// ------------------------------------------------------------------------
	// 2. MOTOR LASER CONTÍNUO (NORITSU CW) E CODIFICAÇÃO DENSA M-ÁRIA
	// ------------------------------------------------------------------------

	// DefaultLaserEngineMode: Modo de operação dos canhões laser contínuos.
	DefaultLaserEngineMode = "Continuous Wave (Noritsu CW Style - Always-ON)"

	// DefaultEncodingFormat: Esquema de codificação de dados espectral/fase.
	DefaultEncodingFormat = "M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)"

	// DefaultBitsPerSymbol: Quantidade de bits codificados por símbolo óptico (8 bits = 1 Byte).
	DefaultBitsPerSymbol = 8

	// DefaultMultiLevelStatesCount: Número de estados ópticos discretos por símbolo (2^8 = 256).
	DefaultMultiLevelStatesCount = 256

	// ------------------------------------------------------------------------
	// 3. HIERARQUIA DE MEMÓRIA FOTÔNICA E PHOTONIC SSD
	// ------------------------------------------------------------------------

	// DefaultCacheL1LatencyPS: Latência de acesso ao Cache L1 (micro-anéis ressonadores) em ps (Alexoudi et al., 2020).
	DefaultCacheL1LatencyPS = 5.0

	// DefaultRamLoopLatencyPS: Latência de recirculação na Photonic RAM em ps (Yao, 1993).
	DefaultRamLoopLatencyPS = 96.73

	// DefaultCacheL1HitRate: Taxa empírica nominal de acerto (Hit Rate) do Cache L1.
	DefaultCacheL1HitRate = 0.92

	// DefaultPhotonicSsdCapacityTB: Capacidade de armazenamento do cubo de vidro (Photonic SSD) em Terabytes.
	DefaultPhotonicSsdCapacityTB = 100.0

	// DefaultPhotonicSsdThroughputTBps: Taxa de leitura paralela WDM do Photonic SSD em TB/s.
	DefaultPhotonicSsdThroughputTBps = 1.2

	// ------------------------------------------------------------------------
	// 4. GPU ÓPTICA (WDM RGB) E ACELERADOR TENSOR DE IA (MZI MESH)
	// ------------------------------------------------------------------------

	// DefaultGpuWdmChannelsCount: Número de canais espectrais RGB (635nm, 532nm, 450nm).
	DefaultGpuWdmChannelsCount = 3

	// DefaultAiTensorDensityTOPS: Densidade de computação do Tensor Core em TOPS/mm^2 (Xu et al., Nature 2021).
	DefaultAiTensorDensityTOPS = 11.0

	// DefaultAiTensorEfficiency: Eficiência energética do Tensor Core em TOPS/W (Shen et al., Nature Phot. 2017).
	DefaultAiTensorEfficiency = 100.0

	// ------------------------------------------------------------------------
	// 5. COMPUTADOR QUÂNTICO FOTÔNICO (LOQC - 298K)
	// ------------------------------------------------------------------------

	// DefaultQuantumQubitsCount: Número de qubits fotônicos dual-rail em temperatura ambiente.
	DefaultQuantumQubitsCount = 16

	// DefaultQuantumHomVisibilityPct: Visibilidade do poço de interferência quântica de 2 fótons (HOM) em % (Crespi et al., 2013).
	DefaultQuantumHomVisibilityPct = 99.4

	// DefaultQuantumCnotFidelityPct: Fidelidade da porta lógica quântica CNOT em % (Carolan et al., Science 2015).
	DefaultQuantumCnotFidelityPct = 98.7

	// ------------------------------------------------------------------------
	// 6. RUÍDO FÍSICO E ATENUAÇÃO DE PROPAGAÇÃO NO VIDRO
	// ------------------------------------------------------------------------

	// DefaultLaserRinDbHz: Ruído de intensidade relativa do laser contínuo (RIN) em dB/Hz.
	DefaultLaserRinDbHz = -155.0

	// DefaultPhaseNoiseSigmaRad: Desvio padrão do ruído de fase eletro-óptica em radianos.
	DefaultPhaseNoiseSigmaRad = 0.012

	// DefaultGlassLossDbPerCm: Perda por atenuação de propagação no guia de onda em sílica fundida (dB/cm).
	DefaultGlassLossDbPerCm = 0.2

	// DefaultMziPhaseErrorRad: Erro de deriva de fase nos interferômetros MZI do Tensor Core em radianos.
	DefaultMziPhaseErrorRad = 0.010
)

// OpticalParams encapsula todos os parâmetros físicos, geométricos e funcionais do SilicaCore.
type OpticalParams struct {
	RefractiveIndex   float64 // Índice de refração efetivo do substrato (ex: 1.4500 para SiO2)
	FastDistanceMM    float64 // d1: Distância em linha reta (mm)
	DelayedDistanceMM float64 // d0: Distância retarda/defletida (mm)
	LaserJitterFwhmPS float64 // Jitter temporal do laser em ps (FWHM)
	SpadJitterFwhmPS  float64 // Jitter do detector SPAD em ps (FWHM)
	TdcResolutionPS   float64 // Resolução LSB do TDC em ps
	WindowWidthPS     float64 // Largura da janela de tempo em ps

	// Motor Laser Contínuo CW (Estilo Noritsu) e Codificação M-ária
	LaserEngineMode       string // Modo de operação ("Continuous Wave (Noritsu CW Style - Always-ON)")
	EncodingFormat        string // Esquema de codificação ("M-ary Multi-Level WDM (Hexadecimal 4-bit / Byte 8-bit)")
	BitsPerSymbol         int    // Bits por símbolo óptico (8 bits = 1 Byte / símbolo)
	MultiLevelStatesCount int    // Número de estados discretos de fase/amplitude (256 para 8-bit Byte)

	// Memória Fotônica & Photonic SSD
	CacheL1LatencyPS          float64 // Latência do Cache L1 (Micro-anéis): <= 5ps (Alexoudi et al., 2020)
	RamLoopLatencyPS          float64 // Latência da RAM Fotônica (Loops de retardo): ~96.73ps (Yao, 1993)
	CacheL1HitRate            float64 // Taxa de acerto (Hit Rate) nominal do Cache L1 (0.92)
	PhotonicSsdThroughputTBps float64 // Taxa de leitura do Photonic SSD em TB/s (1.2 TB/s)
	PhotonicSsdCapacityTB     float64 // Capacidade do Photonic SSD em Terabytes por cubo (100 TB)

	// GPU Óptica & Photonic AI Tensor Core
	GpuWdmChannelsCount int     // Canais espectrais WDM da GPU (RGB: 3 canais)
	AiTensorDensityTOPS float64 // Densidade computacional do Tensor Core (11 TOPS/mm^2)
	AiTensorEfficiency  float64 // Eficiência energética do Tensor Core (100 TOPS/W)

	// Processador Quântico Fotônico LOQC (298K)
	QuantumQubitsCount      int     // Número de qubits dual-rail em temperatura ambiente (16 Qubits)
	QuantumHomVisibilityPct float64 // Visibilidade HOM em % (99.4%) [Crespi et al., 2013]
	QuantumCnotFidelityPct  float64 // Fidelidade da porta CNOT em % (98.7%) [Carolan et al., 2015]

	// Parâmetros Avançados de Ruído Físico e Atenuação
	LaserRinDbHz       float64 // Ruído de intensidade relativa do laser RIN (dB/Hz)
	PhaseNoiseSigmaRad float64 // Ruído de fase eletro-óptica em radianos
	GlassLossDbPerCm   float64 // Atenuação no vidro em dB/cm (0.2 dB/cm)
	MziPhaseErrorRad   float64 // Erro de fase MZI no Tensor Core em radianos
}

// DefaultParams retorna a configuração nominal padrão do SilicaCore utilizando as constantes definidas.
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

// FwhmToSigma converte Full Width at Half Maximum (FWHM) de uma distribuição Gaussiana para desvio padrão (sigma).
func FwhmToSigma(fwhm float64) float64 {
	return fwhm / (2.0 * math.Sqrt(2.0*math.Log(2.0)))
}

// ToFSimulator representa o motor principal de simulação física e modelo de ruído temporal.
type ToFSimulator struct {
	Params            OpticalParams
	MediumSpeedMMps   float64 // Velocidade da luz no substrato em mm/ps
	SpecificDelayPSmm float64 // Taxa de atraso de propagação em ps/mm
	T1NominalPS       float64 // Tempo de propagação nominal da linha rápida t1 (ps)
	T0NominalPS       float64 // Tempo de propagação nominal da linha de retardo t0 (ps)
	DeltaTNominalPS   float64 // Separação temporal t0 - t1 (ps)
	TotalSigmaPS      float64 // Jitter Gaussiano convoluído total (ps)
}

// NewSimulator inicializa o motor de simulação com os parâmetros ópticos fornecidos.
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

// SeparationMarginSigmas calcula a margem de separação temporal em múltiplos de sigma total.
func (s *ToFSimulator) SeparationMarginSigmas() float64 {
	return s.DeltaTNominalPS / s.TotalSigmaPS
}
