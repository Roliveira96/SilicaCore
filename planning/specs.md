# Especificações Consolidadas do Sistema SilicaCore

| Parâmetro | Valor Nominal | Unidade | Referência / Mecanismo Físico |
| :--- | :--- | :--- | :--- |
| Substrato | Sílica Fundida ($SiO_2$) | - | Grau UV/IR de ultra-alta pureza |
| Motor Laser | Continuous Wave (CW) | - | **Lasers Continuous Wave de Estado Sólido (Lasers Always-ON)** |
| Roteamento de Dados | Deflexão EOM/AOM | - | Ativado ao energizar a placa por micro-espelhos 3D |
| Codificação de Dados | M-ária Espectral WDM | - | **Hexadecimal (4-bit) / Byte Completo (8-bit)** |
| Multiplicador de Vazão | $8\times$ | - | Saída direta em palavras de Bytes/Megabytes |
| Comprimentos de Onda ($\lambda$) | 635, 532, 450, 850 | nm | Laser RGB Optoeletrônico + Emissão ToF |
| Índice de Refração ($n$) | 1.4500 | - | Índice efetivo no vidro |
| Velocidade no Meio ($v$) | 0.20675 | mm/ps | Velocidade de fase da luz ($c/n$) |
| Distância Linha Rápida ($d_1$) | 20.000 | mm | Trajetória direta (Estado 1) |
| Distância Linha Atrasada ($d_0$) | 40.675 | mm | Trajetória refletida (Estado 0) |
| Diferencial Temporal ($\Delta t$) | 100.00 | ps | Janelamento entre estados ToF |
| Jitter Total ($\sigma_{\text{total}}$) | 11.24 | ps | Convolução Laser + SPAD + TDC LSB |
| Razão $\Delta t / \sigma$ | 8.90 | - | Separação entre janelas em unidades de jitter |
| Fator de Decisão $Q = \Delta t / 2\sigma$ | 4.45 | - | BER teórica $\approx 4.3 \times 10^{-6}$ (Monte Carlo $\sim 2 \times 10^{-5}$) |
| $Q$ Necessário p/ BER $10^{-12}$ | 7.03 | - | Exige $\sigma_{\text{total}} \le 7.1$ ps |
| Taxa por Canal (slot $\Delta t + W$) | $\approx 5.1$ | GHz | $\approx 20.5$ Gb/s com 4 bits/símbolo; micro-cubo $\approx 55.6$ GHz com $Q = 3.64$ |
| Teto do Detector SPAD | $\le 0.5$ | GHz | Tempo morto $\ge 2$ ns; dados usam fotodiodos UTC ($\sim 150$ Gbaud) |
| Roteamento Adotado | Si₃N₄ + TFLN | - | 1.3 dB/porta, 7 portas em cascata (bloco com espelhos: 46.6 dB/porta, inviável) |
| Latência Cache L1 Óptica | $\sim 25$ | ps | SRAM fotônica com micro-anéis acoplados, 40 GHz, capacidade KB (arXiv:2503.19544) |
| Período do Buffer em Linha de Atraso | $\sim 96.73$ | ps | Buffer/registrador de 619 bits por laço, não RAM principal (*Yao, 1993*) |
| RAM Unificada (HBM via I/O óptico) | $\sim 30$ | ns | Transporte óptico ~130 ps + célula DRAM |
| Photonic SSD Capacid. Volumétrica | 100 | TB em 15.6 cm³ | **Premissa**: $6.4\text{ TB/cm}^3$; vidro publicado é arquivo de escrita única (*Project Silica*) |
| Photonic SSD Vazão de Leitura | 1.2 | TB/s | **Premissa** não demonstrada (Project Silica lê por microscopia) |
| GPU Óptica WDM Canais | 3 | Cores RGB | Red 635nm, Green 532nm, Blue 450nm (*Weng et al., 2020*) |
| GPU Óptica Ray-Tracing | Eletrônico | - | Óptica acelera upscaling/denoise neural, não a geometria da cena |
| Photonic AI Tensor Densidade | 11 | TOPS/mm² | Multiplicação MVM via MZI Mesh (*Xu et al., Nature 2021*) |
| Photonic AI Tensor Eficiência | $> 100$ (núcleo) / $\approx 0.84$ (sistema) | TOPS/W | Núcleo óptico (*Shen et al., 2017*); sistema completo (*Ahmed et al., Nature 2025*) |
| Qubits Quânticos Fotônicos | 16 | Dual-Rail | Circuito em temperatura ambiente; fontes e SNSPDs criogênicos (~1–4 K) (*Kok et al., 2007*) |
| Visibilidade Interferência HOM | 99.4 | % | Interferência quântica de 2 fótons (*Crespi et al., 2013*) |
| Fidelidade Porta CNOT Quântica | 98.7 | % | Portas quânticas fotônicas em vidro (*Carolan et al., 2015*) |
| Latência Média Global de Memória | $\sim 416$ | ps | Monte Carlo: L1 92%, L2/L3 90% das faltas, HBM; dominada pela DRAM |
| Energia por Bit (modelo) | $\sim 1.9$ | pJ/bit | Laser CW dividido pela taxa agregada de 1.313 Gb/s + EOM + detector |
