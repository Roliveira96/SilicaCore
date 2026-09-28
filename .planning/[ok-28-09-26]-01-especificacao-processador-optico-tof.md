# Especificação Técnica: Processador Óptico Tridimensional em Sílica (Arquitetura ToF)

**Documento:** `.planning/01-especificacao-processador-optico-tof.md`  
**Versão:** 1.0  
**Data:** 28-09-26  
**Status:** Em Especificação  

---

## 1. Visão Geral do Sistema

O processador óptico volumétrico em sílica aborda os limites fundamentais da eletrônica semicondutora tradicional:
1. **Dissipação Térmica e Resistência Parasita:** A propagação de elétrons em condutores metálicos gera aquecimento por efeito Joule ($P = I^2 R$) e atraso RC proporcional ao comprimento das interconexões. Fótons em guias de onda dielétricos não interagem resistivamente entre si.
2. **Gargalo de von Neumann:** A separação física entre a Unidade Central de Processamento e a memória primária exige barramentos de cobre com largura de banda limitada e alta latência de barramento. Na presente arquitetura, a memória é contida no próprio substrato de propagação óptica.
3. **Fragilidade da Lógica por Intensidade Luminosa:** Sistemas computacionais fotônicos baseados em amplitude sofrem de degradação severa por dispersão e perda de potência óptica. A **Lógica por Tempo de Voo (*Time-of-Flight Logic*)** transfere a informação do domínio da amplitude para o **domínio temporal determinístico**.

---

## 2. Parâmetros Físicos e Equações de Propagação

### 2.1 Substrato Dielétrico
- **Material Base:** Sílica Fundida (*Fused Silica* - $SiO_2$) de pureza óptica ultra-alta (grau UV/IR, ausência de bolhas e estrias, perda de transmissão $< 0.1 \text{ dB/km}$ nas janelas de telecomunicação).
- **Índice de Refração Efetivo ($n$):**
  - $\lambda = 850 \text{ nm}$: $n \approx 1.4525$
  - $\lambda = 1550 \text{ nm}$: $n \approx 1.4440$
  - Adotaremos valor nominal de projeto: $n = 1.4500$.

### 2.2 Velocidade de Propagação e Atraso Específico
A velocidade de propagação de um pulso luminoso no meio é dada por:
$$v = \frac{c}{n} = \frac{2.9979 \times 10^8 \text{ m/s}}{1.4500} \approx 2.0675 \times 10^8 \text{ m/s} = 0.20675 \text{ mm/ps}$$

A taxa de atraso temporal por unidade de distância física percorrida é:
$$\tau_{\text{prop}} = \frac{1}{v} \approx 4.8367 \text{ ps/mm}$$

### 2.3 Resolução Temporal de Detecção e Janelamento
- **Sensores Primários:** Matrizes de Fotodiodos de Avalanche de Fóton Único (*Single-Photon Avalanche Diodes* - SPAD) integradas em processo CMOS 3D-stacked.
- **Circuitos de Aquisição:** Conversores Tempo-Digital (*Time-to-Digital Converters* - TDC) com resolução temporal de LSB $= 5 \text{ ps}$ a $15 \text{ ps}$ e *timing jitter* total (FWHM) $\sigma_{\text{jitter}} \le 20 \text{ ps}$.
- **Critério de Separação Mínima ($\Delta t$):** Para assegurar uma taxa de erro de bit ($BER < 10^{-12}$), a separação temporal entre os estados lógicos deve satisfazer a regra dos $5\sigma$:
  $$\Delta t \ge 5 \cdot \sigma_{\text{jitter}} \approx 5 \times 20\text{ ps} = 100\text{ ps}$$
  Em implementações com calibração estrita ou detectores SNSPD supercondutores ($\sigma < 3\text{ ps}$), $\Delta t$ pode ser reduzido para $15\text{ ps}$ a $25\text{ ps}$.

### 2.4 Dimensionamento Geométrico do Caminho Diferencial ($\Delta d$)
Para um atraso de chaveamento alvo de $\Delta t = 100 \text{ ps}$:
$$\Delta d = d_0 - d_1 = \Delta t \cdot v = 100 \text{ ps} \times 0.20675 \text{ mm/ps} \approx 20.675 \text{ mm} \approx 2.07 \text{ cm}$$
Em um cubo de sílica com aresta de $25\text{ mm} \times 25\text{ mm} \times 25\text{ mm}$, um caminho direto de $20\text{ mm}$ corresponde à Linha Rápida ($t_1 \approx 96.7\text{ ps}$), e uma reflexão em ângulo duplo (percurso estendido de $40.7\text{ mm}$) atinge a Linha Atrasada ($t_0 \approx 196.7\text{ ps}$).

---

## 3. Arquitetura de Portas Lógicas por Modulação de Trajetória

Em vez de chavear correntes elétricas por transistores de efeito de campo (MOSFET), a lógica ToF utiliza deflexão e atraso guiado:

### 3.1 Porta Inversora (NOT)
- **Princípio:** Um feixe de pulso de sinal (Clock/Probe) é injetado ciclicamente.
- **Entrada $A = 0$ (Ausência de feixe de controle):** O feixe de sinal segue o caminho em linha reta $d_1$ (Linha Rápida). O fotodetector recebe o pulso na janela $t_1$. Saída lógica interpretada = $1$.
- **Entrada $A = 1$ (Presença de feixe de controle):** O feixe de controle induz uma variação de índice de refração local via efeito óptico não-linear (Kerr óptico $\Delta n = n_2 I$) ou acoplador direcional eletro-óptico, desviando o feixe de sinal para a cavidade de reflexão estendida $d_0$. O pulso chega na janela retardada $t_0$. Como a janela $t_1$ fica vazia, o receptor registra ausência na janela rápida, logo saída = $0$.

### 3.2 Porta E (AND)
- Dois feixes de entrada $A$ e $B$ convergem para uma junção de acoplamento não-linear. Apenas a presença simultânea de ambos gera o deslocamento de fase necessário para comutação para o caminho de saída rápida.

### 3.3 Porta OU (OR)
- Multiplexação passiva espacial com guias de onda em Y simétricos desembocando no mesmo detector na janela rápida.

### 3.4 Porta OU-Exclusivo (XOR)
- Acoplador direcional 3dB simétrico onde a interferência destrutiva na saída rápida ocorre quando ambos os sinais estão em fase, desviando a energia para o canal secundário.

---

## 4. O Bloco Volumétrico de Sílica

### 4.1 Estrutura Física
1. **Substrato Central:** Bloco cúbico de $25\text{ mm} \times 25\text{ mm} \times 25\text{ mm}$.
2. **Revestimento Reflexivo:** Três faces externas recebem deposição de espelhos de Bragg dielétricos (Deltas alternados de $TiO_2$ e $SiO_2$, refletividade $> 99.95\%$).
3. **Faces Ativas de Entrada/Saída:**
   - **Face de Emissão:** Matriz bidimensional de lasers VCSEL pulsados (850nm) acoplados via microlentes colimadoras.
   - **Face de Recepção:** Matriz bidimensional de fotodiodos SPAD acoplados a circuitos de leitura CMOS TDC.
4. **Guias de Onda Volumétricos:** Gravados por escrita direta a laser de femtossegundo (FLICE - *Femtosecond Laser Irradiation followed by Chemical Etching* ou modificação de índice induzida).

### 4.2 Divisão em Andares Funcionais
- **Andar 1 (Base - Z = 0 a 5mm):** Barramento óptico principal e distribuidores de clock pulsado mestre.
- **Andar 2 (Central - Z = 5 a 15mm):** Unidade Aritmética e Lógica (ULA ToF) e Registradores de Trabalho.
- **Andar 3 (Intermediário - Z = 15 a 20mm):** Matriz de linhas de atraso recirculantes para armazenamento dinâmico (Memória Transitória).
- **Andar 4 (Topo - Z = 20 a 25mm):** Acelerador de produtos escalares para tensores de IA (Interferometria Óptica) e canais WDM (RGB) para processamento gráfico paralelo.

---

## 5. Memória Recirculante (Delay-Line Memory)

Para manter registradores e memória sem conversão constante para circuitos flip-flop eletrônicos:
- O feixe luminoso é inserido em uma cavidade reflexiva em anel fechado com ganho óptico compensado por micro-amplificadores ópticos semicondutores (SOA) de baixa potência.
- Os bits permanecem circulando como pacotes de pulsos em ordem serial/paralela na velocidade $c/n$.
- A leitura é não-destrutiva: um divisor de feixe 95/5 extrai 5% da potência luminosa para os sensores a cada volta completa do ciclo, enquanto 95% do sinal continua circulando.

---

## 6. Roteiro de Prototipagem em Duas Fases

1. **Fase Bancada Maker/Acadêmica (Macro-escala):**
   - Diodos laser de 850nm / 650nm pulsados por gerador de pulso rápido (< 1ns).
   - Divisores de feixe cúbicos (Beam Splitters 50/50).
   - Espelhos front-surface de alumínio ou dielétricos.
   - Cabo de fibra óptica de sílica monomodo cortado em dois comprimentos ($L_1 = 1\text{ m}$, $L_2 = 2\text{ m}$) para obter $\Delta t \approx 4.8\text{ ns}$, facilmente detectável por osciloscópios convencionais de 100MHz a 1GHz ou CIs TDC comerciais (ex: Texas Instruments TDC7200).
2. **Fase Integrada de Sílica (Meso/Micro-escala):**
   - Bloco monolítico de sílica fundida usinado a laser.
