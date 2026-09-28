# Especificação Técnica de Hardware: Processador Óptico Volumétrico em Sílica

> **Nota de validação (v1.1, 28/09/2026):** este documento foi alinhado ao simulador e aos docs [11](11-roteamento-e-comutacao-optica.md) e [12](12-memoria-unificada-jogos-e-ia-local.md). O bloco de sílica original com espelhos internos e feixe livre perde ~46.6 dB por porta por difração; a plataforma adotada é Si₃N₄ + TFLN em 1550 nm.

## 1. Substrato Físico (Geometria Agnóstica)
O núcleo do processador é um **substrato fotônico integrado**: camadas de guias de onda de Si₃N₄ sobre **sílica fundida de ultra-alta pureza ($SiO_2$)**, com chaves TFLN heterogêneas, empilhadas sobre a eletrônica CMOS de controle. A arquitetura **não depende do formato externo** (retangular, lâmina multicamada ou poligonal):

- **O atraso depende do comprimento do guia, não do invólucro:** $L = v \cdot \Delta t$. Uma linha de 100 ps é uma espiral de ~15 mm de Si₃N₄, que ocupa uma fração de mm² em qualquer formato de chip.
- **Compatível com fabricação industrial:** foundries de fotônica trabalham com wafers planos de 200/300 mm, e um die é limitado pelo retículo de litografia (26 × 33 mm, 8.6 cm²). O conceito original de cubo maciço de 25 mm exigiria lapidação óptica customizada e usaria feixes livres, com ~46.6 dB de perda por porta.
- **Térmica e encapsulamento:** formatos planos maximizam o contato com dissipadores e simplificam a interface com barramentos elétricos (BGA, PCIe).
- **Dimensão de referência:** as distâncias $d_1 = 20$ mm e $d_0 = 40.675$ mm abaixo definem **tempos de voo**, não o tamanho do chip.

- **Índice de Refração Efetivo ($n$):** $1.4500$ (na faixa de $850\text{ nm}$).
- **Velocidade de Propagação no Meio ($v$):**
  $$v = \frac{c}{n} = \frac{2.9979 \times 10^8 \text{ m/s}}{1.4500} \approx 0.20675 \text{ mm/ps}$$
- **Atraso Propagacional Específico ($\tau_{\text{prop}}$):** $4.8367\text{ ps/mm}$.

---

## 2. Motor Laser Contínuo CW em Estado Sólido (Always-ON)
1. **Emissão Contínua (Always-ON CW Lasers):** Diodos/DPSS laser RGB ($\lambda_R = 635\text{ nm}$, $\lambda_G = 532\text{ nm}$, $\lambda_B = 450\text{ nm}$) operam **permanentemente acesos em potência estabilizada**, eliminando surtos térmicos e repetição de chaveamento elétrico (*CW Solid-State Laser Report*).
2. **Roteamento e Comutação:** guias de onda de Si₃N₄ sobre SiO₂ fazem o roteamento fixo (curvas de 50 µm, sem espelhos), chaves Sb₂Se₃ fazem a reconfiguração lenta e chaves de niobato de lítio em filme fino (TFLN) fazem o chaveamento por bit em picossegundos (ver [doc 11](11-roteamento-e-comutacao-optica.md)).
3. **Face de Recepção:** fotodiodos UTC/InGaAs (> 100 GHz) com amostragem TDC para os dados, com leitura em Hexadecimal (4 bits por símbolo). SPADs ficam limitados a $\le 0.5$ GHz pelo tempo morto.
4. **Revestimento Reflexivo:** Faces externas revestidas com espelhos dielétricos multincamadas (Bragg) de titânia/sílica ($TiO_2/SiO_2$) apresentando refletividade $> 99.95\%$.

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
- **Linha Rápida ($d_1$):** Trajetória direta de $20.0\text{ mm}$ ($\Delta t_1 \approx 96.73\text{ ps}$).
- **Linha Atrasada ($d_0$):** Trajetória estendida de $40.675\text{ mm}$ ($\Delta t_0 \approx 196.73\text{ ps}$), implementada como espiral de guia de onda (não por reflexões em espelhos).
- **Decisão temporal:** $Q = \Delta t / 2\sigma = 4.45$, $\text{BER} \approx 4.3 \times 10^{-6}$; taxa de $\approx 5.1$ GHz por canal.
- **Codificação M-ária Densa:** Saída em Hexadecimal (`0x0` a `0xF`) por canal espectral WDM. O modo Byte (256 estados) não é suportado pelo ruído de fase atual.
