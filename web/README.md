# SilicaCore Web Experience | Photonic Race Logic Interactive Simulator

Interactive 3D simulation website and global research portal for **SilicaCore**, designed for IEEE, Nature Photonics, and deep-tech conference audiences.

## Architecture

* **Framework:** React 18 + Vite + TypeScript + Tailwind CSS.
* **3D Engine:** Three.js + React Three Fiber (`@react-three/fiber`) + `@react-three/drei`.
* **Physics & Simulation Engine:** High-performance TypeScript implementation mirroring [`racelogic.go`](file:///home/ricardo/Desktop/SilicaCore/simulations/go/pkg/optical/racelogic.go) with exact Box-Muller Gaussian jitter, component loss budgets, Q-factor complementary error function calculations, and real-time pulse wavefront tracing.

## Sections Included

1. **Brand Identity & Hero:** Global deep-tech narrative, 42.2 ns ToF latency, 25.5M error-free paths benchmark.
2. **Section 1 — Die Architecture & Layer Exploded View:** Interactive 3D layer stack with 0–100% explode slider, wireframe toggle, and physical breakdown of Sb₂Se₃, TFLN, Si₃N₄, and SiO₂ substrate tiers.
3. **Section 2 — The Light Race:** Real-time 3D photonic race logic simulator with clickable source/target nodes, 100 ps / 50 ps delay unit toggles, speed sliders, and live HUD telemetry overlay.
4. **Section 3 — The Maxwell Lab:** FDTD Meep bending loss inspector (10 µm to 50 µm with animated electromagnetic field $E_z$ mode profiles) and adiabatic taper visualizer (Si₃N₄ → TFLN saturation at 25 µm).
5. **Section 4 — Benchmark Arena:** Direct side-by-side comparison between classical x86 CPU Dijkstra (memory latency, min-heap pointer chasing) and SilicaCore passive time-of-flight.

## Getting Started

```bash
cd web
npm install
npm run dev
```

To build for production:

```bash
npm run build
npm run preview
```
