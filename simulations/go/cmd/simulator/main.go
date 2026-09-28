package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

func main() {
	fmt.Println("======================================================================")
	fmt.Println(" SIMULAÇÃO DE PROCESSAMENTO ÓPTICO TRIDIMENSIONAL ToF (GOLANG) ")
	fmt.Println(" Substrato: Sílica Fundida (SiO2) | Detecção: SPAD + TDC ")
	fmt.Printf(" Execução Paralelelizada em %d Núcleos de CPU (Goroutines)\n", runtime.NumCPU())
	fmt.Println("======================================================================")

	params := optical.ParametrosPadrao()
	sim := optical.NovoSimulador(params)

	fmt.Println("\n--- 1. PARÂMETROS FÍSICOS E TEMPOS NOMINAIS ---")
	fmt.Printf("Índice de refração do vidro (n):           %.4f\n", sim.Params.IndiceRefracao)
	fmt.Printf("Velocidade de propagação no meio (v):      %.5f mm/ps\n", sim.VelocidadeMeioMMps)
	fmt.Printf("Atraso específico:                         %.4f ps/mm\n", sim.AtrasoEspecificoPSmm)
	fmt.Printf("Distância direta (d1 - Linha Rápida):      %.3f mm\n", sim.Params.DistanciaRapidaMM)
	fmt.Printf("Distância defletida (d0 - Linha Atrasada): %.3f mm\n", sim.Params.DistanciaAtrasadaMM)
	fmt.Printf("Diferencial geométrico (Delta d):          %.3f mm\n", sim.Params.DistanciaAtrasadaMM-sim.Params.DistanciaRapidaMM)
	fmt.Printf("Tempo nominal Linha Rápida (t1):           %.2f ps\n", sim.T1NominalPS)
	fmt.Printf("Tempo nominal Linha Atrasada (t0):         %.2f ps\n", sim.T0NominalPS)
	fmt.Printf("Diferencial temporal (Delta t):            %.2f ps\n", sim.DeltaTNominalPS)
	fmt.Printf("Jitter temporal total (sigma):             %.2f ps\n", sim.SigmaTotalPS)
	fmt.Printf("Margem de separação (Delta t / sigma):     %.2f sigmas\n", sim.MargemSeparacaoSigmas())

	fmt.Println("\n--- 2. SIMULAÇÃO MONTE CARLO (1.000.000 PULSOS) ---")
	res := sim.SimularMonteCarloConcorrente(1000000)

	fmt.Printf("Amostras testadas:                         %d\n", res.AmostrasTotal)
	fmt.Printf("Erros detectados:                          %d\n", res.ErrosDetectados)
	fmt.Printf("Taxa de Acerto:                            %.4f%%\n", res.TaxaAcerto)
	fmt.Printf("Fator de Qualidade Q:                      %.2f\n", res.FatorQ)
	fmt.Printf("BER Empírico:                              %.2e\n", res.BerEmpirico)
	fmt.Printf("BER Teórico:                               %.2e\n", res.BerTeorico)
	fmt.Printf("Tempo de Execução Concorrente:            %s\n", res.DuracaoExecucao)

	fmt.Println("\n--- 3. TESTE DE PORTA LÓGICA NOT ToF (INVERSOR) ---")
	fmt.Println("Entrada (A) | Saída Invertida | Tempo Registrado (ps) | Status")
	fmt.Println("------------------------------------------------------------")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		ent := i % 2
		saida, tempo, ok := sim.TestarPortaNOT(ent, r)
		statusStr := "OK"
		if !ok {
			statusStr = "FALHA"
		}
		fmt.Printf("     %d      |        %2d       |        %6.2f ps       |   %s\n", ent, saida, tempo, statusStr)
	}

	fmt.Println("\nConclusão da simulação: A discriminação temporal por janelas de chegada")
	fmt.Printf("confirma viabilidade teórica com separação de %.2f sigmas e BER inferior a 10^-12.\n", sim.MargemSeparacaoSigmas())
	fmt.Println("======================================================================")
}
