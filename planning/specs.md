# Especificações Consolidadas do Sistema SilicaCore

| Parâmetro | Valor Nominal | Unidade | Referência / Mecanismo Físico |
| :--- | :--- | :--- | :--- |
| Substrato | Sílica Fundida ($SiO_2$) | - | Grau UV/IR de ultra-alta pureza |
| Comprimento de Onda ($\lambda$) | 850 | nm | Faixa de emissão VCSEL / SPAD |
| Índice de Refração ($n$) | 1.4500 | - | Índice efetivo no vidro |
| Velocidade no Meio ($v$) | 0.20675 | mm/ps | Velocidade de fase da luz ($c/n$) |
| Distância Linha Rápida ($d_1$) | 20.000 | mm | Trajetória direta (Estado 1) |
| Distância Linha Atrasada ($d_0$) | 40.675 | mm | Trajetória refletida (Estado 0) |
| Diferencial Temporal ($\Delta t$) | 100.00 | ps | Janelamento entre estados ToF |
| Jitter Total ($\sigma_{\text{total}}$) | 11.24 | ps | Convolução Laser + SPAD + TDC LSB |
| Margem de Separação | 8.90 | sigmas | Imunidade a ruído ($\text{BER} < 10^{-12}$) |
| Latência Cache L1 Óptica | $\le 5.00$ | ps | Ressonadores de Micro-anéis (*Alexoudi et al., 2020*) |
| Latência RAM Óptica Volátil | $\sim 96.73$ | ps | Linhas de Atraso Recirculantes (*Yao, 1993*) |
| Photonic SSD Capacid. Volumétrica | 100 | TB / cubo | Voxels 3D em $SiO_2$ ($6.4\text{ TB/cm}^3$) (*Project Silica*) |
| Photonic SSD Vazão de Leitura | 1.2 | TB/s | Multiplexação WDM paralela |
| Photonic SSD Durabilidade | $> 10^9$ | Anos | Estabilidade de nanofilamentos em vidro |
| Latência Média Global de Memória | $\sim 37.60$ | ps | Simulação Monte Carlo |
