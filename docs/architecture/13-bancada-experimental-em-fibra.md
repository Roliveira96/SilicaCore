# Bancada Experimental em Fibra: Validação da Lógica ToF e do Race Logic

> **Substitui** a bancada original da tarefa 17 (divisores de feixe cúbicos, feixe livre e detectores SPAD). Feixe livre foi descartado pela difração (doc [11](11-roteamento-e-comutacao-optica.md)); SPADs ficam cegos por ≥ 1,5–2 ns após cada detecção e são desnecessários quando cada pulso carrega muitos fótons.

## 1. Objetivo

Medir, com componentes comerciais de telecomunicações em 1550 nm, os parâmetros que o simulador hoje trata como **premissa** (ASSUMPTION), e reproduzir em escala de bancada os dois resultados centrais do [artigo preliminar](../papers/artigo-preliminar.md):

1. o fator de decisão $Q = \Delta t / 2\sigma$ de uma porta ToF;
2. a correção do race logic para menor caminho e a curva de erro por número de saltos.

A bancada **não** demonstra a integração em Si₃N₄/TFLN, a área nem a energia por bit do chip. Ela valida o princípio físico e o modelo de ruído.

## 2. Escala: por que a bancada usa atrasos em nanossegundos

| Grandeza | Chip (modelo) | Bancada |
| :--- | :--- | :--- |
| Atraso de 100 ps | ~15 mm de espiral Si₃N₄ ($n_g = 2{,}0$) | **20,4 mm de fibra SMF-28** ($n_g = 1{,}4682$) |
| Latência do nó (detectar + reemitir) | 20 ps (premissa) | **~1–3 ns** (fotodiodo, comparador, driver e cabos discretos) |
| Unidade de atraso $\tau$ | 100 ps | **~5 ns ≈ 1,02 m de fibra** (precisa ser maior que a latência do nó) |
| Tolerância de corte da fibra | — | ±1 mm ≈ ±4,9 ps (pode ser calibrada) |

Como $\tau \gg \sigma$ na bancada, o ruído natural não produz erros de decodificação. Para medir a curva de erro por número de saltos, a Fase C injeta jitter calibrado (Seção 4.3).

## 3. Componentes

| Função | Componente comercial | Especificação alvo | Observação |
| :--- | :--- | :--- | :--- |
| Fonte | Laser DFB CW 1550 nm | Linha estreita, ≥ 10 mW | Uma fonte alimenta todos os nós via divisores |
| Geração de pulsos | Modulador de intensidade de niobato de lítio (LiNbO₃ bulk) + gerador de pulsos | Pulsos de 20–50 ps | Alternativa: laser de fibra mode-locked. TFLN comercial é opcional; LiNbO₃ bulk tem o mesmo efeito Pockels |
| Chave ToF (Fase A) | Chave eletro-óptica 2×2 de LiNbO₃ | Tempo de subida < 50 ps | Faz o papel da chave TFLN do chip |
| Atrasos | Cordões de fibra SMF-28 cortados sob medida | Precisão ±1 mm | Pesos $w \cdot \tau$ como comprimentos de fibra |
| Programação dos pesos (Fase C2) | Chaves de fibra MEMS 1×2 / 2×2 | ms de comutação | Faz o papel das chaves Sb₂Se₃ (reconfiguração lenta) |
| Detecção | Fotodiodo rápido InGaAs (PIN ou UTC) | ≥ 20 GHz | Muitos fótons por pulso; sem SPAD |
| Decisão do nó | Amplificador limitador + comparador rápido (ECL/CML) | Latência ~ns, jitter de poucos ps | Primeira chegada = primeiro cruzamento de limiar |
| Reemissão do nó | Modulador LiNbO₃ acionado pelo comparador | — | Recorta um pulso da fonte CW compartilhada |
| Amplificação | EDFA | Compensar divisões de fan-out | Na bancada, divisão passiva + EDFA substitui o "modulador por aresta" do chip |
| Leitura (TDC) | Time tagger | **1,5–2 ps rms por canal, até 20 canais** (ex.: Swabian Time Tagger X) | Um canal por nó + referência de disparo |
| Jitter controlado (Fase C3) | Gerador de forma de onda arbitrária (AWG) | Ruído somado ao limiar do comparador | Permite varrer $\sigma$ por nó |

## 4. Fases

### 4.1 Fase A: porta ToF isolada

**Montagem:** pulso → chave LiNbO₃ 2×2 → braço rápido ($L_1$) ou braço atrasado ($L_1$ + 20,4 mm, $\Delta t$ = 100 ps) → combinador → fotodiodo → time tagger.

**Medir:** $\sigma_{\text{total}}$ (jitter de pulso, fotodiodo, time tagger), $Q = \Delta t / 2\sigma$ e a BER por contagem, com a chave alternando o bit.

**Comparar com o simulador:** `ComputeTimingBudget` com os jitters medidos no lugar das premissas. **Critério de sucesso:** BER medida dentro do intervalo de confiança previsto pelo Q medido.

### 4.2 Fase B: um nó de race logic

**Montagem:** duas entradas com atrasos diferentes → fotodiodos → OR eletrônico → comparador → modulador de reemissão → saída para o time tagger.

**Medir:**
- **Latência do nó** $L_{\text{nó}}$ (premissa atual: 20 ps no chip).
- **Jitter do nó** $\sigma_{\text{nó}}$ (premissa atual: 1,5 ps rms).
- **Discriminação da primeira chegada:** a menor diferença entre as duas entradas que o nó separa corretamente.

**Critério de sucesso:** o nó dispara pela entrada mais adiantada em ≥ 99,9% dos eventos quando a diferença excede $2 \times 3\sigma_{\text{nó}}$.

### 4.3 Fase C: grafo pequeno resolvendo menor caminho

**Grafo:** grade 3×3 (9 nós, 24 arestas dirigidas, pesos 1–3) ou um grafo irregular de 6 nós. Nove nós cabem em um time tagger de 10 canais.

- **C1, pesos fixos:** cordões de fibra de comprimento $w\tau - L_{\text{nó}}$ trocados à mão a cada mapa. Medir o tempo de disparo de cada nó, decodificar $\hat d = \text{round}(t/\tau)$ e comparar com Dijkstra para todas as origens.
- **C2, pesos programáveis:** em um subconjunto de arestas, estágios binários de fibra selecionados por chaves MEMS. Isso demonstra a reconfiguração lenta e confirma que trocar a origem não exige reprogramar.
- **C3, curva de erro:** o AWG soma ruído calibrado ao limiar dos comparadores para varrer $\sigma_{\text{nó}}$. Com $\tau$ fixo, mede-se a taxa de erro por distância em função do número de saltos e compara-se com $P_{\text{err}}(h) = \text{erfc}\!\left(\frac{\tau/2}{\sigma\sqrt{2h}}\right)$ e com o Monte Carlo (`cmd/racestats` com os parâmetros medidos). Este é o análogo experimental da Figura 1 do artigo.

**Critério de sucesso:** C1 com 0 erros em todas as origens; C3 com a curva medida dentro do intervalo de confiança da previsão.

## 5. O que cada fase valida no simulador

| Premissa do simulador | Onde é medida |
| :--- | :--- |
| Jitter total da porta ToF (laser, detector, TDC) | Fase A |
| Latência de regeneração do nó (20 ps no chip) | Fase B (valor de bancada; o valor integrado exige chip) |
| Jitter por nó (1,5 ps rms) | Fase B |
| Erro estático por aresta após calibração (0,5 ps rms) | Fase C1 (tolerância de corte + calibração) |
| Modelo gaussiano de ruído acumulado por salto | Fase C3 |
| Tempo de programação dos pesos | Fase C2 (MEMS; Sb₂Se₃ exige chip) |

## 6. Riscos

- **Deriva térmica das fibras:** ~40 ps/km/°C em SMF. Em alguns metros é desprezível frente a $\tau$ = 5 ns, mas importa na Fase A ($\Delta t$ = 100 ps). Mitigação: bancada termicamente estável e medidas diferenciais.
- **Reflexões e cabos:** a latência elétrica dos nós depende do comprimento dos cabos. Mitigação: cabos iguais por nó e calibração de $L_{\text{nó}}$ nó a nó.
- **Número de canais:** grafos maiores que ~18 nós exigem mais de um time tagger sincronizado.
- **Custo:** os itens mais caros são o time tagger, os moduladores LiNbO₃ e o gerador de pulsos. O orçamento deve ser levantado com fornecedores; laboratórios de óptica e telecomunicações costumam ter parte desse equipamento.

## 7. Referências

1. Corning SMF-28 Optical Fiber — índice de grupo efetivo 1,4682 em 1550 nm. [Datasheet](https://www.tlc.unipr.it/cucinotta/cfa/datasheet_SMF28e.pdf)
2. Swabian Instruments — Time Tagger Series (jitter de 1,5–2 ps rms, até 20 canais). [Especificações](https://www.swabianinstruments.com/time-tagger/)
3. Madhavan, A., Sherwood, T., & Strukov, D. (2014). Race Logic. *ISCA 2014*.
4. Artigo preliminar do SilicaCore: [artigo-preliminar.md](../papers/artigo-preliminar.md).
