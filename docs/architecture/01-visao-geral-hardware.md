# Especificação Técnica de Hardware: do Processador Óptico em Sílica à Plataforma Si₃N₄ + TFLN

> **Nota de validação (v1.1, 28/09/2026):** este documento foi alinhado ao simulador e aos docs [11](11-roteamento-e-comutacao-optica.md) e [12](12-memoria-unificada-jogos-e-ia-local.md). O bloco de sílica original com espelhos internos e feixe livre perde ~46.6 dB por porta por difração; a plataforma adotada é Si₃N₄ + TFLN em 1550 nm.

> **Revisão v1.2 (29/09/2026):** lasers RGB visíveis, espelhos de Bragg nas faces e codificação hexadecimal saíram deste documento (ver doc [09](09-canhoes-laser-continuos-cw-e-codificacao-multi-nivel.md)). Os tempos da porta ToF foram recalculados para Si₃N₄ com fotodiodo (`go run ./cmd/tofplatform`).

## 1. Substrato Físico (Geometria Agnóstica)
O núcleo do processador é um **substrato fotônico integrado**: camadas de guias de onda de Si₃N₄ sobre **sílica fundida de ultra-alta pureza ($SiO_2$)**, com chaves TFLN heterogêneas, empilhadas sobre a eletrônica CMOS de controle. A arquitetura **não depende do formato externo** (retangular, lâmina multicamada ou poligonal):

- **O atraso depende do comprimento do guia, não do invólucro:** $L = v \cdot \Delta t$. Uma linha de 100 ps é uma espiral de ~15 mm de Si₃N₄, que ocupa uma fração de mm² em qualquer formato de chip.
- **Compatível com fabricação industrial:** foundries de fotônica trabalham com wafers planos de 200/300 mm, e um die é limitado pelo retículo de litografia (26 × 33 mm, 8.6 cm²). O conceito original de cubo maciço de 25 mm exigiria lapidação óptica customizada e usaria feixes livres, com ~46.6 dB de perda por porta.
- **Térmica e encapsulamento:** formatos planos maximizam o contato com dissipadores e simplificam a interface com barramentos elétricos (BGA, PCIe).
- **Dimensão de referência:** os tempos $t_1 = 96.73$ ps e $t_0 = 196.73$ ps definem a porta ToF, não o tamanho do chip.

- **Índice de grupo do guia Si₃N₄ ($n_g$):** $2.0$ no modelo ($2.06$ no FDTD para o guia de 800 nm × 0,7 µm) em 1550 nm.
- **Velocidade de grupo:**
  $$v_g = \frac{c}{n_g} = \frac{2.9979 \times 10^8 \text{ m/s}}{2.0} \approx 0.150 \text{ mm/ps}$$
- **Atraso por comprimento:** $\approx 6.67\text{ ps/mm}$. O conceito original usava sílica ($n = 1.45$, $4.84$ ps/mm).

---

## 2. Motor Laser Contínuo CW em Estado Sólido (Always-ON)
1. **Laser sempre aceso em 1550 nm:** um laser CW estabilizado alimenta um pente de frequências ou laser mode-locked, que fornece os pulsos e os canais WDM. A ideia de manter a fonte sempre acesa vem das máquinas de exposição da Grafis. Os lasers RGB visíveis do conceito original foram descartados: espalham ~140× mais que 1550 nm, e os fotodiodos InGaAs usados na leitura não detectam luz visível.
2. **Roteamento e Comutação:** guias de onda de Si₃N₄ sobre SiO₂ fazem o roteamento fixo (curvas de 50 µm, sem espelhos), chaves Sb₂Se₃ fazem a reconfiguração lenta e chaves de niobato de lítio em filme fino (TFLN) fazem o chaveamento por bit em picossegundos (ver [doc 11](11-roteamento-e-comutacao-optica.md)).
3. **Face de Recepção:** fotodiodos UTC/InGaAs (> 100 GHz) com TDC, leitura binária por janela de tempo. SPADs ficam limitados a $\le 0.5$ GHz pelo tempo morto e restritos ao núcleo quântico.
4. **Entrada e saída de luz:** acopladores de borda ligam os guias a fibras (≤ 1,5 dB por face). Não há espelhos nas faces: a luz fica confinada nos guias.

### 2.1 Mecanismos de Roteamento Dinâmico: Estado Sólido Eletro-Óptico vs MOEMS Mecânicos

A arquitetura do SilicaCore adota exclusivamente **mecanismos de comutação em estado sólido (sem partes móveis)**. A validação física (simulador, seção 10) descartou as duas opções originais em sílica:

- **Descartados após validação:**
  - **Moduladores Acusto-Ópticos (AOM):** o som leva ~16.8 ns para atravessar um feixe de 100 µm em sílica. É uma escala de nanossegundos, não de picossegundos.
  - **Reflexão interna total por campo elétrico no $SiO_2$:** a sílica é amorfa e não tem efeito Pockels; o efeito Kerr DC gera $\Delta n \lesssim 10^{-9}$, que só produziria TIR com o feixe a 0.002° da superfície.
- **Adotado:** chaves eletro-ópticas de **niobato de lítio em filme fino (TFLN)** integradas heterogeneamente sobre guias de Si₃N₄ (> 67 GHz, transições < 0.1 dB; Churaev et al., *Nat. Commun.* 2023), com chaves não-voláteis Sb₂Se₃ para reconfiguração.
- **Microestruturas Mecânicas (MOEMS / FLICE - Descartado para ToF High-Speed):**
  - Cavidades gravadas por laser femtossegundo e corrosão química (*FLICE*) com micro-palhetas suspensas (estilo DMD/DLP) apresentam inércia mecânica intrínseca, limitando a frequência de oscilação às faixas de kHz a MHz (incompatíveis com a computação fotônica em picossegundos).

---

## 3. Disposição Geométrica e Codificação Densa
- **Linha Rápida ($d_1$):** $14.5\text{ mm}$ de Si₃N₄ ($t_1 \approx 96.73\text{ ps}$); no conceito em sílica eram 20,0 mm.
- **Linha Atrasada ($d_0$):** espiral de $29.5\text{ mm}$ de Si₃N₄ ($t_0 \approx 196.73\text{ ps}$); no conceito em sílica eram 40,675 mm.
- **Decisão temporal:** com fotodiodo InGaAs e TDC, $\sigma_{\text{total}} \approx 2.1$ ps e $Q = \Delta t / 2\sigma \approx 23.5$. O conceito original, com SPAD de 25 ps FWHM, tinha $\sigma = 11.24$ ps, $Q = 4.45$ e $\text{BER} \approx 4.3 \times 10^{-6}$.
- **Taxa por canal:** limitada pela janela de decisão (~5,1 GHz com o SPAD original; até ~17 GHz de limite temporal com fotodiodo e $\Delta t = 30$ ps), antes de contar a banda dos moduladores e detectores.
- **Codificação:** binária por tempo de chegada. A codificação hexadecimal por cores e fases RGB do conceito original não tem sustentação física (doc 09).
