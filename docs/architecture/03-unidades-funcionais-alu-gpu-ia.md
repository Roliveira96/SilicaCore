# Divisão Funcional dos Andares Volumétricos (ULA, IA, GPU e Memória)

## 1. Organização por Andares Espaciais (Eixo Z)

O cubo de sílica é dividido em quatro zonas funcionais ao longo da profundidade física:

- **Andar 1 (Base - Z = 0 a 5mm): Barramento Óptico Mestre**
  - Distribuição de relógio pulsado síncrono para todos os detectores SPAD.
  - Guias de onda para roteamento de controle entre andares.

- **Andar 2 (Central - Z = 5 a 15mm): Unidade Aritmética e Lógica (ULA)**
  - Execução de operações lógicas fundamentais (NOT, AND, OR, XOR) por tempo de voo (*Time-of-Flight*).
  - Comutadores ópticos baseados em efeito Kerr não-linear.

- **Andar 3 (Intermediário - Z = 15 a 20mm): Memória Recirculante em Linha de Atraso**
  - Pulsos retidos em cavidades reflexivas de anel fechado a $v = 0.20675\text{ mm/ps}$.
  - Leitura não-destrutiva com extração parcial de 5% por divisores de feixe.

- **Andar 4 (Topo - Z = 20 a 25mm): Acelerador Tensor de IA & Processamento Gráfico (GPU)**
  - Interferometria de malha Mesh de Mach-Zehnder para multiplicação paralela de matrizes.
  - Multiplexação por comprimento de onda (WDM RGB) para processamento paralelo de canais gráficos.
