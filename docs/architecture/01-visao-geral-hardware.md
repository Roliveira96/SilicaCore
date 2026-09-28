# Especificação Técnica de Hardware: Processador Óptico Volumétrico em Sílica

## 1. Substrato Físico
O núcleo do processador consiste em um bloco cúbico de **sílica fundida de ultra-alta pureza ($SiO_2$)** com dimensões nominais de $25\text{ mm} \times 25\text{ mm} \times 25\text{ mm}$.

- **Índice de Refração Efetivo ($n$):** $1.4500$ (na faixa de $850\text{ nm}$).
- **Velocidade de Propagação no Meio ($v$):**
  $$v = \frac{c}{n} = \frac{2.9979 \times 10^8 \text{ m/s}}{1.4500} \approx 0.20675 \text{ mm/ps}$$
- **Atraso Propagacional Específico ($\tau_{\text{prop}}$):** $4.8367\text{ ps/mm}$.

---

## 2. Motor Laser Contínuo CW em Estado Sólido (Always-ON)
1. **Emissão Contínua (Always-ON CW Lasers):** Diodos/DPSS laser RGB ($\lambda_R = 635\text{ nm}$, $\lambda_G = 532\text{ nm}$, $\lambda_B = 450\text{ nm}$) operam **permanentemente acesos em potência estabilizada**, eliminando surtos térmicos e repetição de chaveamento elétrico (*CW Solid-State Laser Report*).
2. **Deflexão e Roteamento Eletro-Óptico:** Ao energizar a placa, moduladores AOM/EOM e micro-espelhos 3D gravados no vidro iniciam o espelhamento e o direcionamento contínuo dos feixes pelas rotas lógicas.
3. **Face de Recepção (Detectores M-ários):** Matriz 2D de Fotodiodos SPAD/CMOS operando amostragem síncrona TDC para leitura de estados em Hexadecimal (4 bits) e Bytes (8 bits por símbolo).
4. **Revestimento Reflexivo:** Faces externas revestidas com espelhos dielétricos multincamadas (Bragg) de titânia/sílica ($TiO_2/SiO_2$) apresentando refletividade $> 99.95\%$.

### 2.1 Mecanismos de Roteamento Dinâmico: Estado Sólido Eletro-Óptico vs MOEMS Mecânicos

A arquitetura do SilicaCore adota exclusivamente **mecanismos de comutação em estado sólido (sem partes móveis)** para preservação da escala temporal de picossegundos e imunidade à fadiga mecânica:

- **Espelhos Dinâmicos em Estado Sólido (Adotado):**
  - **Moduladores Acusto-Ópticos (AOM / Grade de Bragg por Som):** Transdutores piezoelétricos de micro-ondas geram ondas acústicas no vidro, criando redes periódicas de difração que defletem os feixes em nanosegundos/picossegundos.
  - **Moduladores Eletro-Ópticos (EOM / Efeito Kerr e Pockels):** A aplicação de um campo elétrico modulado altera instantaneamente o índice de refração local do $SiO_2$, criando superfícies de reflexão interna total (TIR) sob demanda sem desgaste físico.
- **Microestruturas Mecânicas (MOEMS / FLICE - Descartado para ToF High-Speed):**
  - Cavidades gravadas por laser femtossegundo e corrosão química (*FLICE*) com micro-palhetas suspensas (estilo DMD/DLP) apresentam inércia mecânica intrínseca, limitando a frequência de oscilação às faixas de kHz a MHz (incompatíveis com a computação fotônica em picossegundos).

---

## 3. Disposição Geométrica e Codificação Densa
- **Linha Rápida ($d_1$):** Trajetória direta de $20.0\text{ mm}$ ($\Delta t_1 \approx 96.73\text{ ps}$).
- **Linha Atrasada ($d_0$):** Trajetória com reflexão estendida de $40.675\text{ mm}$ ($\Delta t_0 \approx 196.73\text{ ps}$).
- **Codificação M-ária Densa:** Saída direta em formato Hexadecimal (`0x0` a `0xF`) e Bytes completos (0 a 255) por canal espectral WDM.
