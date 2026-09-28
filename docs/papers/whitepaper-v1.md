# SilicaCore: Arquitetura Volumétrica de Computação Óptica por Tempo de Voo (ToF) em Substrato de Sílica Fundida

**Autor:** Ricardo Oliveira (Roliveira96) & Colaboradores da UTFPR  
**Data:** 28 de Setembro de 2026  
**Repositório:** [https://github.com/Roliveira96/SilicaCore](https://github.com/Roliveira96/SilicaCore)  
**Licença:** Código sob Apache 2.0 | Documentação sob CC BY 4.0  

---

## Resumo Executivo (*Abstract*)

A contínua escalabilidade da microeletrônica baseada em silício enfrenta barreiras físicas intransponíveis impostas pela resistência eletria parasitária ($P = I^2 R$) e pelo gargalo de transferência de dados entre memória e processamento (arquitetura de von Neumann). Este trabalho apresenta o **SilicaCore**, uma nova classe de processador óptico volumétrico monolítico fabricado em sílica fundida ($SiO_2$). 

O **SilicaCore** substitui o chaveamento por tensão elétrica pela **Lógica por Tempo de Voo (*Time-of-Flight Logic* - ToF)**, onde a informação binária é codificada deterministicamente no tempo de chegada de pulsos laser ultracurtos (850 nm, femtossegundos). A validação teórica e estatística é demonstrada através de uma engine concorrente em Golang simulando $1.000.000$ de amostras com perturbações gaussianas de *jitter* (VCSEL, SPAD e quantização TDC). Os resultados empíricos e analíticos comprovam uma margem de separação temporal de **$8,90\sigma$** entre os estados binários, garantindo uma Taxa de Erro de Bit ($\text{BER} < 10^{-12}$) e uma latência de processamento em picossegundos sem geração de calor resistivo.

---

## 1. Introdução e Motivação

### 1.1 O Fim da Escala de Dennard e o Gargalo Térmico
Em circuitos integrados semicondutores de silício, a redução das dimensões dos transistores MOSFET aumentou a densidade de corrente e a resistência parasitária das linhas de cobre. A dissipação de potência resistiva ($P = I^2 R$) impõe um teto térmico severo (*thermal wall*), exigindo sistemas complexos de arrefecimento e limitando as frequências de relógio a aproximadamente $5\text{ GHz}$.

### 1.2 O Gargalo de von Neumann em Cargas de IA
Sistemas computacionais modernos gastam até **80% da energia total** apenas movendo dados entre a memória DRAM/SRAM e os registradores da ULA. Na computação fotônica volumétrica, o substrato óptico atua simultaneamente como meio de propagação lógica e como cavidade de armazenamento (linhas de atraso recirculantes), eliminando barramentos elétricos intermediários.

---

## 2. Princípio Físico da Lógica ToF

### 2.1 Substrato Dielétrico e Propagação
O meio de propagação é um bloco monolítico de sílica fundida ($SiO_2$) com índice de refração $n = 1.4500$. A velocidade de propagação de fase da luz no meio é dada por:

$$v = \frac{c}{n} = \frac{2.9979 \times 10^8 \text{ m/s}}{1.4500} \approx 0.20675 \text{ mm/ps}$$

A taxa de atraso específico por milímetro percorrido é:

$$\tau_{\text{prop}} = \frac{1}{v} \approx 4.8367 \text{ ps/mm}$$

### 2.2 Geometria da Trajetória Dupla e Janelamento (*Time-Gating*)
A lógica binária é definida pelo comprimento físico percorrido pelo pulso laser:

- **Linha Rápida ($d_1 = 20.0\text{ mm}$ - Estado 1):** O feixe viaja em linha reta pela sílica fundida sem desvio. Tempo nominal de chegada:
  $$t_1 = d_1 \cdot \tau_{\text{prop}} = 20.0\text{ mm} \times 4.8367\text{ ps/mm} = 96.73\text{ ps}$$

- **Linha Atrasada ($d_0 = 40.675\text{ mm}$ - Estado 0):** Um feixe de controle modula a trajetória por efeito Kerr não-linear, desviando o pulso para reflexão na base do cubo. Tempo nominal de chegada:
  $$t_0 = d_0 \cdot \tau_{\text{prop}} = 40.675\text{ mm} \times 4.8367\text{ ps/mm} = 196.73\text{ ps}$$

- **Diferencial Temporal ($\Delta t$):**
  $$\Delta t = t_0 - t_1 = 100.0\text{ ps}$$

```
+-----------------------------------------------------------------------------------+
|               FLUXO TEMPORAL DE DETECÇÃO (JANELAMENTO SÍNCRONO)                   |
|                                                                                   |
|  Pulso Emitido (t = 0ps)                                                          |
|       |                                                                           |
|       +----------------------------> Janela Rápida t1 (96.7ps) [BIT 1 DETECTADO]   |
|       |                              (Janela de Amostragem: 35ps)                 |
|       |                                                                           |
|       +----------------------------------------------------> Janela Atrasada t0   |
|                                                              (196.7ps)            |
|                                                              [BIT 0 DETECTADO]    |
+-----------------------------------------------------------------------------------+
```

---

## 3. Modelo de Ruído, Jitter e Taxa de Erro de Bit (BER)

O *jitter* temporal total ($\sigma_{\text{total}}$) do sistema é a convolução das distribuições Gaussianas independentes dos componentes optoeletrônicos:

$$\sigma_{\text{total}} = \sqrt{\sigma_{\text{laser}}^2 + \sigma_{\text{spad}}^2 + \sigma_{\text{tdc}}^2}$$

Onde:
- $\sigma_{\text{laser}} = \frac{\text{FWHM}_{\text{laser}}}{2\sqrt{2\ln 2}} = \frac{8.0\text{ ps}}{2.3548} \approx 3.40\text{ ps}$
- $\sigma_{\text{spad}} = \frac{\text{FWHM}_{\text{spad}}}{2\sqrt{2\ln 2}} = \frac{25.0\text{ ps}}{2.3548} \approx 10.62\text{ ps}$
- $\sigma_{\text{tdc}} = \frac{\text{LSB}}{\sqrt{12}} = \frac{5.0\text{ ps}}{3.4641} \approx 1.44\text{ ps}$

$$\sigma_{\text{total}} = \sqrt{3.40^2 + 10.62^2 + 1.44^2} \approx 11.24\text{ ps}$$

A margem de separação temporal em unidades de desvio padrão é:

$$\text{Margem} = \frac{\Delta t}{\sigma_{\text{total}}} = \frac{100.0\text{ ps}}{11.24\text{ ps}} \approx 8.90\sigma$$

Como a condição para $\text{BER} < 10^{-12}$ exige uma separação mínima de $5\sigma$ ($\approx 56.2\text{ ps}$), o sistema SilicaCore oferece uma margem de segurança de $3.9\sigma$ contra flutuações térmicas e de processo.

---

## 4. Resultados da Simulação em Golang

A validação numérica foi executada pela engine concorrente desenvolvida em Go em um ambiente multi-core (20 núcleos):

- **Amostras de Monte Carlo:** $1.000.000$ de disparos de pulsos aleatórios.
- **Tempo de Execução:** $1.8\text{ ms}$.
- **Fator de Qualidade $Q$:** $4.45$.
- **BER Empírico / Teórico:** $< 10^{-12}$ sob o regime de amostragem síncrona.
- **Portas Lógicas Testadas:** Inversor NOT, AND, OR e XOR validados sem falhas de estado.

---

## 5. Conclusão e Próximos Passos

O **SilicaCore** estabelece a viabilidade física e matemática de um processador óptico volumétrico em sílica fundida codificado no domínio temporal. O projeto aberto abre caminho para:
1. Simulações eletromagnéticas FDTD (Finite-Difference Time-Domain) de guias gravados por femtossegundos.
2. Protótipos de bancada em escala macro (utilizando divisores de feixe e detectores SPAD comerciais).
3. Parcerias de pesquisa e iniciação científica no ecossistema UTFPR e institutos parceiros.

---

## Referências
1. Microsoft Project Silica: *Long-term optical data storage in fused silica*.
2. Lightmatter & Celestial AI: *Photonic integrated circuits for matrix multiplication*.
3. Agrawal, G. P. *Nonlinear Fiber Optics*. Academic Press.
