# SilicaCore · site de exibição

Site de divulgação do SilicaCore: race logic fotônica em Si₃N₄/TFLN. Todos os números vêm do simulador aberto em `simulations/go` e das simulações FDTD em `simulations/fdtd`. O chip **não foi fabricado**.

## Seções

1. **Hero:** 42,2 ns por consulta (16×16), 0 erros em 25,5 M distâncias, ~42× contra o melhor algoritmo numa CPU atual (estimado) e 648 mm².
2. **How it works:** os três passos de uma consulta e o conceito de qubit fotônico, marcado como pesquisa de longo prazo.
2b. **Decision:** modelo 3D de como um bit é decidido (janelas de tempo na plataforma Si₃N₄) e do qubit dual-rail.
2c. **Comparisons:** tabela com abas de desempenho, potência, temperatura e IA, com a base de cada número.
3. **The chip:** modelo 3D em escala física (`ChipModel3D.tsx`): die de 25,5 mm com 960 espirais, wire bonds, encapsulamento BGA e array de fibras. A corrida 16×16 é a simulada pelo modelo, em falsa cor. Tem vista explodida das camadas e uma corrida 4×4 comparada com Dijkstra medido no navegador.
4. **Waveguide physics:** resultados FDTD 2D (Meep) de curvas e da transição Si₃N₄ → TFLN.
5. **Energy:** para onde vai a energia e o que ainda não foi calculado.
6. **Status:** resultados validados, premissas, limites e próximos passos.

## Rodar

```bash
cd web
npm install
npm run dev      # http://localhost:3000
npm run build    # gera dist/ e copia para build/
```

## Stack

React 18, Vite, TypeScript, Tailwind CSS, Three.js com React Three Fiber, drei e postprocessing (bloom).
