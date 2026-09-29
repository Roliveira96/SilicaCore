import React, { useState } from 'react';
import {
  Cpu,
  Zap,
  ArrowRight,
  Clock,
  Layers,
  Flame,
  CheckCircle2,
  XCircle,
  Activity,
  GitBranch,
  ShieldAlert,
  Sparkles,
  Server,
} from 'lucide-react';

interface MetricComparison {
  category: string;
  siliconTitle: string;
  siliconDesc: string;
  siliconDetail: string;
  photonicTitle: string;
  photonicDesc: string;
  photonicDetail: string;
  highlightAdvantage: string;
}

const COMPARISONS: MetricComparison[] = [
  {
    category: 'Architecture Paradigm',
    siliconTitle: 'Von Neumann Architecture',
    siliconDesc: 'Separated compute and memory units. Constantly moves instructions and graph matrices across buses.',
    siliconDetail: 'Heavy L1/L2/L3 cache hierarchies; memory bus arbitration bottlenecks.',
    photonicTitle: 'Photonic Race Logic (In-Hardware Graph)',
    photonicDesc: 'The graph topology is directly mapped onto physical planar delay waveguides. Zero instruction fetch.',
    photonicDetail: 'Sb₂Se₃ non-volatile phase change switches set weights with 0 static power.',
    highlightAdvantage: 'Zero Instruction Overhead',
  },
  {
    category: 'Signal Medium',
    siliconTitle: 'Electrons through Copper Traces',
    siliconDesc: 'Electron transport governed by drift velocity and parasitic capacitance (RC delay).',
    siliconDetail: 'Drift velocity ~1 mm/s; wire propagation ~0.3c with severe resistive signal degradation.',
    photonicTitle: 'Photons in Si₃N₄ Waveguides',
    photonicDesc: 'Light traveling in dielectric waveguides at physical phase velocity c/n_g (~160,000 km/s).',
    photonicDetail: 'Continuous-wave 1550 nm carrier; 0.05 dB/cm ultra-low loss Si₃N₄ spirals.',
    highlightAdvantage: 'Speed of Light Propagation',
  },
  {
    category: 'Scaling & Query Complexity',
    siliconTitle: 'Sequential Dijkstra: O(E + V log V)',
    siliconDesc: 'Central CPU must pop nodes from a priority min-heap one-by-one, relaxing edges in serialized loops.',
    siliconDetail: 'Heavy branch mispredictions and pointer-chasing memory stalls scale with graph density.',
    photonicTitle: 'Wavefront Transit + Readout: O(D·τ + N/R)',
    photonicDesc: 'Not O(1): transit scales with maximum graph distance (D·τ), and electronic TDC readout scales with node count (N).',
    photonicDetail: 'At 16×16, electronic readout (30.7 ns) dominates light transit (11.5 ns), delivering ~150–800× speedup over modern desktop CPUs.',
    highlightAdvantage: 'Parallel Physical Wavefront',
  },
  {
    category: 'Latency per Query (16×16)',
    siliconTitle: '~50 – 100 µs (~84,000 ns)',
    siliconDesc: 'Millions of clock cycles spent managing heaps and reading DRAM state sequentially.',
    siliconDetail: 'Measured ~84 µs baseline on single-core 2012 x86 CPU (Intel i3-3217U).',
    photonicTitle: '42.2 Nanoseconds',
    photonicDesc: '11.5 ns optical transit race + 30.7 ns high-speed Time-to-Digital Converter (TDC) readout.',
    photonicDetail: 'Measured ~2,000× vs 2012 baseline; estimated ~150–800× vs modern high-end desktop CPUs.',
    highlightAdvantage: '~150–800× Modern Speedup',
  },
  {
    category: 'Dynamic Heat Dissipation',
    siliconTitle: 'Severe Joule Heating (P = I²R)',
    siliconDesc: '> 80% of supplied energy is lost as resistive thermal dissipation across copper interconnects.',
    siliconDetail: 'Requires high-maintenance 360mm liquid cooling loops and vapor chambers.',
    photonicTitle: '0 Joules Dynamic Heating in Waveguides',
    photonicDesc: 'Photons travel passively through dielectric SiO₂/Si₃N₄ cores without electric resistance.',
    photonicDetail: 'Heat is strictly limited to laser diodes and receiver TDCs at chip edges.',
    highlightAdvantage: 'Zero Interconnect Heat',
  },
  {
    category: 'Memory Wall & Bandwidth',
    siliconTitle: 'DDR5 / HBM Memory Bottleneck',
    siliconDesc: 'Moving 1 bit from external DRAM costs 10–100 pJ, causing frequent ALU pipeline starvation.',
    siliconDetail: 'DDR5: ~90 GB/s; HBM3: ~3.3 TB/s with high power draw.',
    photonicTitle: 'Hardware Delays Replace Memory',
    photonicDesc: 'Edge weights are stored as physical optical delays, eliminating pointer lookups entirely.',
    photonicDetail: 'Unified optical I/O bandwidth ready for terabit-scale photonic packaging.',
    highlightAdvantage: 'Bypasses Memory Wall',
  },
];

export const SectionComparison: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'matrix' | 'pipeline'>('matrix');

  return (
    <section id="section-comparison" className="relative py-24 bg-[#03060d] border-t border-slate-800/80">
      {/* Background ambient lighting */}
      <div className="absolute top-1/2 left-1/4 -translate-y-1/2 w-96 h-96 bg-cyan-500/10 rounded-full blur-[140px] pointer-events-none" />
      <div className="absolute top-1/2 right-1/4 -translate-y-1/2 w-96 h-96 bg-rose-500/10 rounded-full blur-[140px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-14">
          <div className="inline-flex items-center space-x-2 px-3.5 py-1.5 rounded-full bg-slate-900/90 border border-cyan-500/30 text-cyan-300 text-xs font-mono mb-4 shadow-[0_0_15px_rgba(6,182,212,0.2)]">
            <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
            <span>Architectural Paradigm Shift</span>
          </div>
          <h2 className="text-3xl sm:text-5xl font-extrabold text-white tracking-tight">
            Conventional Processor vs.{' '}
            <span className="bg-gradient-to-r from-cyan-300 via-teal-300 to-blue-400 bg-clip-text text-transparent glow-text-cyan">
              SilicaCore
            </span>
          </h2>
          <p className="mt-4 text-base text-slate-300 leading-relaxed">
            Conventional microprocessors are trapped by the Von Neumann bottleneck, sequential clock cycles, and copper Joule heating. SilicaCore replaces instructions with optical propagation physics.
          </p>

          {/* Tab Selector */}
          <div className="mt-8 inline-flex items-center p-1 rounded-xl bg-slate-900/80 border border-slate-800 backdrop-blur-md">
            <button
              onClick={() => setActiveTab('matrix')}
              className={`px-5 py-2 rounded-lg text-xs font-mono font-medium transition-all ${
                activeTab === 'matrix'
                  ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/50 shadow-[0_0_12px_rgba(6,182,212,0.3)]'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Head-to-Head Comparison Matrix
            </button>
            <button
              onClick={() => setActiveTab('pipeline')}
              className={`px-5 py-2 rounded-lg text-xs font-mono font-medium transition-all ${
                activeTab === 'pipeline'
                  ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/50 shadow-[0_0_12px_rgba(6,182,212,0.3)]'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Execution Pipeline Walkthrough
            </button>
          </div>
        </div>

        {/* Tab 1: Comprehensive Comparison Matrix */}
        {activeTab === 'matrix' && (
          <div className="space-y-4">
            {/* Header Columns */}
            <div className="hidden lg:grid grid-cols-12 gap-4 px-6 py-3 text-xs font-mono uppercase tracking-wider text-slate-400 border-b border-slate-800">
              <div className="col-span-3">Architectural Metric</div>
              <div className="col-span-4 flex items-center gap-1.5 text-rose-400">
                <Cpu className="w-3.5 h-3.5" />
                <span>Conventional Silicon (x86/ARM/GPU)</span>
              </div>
              <div className="col-span-5 flex items-center gap-1.5 text-cyan-400">
                <Zap className="w-3.5 h-3.5" />
                <span>SilicaCore Photonic Accelerator</span>
              </div>
            </div>

            {/* Comparison Cards */}
            {COMPARISONS.map((item, idx) => (
              <div
                key={idx}
                className="glass-panel p-5 sm:p-6 rounded-2xl border border-slate-800/80 hover:border-slate-700 transition-all duration-300 grid grid-cols-1 lg:grid-cols-12 gap-6 items-center"
              >
                {/* Metric Column */}
                <div className="lg:col-span-3">
                  <div className="text-xs font-mono uppercase text-cyan-400 font-semibold mb-1">
                    {item.category}
                  </div>
                  <span className="inline-block px-2 py-0.5 rounded text-[10px] font-mono bg-cyan-950/80 border border-cyan-500/30 text-cyan-300 font-medium">
                    {item.highlightAdvantage}
                  </span>
                </div>

                {/* Silicon Processor Column */}
                <div className="lg:col-span-4 p-4 rounded-xl bg-slate-950/70 border border-rose-500/20 space-y-1.5">
                  <div className="flex items-center gap-2 text-rose-400 font-semibold text-sm">
                    <XCircle className="w-4 h-4 flex-shrink-0" />
                    <span>{item.siliconTitle}</span>
                  </div>
                  <p className="text-xs text-slate-300 leading-relaxed">{item.siliconDesc}</p>
                  <p className="text-[11px] font-mono text-slate-500 pt-1">{item.siliconDetail}</p>
                </div>

                {/* SilicaCore Column */}
                <div className="lg:col-span-5 p-4 rounded-xl bg-gradient-to-br from-cyan-950/40 via-slate-950/70 to-blue-950/30 border border-cyan-500/40 shadow-[0_0_20px_rgba(6,182,212,0.1)] space-y-1.5 ring-1 ring-cyan-500/20">
                  <div className="flex items-center gap-2 text-cyan-300 font-semibold text-sm">
                    <CheckCircle2 className="w-4 h-4 text-cyan-400 flex-shrink-0" />
                    <span>{item.photonicTitle}</span>
                  </div>
                  <p className="text-xs text-slate-200 leading-relaxed">{item.photonicDesc}</p>
                  <p className="text-[11px] font-mono text-cyan-300/80 pt-1">{item.photonicDetail}</p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Tab 2: Visual Step-by-Step Execution Pipeline */}
        {activeTab === 'pipeline' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-stretch">
            {/* Silicon Pipeline */}
            <div className="glass-panel p-6 sm:p-8 rounded-2xl border border-rose-500/30 relative overflow-hidden flex flex-col justify-between">
              <div className="flex items-center justify-between mb-6 pb-4 border-b border-slate-800">
                <div className="flex items-center space-x-3">
                  <div className="p-2.5 rounded-xl bg-rose-500/10 border border-rose-500/30">
                    <Cpu className="w-6 h-6 text-rose-400" />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-white">x86-64 Classical Pipeline</h3>
                    <p className="text-xs font-mono text-rose-400">Sequential Software Execution</p>
                  </div>
                </div>
                <span className="px-2.5 py-1 rounded bg-rose-950/60 border border-rose-500/30 text-xs font-mono text-rose-300">
                  ~84,000 ns
                </span>
              </div>

              <div className="space-y-4">
                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-slate-800">
                  <div className="text-xs font-mono text-slate-400 font-semibold mb-1">
                    Step 1: Graph Loading from DRAM
                  </div>
                  <p className="text-xs text-slate-300">
                    Reads adjacency matrix from DRAM into L3/L2 caches. Incurs ~150–200 clock cycles stall on initial cache misses.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-slate-800">
                  <div className="text-xs font-mono text-slate-400 font-semibold mb-1">
                    Step 2: Min-Heap Queue Balancing
                  </div>
                  <p className="text-xs text-slate-300">
                    Every node relaxation pushes elements into a priority queue. Re-sifting heap arrays causes heavy pointer chasing and branching.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-slate-800">
                  <div className="text-xs font-mono text-slate-400 font-semibold mb-1">
                    Step 3: Sequential Edge Relaxation
                  </div>
                  <p className="text-xs text-slate-300">
                    ALUs calculate tentative distances node-by-node. For a 16×16 graph, tens of thousands of instructions serialize through pipeline hazards.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-rose-950/30 border border-rose-500/30 text-xs text-rose-300 font-mono">
                  Bottleneck: High latency, memory bus contention, and exponential clock cycle scaling.
                </div>
              </div>
            </div>

            {/* SilicaCore Pipeline */}
            <div className="glass-panel p-6 sm:p-8 rounded-2xl border border-cyan-500/40 relative overflow-hidden flex flex-col justify-between shadow-[0_0_30px_rgba(6,182,212,0.15)] ring-1 ring-cyan-500/20">
              <div className="flex items-center justify-between mb-6 pb-4 border-b border-slate-800">
                <div className="flex items-center space-x-3">
                  <div className="p-2.5 rounded-xl bg-cyan-500/10 border border-cyan-500/40">
                    <Zap className="w-6 h-6 text-cyan-400" />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-white">SilicaCore Photonic Pipeline</h3>
                    <p className="text-xs font-mono text-cyan-300">Time-of-Flight Physics Execution</p>
                  </div>
                </div>
                <span className="px-2.5 py-1 rounded bg-cyan-950/60 border border-cyan-500/30 text-xs font-mono text-cyan-300 font-bold">
                  42.2 ns
                </span>
              </div>

              <div className="space-y-4">
                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-cyan-500/30">
                  <div className="text-xs font-mono text-cyan-300 font-semibold mb-1">
                    Step 1: Electro-Optic Pulse Release (t = 0)
                  </div>
                  <p className="text-xs text-slate-300">
                    A TFLN electro-optic modulator fires a 1550 nm optical pulse into the source node. No instructions or memory fetches needed.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-cyan-500/30">
                  <div className="text-xs font-mono text-cyan-300 font-semibold mb-1">
                    Step 2: Passive Waveguide Transit (~11.5 ns)
                  </div>
                  <p className="text-xs text-slate-300">
                    Photons branch through binary-weighted Si₃N₄ spiral delay lines at physical speed of light. All possible paths compete concurrently.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-slate-950/80 border border-cyan-500/30">
                  <div className="text-xs font-mono text-cyan-300 font-semibold mb-1">
                    Step 3: First-Arrival Latching &amp; TDC Readout (~30.7 ns)
                  </div>
                  <p className="text-xs text-slate-300">
                    The fastest wavefront triggers high-speed UTC photodiodes and latches the TDC. The earliest arrival time *is* the shortest path.
                  </p>
                </div>

                <div className="p-3.5 rounded-xl bg-cyan-950/30 border border-cyan-500/30 text-xs text-cyan-300 font-mono">
                  Verified Reality: Zero memory lookups and zero branch hazards; latency bounded by physical graph diameter (11.5 ns) and TDC electronic readout (30.7 ns).
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </section>
  );
};
