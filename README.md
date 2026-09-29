<div align="center">

# SilicaCore

### Computação por tempo de voo: a primeira luz a chegar é a resposta.

*Um acelerador fotônico de pesquisa aberta que resolve menor caminho em grafos deixando pulsos de luz apostarem corrida, com cada número checado contra a física.*

[![Site Oficial](https://img.shields.io/badge/🌐%20Site%20Oficial-silicacore.rmo.dev.br-00f2fe?style=for-the-badge)](https://silicacore.rmo.dev.br/)

[![Licença: Apache 2.0](https://img.shields.io/badge/c%C3%B3digo-Apache%202.0-blue.svg)](LICENSE) [![Docs: CC BY 4.0](https://img.shields.io/badge/docs-CC%20BY%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/) [![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-00ADD8.svg)](https://go.dev/) [![FDTD: Meep](https://img.shields.io/badge/FDTD-Meep-7a3ea1.svg)](https://meep.readthedocs.io/) [![Status](https://img.shields.io/badge/status-pesquisa%20aberta-orange.svg)](planning/roadmap.md)

### 🔗 **Site Oficial do Projeto:** [https://silicacore.rmo.dev.br/](https://silicacore.rmo.dev.br/)

**[🌐 Acessar Demonstração Online](https://silicacore.rmo.dev.br/)** · **[A ideia](#a-ideia-em-30-segundos)** · **[Demonstração Online](#-site-oficial--demonstração-interativa-live)** · **[Resultados](#o-que-os-números-dizem)** · **[O que a física mudou](#o-que-a-física-nos-obrigou-a-mudar)** · **[Rodar agora](#experimente-em-um-minuto)** · **[Artigo](docs/papers/artigo-preliminar.md)**

</div>

> [!IMPORTANT]
> ### 🌐 Site Oficial & Demonstração Interativa ao Vivo
> Acesse a experiência completa e simulação interativa em:  
> 👉 **[https://silicacore.rmo.dev.br/](https://silicacore.rmo.dev.br/)**  
> Veja o chip fotônico em 3D, a corrida óptica num mapa 4×4, o mecanismo de decisão do bit e a comparação com processadores atuais, com a origem de cada número.

---

## A ideia em 30 segundos

Um processador comum calcula o menor caminho num mapa passo a passo: lê um nó, compara distâncias, reordena uma fila, repete milhares de vezes.

O SilicaCore faz outra pergunta: **e se o próprio mapa fosse o computador?**

- Cada **estrada** do mapa vira um trecho de guia de onda, com comprimento proporcional ao seu custo.
- Cada **cruzamento** vira um detector que, ao receber o primeiro pulso, dispara um novo pulso para todas as saídas.
- Um único pulso sai da origem e se espalha pelo grafo. **O momento em que a luz chega a cada ponto é a distância até ele.**

Ninguém precisa comparar nada: quem chega primeiro já é, por definição, o caminho mais curto. É a ideia de *race logic* (Madhavan, Sherwood e Strukov, ISCA 2014) transposta para fótons em chips de nitreto de silício e niobato de lítio.

```mermaid
flowchart LR
    S(("Origem<br/>t = 0")) -- "custo 2 · 200 ps" --> A(("A<br/>200 ps"))
    S -- "custo 5 · 500 ps" --> B(("B<br/>300 ps"))
    A -- "custo 1 · 100 ps" --> B
    B -- "custo 3 · 300 ps" --> T(("Destino<br/>600 ps"))
    A -- "custo 7 · 700 ps" --> T
```

*O pulso chega a B em 300 ps pelo desvio S→A→B, antes dos 500 ps do caminho direto: a distância de B é 3. No destino, o primeiro pulso chega em 600 ps: distância 6.*

---

## 🌐 Site Oficial & Demonstração Interativa (Live)

O SilicaCore possui um portal web completo e interativo em produção para validação científica e demonstração visual:

👉 **Acesse online:** **[https://silicacore.rmo.dev.br/](https://silicacore.rmo.dev.br/)**

### O que o site mostra
- **Chip em 3D em escala física:** die de 25,5 mm com as 960 espirais, camadas ($\text{Sb}_2\text{Se}_3$, TFLN, $\text{Si}_3\text{N}_4$, $\text{SiO}_2$) em vista explodida e a corrida 16×16 simulada em falsa cor.
- **Corrida 4×4 no navegador:** a corrida óptica ao lado de um Dijkstra executado pelo próprio navegador.
- **Mecanismo de decisão em 3D:** como um bit vira 0 ou 1 pelo tempo de chegada (janelas e limiar), e o qubit dual-rail como conceito.
- **Comparação com processadores atuais:** uma tabela com abas de desempenho, potência, temperatura e IA (Apple M5 Max, Intel Core Ultra 9 285K, AMD Ryzen 9 9950X3D, EPYC 9965, Xeon 6980P, RTX 5090, B300, Rubin), cada número marcado como simulado, estimado ou do fabricante.
- **Física dos guias e energia:** resultados FDTD 2D e o orçamento de energia do chip.

Para rodar o portal web localmente:
```bash
cd web
npm install
npm run dev
```

---

## O que os números dizem

Resultados do simulador em Go (física de perdas, jitter, leitura e área), com a estatística feita para aguentar uma banca:

| | Resultado | Como foi obtido |
| :--- | :--- | :--- |
| **Correção** | **0 erros em 25,5 milhões de distâncias** (< 1,2×10⁻⁷ com 95% de confiança) | Mapa 16×16, 10 chips simulados × 10⁴ consultas, unidade de atraso de 100 ps |
| **Velocidade** | **42 ns por consulta** (11,5 ns de corrida + 30,7 ns de leitura) | Contra o melhor algoritmo para pesos inteiros (Dial): 25,4 µs numa CPU de 2012 (~600×) e **~42× estimado** contra uma CPU atual. Num mapa fixo, uma tabela pré-calculada é mais rápida que o chip |
| **Área** | **~864 mm²** para o bloco 16×16 | Espirais com passo de 4 µm, exigido pela diafonia entre voltas (com 3 µm eram 648 mm², mas o pulso se embaralha). Passa um pouco do retículo (858 mm²): pede duas camadas de guias (~432 mm²) |
| **Escala** | **64×64 em 16 chips, 0 erros, 31 ns** origem→destino | Exige unidade de 150 ps e acoplamento entre chips ≤ 1,5 dB por face; com 150 ps e passo de 3 µm, cada bloco precisa cair para ~12×12 para caber no retículo (menos ainda com o passo de 4 µm) |
| **Roteamento** | **1,3 dB por porta** em Si₃N₄ + TFLN | Contra 46,6 dB dos espelhos internos do conceito original |
| **Física dos guias** | Curva de 50 µm: **0,003 dB**; transição para o niobato: **< 0,003 dB** | Simulação eletromagnética FDTD 2D (Meep), confirmando que as premissas do modelo são conservadoras |

<p align="center">
  <img src="docs/papers/figuras/race_logic_erro_por_salto.svg" alt="Erro medido contra previsto por número de saltos" width="720">
</p>
<p align="center"><em><b>Figura 1.</b> O erro medido (pontos) segue a previsão teórica de ruído acumulado (linhas) em cada unidade de atraso. Com 100 ps, nenhum erro em 2,55×10⁷ distâncias.</em></p>

> [!NOTE]
> A leitura eletrônica dos resultados, e não a luz, é o que domina o tempo. O ganho não é O(1): cresce com a distância máxima do mapa. Esses limites estão quantificados no [artigo preliminar](docs/papers/artigo-preliminar.md), não escondidos.

---

## De onde veio a ideia

Na empresa **Grafis**, o autor (Ricardo Oliveira) trabalhava com revelação fotográfica em máquinas industriais: lasers **sempre acesos** varriam papel fotossensível guiados por um prisma rotativo e espelhos de precisão. Nos reparos dessas máquinas, o colega **Valmor Moreira** explicava como funcionava o sistema a laser.

Foi daí que o autor tirou a pergunta que deu origem ao projeto: se um feixe contínuo e bem guiado grava uma imagem com precisão micrométrica, por que não usá-lo para **computar**?

O primeiro desenho era um cubo de vidro com espelhos internos. A física respondeu. E a resposta levou a algo mais interessante.

---

## O que a física nos obrigou a mudar

Este projeto trata cada afirmação como hipótese. Quando o modelo desmentiu uma ideia, ela foi trocada, e o motivo ficou registrado:

| Ideia original | O que a física mostrou | Como ficou |
| :--- | :--- | :--- |
| Cubo maciço de sílica com espelhos internos | Feixe livre abre para ~2,8 mm e perde **46,6 dB por porta** | Substrato fotônico planar multicamada, compatível com wafers industriais |
| Espelho eletro-óptico na própria sílica | Sílica não tem efeito Pockels; reflexão exigiria incidência de 0,002° | Chaves de **niobato de lítio em filme fino (TFLN)**, acima de 67 GHz |
| "206 GHz" de clock pelo tempo de voo | Tempo de voo é latência, não período de clock | Taxa limitada pela janela de decisão: ~5,1 GHz com o SPAD original, até ~17 GHz de limite temporal com fotodiodo em Si₃N₄ |
| Margem de 8,9σ e erro < 10⁻¹² | O critério correto é Q = Δt/2σ = 4,45 (com o SPAD de 25 ps) | Detector trocado por fotodiodo InGaAs: σ = 2,1 ps e Q = 23,5 na plataforma Si₃N₄ (`cmd/tofplatform`) |
| Gigabytes de RAM em linhas de atraso | 16 GB exigiriam **~3.000 km** de guia de onda | Memória unificada: pSRAM fotônica (~25 ps) + HBM via I/O óptico |
| 0,05 pJ/bit e >100 TOPS/W | Lasers, conversores e controle dominam o consumo | ~7,6 pJ/bit com codificação binária; ~9,2 W para o chip 16×16; 0,84 TOPS/W é o melhor sistema fotônico medido |
| Comparar com Dijkstra comum | Pesos inteiros permitem fila de baldes (Dial), 2,6× mais rápida; mapa fixo admite tabela pré-calculada | Ganho reportado contra o melhor algoritmo (~42×) e nicho restrito a mapas que mudam |
| Espirais com passo de 3 µm | Voltas vizinhas acoplam com L_c ≈ 27 mm; o pulso se embaralha ao longo de 222 mm (conferido em FDTD) | Passo de 4 µm: eco adiantado de ~−30 dB, área ×4/3 |
| Race logic óptica contra CPU | Uma race logic síncrona em CMOS atinge a mesma latência com área e energia ordens de grandeza menores | O chip passa a ser tratado como veículo de validação da temporização óptica |
| Race logic "sem erros" com 50 ps | Com 10⁵ consultas aparecem 1,2×10⁻⁴ erros por distância | Unidade de 100 ps, validada estatisticamente |

O resultado é uma proposta menor que o sonho inicial, e por isso mesmo **defensável e publicável**.

---

## Arquitetura

```mermaid
flowchart TB
    Laser["Pente de frequências 1550 nm<br/>laser sempre aceso (herança da Grafis)"] --> Chip
    subgraph Chip["Substrato fotônico Si3N4 + TFLN (geometria livre, planar multicamada)"]
        direction LR
        Race["Race logic<br/>menor caminho por tempo de voo"]
        Tensor["Núcleo tensorial<br/>multiplicação matriz-vetor"]
        Switch["Chaves TFLN (ps)<br/>+ Sb2Se3 (reconfiguração)"]
    end
    Chip --> Read["Fotodiodos + TDC<br/>leitura eletrônica"]
    Chip <--> Mem["Memória unificada<br/>pSRAM + HBM via I/O óptico"]
    CPU["CPU / GPU eletrônicas"] <--> Mem
```

- **Guias de Si₃N₄** (800 nm × 0,7 µm, monomodo) formam as espirais de atraso, com perda < 0,1 dB/cm.
- **TFLN** faz a comutação em picossegundos; **Sb₂Se₃** guarda os pesos do grafo sem consumir energia parado.
- O SilicaCore é um **coprocessador**: lógica booleana comum, shading de jogos e controle continuam na eletrônica.

---

## Experimente em um minuto

Requer [Go 1.26+](https://go.dev/dl/).

```bash
git clone https://github.com/Roliveira96/SilicaCore.git
cd SilicaCore/simulations/go

go test ./...                     # suíte de testes
go run ./cmd/simulator            # relatório completo: física, memória, energia e race logic
go run ./cmd/racestats            # campanha de 10⁵ consultas com intervalo de confiança
go run ./cmd/racemultichip        # composição em vários chips: modo exato vs hierárquico
go run ./cmd/dijkstrabench         # baselines da CPU: heap, Dial e tabela pré-calculada
go run ./cmd/tofplatform          # porta ToF: conceito original (sílica + SPAD) vs Si3N4 + fotodiodo
go run ./cmd/cmosrace             # a mesma corrida em CMOS síncrono
```

<details>
<summary><b>Regenerar a Figura 1 e rodar as simulações eletromagnéticas</b></summary>

```bash
# Figura 1 a partir do CSV da campanha (sem dependências)
go run ./cmd/racestats -chips 10 -queries 10000 -csv ../results/race_logic_hops_16x16.csv
python3 ../results/plot_race_logic_hops.py

# FDTD 2D com Meep (instalado via conda-forge: pymeep)
cd ../fdtd
bash run_sweeps.sh bends          # curvas de 90° em Si3N4
bash run_sweeps.sh transition     # transição Si3N4 -> TFLN
```
</details>

---

## Onde estamos

- [x] Modelo físico: orçamento de perdas, jitter, taxa de erro, energia e memória
- [x] Race logic fotônica com validação estatística e composição multi-chip
- [x] Artigo preliminar com Figura 1 ([rascunho](docs/papers/artigo-preliminar.md))
- [x] Bancada experimental em fibra desenhada ([doc 13](docs/architecture/13-bancada-experimental-em-fibra.md))
- [x] Simulação eletromagnética FDTD 2D de curvas e transições ([doc 11, seção 2.1](docs/architecture/11-roteamento-e-comutacao-optica.md))
- [x] Portal web interativo e simulador 3D em produção ([silicacore.rmo.dev.br](https://silicacore.rmo.dev.br/))
- [ ] FDTD 3D das mesmas estruturas
- [x] Plano de fabricação e roteiro de protótipos G0–G4 ([doc 14](docs/architecture/14-fabricacao-e-prototipagem.md))
- [x] Andar de IA: unidade tensorial fotônica empilhada com o processador, com orçamento de ruído e energia ([doc 15](docs/architecture/15-andar-de-ia-unidade-tensorial-fotonica.md))
- [ ] Simulação T0 da unidade tensorial: rede pequena com o ruído do chip
- [ ] Montagem da bancada: porta ToF, nó de race logic e grafo 3×3
- [ ] Submissão (WSCAD/SBESC) e editais de iniciação científica

Detalhes em [roadmap](planning/roadmap.md) e [tarefas](planning/tasks.md).

---

## Limites conhecidos

Porque um projeto sério diz onde não funciona:

- **Não substitui a CPU.** É um acelerador para problemas específicos (grafos, multiplicação de matrizes).
- **A vantagem depende do mapa mudar.** Num mapa fixo, uma tabela de todas as distâncias responde em nanossegundos numa CPU atual. O chip ganha quando o mapa muda entre lotes de consultas (~34× com 100 consultas por mudança, se a reprogramação levar 1 µs).
- **Uma race logic em CMOS faz o mesmo.** Um circuito síncrono a 3 GHz (um ciclo por unidade) resolve o 16×16 em ~40 ns, sem erros, com ~0,007 mm² e ~0,8 nJ por consulta (estimativa), contra 42 ns, ~864 mm² e 386 nJ do RL-16 (`go run ./cmd/cmosrace`). A corrida óptica só é ~3× mais rápida na fase da corrida. Para superar a eletrônica, a unidade de atraso precisaria cair bem abaixo de um ciclo de relógio.
- **O RL-16 é um veículo de pesquisa.** Ele serve para validar a física da temporização óptica (jitter, diafonia, perdas) com estatística, não para vencer um circuito digital dedicado.
- **Várias premissas ainda não foram medidas**, como a latência e o jitter do nó e a perda dos enlaces entre chips. Todas estão marcadas como `ASSUMPTION` no código, e a [bancada](docs/architecture/13-bancada-experimental-em-fibra.md) foi desenhada para medi-las.
- **Área limita a escala:** o bloco 16×16 ocupa quase um retículo inteiro.
- **O núcleo quântico é conceitual**, e seus detectores exigiriam criogenia (1–4 K).

---

## Documentação

| Tema | Documentos |
| :--- | :--- |
| **Comece por aqui** | [Artigo preliminar](docs/papers/artigo-preliminar.md) · [Roteamento, materiais e race logic](docs/architecture/11-roteamento-e-comutacao-optica.md) |
| **Hardware e lógica** | [Visão geral](docs/architecture/01-visao-geral-hardware.md) · [Lógica por tempo de voo](docs/architecture/02-logica-tempo-de-voo.md) · [Camadas funcionais](docs/architecture/03-unidades-funcionais-alu-gpu-ia.md) · [Lasers e codificação](docs/architecture/09-canhoes-laser-continuos-cw-e-codificacao-multi-nivel.md) |
| **Memória e IA** | [Hierarquia de memória](docs/architecture/04-hierarquia-de-memoria-optica.md) · [Armazenamento em vidro](docs/architecture/05-armazenamento-em-vidro-disco-optico-ssd.md) · [Memória unificada, jogos e IA local](docs/architecture/12-memoria-unificada-jogos-e-ia-local.md) · [Acelerador tensorial](docs/architecture/07-acelerador-tensor-ia-fototectonico.md) · [Andar de IA (PTU)](docs/architecture/15-andar-de-ia-unidade-tensorial-fotonica.md) |
| **Outros núcleos** | [GPU óptica](docs/architecture/06-processamento-de-video-gpu-optica.md) · [Quântico LOQC](docs/architecture/08-processamento-quantico-fotonico-loqc.md) · [Energia](docs/architecture/10-consumo-energetico-e-comparativo-silicio.md) |
| **Experimento e fabricação** | [Bancada em fibra 1550 nm](docs/architecture/13-bancada-experimental-em-fibra.md) · [Fabricação e prototipagem](docs/architecture/14-fabricacao-e-prototipagem.md) |
| **Histórico** | [Whitepaper v1.1](docs/papers/whitepaper-v1.md) · [Especificações](planning/specs.md) |

<details>
<summary><b>Estrutura do repositório</b></summary>

```text
.
├── docs/
│   ├── architecture/        # 15 documentos de arquitetura (01, 05, 06 e 09 marcados como conceito superado)
│   └── papers/              # artigo preliminar, whitepaper e figuras
├── simulations/
│   ├── go/                  # simulador (pkg/optical) e comandos (simulator, racestats, racemultichip, dijkstrabench, tofplatform)
│   ├── ptu/                 # orçamento da unidade tensorial fotônica (doc 15)
│   ├── fdtd/                # simulações eletromagnéticas 2D com Meep
│   └── results/             # CSVs de resultados e gerador da Figura 1
├── web/                     # Portal web e simulador 3D (React + Three.js) -> https://silicacore.rmo.dev.br/
└── planning/                # roadmap, tarefas e especificações
```
</details>

---

## Como citar

```bibtex
@misc{silicacore2026,
  author       = {Oliveira, Ricardo},
  title        = {SilicaCore: race logic fotônica em Si3N4/TFLN com orçamento físico e validação estatística},
  year         = {2026},
  howpublished = {\url{https://github.com/Roliveira96/SilicaCore}}
}
```

## Licença

Código sob [Apache 2.0](LICENSE). Documentação e artigos sob [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).

## Agradecimentos

A **Valmor Moreira**, colega de trabalho na **Grafis**, que explicou ao autor o funcionamento das máquinas de revelação fotográfica a laser durante os reparos que faziam juntos. A ideia e a concepção do SilicaCore são do autor.

<div align="center">
<em>Ciência aberta: cada número deste README é reproduzível com os comandos acima.</em>
</div>
