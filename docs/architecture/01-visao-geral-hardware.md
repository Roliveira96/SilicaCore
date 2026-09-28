# Especificação Técnica de Hardware: Processador Óptico Volumétrico em Sílica

## 1. Substrato Físico
O núcleo do processador consiste em um bloco cúbico de **sílica fundida de ultra-alta pureza ($SiO_2$)** com dimensões nominais de $25\text{ mm} \times 25\text{ mm} \times 25\text{ mm}$.

- **Índice de Refração Efetivo ($n$):** $1.4500$ (na faixa de $850\text{ nm}$).
- **Velocidade de Propagação no Meio ($v$):**
  $$v = \frac{c}{n} = \frac{2.9979 \times 10^8 \text{ m/s}}{1.4500} \approx 0.20675 \text{ mm/ps}$$
- **Atraso Propagacional Específico ($\tau_{\text{prop}}$):** $4.8367\text{ ps/mm}$.

---

## 2. Faces Ativas e Revestimento Optoeletrônico
1. **Face Frontal (Emissão):** Matriz 2D de lasers emissores de superfície com cavidade vertical (VCSEL) operando em $850\text{ nm}$ com duração de pulso em femtossegundos.
2. **Face Posterior (Recepção):** Matriz 2D de Fotodiodos de Avalanche de Fóton Único (SPAD) integrados em tecnologia CMOS 3D com conversores Tempo-Digital (TDC).
3. **Faces Laterais e Inferior:** Revestidas com espelhos dielétricos multincamadas (Bragg) de titânia/sílica ($TiO_2/SiO_2$) apresentando refletividade $> 99.95\%$.

---

## 3. Disposição Geométrica dos Caminhos Lógicos
- **Linha Rápida ($d_1$):** Trajetória direta de $20.0\text{ mm}$ ($\Delta t_1 \approx 96.73\text{ ps}$).
- **Linha Atrasada ($d_0$):** Trajetória com reflexão estendida de $40.675\text{ mm}$ ($\Delta t_0 \approx 196.73\text{ ps}$).
- **Diferencial Temporal ($\Delta t$):** $100.0\text{ ps}$.
