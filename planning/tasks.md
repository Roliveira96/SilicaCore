# Lista de Tarefas Técnicas Consolidadas (SilicaCore)

## Tarefas Concluídas
- [x] 1. Especificação técnica detalhada da arquitetura ToF em sílica fundida (`docs/architecture/01-visao-geral-hardware.md`) [complexidade: media] [concluída em 28-09-26]
- [x] 2. Desenvolvimento do motor de simulação estatística Monte Carlo em Go paralelizado em Goroutines (`simulations/go/`) [complexidade: alta] [concluída em 28-09-26]
- [x] 3. Internacionalização do código Go e relatório CLI para Inglês técnico [complexidade: baixa] [concluída em 28-09-26]
- [x] 4. Modelagem da Hierarquia de Memória Fotônica (Cache L1 micro-anéis, RAM recirculante e ROM em $SiO_2$) [complexidade: alta] [concluída em 28-09-26]
- [x] 5. Modelagem do Photonic SSD em vidro de sílica ($100\text{ TB}$ / $1.2\text{ TB/s}$) [complexidade: alta] [concluída em 28-09-26]
- [x] 6. Arquitetura da GPU Óptica por Multiplexação WDM RGB (635nm/532nm/450nm) e Ray-Tracing Nativo [complexidade: altíssima] [concluída em 28-09-26]
- [x] 7. Arquitetura do Photonic AI Tensor Core (Multiplicação MVM por malha MZI e PCM $GST$) [complexidade: altíssima] [concluída em 28-09-26]
- [x] 8. Arquitetura do Processador Quântico Fotônico Híbrido (LOQC, qubits dual-rail, interferência HOM) [complexidade: altíssima] [concluída em 28-09-26]
- [x] 9. Arquitetura do Motor Laser Contínuo CW e Codificação Densa M-ária em Byte (8-bit) [complexidade: alta] [concluída em 28-09-26]
- [x] 10. Redação completa do Whitepaper v1.0 com 11 referências científicas peer-reviewed (*Nature*, *Science*, *PRL*, *IEEE*) e agradecimentos à empresa Grafis e Valmor Moreira (`docs/papers/whitepaper-v1.md`) [complexidade: altíssima] [concluída em 28-09-26]

## Próximos Passos (Roadmap de Pesquisa)
- [x] 11. Esclarecimento da inspiração conceitual do motor laser CW na documentação e aprimoramento físico do simulador Go (RIN noise, interferência quântica HOM, M-ary encoding) [complexidade: alta] [concluída em 28-09-26]
- [x] 12. Calibração estatística automática da janela de tempo (Time-Gating Window $W_{\text{gate}} = 6\sigma$) com 100+ execuções de teste para zerar erros de leitura ToF [complexidade: média] [concluída em 28-09-26]
- [x] 13. Modelagem técnica e simulação de consumo energético (fJ/bit, TOPS/W, TDP Watts) e comparativo com processadores de silício de alta performance (Intel i9, AMD EPYC, NVIDIA H100) [complexidade: alta] [concluída em 28-09-26]
- [x] 14. Ajuste da codificação M-ária para 16 estados (4 bits/símbolo - Hexadecimal / 16-QAM) e inclusão da tabela comparativa de latência L1/Ciclo no main.go [complexidade: baixa] [concluída em 28-09-26]
- [x] 15. Redução geométrica do bloco fotônico (d1 = 2.0 mm / latência ~9.67 ps) e expansão DWDM massiva (64 canais espectrais / 256 bits/pulso) no simulador Go [complexidade: alta] [concluída em 28-09-26]
- [x] 19. Orçamento físico de roteamento (espelhos internos vs Si₃N₄/TFLN), BER corrigida e modelo de memória unificada / IA local com literatura 2024–2026 (`budget.go`, `memory.go`, docs 11 e 12) [complexidade: alta] [concluída em 28-09-26]
- [x] 20. Revisar afirmações do README/whitepaper/specs (8.9σ → Q=4.45, 206 GHz → ~5 GHz/canal, >100 TOPS/W no sistema, ray-tracing nativo, quântico sem criogenia) [complexidade: média] [concluída em 28-09-26]
- [x] 22. Alinhar docs de arquitetura 01–10 com a simulação (EOM/TIR em SiO₂, espelhos internos, ray-tracing, 298 K, 18.5 W, 0.05 pJ/bit, L1 ≤5 ps) e geometria agnóstica do substrato (cubo → substrato fotônico integrado) [complexidade: média] [concluída em 28-09-26]
- [x] 21. Protótipo em Go de race logic fotônica (menor caminho em grafo por rede de atrasos programáveis) [complexidade: alta] [concluída em 28-09-26]
- [x] 23. Campanha estatística de race logic (10 chips × 10⁴ consultas, IC 95%, validação por número de saltos); unidade padrão elevada para 100 ps (`cmd/racestats`, `simulations/results/`) [complexidade: média] [concluída em 28-09-26]
- [x] 24. Alinhar núcleo do simulador (`core.go`/`tof.go`) aos docs: L1 pSRAM 25 ps, hierarquia L1→L2/L3→HBM (~416 ps), energia pela taxa agregada (~1.9 pJ/bit), premissas de TDP/SSD, suposições do modelo quântico e correção da perda de guia (0.4 dB) [complexidade: média] [concluída em 28-09-26]
- [x] 25. Rascunho do artigo preliminar (race logic fotônica em Si₃N₄/TFLN) com Figura 1 de erro medido vs previsto gerada do CSV (`docs/papers/artigo-preliminar.md`, `simulations/results/plot_race_logic_hops.py`) [complexidade: alta] [concluída em 28-09-26]
- [ ] 16. Simulações eletromagnéticas de propagação de ondas em FDTD (Meep / Lumerical)
- [ ] 17. Validação de protótipo de bancada em macro-escala com divisores de feixe cúbicos e detectores SPAD comerciais
- [ ] 18. Submissão do projeto a edital de iniciação científica (PIBIC/PIBITI UTFPR) e congressos (SBC/SBESC, IEEE Photonics)
