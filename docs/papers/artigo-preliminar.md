# Race Logic Fotônica para Menor Caminho em Plataforma Si₃N₄/TFLN: Orçamento Físico e Validação Estatística por Simulação

**Autor:** Ricardo Oliveira — UTFPR *[afiliação e coautores a confirmar]*  
**Status:** rascunho preliminar v0.1 (29/09/2026, revisão de baselines e da porta ToF) — alvo: SBESC / WSCAD (SBC) ou IEEE Photonics Conference  
**Código e dados:** [github.com/Roliveira96/SilicaCore](https://github.com/Roliveira96/SilicaCore) (`simulations/go`, `simulations/results`)

---

## Resumo

Propostas de processadores ópticos frequentemente reportam ganhos derivados apenas do tempo de voo da luz, sem orçamento de potência, sem limites de detecção e sem estatística de erro. Este trabalho parte de uma arquitetura conceitual de computação por tempo de voo (SilicaCore) e a submete a um modelo físico com parâmetros da literatura de 2023–2026. Mostramos que o roteamento por espelhos internos em bloco de sílica perde ~46,6 dB por porta por difração e não tem comutador em picossegundos, enquanto guias de Si₃N₄ com chaves de niobato de lítio em filme fino (TFLN) perdem ~1,3 dB por porta e permitem 7 portas em cascata antes de regenerar o sinal. Sobre essa plataforma, mapeamos *race logic* — computação em que o valor é o tempo de chegada de um pulso — para resolver menor caminho em grafos: arestas são linhas de atraso programáveis por chaves Sb₂Se₃ e nós detectam a primeira chegada e reemitem. Uma campanha Monte Carlo com 10 chips simulados × 10⁴ consultas (2,55×10⁷ distâncias por configuração) em um mapa 16×16 mostra que a unidade de atraso de 50 ps produz 1,23×10⁻⁴ erros por distância, enquanto 100 ps produz 0 erros (limite superior de 1,2×10⁻⁷ com 95% de confiança). O erro medido por número de saltos segue a previsão gaussiana de ruído acumulado. Cada consulta leva 42,2 ns, dominados pela leitura eletrônica dos tempos (30,7 ns), e o circuito ocupa 648 mm², dentro de um retículo de litografia. O ganho não é O(1): cresce com a distância máxima do grafo e é limitado pela leitura e pela área. Mapas maiores podem ser compostos por vários chips numa corrida única (64×64: 30,7 ns origem→destino com unidade de 150 ps e 0 erros), desde que o acoplamento entre chips fique em ≤ 1,5–2 dB por face. Contra o melhor algoritmo eletrônico para pesos inteiros (fila de baldes de Dial), a vantagem estimada sobre uma CPU atual é de ~42× por consulta; num mapa fixo, uma tabela pré-calculada responde mais rápido que o chip. O nicho favorável é o de mapas que mudam entre lotes de consultas: ~34× com 100 consultas por mudança, sob a premissa de reprogramação em 1 µs. Duas verificações adicionais limitam a proposta: uma race logic síncrona em CMOS a 3 GHz resolve o mesmo mapa em ~40 ns com área e energia ordens de grandeza menores, e a diafonia entre voltas das espirais exige passo ≥ 4 µm, o que leva o bloco 16×16 a ~864 mm² (duas camadas de guias para caber no retículo). Nesta forma, o RL-16 é um veículo para validar a física da temporização óptica, não um acelerador superior à eletrônica.

**Palavras-chave:** computação fotônica, race logic, niobato de lítio em filme fino, nitreto de silício, menor caminho, materiais de mudança de fase.

## Abstract

*[Tradução para o inglês a ser feita após revisão do resumo.]*

---

## 1. Introdução

A escala dos processadores CMOS esbarra no custo energético de mover dados e na dissipação resistiva das interconexões. A fotônica integrada é candidata natural para reduzir esse custo, e aceleradores fotônicos de multiplicação matriz-vetor já foram demonstrados em escala de sistema [Ahmed et al. 2025; Hua et al. 2025]. Esses mesmos trabalhos, porém, mostram que a eficiência do sistema completo (≈0,84 TOPS/W em [Ahmed et al. 2025]) fica muito abaixo da eficiência do núcleo óptico isolado, porque lasers, conversores e controle eletrônico dominam o consumo.

Arquiteturas de "computação na velocidade da luz" costumam derivar desempenho do tempo de voo em milímetros de guia (dezenas de picossegundos) e tratá-lo como período de clock. Essa leitura ignora três restrições: (i) a lógica óptica precisa de cascateamento e restauração de nível [Miller 2010]; (ii) cada símbolo temporal precisa conter as janelas de decisão, e o jitter acumulado define a taxa de erro; (iii) o resultado precisa ser lido pela eletrônica.

Este artigo tem três contribuições:

1. **Orçamento físico de roteamento** comparando espelhos internos em bloco, guias gravados por laser de femtossegundo em vidro e guias Si₃N₄ com chaves TFLN heterogêneas (Seção 5.1).
2. **Correção da métrica de decisão temporal** de uma arquitetura por tempo de voo: fator Q, taxa de erro de bit e taxa real por canal (Seção 5.2).
3. **Mapeamento e validação estatística de race logic fotônica** para menor caminho, com leitura, área e potência incluídas no modelo (Seções 3.3 e 5.3–5.4).

## 2. Trabalhos Relacionados

**Race logic.** Madhavan, Sherwood e Strukov [2014] propuseram codificar valores no tempo de subida de sinais digitais e computar por atrasos: `MIN` é a primeira chegada, `MAX` a última, somar uma constante é atrasar. A formulação resolve programação dinâmica (alinhamento de sequências, menor caminho) de forma nativa. Este trabalho transpõe a ideia para pulsos ópticos em guias.

**Aceleradores fotônicos.** Lightmatter [Ahmed et al. 2025] integrou quatro núcleos tensoriais 128×128 em encapsulamento 3D e executou ResNet e BERT com precisão próxima de FP32, com 65,5 TOPS a 78 W elétricos. O PACE [Hua et al. 2025] integrou >16.000 componentes e resolveu problemas de Ising com latência de ~5 ns. O Taichi [Xu et al. 2024] reporta 160 TOPS/W no chiplet.

**Lógica óptica de propósito geral.** Kissner et al. [2024] demonstraram uma CPU totalmente óptica com registradores em linha de atraso, memória PCM e regeneração 2R, operando abaixo de 1 GHz. É o trabalho mais próximo em espírito deste.

**Materiais.** TFLN heterogêneo sobre Si₃N₄ atinge transições <0,1 dB e guias <0,1 dB/cm [Churaev et al. 2023]. Chaves de Sb₂Se₃ sobre Si₃N₄ atingem >1,4×10⁸ ciclos, 25 dB de extinção e >6 bits multinível [Yu et al. 2026].

## 3. Arquitetura

### 3.1 Substrato e roteamento

A arquitetura original previa um cubo maciço de sílica fundida com feixes livres desviados por espelhos internos, moduladores acusto-ópticos e reflexão interna total induzida por campo elétrico. A validação física (Seção 5.1) levou a três funções separadas:

| Função | Velocidade | Implementação adotada |
| :--- | :--- | :--- |
| Roteamento fixo | — | Guias Si₃N₄ multicamada sobre SiO₂ (curvas de 50 µm, acopladores verticais de 0,01 dB) [Shang et al. 2015] |
| Reconfiguração | µs | Chaves não-voláteis Sb₂Se₃ [Yu et al. 2026] |
| Chaveamento por bit | ps | Chaves eletro-ópticas TFLN [Churaev et al. 2023] |

O substrato é agnóstico à geometria externa: o atraso depende do comprimento do guia ($L = v\,\Delta t$), e o formato plano segue o padrão de wafers de 200/300 mm.

### 3.2 Lógica por tempo de voo

Um pulso segue a linha rápida ($t_1$) ou é desviado por uma chave TFLN para a linha atrasada ($t_0$). O detector amostra duas janelas de largura $W$; a decisão tem fator

$$Q = \frac{\Delta t}{2\,\sigma_{\text{total}}}, \qquad \text{BER} = \tfrac{1}{2}\,\text{erfc}\!\left(\frac{Q}{\sqrt{2}}\right),$$

com $\sigma_{\text{total}}$ a convolução gaussiana dos jitters de laser, detector e quantização do TDC.

### 3.3 Mapeamento de race logic para menor caminho

- **Aresta** $(u, v)$ com peso inteiro $w \in [1, 15]$: linha de atraso de 4 estágios binários de espiral Si₃N₄, cada estágio selecionado por um par de chaves Sb₂Se₃. Atraso programado $w\,\tau - L_{\text{nó}}$, com $\tau$ a unidade de atraso.
- **Nó**: um fotodiodo por aresta de entrada, com OR eletrônico (evita a perda de um combinador passivo), comparador e um modulador TFLN por aresta de saída. O nó dispara uma única vez, na primeira chegada, após a latência de regeneração $L_{\text{nó}}$.
- **Compensação de latência**: descontar $L_{\text{nó}}$ do atraso de cada aresta faz o tempo de disparo do nó $v$ valer exatamente $d(s, v)\,\tau$, independentemente do número de saltos. Sem essa compensação, caminhos com mais saltos seriam penalizados e a corrida retornaria a resposta errada (verificado em teste unitário).
- **Leitura**: um TDC por nó captura o tempo de disparo; a distância decodificada é $\hat d(v) = \text{round}(t_v / \tau)$.
- **Modulador por aresta**: com um único modulador por nó, a divisão de potência para 4 saídas custa 6 dB sobre o pulso de dados e estoura a margem. Com um modulador por aresta, a divisão ocorre no laser CW antes da modulação, que tem orçamento próprio.

## 4. Metodologia

O modelo foi implementado em Go (`simulations/go/pkg/optical`) e é reproduzível com `go test ./...`, `go run ./cmd/simulator` e `go run ./cmd/racestats`.

**Parâmetros físicos** (Tabela 1). Valores sem número publicado estão marcados como premissa (P) no código.

| Parâmetro | Valor | Fonte |
| :--- | :--- | :--- |
| Potência por canal / sensibilidade do fotodiodo | 0 dBm / −10 dBm (margem de 10 dB) | — |
| Perda de propagação Si₃N₄ | 0,1 dB/cm | Churaev et al. 2023 |
| Índice de grupo Si₃N₄ | 2,0 | — |
| Perda de chave TFLN (com transições) | 1,0 dB | Churaev et al. 2023 |
| Perda de chave Sb₂Se₃ | 0,25 dB | Yu et al. 2026 (0,1–0,5 dB/π) |
| Latência de regeneração do nó | 20 ps | P |
| Jitter por nó (rms) | 1,5 ps | P |
| Erro estático de atraso por aresta (rms, pós-calibração) | 0,5 ps | P |
| TDC / enlace de leitura | 12 bits / 100 Gb/s | P |
| Espaçamento entre espirais | 3 µm | P |
| Tempo de programação das chaves Sb₂Se₃ | 1 µs | P |

**Grafo.** Mapa de jogo em grade 4-vizinhos, custos de terreno uniformes em [1, 15], simétricos.

**Simulação da corrida.** Simulação por eventos: cada nó dispara uma vez, na primeira chegada mais $L_{\text{nó}}$ mais jitter gaussiano; cada aresta soma seu atraso programado e seu erro estático. O erro estático é sorteado **uma vez por chip**, como num dispositivo fabricado.

**Campanha estatística.** Para cada unidade de atraso $\tau \in \{35, 50, 75, 100\}$ ps: 10 chips × 10⁴ consultas de origem aleatória, com todas as 255 distâncias de cada consulta comparadas com Dijkstra exato [Dijkstra 1959]. Limites superiores de 95% por regra de 3/N (sem erros) ou intervalo de Wilson.

**Previsão teórica.** Um nó a $h$ saltos erra quando o ruído acumulado excede $\tau/2$: $P_{\text{err}}(h) = \text{erfc}\!\left(\frac{\tau/2}{\sigma\sqrt{2h}}\right)$, com $\sigma = \sqrt{\sigma_{\text{nó}}^2 + \sigma_{\text{aresta}}^2}$.

**Baselines eletrônicos.** Três implementações em Go, thread única, em Intel Core i3-3217U (2012, 1,8 GHz), todas devolvendo as 256 distâncias de uma origem (`go run ./cmd/dijkstrabench`): (i) Dijkstra com heap binário tipado; (ii) Dijkstra com fila de baldes circular [Dial 1969], o algoritmo adequado a pesos inteiros pequenos; (iii) cópia de uma linha de uma tabela com todas as distâncias, pré-calculada uma vez (mapa fixo). Para CPUs atuais, os tempos são escalados pela pontuação single-core do Geekbench 6 (i3-3217U: 307; Apple M5 Max: 4349), uma estimativa de primeira ordem para um problema que cabe na cache.

**Verificação eletromagnética.** As perdas de curva e de transição Si₃N₄→TFLN do modelo foram verificadas por FDTD no Meep 1.34, em 2D pelo método do índice efetivo (a curva no plano do chip; a transição num corte vertical, com a largura do taper convertida em índice de camada). A transmissão é medida por decomposição no modo fundamental; as geometrias usam suavização subpixel e cada caso foi repetido com resolução maior (20→30 px/µm nas curvas, 50→70 px/µm na transição), com variação ≤ 8%.

## 5. Resultados

### 5.1 Orçamento de roteamento

| Plataforma | Perda por porta | Portas em cascata | Comutador em ps |
| :--- | :---: | :---: | :---: |
| Bloco SiO₂ + espelhos internos, feixe livre | 46,6 dB (45,8 dB de difração) | 0 | Não (AOM: 16,8 ns; sem efeito Pockels em SiO₂) |
| Guias fs em vidro + micro-espelhos TIR + TFLN | 6,2 dB | 1 | Sim |
| **Si₃N₄ multicamada + TFLN** | **1,3 dB** | **7** | **Sim** |

Um feixe com cintura de 5 µm chega com raio de ~2,8 mm após 40,7 mm de sílica. A reflexão interna total induzida por campo em sílica exigiria incidência rasante de 0,002°.

### 5.1.1 Verificação eletromagnética (FDTD 2D)

| Estrutura | Resultado | Premissa do modelo |
| :--- | :--- | :--- |
| Curva de 90°, Si₃N₄ 800 nm × 0,7 µm (monomodo) | 0,138 / 0,029 / 0,012 / 0,0026 dB para R = 10 / 20 / 30 / 50 µm | 0,01 dB por curva |
| Curva de 90°, Si₃N₄ 800 nm × 1,2 µm (multimodo) | 0,10 dB em R = 20 µm | — |
| Transição Si₃N₄ → filme de LN | 1,05 / 0,14 / 0,0017 / 0,0005 dB para taper de 5 / 10 / 25 / 100 µm | 0,1 dB por transição |

Os resultados indicam, em aproximação planar, que a premissa de 0,01 dB por curva é conservadora para raios ≥ ~31 µm no guia monomodo, e que a transição satura a partir de ~25 µm. O guia multimodo perde ~3,5× mais na mesma curva, pela conversão para o segundo modo lateral, o que fixa o guia de 0,7 µm como escolha de projeto. O modelo mantém 0,1 dB por transição para cobrir efeitos 3D e de fabricação fora do alcance da simulação 2D; esse valor é compatível com o medido em chips reais [Churaev et al. 2023].

### 5.2 Decisão temporal

Com $\Delta t = 100$ ps e $\sigma_{\text{total}} = 11,24$ ps, $Q = 4,45$ e BER ≈ 4,3×10⁻⁶ (Monte Carlo com 10⁶ amostras: ~2×10⁻⁵). A razão $\Delta t/\sigma = 8,9$, anteriormente reportada como margem com BER < 10⁻¹², é o dobro de Q. BER de 10⁻¹² exige $Q = 7,03$, isto é, $\sigma_{\text{total}} \le 7,1$ ps. O slot mínimo $\Delta t + W \approx 195$ ps limita a taxa a ~5,1 GHz por canal.

A diferença entre a BER teórica (4,3×10⁻⁶) e a de Monte Carlo (~2×10⁻⁵) vem do critério de decisão: a fórmula acima usa um limiar no meio entre $t_1$ e $t_0$, enquanto a simulação rejeita pulsos fora da janela de ±47,5 ps. Nesse caso contam as duas caudas, e a taxa esperada é $\text{erfc}\!\left(\tfrac{W/2}{\sigma\sqrt2}\right) \approx 2{,}4\times10^{-5}$.

**Estes números são do conceito original** (caminhos em sílica com $n = 1{,}45$, laser de 8 ps FWHM, SPAD de 25 ps FWHM e TDC de 5 ps). Na plataforma adotada, com espirais de Si₃N₄ ($n_g = 2{,}0$) e o mesmo fotodiodo do nó de race logic (1,5 ps rms, premissa), a porta foi recalculada (`go run ./cmd/tofplatform`):

| Porta ToF | Caminhos $d_1$ / $d_0$ | $\sigma_{\text{total}}$ | $Q$ ($\Delta t$ = 100 ps) | $\Delta t$ para BER 10⁻¹² | Taxa por canal (só temporização) |
| :--- | :---: | :---: | :---: | :---: | :---: |
| Conceito original (sílica + SPAD) | 20,0 / 40,7 mm | 11,24 ps | 4,45 | 158 ps | 5,1 GHz |
| **Plataforma (Si₃N₄ + fotodiodo InGaAs)** | **14,5 / 29,5 mm** | **2,12 ps** | **23,5** | **30 ps (4,5 mm)** | **17 GHz** |

A taxa da plataforma é um limite de temporização; a banda dos moduladores, dos detectores e do TDC ainda precisa entrar na conta.

### 5.3 Correção do race logic

| $\tau$ (ps) | Q no pior caminho (32 saltos) | Distâncias erradas / 2,55×10⁷ | Taxa (limite 95%) | Previsto | Consultas com erro |
| :---: | :---: | :---: | :---: | :---: | :---: |
| 35 | 1,96 | 101.390 | 3,98×10⁻³ (4,00×10⁻³) | 3,68×10⁻³ | 17,6% |
| 50 | 2,80 | 3.136 | 1,23×10⁻⁴ (1,27×10⁻⁴) | 1,08×10⁻⁴ | 0,83% |
| 75 | 4,19 | 3 | 1,2×10⁻⁷ (3,5×10⁻⁷) | 9,3×10⁻⁸ | 0,002% |
| **100** | **5,59** | **0** | **< 1,2×10⁻⁷** | 2,0×10⁻¹¹ | **< 3×10⁻⁵** |

![Erro de decodificação por número de saltos](figuras/race_logic_erro_por_salto.svg)

**Figura 1.** Taxa de erro por distância em função do número de saltos no menor caminho. Linhas: previsão gaussiana $P_{\text{err}}(h)$; pontos: taxa medida (apenas saltos com ≥1 erro). Com 75 ps houve só 3 erros, então os pontos têm grande incerteza. Com 100 ps não houve erros. Gerada por `simulations/results/plot_race_logic_hops.py`.

Dados da Figura 1 para $\tau = 50$ ps:

| Saltos | Distâncias | Medido | Previsto |
| :---: | ---: | :---: | :---: |
| 12 | 1.530.455 | 7,2×10⁻⁶ | 5,0×10⁻⁶ |
| 16 | 1.272.043 | 9,4×10⁻⁵ | 7,7×10⁻⁵ |
| 20 | 627.896 | 4,8×10⁻⁴ | 4,1×10⁻⁴ |
| 24 | 182.707 | 1,31×10⁻³ | 1,25×10⁻³ |
| 28 | 31.964 | 2,60×10⁻³ | 2,81×10⁻³ |
| 32 | 1.558 | 5,1×10⁻³ | 5,2×10⁻³ |

O erro medido segue a previsão de ruído acumulado. O excesso agregado de ~14% é esperado: quando dois caminhos têm comprimentos quase iguais, o nó dispara pelo mais adiantado dos dois ruídos, o que desloca o tempo de disparo para cedo. Uma amostra piloto de 200 consultas não mostrou erros com 50 ps, o que ilustra por que amostras pequenas não bastam para afirmar correção.

### 5.4 Desempenho, leitura e área ($\tau = 100$ ps)

| Mapa | Corrida | Leitura TDC | Total | Dijkstra com heap (Go, i3-3217U) | Área das espirais |
| :--- | ---: | ---: | ---: | ---: | :--- |
| 8×8 | 6,0 ns | 7,7 ns | 13,7 ns | ~20 µs | 151 mm² |
| 16×16 | 11,5 ns | 30,7 ns | 42,2 ns | ~50–100 µs | 648 mm² |
| 32×32 | 20,0 ns | 122,9 ns | 142,9 ns | ~0,3 ms | 2.677 mm² |
| 64×64 | 41,9 ns | 491,5 ns | 533,5 ns | ~2 ms | 10.879 mm² |

O circuito 16×16 usa 960 fotodiodos, 960 moduladores TFLN, 7.680 chaves Sb₂Se₃ e 256 TDCs; a aresta mais longa tem 221,8 mm de espiral e perde 5,2 dB (margem de 10 dB). As áreas da tabela usam passo de espiral de 3 µm; a verificação de diafonia (Seção 5.4.2) pede ≥ 4 µm, o que multiplica as áreas por 4/3 (16×16: ~864 mm², acima do retículo de 858 mm²; com duas camadas de guias, ~432 mm²).

### 5.4.1 Baselines eletrônicos e o nicho real

| Baseline (uma origem, todas as distâncias, 16×16) | i3-3217U (medido) | Ganho da corrida (42,2 ns) | CPU atual (estimado, M5 Max) | Ganho da corrida |
| :--- | ---: | ---: | ---: | ---: |
| Dijkstra com heap binário | 65,5 µs | 1553× | 4,6 µs | ~110× |
| **Dijkstra com fila de baldes (Dial)** | **25,4 µs** | **601×** | **1,79 µs** | **~42×** |
| Tabela pré-calculada, cópia de uma linha (mapa fixo) | 0,11 µs | 2,7× | ~8 ns | **~0,2× (a CPU vence)** |

A tabela de todas as distâncias do 16×16 ocupa 256 KB e é calculada em 6,5 ms no i3 (~0,46 ms estimado numa CPU atual). **Num mapa fixo, ela responde mais rápido que o chip.** A corrida só compensa quando o mapa muda entre consultas. Com reprogramação do chip em 1 µs (premissa), contra a melhor opção da CPU atual (Dial por consulta ou recalcular a tabela):

| Consultas por mudança do mapa | Chip (reprogramar + consultas) | CPU atual (melhor opção) | Ganho do chip |
| ---: | ---: | ---: | ---: |
| 1 | 1,04 µs | 1,79 µs | 1,7× |
| 10 | 1,42 µs | 17,9 µs | 12,6× |
| 100 | 5,22 µs | 179 µs | 34× |
| 1000 | 43,2 µs | 464 µs | 10,7× |
| 100000 | 4,22 ms | 1,25 ms | 0,3× |

O ganho depende diretamente do tempo de reprogramação: se a gravação das chaves Sb₂Se₃ levar 10 µs em vez de 1 µs, o caso de uma consulta por mudança passa a favorecer a CPU.

### 5.4.2 Diafonia entre voltas das espirais

Voltas vizinhas de uma espiral são guias paralelos e idênticos, portanto em casamento de fase: trocam potência com comprimento de acoplamento $L_c = \lambda / 2(n_{\text{par}} - n_{\text{ímpar}})$. Os supermodos par e ímpar do guia de 800 nm × 0,7 µm foram calculados pelo índice efetivo com matriz de transferência (`simulations/fdtd/spiral_crosstalk.py`) e conferidos em FDTD 2D (`coupler_2d.py`: $L_c$ de 21,2 e 44,8 µm contra 21,9 e 46,2 µm analíticos, para vãos de 0,4 e 0,6 µm). A luz que passa de uma volta para a seguinte pula aquela volta e chega ao nó **adiantada**; um eco acima do limiar do comparador dispara o nó antes da hora.

| Passo da espiral | $L_c$ | Pior volta | Eco adiantado total (aresta de 222 mm) |
| :---: | ---: | ---: | ---: |
| 3,0 µm (premissa do modelo) | 27 mm | −14 dB | ~0 dB (a espiral embaralha o pulso) |
| 3,5 µm | 173 mm | −30 dB | −13,6 dB |
| **4,0 µm** | **1,1 m** | **−46 dB** | **−29,6 dB** |
| 5,0 µm | 47 m | −77 dB | −62 dB |

O passo de 3 µm usado nas áreas é inviável: o guia de alto confinamento ainda tem cauda evanescente suficiente para acoplar ao longo de centenas de milímetros. Com passo de 4 µm o eco fica ~30 dB abaixo do pulso principal. Esse é o novo valor de projeto, e ele aumenta a área em 4/3.

### 5.5 Composição multi-chip (mapas maiores que um retículo)

Mapas maiores que 16×16 foram compostos por blocos de 16×16, cada um em um chip (`go run ./cmd/racemultichip`). Enlace entre chips: 1,5 dB por face (duas faces por cruzamento), 25 ps de voo compensados no atraso programado, 1 ps de jitter e 1 ps de erro estático rms por cruzamento (premissas). Dois modos foram comparados em 3.000 consultas origem→destino por configuração.

**Modo exato (uma corrida óptica atravessando os chips).** Os erros foram contados em todas as distâncias de cada corrida:

| Mapa | Unidade | Chips | Saltos máx. | Perda na aresta de borda | Erro por distância (limite 95%) | Bordas no caminho | Origem→destino | Todas as distâncias |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | ---: | ---: |
| 32×32 | 100 ps | 4 | 61 | 8,2 dB | 0 (< 9,8×10⁻⁷) | 1,2 | 10,9 ns | 50,9 ns |
| 64×64 | 100 ps | 16 | 120 | 8,2 dB | **1,52×10⁻⁴** | 2,8 | 20,5 ns | 72,6 ns |
| 64×64 | 150 ps | 16 | 120 | 9,3 dB | 0 (< 2,4×10⁻⁷) | 2,8 | 30,7 ns | 93,5 ns |

A leitura "todas as distâncias" é paralela: cada chip lê seus 256 nós pelo próprio enlace (30,7 ns). O orçamento da aresta de borda fecha até 2,0 dB por face com 100 ps (9,2 dB) e estoura com 2,5 dB (10,2 dB).

**Conflito de escala.** Mapas maiores têm caminhos com mais saltos, e o jitter acumulado exige unidade de atraso maior: com 100 ps o 64×64 erra 1,52×10⁻⁴ por distância. Com 150 ps os erros somem, mas (i) a espiral de cada aresta cresce 1,5×, o bloco 16×16 passa a ocupar 971 mm² e deixa de caber no retículo (um bloco 12×12 cabe, com 534 mm²), e (ii) a aresta de borda chega a 9,3 dB, o que só fecha a margem com acoplamento ≤ 1,5 dB por face. **Mapas maiores exigem unidade maior, que força blocos menores e mais cruzamentos entre chips.**

**Modo hierárquico (estilo HPA\*).** Cada chip calcula uma vez, com a própria corrida óptica, as distâncias internas entre pontos de passagem na borda; a consulta faz duas corridas locais e um Dijkstra eletrônico sobre esses pontos. Dispensa enlaces ópticos entre chips:

| Mapa | Espaçamento dos pontos de passagem | Pontos | Rotas ótimas | Excesso médio | Excesso p99 | Excesso máx. | Tempo por consulta* |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | ---: |
| 32×32 | 8 | 8 | 30,7% | 10,8% | 69% | 600% | ~15 µs |
| 32×32 | 4 | 16 | 38,6% | 6,5% | 41% | 310% | ~30 µs |
| 32×32 | 2 | 31 | 57,1% | 3,0% | 25% | 267% | ~73 µs |
| 64×64 | 8 | 48 | 9,2% | 13,2% | 53% | 1.160% | ~67 µs |
| 64×64 | 4 | 96 | 10,0% | 9,8% | 36% | 420% | ~188 µs |
| 64×64 | 2 | 187 | 26,7% | 4,1% | 19% | 150% | ~464 µs |

\*Dominado pela busca eletrônica no grafo abstrato, medida em Go no i3-3217U.

Com todos os cruzamentos como pontos de passagem e sem ruído, o modo hierárquico reproduz exatamente o Dijkstra (teste unitário), o que confirma que o excesso vem da esparsidade dos pontos. Os excessos máximos grandes ocorrem em pares origem-destino próximos, separados por uma borda, cujo ponto de passagem mais próximo fica distante. Nenhuma estimativa ficou abaixo do ótimo.

**Conclusão da composição.** O modo exato preserva a vantagem da corrida óptica (dezenas de ns) desde que o acoplamento entre chips fique em ≤ 1,5–2 dB por face e a unidade de atraso cresça com a profundidade do mapa. O modo hierárquico evita óptica entre chips, mas perde exatidão e passa a ser limitado pela busca eletrônica.

## 6. Discussão e Limitações

1. **Não é O(1).** O tempo de corrida cresce com a maior distância do grafo ($D\,\tau$), e a leitura cresce com o número de nós. A partir de 8×8, a leitura eletrônica domina.
2. **Baseline.** Contra o melhor algoritmo para pesos inteiros (Dial), o ganho é de ~600× numa CPU de 2012 e de ~42× estimado numa CPU atual (Seção 5.4.1). Num mapa fixo, uma tabela pré-calculada vence o chip. A* [Hart et al. 1968] ajuda quando só um destino interessa, mas a corrida entrega todas as distâncias. Uma comparação com GPU em lote está pendente.
3. **Área e composição.** Estágios binários fazem cada aresta carregar a espiral completa, limitando o bloco a 16×16 com 100 ps. Acima disso, a composição exata entre chips funciona (Seção 5.5), mas mapas mais profundos exigem unidade maior, blocos menores e acoplamento entre chips ≤ 1,5–2 dB por face; a composição hierárquica perde exatidão (excesso médio de 3–13%).
4. **Premissas.** Latência e jitter do nó, erro estático por aresta, tempo de programação e taxa de leitura não têm valor publicado para esta integração e precisam de medição.
5. **Modelo.** O ruído é gaussiano e independente por salto; deriva térmica lenta (1,8–3,7 rad/K de fase, irrelevante para ToF mas relevante para atrasos longos) e perdas por cruzamento em layout real não foram modeladas. As perdas de curva e de transição foram verificadas apenas em FDTD 2D (índice efetivo); a verificação 3D e a medição em chip ficam pendentes.
6. **Nicho.** O caso favorável é um mapa que muda entre lotes de consultas (p. ex. custos de terreno ou de tráfego atualizados várias vezes por segundo, com dezenas a milhares de consultas entre atualizações). Com mapa fixo, a tabela pré-calculada é melhor; com uma consulta por mudança, o ganho depende do tempo real de reprogramação.
7. **Race logic em CMOS.** A race logic original é eletrônica e síncrona [Madhavan et al. 2014]: uma unidade de atraso é um ciclo de relógio, cada aresta é um contador e cada nó é um latch. O resultado é exato, sem erro de temporização. Para o mesmo mapa 16×16 (`go run ./cmd/cmosrace`), a corrida leva em média 116 ciclos:

| Implementação | Corrida | Total com leitura | Área | Energia por consulta |
| :--- | ---: | ---: | ---: | ---: |
| RL-16 fotônico (simulado) | 11,4 ns | 42,1 ns (enlace de 100 Gb/s) | 648 mm² (864 mm² com passo de 4 µm) | 386 nJ |
| CMOS síncrono a 3 GHz | 38,8 ns | 40,1 ns (barramento no próprio die) | ~0,007 mm² (~7100 flip-flops) | ~0,8 nJ |
| CMOS síncrono a 4 GHz | 29,1 ns | 30,1 ns | ~0,007 mm² | ~0,8 nJ |

Área e energia do CMOS são estimativas de ordem de grandeza (1 µm² e 1 fJ por flip-flop por ciclo, num processo da classe de 7 nm). Mesmo que errem por 10×, a conclusão não muda: **a corrida fotônica é ~3× mais rápida só na fase da corrida (100 ps por unidade contra 250–330 ps de um ciclo), e perde por ordens de grandeza em área e energia.** Para que a óptica ganhe, a unidade de atraso teria de cair muito abaixo de um ciclo de relógio, o que exige jitter por salto bem menor que o assumido, ou a aplicação teria de pedir atrasos longos e estáveis que o CMOS não oferece.
8. **Riscos não modelados.** (i) Diafonia: a Seção 5.4.2 mostra que o passo de 3 µm é inviável e que 4 µm deixa o eco em ~−30 dB; o valor precisa ser medido no protótipo G1. (ii) Programação: gravar as 7680 chaves Sb₂Se₃ em 1 µs, todas em paralelo, é otimista; a cristalização costuma exigir pulsos mais longos, e o pico de potência dos aquecedores não está no modelo. (iii) Latência do nó: 20 ps é otimista para fotodiodo, comparador e driver; o resultado continua válido enquanto $L_{\text{nó}} < \tau$, porque a latência é descontada de cada aresta, mas o jitter do nó pode crescer com ela.

## 7. Conclusão e Trabalhos Futuros

Um modelo físico com parâmetros publicados descarta o roteamento por espelhos internos em bloco e aponta Si₃N₄ + TFLN como plataforma viável para lógica por tempo de voo. Sobre ela, race logic fotônica resolve menor caminho num mapa 16×16 com 0 erros em 2,55×10⁷ distâncias usando unidade de 100 ps, em 42,2 ns por consulta, com área de 648 mm² que sobe para ~864 mm² quando o passo das espirais respeita a diafonia. A estatística por número de saltos valida o modelo de ruído e mostra que amostras pequenas superestimam a correção. As comparações honestas estreitam a proposta: contra o melhor algoritmo em CPU o ganho é de ~42×, uma tabela pré-calculada vence em mapa fixo, e uma race logic síncrona em CMOS alcança latência parecida com área e energia ordens de grandeza menores. O valor do trabalho está no método (orçamento físico, estatística de erro e verificação eletromagnética) e na plataforma de temporização óptica, que só se justifica como acelerador se a unidade de atraso cair bem abaixo de um ciclo de relógio.

Trabalhos futuros:
- **Simulação eletromagnética 3D** de curvas e transições Si₃N₄/TFLN (a verificação 2D por índice efetivo já foi feita, Seção 5.1.1), idealmente em máquina com mais memória ou por expansão em modos próprios (EME).
- **Demonstrador em fibra** (desenho em `docs/architecture/13-bancada-experimental-em-fibra.md`): atrasos em SMF-28 (20,4 mm por 100 ps), moduladores de niobato de lítio comerciais, fotodiodos rápidos e time tagger de 1,5–2 ps rms, em três fases: porta ToF, nó de race logic e grafo 3×3 com curva de erro por salto sob jitter injetado.
- **Baseline moderno**: medir Dial e a tabela pré-calculada em CPUs atuais (hoje escalados pelo Geekbench) e comparar com GPU em lote; sintetizar a race logic em CMOS numa biblioteca real para substituir as estimativas de área e energia.
- **Jitter por salto**: identificar se algum caminho físico permite unidade de atraso de ~10–30 ps, condição para a corrida óptica superar a eletrônica.
- **Escala**: atrasos compartilhados entre arestas para reduzir área, e medição real do acoplamento entre chips (faceta ou interposer) que viabiliza o modo exato.

## Agradecimentos

Ao colega de trabalho Valmor Moreira, da empresa Grafis, que explicou ao autor o funcionamento das máquinas de exposição fotográfica a laser contínuo durante os reparos que faziam juntos; a experiência com essas máquinas inspirou o autor a conceber o SilicaCore.

## Referências

1. Ahmed, S. R., et al. (2025). Universal photonic artificial intelligence acceleration. *Nature*, 640, 368–374.
2. Churaev, M., et al. (2023). A heterogeneously integrated lithium niobate-on-silicon nitride photonic platform. *Nature Communications*, 14, 3499.
3. Dial, R. B. (1969). Algorithm 360: Shortest-path forest with topological ordering. *Communications of the ACM*, 12(11), 632–633.
4. Dijkstra, E. W. (1959). A note on two problems in connexion with graphs. *Numerische Mathematik*, 1, 269–271.
5. Hart, P. E., Nilsson, N. J., & Raphael, B. (1968). A formal basis for the heuristic determination of minimum cost paths. *IEEE Transactions on Systems Science and Cybernetics*, 4(2), 100–107.
6. Hua, S., et al. (2025). An integrated large-scale photonic accelerator with ultralow latency. *Nature*, 640, 361–367.
7. Kissner, M., et al. (2024). An all-optical general-purpose CPU and optical computer architecture. *Journal of Lightwave Technology*, 42, 7999. arXiv:2403.00045.
8. Madhavan, A., Sherwood, T., & Strukov, D. (2014). Race logic: A hardware acceleration for dynamic programming algorithms. *ISCA 2014*.
9. Miller, D. A. B. (2010). Are optical transistors the logical next step? *Nature Photonics*, 4, 3–5.
10. Shang, K., et al. (2015). Low-loss compact multilayer silicon nitride platform for 3D photonic integrated circuits. *Optics Express*, 23(16), 21334.
11. Xu, Z., et al. (2024). Large-scale photonic chiplet Taichi empowers 160-TOPS/W artificial general intelligence. *Science*, 384, 202–209.
12. Yu, X., et al. (2026). High-endurance, low-loss Sb₂Se₃ optical switches on silicon nitride using transparent conductive heaters. arXiv:2604.11649.
