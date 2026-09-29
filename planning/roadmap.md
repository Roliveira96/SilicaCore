# Roadmap de Desenvolvimento: Processador Óptico ToF

## Fases do Projeto Deep Tech

### Fase 1: Especificação & Simulação Estatística Monte Carlo (Concluída)
- [x] Modelo físico em Go com execução concorrente em goroutines (`simulations/go/`).
- [x] Especificação técnica de propagação, jitter (VCSEL, SPAD, TDC) e taxa de erro de bit (BER).
- [x] Validação teórica de separação temporal em $8.9\sigma$.

### Fase 2: Simulação Eletromagnética de Ondas (FDTD) (2D concluída)
- [x] FDTD 2D (Meep, índice efetivo) de curvas em Si₃N₄ e da transição Si₃N₄ → TFLN (`simulations/fdtd/`, doc 11 seção 2.1).
- [ ] FDTD 3D ou EME das mesmas estruturas, em máquina com mais memória.
- [ ] Perdas por cruzamento de guias e rugosidade em layout real.

### Fase 3: Validação Acadêmica & Artigo Científico (Em Andamento)
- [ ] Redação do artigo técnico rigoroso em LaTeX/Markdown para submissão a conferências (SBC/SBESC, IEEE Photonics).
- [ ] Apresentação em editais de Iniciação Científica (PIBIC/PIBITI UTFPR).

### Fase 4: Prototipagem Macro-Escala em Bancada
- [ ] Bancada em fibra 1550 nm com moduladores LiNbO₃, atrasos em SMF-28, fotodiodos rápidos e time tagger (ver `docs/architecture/13-bancada-experimental-em-fibra.md`).
