package optical

import (
	"math"
)

// CVacuo é a velocidade da luz no vácuo em metros por segundo (m/s).
const CVacuo = 299792458.0

// ParametrosOpticos encapsula as variáveis físicas e de hardware do sistema.
type ParametrosOpticos struct {
	IndiceRefracao      float64 // Índice de refração efetivo do substrato (ex: 1.4500 para SiO2)
	DistanciaRapidaMM   float64 // d1: Distância do percurso em linha reta (mm)
	DistanciaAtrasadaMM float64 // d0: Distância do percurso refletido/estendido (mm)
	JitterLaserFWHMps   float64 // Jitter temporal do pulso VCSEL em ps (FWHM)
	JitterSpadFWHMps    float64 // Jitter do detector SPAD em ps (FWHM)
	ResolucaoTDCps      float64 // Menor bit significativo (LSB) do TDC em ps
	LarguraJanelaPS     float64 // Largura da janela de amostragem temporal em ps
}

// ParametrosPadrao retorna a configuração nominal do processador ToF.
func ParametrosPadrao() ParametrosOpticos {
	return ParametrosOpticos{
		IndiceRefracao:      1.4500,
		DistanciaRapidaMM:   20.0,
		DistanciaAtrasadaMM: 40.675,
		JitterLaserFWHMps:   8.0,
		JitterSpadFWHMps:    25.0,
		ResolucaoTDCps:      5.0,
		LarguraJanelaPS:     35.0,
	}
}

// FwhmParaSigma converte a largura a meia altura (FWHM) de uma distribuição Gaussiana para desvio padrão (sigma).
func FwhmParaSigma(fwhm float64) float64 {
	return fwhm / (2.0 * math.Sqrt(2.0*math.Log(2.0)))
}

// SimuladorToF representa o motor de simulação contendo os tempos calculados e perturbações de jitter.
type SimuladorToF struct {
	Params               ParametrosOpticos
	VelocidadeMeioMMps   float64 // Velocidade no meio em mm/ps
	AtrasoEspecificoPSmm float64 // Taxa de atraso em ps/mm
	T1NominalPS          float64 // Tempo nominal da Linha Rápida (ps)
	T0NominalPS          float64 // Tempo nominal da Linha Atrasada (ps)
	DeltaTNominalPS      float64 // Diferença temporal t0 - t1 (ps)
	SigmaTotalPS         float64 // Jitter Gaussiano total convoluído (ps)
}

// NovoSimulador inicializa o motor de simulação com os parâmetros fornecidos.
func NovoSimulador(params ParametrosOpticos) *SimuladorToF {
	// Velocidade da luz no meio (mm/ps): (c / n) * 1e3 mm / 1e12 ps = (c / n) * 1e-9 mm/ps
	vMeio := (CVacuo / params.IndiceRefracao) * 1e-9
	atrasoEspec := 1.0 / vMeio

	t1 := params.DistanciaRapidaMM * atrasoEspec
	t0 := params.DistanciaAtrasadaMM * atrasoEspec
	deltaT := t0 - t1

	sigmaLaser := FwhmParaSigma(params.JitterLaserFWHMps)
	sigmaSpad := FwhmParaSigma(params.JitterSpadFWHMps)
	// O erro de quantização uniforme do TDC [-LSB/2, LSB/2] tem desvio padrão LSB / sqrt(12)
	sigmaTDC := params.ResolucaoTDCps / math.Sqrt(12.0)

	sigmaTotal := math.Sqrt(sigmaLaser*sigmaLaser + sigmaSpad*sigmaSpad + sigmaTDC*sigmaTDC)

	return &SimuladorToF{
		Params:               params,
		VelocidadeMeioMMps:   vMeio,
		AtrasoEspecificoPSmm: atrasoEspec,
		T1NominalPS:          t1,
		T0NominalPS:          t0,
		DeltaTNominalPS:      deltaT,
		SigmaTotalPS:         sigmaTotal,
	}
}

// MargemSeparacaoSigmas calcula a razão de separação temporal em múltiplos de sigma total.
func (s *SimuladorToF) MargemSeparacaoSigmas() float64 {
	return s.DeltaTNominalPS / s.SigmaTotalPS
}
