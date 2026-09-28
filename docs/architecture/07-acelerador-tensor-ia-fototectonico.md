# Acelerador Tensor de IA Fotônico (Photonic AI Tensor Core & MVM)

> **Nota de validação (v1.1, 28/09/2026):** a eficiência > 100 TOPS/W vale só no núcleo óptico; o sistema completo publicado mais avançado entrega ~0.84 TOPS/W. O GST foi substituído por Sb₂Se₃ e o limite de capacidade de pesos no chip foi registrado (doc [12](12-memoria-unificada-jogos-e-ia-local.md)).

## 1. Visão Geral do Photonic Tensor Core

O **SilicaCore** incorpora um motor dedicado de Inteligência Artificial denominado **Photonic AI Tensor Engine**, localizado na **Camada 4 da pilha fotônica**. O acelerador é projetado para resolver a operação matemática computacionalmente mais exigente da Inteligência Artificial moderna (Transformers, LLMs, Visão Computacional e CNNs): a **Multiplicação Matriz-Vetor (MVM - *Matrix-Vector Multiplication*)**.

```mermaid
flowchart TD
    subgraph PhotonicTensorCore["Acelerador Tensor de IA Fotônico (Camada 4)"]
        direction TB
        VetorEntrada["Vetor de Entrada Óptico X (Amplitude / Fase dos Pulsos)"]
        MZIMesh["Malha de Interferômetros Mach-Zehnder (MZI Mesh - Rotações Unitárias)"]
        PesosPCM["Matriz de Pesos Não-Volátil (PCM Sb2Se3)"]
        DetecaoSaida["Vetor de Saída Y = W * X (Fotodetecção SPAD Integrada)"]
        
        VetorEntrada --> MZIMesh
        MZIMesh --> PesosPCM
        PesosPCM --> DetecaoSaida
    end
```

---

## 2. Computação In-Memory e Multiplicação Matriz-Vetor (MVM)

### 2.1 Malha de Interferômetros Mach-Zehnder (MZI Mesh)
A arquitetura utiliza uma malha de interferômetros Mach-Zehnder acoplados (*Shen et al., Nature Photonics 2017*). Cada nó MZI ajusta a fase e a divisão de amplitude do pulso óptico através de deslocadores de fase eletro-ópticos ou térmicos:

$$Y = U \cdot \Sigma \cdot V^\dagger \cdot X$$

A decomposição em valores singulares (SVD) permite que qualquer matriz de pesos arbitrária $W$ de um modelo de Inteligência Artificial seja executada no domínio óptico em uma **única passagem contínua da luz (*single-pass*)** com latência $< 10\text{ ps}$.

### 2.2 Pesos Não-Voláteis em Materiais de Mudança de Fase (PCM $GST$)
- **Armazenamento de Pesos em PCM:** filmes de mudança de fase sobre os guias armazenam pesos por estados de cristalização multinível (*Feldmann et al., Nature 2019*). **Sb₂Se₃** substitui o GST (transparente em 1550 nm, > 6 bits por célula, > 1.4×10⁸ ciclos; Yu et al., 2026).
- **Limite de capacidade:** um LLM de 8B parâmetros em 4 bits exigiria ~2.000 cm² de células PCM, contra 8.6 cm² de um retículo. Pesos de LLMs ficam em RAM unificada de alta banda e chegam por I/O óptico; a PCM no chip guarda adaptadores e modelos pequenos (doc 12).
- **Zero Consumo Específico:** Os pesos do modelo de IA permanecem gravados sem consumo elétrico contínuo, eliminando a busca repetitiva na memória RAM ou Flash.

---

## 3. Eficiência Energética e Métricas de Desempenho

| Métrica de Desempenho | GPU Eletrônica de IA (Nvidia H100) | Photonic AI Tensor Core (SilicaCore) |
| :--- | :--- | :--- |
| **Operação Primária** | Multiplicação digital por portas lógica CMOS | Produtos escalares por interferência e atenuação PCM |
| **Densidade de Processamento** | $\sim 0.5 \text{ TOPS/mm}^2$ | **Up to 11 TOPS/mm²** (*Xu et al., Nature 2021*) |
| **Eficiência Energética (núcleo)** | $\sim 0.5 \text{ a } 2.0 \text{ TOPS/W}$ | **$> 100 \text{ TOPS/W}$** só no núcleo óptico |
| **Eficiência Energética (sistema)** | $\sim 0.5 \text{ a } 2.0 \text{ TOPS/W}$ | **$\approx 0.84 \text{ TOPS/W}$** no estado da arte (Lightmatter, *Nature* 2025: 65.5 TOPS com 78 W + 1.6 W ópticos); Taichi: 160 TOPS/W no chiplet, sem periféricos |
| **Latência por Multiplicação MVM** | $10\text{ ns}$ a $100\text{ ns}$ (Ciclos de clock de GPU) | **$< 10\text{ ps}$** (Tempo de voo na malha MZI) |
| **Transferência de Memória** | Gargalo na VRAM HBM3 | In-memory para pesos pequenos; LLMs grandes continuam limitados pela banda da RAM (8B em 4 bits: ~840 tok/s com 3.35 TB/s) |
| **Precisão** | FP32/FP16/INT8 | Precisão efetiva comparável a FP32 em ResNet/BERT com ABFP16 (Lightmatter, 2025) |

---

## 4. Referências Bibliográficas Científicas

1. **Shen, Y., et al. (2017).** "Deep learning with coherent photonic circuits." *Nature Photonics*, 11(7), 441–446. [DOI: 10.1038/nphoton.2017.93](https://doi.org/10.1038/nphoton.2017.93)
2. **Feldmann, J., et al. (2021).** "Parallel convolutional processing using an integrated photonic tensor core." *Nature*, 595(7867), 373–378. [DOI: 10.1038/s41586-021-03597-z](https://doi.org/10.1038/s41586-021-03597-z)
3. **Xu, X., et al. (2021).** "11 TOPS mm⁻² photonic tensor core for optical neural networks." *Nature*, 589(7840), 44–51. [DOI: 10.1038/s41586-020-03070-1](https://doi.org/10.1038/s41586-020-03070-1)
