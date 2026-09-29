# Codificação Temporal e Roteamento de Trajetória (*Time-of-Flight Logic*)

> **Nota de validação (v1.1, 28/09/2026):** a lógica ToF foi mantida, mas o desvio por espelhos internos e efeito Kerr foi substituído por chave TFLN e linha de atraso em espiral Si₃N₄. A margem de 8.9σ com BER < 10⁻¹² foi corrigida para Q = 4.45.

> **Revisão v1.2 (29/09/2026):** os comprimentos e o detector foram atualizados para a plataforma adotada: espirais de Si₃N₄ ($n_g = 2.0$) e fotodiodo InGaAs com TDC. Os números antigos (sílica e SPAD) ficam indicados como conceito original (`go run ./cmd/tofplatform`).

## 1. Visão Geral da Lógica ToF

A lógica por tempo de voo (*Time-of-Flight* ou modulação por atraso de fase/propagação) resolve um problema crítico da computação óptica tradicional: a necessidade de detectores hiper-sensíveis a variações contínuas de intensidade luminosa. Em vez de medir apenas se a luz chegou forte ou fraca, o sistema mede **quando** ela atingiu a matriz de detectores.

A determinação dos estados binários não depende da absorção ou interrupção do feixe, mas sim do **comprimento do caminho óptico percorrido** dentro do substrato fotônico integrado (independente do formato externo do chip).

---

## 2. Arquitetura de Trajetória Dupla

### 2.1 Caminho Direto (Linha Rápida - Estado 1)
- O pulso emitido (pente de frequências ou laser mode-locked; plataforma recomendada em 1550 nm) viaja por um guia de onda direto.
- Incide no detector da face oposta.
- Como percorre a menor distância física ($d_1 = 14.5\text{ mm}$ de Si₃N₄), atinge o detector no tempo mínimo de propagação:

$$t_1 = \frac{d_1 \cdot n_g}{c}$$

Para $n_g = 2.0$, $t_1 \approx 96.73\text{ ps}$ (no conceito original em sílica, $n = 1.45$ e $d_1 = 20$ mm davam o mesmo tempo).

### 2.2 Caminho Desviado por Chave TFLN (Linha Atrasada - Estado 0)
- Uma chave eletro-óptica de niobato de lítio em filme fino (TFLN) desvia o pulso para a linha atrasada. O efeito Kerr óptico na sílica é fraco demais para essa função, e espelhos internos com feixe livre perdem ~46.6 dB por difração (ver [doc 11](11-roteamento-e-comutacao-optica.md)).
- A linha atrasada é uma espiral de guia Si₃N₄ de $d_0 \approx 29.5$ mm (40,675 mm no conceito em sílica).
- O tempo total de chegada sofre um atraso mensurável ($\Delta t$):

$$t_0 = \frac{d_0 \cdot n_g}{c} \quad \implies \quad \Delta t = t_0 - t_1 \approx 100.0\text{ ps}$$

---

## 3. Discriminação Lógica por Janela de Tempo (*Time-Gating*)

Um fotodiodo InGaAs rápido, seguido de comparador e TDC sincronizado com o relógio óptico pulsado, marca o tempo de chegada (cada pulso tem muitos fótons, então não é preciso SPAD):

- **Pulso registrado na janela rápida ($t_1 \pm \text{janela}/2$):** Interpretado como bit **`1`**.
- **Pulso ausente em $t_1$, mas detectado na janela atrasada ($t_0 \pm \text{janela}/2$):** Interpretado como bit **`0`**.
- **Pulso fora de ambas as janelas:** Sinal de ruído ou perda de coerência, descartado pelo circuito de amostragem.

### 3.1 Taxa de Erro e Taxa por Canal

- **Plataforma adotada** (laser de 1 ps FWHM, fotodiodo com 1,5 ps rms, TDC de 5 ps): $\sigma_{\text{total}} \approx 2.12$ ps e $Q = \Delta t / 2\sigma \approx 23.5$. Com janela de ±7,13σ (±15 ps) o erro fica em $10^{-12}$; BER $10^{-12}$ seria possível já com $\Delta t \approx 30$ ps.
- **Conceito original** (SPAD de 25 ps FWHM): $\sigma_{\text{total}} = 11.24$ ps, $Q = 4.45$, $\text{BER} \approx 4.3 \times 10^{-6}$ com limiar no meio. O Monte Carlo com janela de ±47,5 ps mede $\sim 2 \times 10^{-5}$, porque um pulso que sai da janela conta como erro nas duas caudas: $\text{erfc}(W/2\sigma\sqrt2) \approx 2.4 \times 10^{-5}$.
- Cada símbolo ocupa as duas janelas ($\Delta t + W$), então a taxa por canal fica em ~5,1 GHz no conceito original e em até ~17 GHz de limite temporal na plataforma, antes de contar a banda de moduladores e detectores. O tempo de voo $t_1$ é latência, não período de clock.

---

## 4. Vantagens do Modelo Temporal

1. **Imunidade à Atenuação de Amplitude:** Variações moderadas de intensidade por espalhamento no vidro não corrompem o dado, pois a leitura baseia-se exclusivamente na posição da borda de subida do pulso no tempo.
2. **Exequibilidade com Detectores Existentes:** circuitos TDC (*Time-to-Digital Converters*) oferecem resoluções de picossegundos (LSB $\le 5\text{ ps}$). Para dados, fotodiodos UTC/InGaAs (> 100 GHz); SPADs têm tempo morto $\ge 2$ ns ($\le 0.5$ GHz).
3. **Base para Race Logic:** a mesma codificação temporal permite resolver menor caminho por corrida de pulsos; com unidade de 100 ps, o protótipo simulado teve 0 erros em 2,55×10⁷ distâncias num mapa 16×16 (taxa < 1.2×10⁻⁷ com 95% de confiança; com 50 ps, 1.2×10⁻⁴) ([doc 11, seção 5.1](11-roteamento-e-comutacao-optica.md)).
