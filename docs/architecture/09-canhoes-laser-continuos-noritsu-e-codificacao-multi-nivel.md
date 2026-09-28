# Emissão Laser Contínua (Estilo Minilab Noritsu) e Codificação M-ária (Hexadecimal / Byte)

## 1. Visão Geral do Motor Laser Contínuo (Noritsu CW Engine)

Inspirado nos sistemas de exposição a laser de ultra-alta precisão dos **minilabs fotográficos Noritsu**, o **SilicaCore** substitui a modulação por pulsação liga/desliga de diodo laser por um **Motor Laser de Onda Contínua (*Continuous Wave - CW Laser Engine*)**.

### 1.1 Emissores Permanentes (*Always-ON Lasers*)
- **Operação de Alta Estabilidade:** Os canhões laser sólidos RGB ($\lambda_R = 635\text{ nm}$, $\lambda_G = 532\text{ nm}$, $\lambda_B = 450\text{ nm}$) permanecem **constantemente ligados e alimentados** em nível de potência e fase calibrados.
- **Eliminação de Surtos Térmicos:** Ao eliminar o chaveamento elétrico de alta frequência nos diodos, erradica-se a fadiga dos semicondutores, oscilações de relaxamento óptico e ruídos térmicos associados ao liga/desliga.

### 1.2 Roteamento e Espelhamento Dinâmico na Inicialização
Assim que a placa do SilicaCore é energizada, o processador inicia o direcionamento dos feixes contínuos através de:
- **Moduladores Acusto-Ópticos / Eletro-Ópticos (AOM/EOM):** Deflexão angular ultra-rápida sem partes mecânicas.
- **Matriz de Micro-Espelhos 3D Gravados em Sílica:** Condução e combinação dos feixes pelas trajetórias do substrato até a matriz de fotodetectores.

```mermaid
flowchart TD
    subgraph NoritsuCWEngine["Motor Laser Contínuo CW (Estilo Minilab Noritsu - Always-ON)"]
        direction TB
        LaserR["Laser Vermelho CW (635nm - Potência Constante)"]
        LaserG["Laser Verde CW (532nm - Potência Constante)"]
        LaserB["Laser Azul CW (450nm - Potência Constante)"]
        LaserR --- LaserG --- LaserB
    end

    subgraph CuboSilica["Substrato de Sílica Fundida (Roteamento & Espelhamento Dinâmico)"]
        direction TB
        Moduladores["Moduladores Eletro-Ópticos (EOM/AOM): Ativados ao ligar a placa"]
        Espelhos3D["Matriz de Micro-Espelhos 3D: Deflexão e Combinação de Feixes Espectrais"]
        Moduladores --> Espelhos3D
    end

    subgraph DetectoresMArio["Matriz de Detectores Multi-Nível (Hexadecimal / Byte)"]
        direction TB
        SensorHex["Leitura Espectral WDM -> Saída em Hexadecimal (4-bit) / Byte (8-bit) Direct"]
    end

    NoritsuCWEngine -- "Feixes RGB Contínuos de Alta Estabilidade" --> CuboSilica
    CuboSilica -- "Combinação de Cores e Rotas Defletidas" --> DetectoresMArio
```

---

## 2. Codificação Densa M-ária (Hexadecimal e Byte por Símbolo)

Em vez de limitar a transmissão a um sinal binário simples (`0` ou `1`, $1\text{ bit}$ por feixe), o SilicaCore explora a **combinação de cores RGB e níveis de fase/amplitude** para transmitir valores densos por canal espacial:

### 2.1 Codificação Hexadecimal (4 Bits por Símbolo - 16 Estados)
- **Modulação:** 16 combinações discretas de amplitude e fase entre os 3 feixes RGB.
- **Saída:** Cada canal espacial entrega diretamente um caractere hexadecimal (`0x0` a `0xF`).

### 2.2 Codificação em Byte Completo (8 Bits por Símbolo - 256 Estados)
- **Modulação:** 256 combinações espectrais discretas codificadas na interferência WDM.
- **Saída:** Cada medição do detector lê diretamente um **Byte de dados (0 a 255)** já processado pelo bloco óptico.

---

## 3. Vantagens Computacionais da Saída em Byte/Megabyte

1. **Eliminação do Bottleneck de Decodificação Binária:** A placa não precisa gastar ciclos de clock agrupando bits individuais em bytes; o dado processado sai no formato nativo de palavras de memória (8, 16, 32 ou 64 bits).
2. **Vazão de Dados Multiplicada:** A taxa de transferência efetiva do processador é multiplicada por **$8\times$ (em modo Byte)** em relação a sistemas ópticos binários convencionais.

---

## 4. Referências Bibliográficas Científicas

1. **Noritsu Koki Co., Ltd.** "Precision Laser Exposure Engine Technology for Photofinishing Systems." *Technical Whitepaper Report*.
2. **Weng, L., et al. (2020).** "Wavelength-division multiplexed photonic computing for high-throughput matrix processing." *IEEE JSTQE*, 26(5), 1–12.
3. **Agrawal, G. P. (2012).** *Fiber-Optic Communication Systems: M-ary Modulation Techniques*. John Wiley & Sons.
