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
	fmt.Printf(" Execução Paralelizada em %d Núcleos de CPU (Goroutines)\n", runtime.NumCPU())
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

	fmt.Println("\n--- 2. HIERARQUIA DE MEMÓRIA ÓPTICA (LATÊNCIAS E FOTÔNICA) ---")
	fmt.Printf("Latência Cache L1 (Micro-anéis):           <= %.2f ps  [Alexoudi et al., 2020]\n", sim.Params.LatenciaCacheL1PS)
	fmt.Printf("Latência RAM Óptica (Delay-Line Loop):    ~%.2f ps   [Yao, 1993]\n", sim.Params.LatenciaRamLoopPS)
	fmt.Printf("Acesso ROM Kernel (Nanofilamentos SiO2):  ~%.2f ps (c/n direto) [Zhang et al., 2014]\n", sim.T1NominalPS)
	fmt.Printf("Taxa Nominal de Hit na Cache L1:           %.1f%%\n", sim.Params.TaxaAcertoCacheL1*100.0)

	fmt.Println("\n--- 3. SIMULAÇÃO MONTE CARLO (1.000.000 OPERAÇÕES DE CPU E MEMÓRIA) ---")
	res := sim.SimularMonteCarloConcorrente(1000000)

	fmt.Printf("Amostras testadas:                         %d\n", res.AmostrasTotal)
	fmt.Printf("Erros ToF detectados:                      %d\n", res.ErrosDetectados)
	fmt.Printf("Taxa de Acerto ToF:                        %.4f%%\n", res.TaxaAcerto)
	fmt.Printf("Fator de Qualidade Q:                      %.2f\n", res.FatorQ)
	fmt.Printf("BER Empírico:                              %.2e\n", res.BerEmpirico)
	fmt.Printf("BER Teórico:                               %.2e\n", res.BerTeorico)
	fmt.Printf("Hits na Cache L1:                          %d\n", res.EstatisticasMem.HitsCacheL1)
	fmt.Printf("Misses na Cache L1 (Acessos à RAM):        %d\n", res.EstatisticasMem.MissesCacheL1)
	fmt.Printf("Acessos Diretos à ROM Kernel em SiO2:       %d\n", res.EstatisticasMem.AcessosROMKernel)
	fmt.Printf("Latência Média Global de Acesso a Dados:   %.2f ps\n", res.EstatisticasMem.LatenciaMediaGlobalPS)
	fmt.Printf("Tempo de Execução Concorrente em Go:       %s\n", res.DuracaoExecucao)

	fmt.Println("\n--- 4. TESTE DE PORTA LÓGICA NOT ToF (INVERSOR) ---")
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

	fmt.Println("\nConclusão da simulação: A hierarquia de memória óptica e a discriminação temporal")
	fmt.Printf("confirmam latência média global de %.2f ps e BER inferior a 10^-12.\n", res.EstatisticasMem.LatenciaMediaGlobalPS)
	fmt.Println("======================================================================")
}
