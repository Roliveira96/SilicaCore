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
- [ ] 16. Simulações eletromagnéticas de propagação de ondas em FDTD (Meep / Lumerical)
- [ ] 17. Validação de protótipo de bancada em macro-escala com divisores de feixe cúbicos e detectores SPAD comerciais
- [ ] 18. Submissão do projeto a edital de iniciação científica (PIBIC/PIBITI UTFPR) e congressos (SBC/SBESC, IEEE Photonics)
