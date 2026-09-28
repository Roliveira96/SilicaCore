# Codificação Temporal e Roteamento de Trajetória (*Time-of-Flight Logic*)

## 1. Visão Geral da Lógica ToF

A lógica por tempo de voo (*Time-of-Flight* ou modulação por atraso de fase/propagação) resolve um problema crítico da computação óptica tradicional: a necessidade de detectores hiper-sensíveis a variações contínuas de intensidade luminosa. Em vez de medir apenas se a luz chegou forte ou fraca, o sistema mede **quando** ela atingiu a matriz de detectores.

A determinação dos estados binários não depende da absorção ou interrupção do feixe, mas sim do **comprimento do caminho óptico percorrido** dentro do substrato cúbico de sílica.

---

## 2. Arquitetura de Trajetória Dupla

### 2.1 Caminho Direto (Linha Rápida - Estado 1)
- O feixe emitido pelo laser lateral (VCSEL 850 nm) viaja em linha reta pela sílica sem desvio.
- Incide no detector da face oposta.
- Como percorre a menor distância física ($d_1 = 20.0\text{ mm}$), atinge o detector no tempo mínimo de propagação:

$$t_1 = \frac{d_1 \cdot n}{c}$$

Para $n = 1.4500$, $t_1 \approx 96.73\text{ ps}$.

### 2.2 Caminho Defletido por Espelho Interno (Linha Atrasada - Estado 0)
- Ao ocorrer colisão ou modulação com feixe de controle (efeito Kerr óptico ou acoplador eletro-óptico), o pulso é desviado para a base do cubo.
- Na base, espelhos internos adicionais impõem uma trajetória reflexiva em zigue-zague ou percurso estendido ($d_0 = 40.675\text{ mm}$).
- O tempo total de chegada sofre um atraso mensurável ($\Delta t$):

$$t_0 = \frac{d_0 \cdot n}{c} \quad \implies \quad \Delta t = t_0 - t_1 \approx 100.0\text{ ps}$$

---

## 3. Discriminação Lógica por Janela de Tempo (*Time-Gating*)

O fotodetector SPAD opera acoplado a um circuito TDC de amostragem síncrona com o relógio óptico pulsado:

- **Pulso registrado na janela rápida ($t_1 \pm \text{janela}/2$):** Interpretado como bit **`1`**.
- **Pulso ausente em $t_1$, mas detectado na janela atrasada ($t_0 \pm \text{janela}/2$):** Interpretado como bit **`0`**.
- **Pulso fora de ambas as janelas:** Sinal de ruído ou perda de coerência, descartado pelo circuito de amostragem.

---

## 4. Vantagens do Modelo Temporal

1. **Imunidade à Atenuação de Amplitude:** Variações moderadas de intensidade por espalhamento no vidro não corrompem o dado, pois a leitura baseia-se exclusivamente na posição da borda de subida do pulso no tempo.
2. **Exequibilidade com Detectores Existentes:** Sensores do tipo SPAD (*Single-Photon Avalanche Diode*) e circuitos TDC (*Time-to-Digital Converters*) comerciais oferecem resoluções na faixa de picosegundos (LSB $\le 5\text{ ps}$), viabilizando leituras em distâncias micrométricas ou milimétricas.
