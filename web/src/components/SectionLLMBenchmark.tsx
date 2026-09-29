import React, { useState } from 'react';
import {
  Brain,
  Zap,
  Clock,
  Sparkles,
  GitBranch,
  Network,
  Cpu,
  Layers,
  CheckCircle2,
  TrendingDown,
  Server,
  ArrowRight,
} from 'lucide-react';

interface LLMWorkload {
  id: 'mcts' | 'graph-rag' | 'moe';
  name: string;
  tag: string;
  description: string;
  metricLabel: string;
  unit: string;
  silicaValue: number;
  silicaFormatted: string;
  appleValue: number;
  appleFormatted: string;
  nvidiaValue: number;
  nvidiaFormatted: string;
  amdValue: number;
  amdFormatted: string;
  intelValue: number;
  intelFormatted: string;
  explanation: string;
}

const LLM_WORKLOADS: LLMWorkload[] = [
  {
    id: 'mcts',
    name: 'Tree-of-Thought & Reasoning Graphs (MCTS)',
    tag: 'DeepSeek-R1 / OpenAI o1 / Reasoning Models',
    description:
      'Modern reasoning models generate dynamic trees of thousands of possible thought branches. Evaluating and pruning the best reasoning trajectory requires continuous graph search.',
    metricLabel: 'Reasoning Tree Traversal Latency',
    unit: 'ns (Lower is better)',
    silicaValue: 42.2,
    silicaFormatted: '42.2 ns (0.042 µs)',
    appleValue: 16800,
    appleFormatted: '16,800 ns (16.8 µs)',
    nvidiaValue: 12500,
    nvidiaFormatted: '12,500 ns (12.5 µs)',
    amdValue: 21200,
    amdFormatted: '21,200 ns (21.2 µs)',
    intelValue: 22600,
    intelFormatted: '22,600 ns (22.6 µs)',
    explanation:
      'SilicaCore propagates the entire reasoning graph search at the physical speed of light. Electronic CPUs and GPUs stall on DRAM pointer chases and min-heap rebalancing for every branch evaluation.',
  },
  {
    id: 'graph-rag',
    name: 'Graph RAG & Knowledge Graph Entity Retrieval',
    tag: 'Enterprise AI & Hallucination Elimination',
    description:
      'Traversing billions of semantic relationships in enterprise knowledge graphs to provide verified factual context for LLM generation.',
    metricLabel: 'Entity Pathfinding Query Time',
    unit: 'µs (Lower is better)',
    silicaValue: 0.042,
    silicaFormatted: '0.042 µs (42.2 ns)',
    appleValue: 24.2,
    appleFormatted: '24.2 µs',
    nvidiaValue: 18.5,
    nvidiaFormatted: '18.5 µs',
    amdValue: 29.0,
    amdFormatted: '29.0 µs',
    intelValue: 31.5,
    intelFormatted: '31.5 µs',
    explanation:
      'In Graph RAG, finding shortest multi-hop entity connections is solved in real Time-of-Flight (42 ns), delivering up to 23.7 Million verified entity paths per second.',
  },
  {
    id: 'moe',
    name: 'Mixture-of-Experts (MoE) Token Dispatch',
    tag: 'DeepSeek-V3 / Mixtral / Large MoE Clusters',
    description:
      'Dynamically routing token embeddings to top-k expert networks across distributed multi-chip clusters with zero buffer bloat.',
    metricLabel: 'Expert Dispatch Switch Latency',
    unit: 'ns (Lower is better)',
    silicaValue: 0.8,
    silicaFormatted: '0.8 ns (Photonic Switch)',
    appleValue: 140,
    appleFormatted: '140 ns (Bus Crossbar)',
    nvidiaValue: 95,
    nvidiaFormatted: '95 ns (NVLink Fabric)',
    amdValue: 165,
    amdFormatted: '165 ns (Infinity Fabric)',
    intelValue: 180,
    intelFormatted: '180 ns (UPI Interconnect)',
    explanation:
      'Optical routing in thin-film lithium niobate (TFLN) switches token streams in sub-nanosecond delay, eliminating all-to-all GPU communication bottlenecks.',
  },
];

export const SectionLLMBenchmark: React.FC = () => {
  const [selectedWorkload, setSelectedWorkload] = useState<LLMWorkload>(LLM_WORKLOADS[0]);

  return (
    <section id="section-llm" className="relative py-28 bg-black border-t border-white/5 overflow-hidden">
      {/* Background glow */}
      <div className="absolute top-1/3 left-1/2 -translate-x-1/2 w-[750px] h-[400px] bg-gradient-to-tr from-purple-600/10 via-cyan-600/10 to-transparent blur-[180px] rounded-full pointer-events-none -z-10" />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          <div className="inline-flex items-center space-x-2 px-3.5 py-1.5 rounded-full bg-purple-950/60 border border-purple-500/30 text-purple-300 text-xs font-mono mb-4">
            <Brain className="w-3.5 h-3.5 text-purple-400" />
            <span>AI &amp; LARGE LANGUAGE MODEL (LLM) ACCELERATION</span>
          </div>
          <h2 className="text-4xl sm:text-5xl font-black text-white tracking-tight">
            Accelerating Modern Reasoning LLMs.
          </h2>
          <p className="mt-4 text-base sm:text-lg text-slate-300 font-light leading-relaxed">
            The biggest bottleneck in next-generation reasoning models (like DeepSeek-R1 and OpenAI o1) isn’t matrix multiplication — it is <strong>searching massive trees of thought</strong> and <strong>routing through complex knowledge graphs</strong>.
          </p>

          {/* Workload Selector Tabs */}
          <div className="mt-10 inline-flex flex-wrap justify-center items-center p-1.5 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl gap-1">
            {LLM_WORKLOADS.map((w) => (
              <button
                key={w.id}
                onClick={() => setSelectedWorkload(w)}
                className={`px-5 py-2.5 rounded-full text-xs font-semibold transition-all ${
                  selectedWorkload.id === w.id
                    ? 'bg-gradient-to-r from-cyan-400 to-purple-400 text-slate-950 shadow-[0_0_15px_rgba(168,85,247,0.4)]'
                    : 'text-slate-300 hover:text-white'
                }`}
              >
                {w.name.split('(')[0]}
              </button>
            ))}
          </div>
        </div>

        {/* Main Workload Card */}
        <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-b from-white/[0.04] to-black border border-white/10 shadow-2xl mb-12">
          {/* Card Header */}
          <div className="flex flex-col md:flex-row md:items-center justify-between pb-6 border-b border-white/10 gap-4 mb-8">
            <div>
              <div className="flex items-center space-x-2">
                <span className="px-3 py-1 rounded-full text-xs font-mono bg-purple-950/80 border border-purple-500/40 text-purple-300 font-bold">
                  {selectedWorkload.tag}
                </span>
                <span className="text-xs font-mono text-cyan-400 font-bold">100% Optical Acceleration</span>
              </div>
              <h3 className="text-2xl font-bold text-white mt-2">{selectedWorkload.name}</h3>
              <p className="text-sm text-slate-300 font-normal mt-1 max-w-3xl leading-relaxed">
                {selectedWorkload.description}
              </p>
            </div>

            <div className="text-right">
              <span className="text-xs font-mono text-slate-400 block">{selectedWorkload.metricLabel}</span>
              <span className="text-2xl font-black text-cyan-400 font-mono">
                {selectedWorkload.silicaFormatted}
              </span>
            </div>
          </div>

          {/* Comparative Benchmark Bars */}
          <div className="space-y-4 mb-8">
            {/* 1. SilicaCore */}
            <div className="p-4 rounded-2xl bg-cyan-950/30 border border-cyan-400/50 shadow-[0_0_20px_rgba(0,242,254,0.15)] ring-1 ring-cyan-400/30">
              <div className="flex items-center justify-between text-xs font-mono mb-2">
                <div className="flex items-center space-x-2">
                  <span className="w-3 h-3 rounded-full bg-cyan-400 shadow-[0_0_8px_#00f2fe]" />
                  <span className="font-bold text-cyan-300 text-sm">SilicaCore Photonic Accelerator</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-cyan-950 border border-cyan-400 text-cyan-300 font-bold">
                    🏆 #1 WINNER
                  </span>
                </div>
                <div className="text-base font-extrabold text-cyan-400 font-mono">
                  {selectedWorkload.silicaFormatted}
                </div>
              </div>
              <div className="w-full bg-slate-900 h-2.5 rounded-full overflow-hidden p-0.5 border border-white/10">
                <div className="bg-cyan-400 h-full w-[4%] rounded-full shadow-[0_0_12px_rgba(0,242,254,0.9)] animate-pulse" />
              </div>
              <div className="flex justify-between items-center mt-2 text-xs font-mono text-slate-400">
                <span className="text-cyan-300 font-semibold">Speed of Light Traversal (Zero Memory Stalls)</span>
                <span className="text-white">1.0× Reference</span>
              </div>
            </div>

            {/* 2. NVIDIA Blackwell B200 */}
            <div className="p-4 rounded-2xl bg-white/[0.02] hover:bg-white/[0.04] border border-white/5 transition-all">
              <div className="flex items-center justify-between text-xs font-mono mb-2">
                <div className="flex items-center space-x-2">
                  <span className="w-3 h-3 rounded-full bg-emerald-400" />
                  <span className="font-bold text-white text-sm">NVIDIA Blackwell B200 GPU</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-white/5 text-slate-400">
                    TSMC 4NP • 1,000 W
                  </span>
                </div>
                <div className="text-base font-extrabold text-emerald-400 font-mono">
                  {selectedWorkload.nvidiaFormatted}
                </div>
              </div>
              <div className="w-full bg-slate-900 h-2 rounded-full overflow-hidden border border-white/5">
                <div className="bg-emerald-500 h-full w-[60%] rounded-full" />
              </div>
              <div className="flex justify-between items-center mt-2 text-xs font-mono text-slate-400">
                <span>CUDA Kernel Launch Overhead + Divergent Branching</span>
                <span className="text-emerald-300">~296× slower than SilicaCore</span>
              </div>
            </div>

            {/* 3. Apple M5 Max */}
            <div className="p-4 rounded-2xl bg-white/[0.02] hover:bg-white/[0.04] border border-white/5 transition-all">
              <div className="flex items-center justify-between text-xs font-mono mb-2">
                <div className="flex items-center space-x-2">
                  <span className="w-3 h-3 rounded-full bg-purple-400" />
                  <span className="font-bold text-white text-sm">Apple M5 Max</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-white/5 text-slate-400">
                    TSMC 2nm GAAFET • 80 W
                  </span>
                </div>
                <div className="text-base font-extrabold text-purple-400 font-mono">
                  {selectedWorkload.appleFormatted}
                </div>
              </div>
              <div className="w-full bg-slate-900 h-2 rounded-full overflow-hidden border border-white/5">
                <div className="bg-purple-500 h-full w-[78%] rounded-full" />
              </div>
              <div className="flex justify-between items-center mt-2 text-xs font-mono text-slate-400">
                <span>Unified Memory Bus Contention &amp; Min-Heap Serialization</span>
                <span className="text-purple-300">~398× slower than SilicaCore</span>
              </div>
            </div>

            {/* 4. AMD Ryzen 9 9950X */}
            <div className="p-4 rounded-2xl bg-white/[0.02] hover:bg-white/[0.04] border border-white/5 transition-all">
              <div className="flex items-center justify-between text-xs font-mono mb-2">
                <div className="flex items-center space-x-2">
                  <span className="w-3 h-3 rounded-full bg-orange-400" />
                  <span className="font-bold text-white text-sm">AMD Ryzen 9 9950X</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-white/5 text-slate-400">
                    TSMC 4nm • 200 W
                  </span>
                </div>
                <div className="text-base font-extrabold text-orange-400 font-mono">
                  {selectedWorkload.amdFormatted}
                </div>
              </div>
              <div className="w-full bg-slate-900 h-2 rounded-full overflow-hidden border border-white/5">
                <div className="bg-orange-500 h-full w-[88%] rounded-full" />
              </div>
              <div className="flex justify-between items-center mt-2 text-xs font-mono text-slate-400">
                <span>DDR5-6000 Random Access Latency (~65 ns/hop)</span>
                <span className="text-orange-300">~502× slower than SilicaCore</span>
              </div>
            </div>

            {/* 5. Intel Core Ultra 9 285K */}
            <div className="p-4 rounded-2xl bg-white/[0.02] hover:bg-white/[0.04] border border-white/5 transition-all">
              <div className="flex items-center justify-between text-xs font-mono mb-2">
                <div className="flex items-center space-x-2">
                  <span className="w-3 h-3 rounded-full bg-blue-400" />
                  <span className="font-bold text-white text-sm">Intel Core Ultra 9 285K</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-white/5 text-slate-400">
                    Intel Foveros 3D • 250 W
                  </span>
                </div>
                <div className="text-base font-extrabold text-blue-400 font-mono">
                  {selectedWorkload.intelFormatted}
                </div>
              </div>
              <div className="w-full bg-slate-900 h-2 rounded-full overflow-hidden border border-white/5">
                <div className="bg-blue-500 h-full w-[96%] rounded-full" />
              </div>
              <div className="flex justify-between items-center mt-2 text-xs font-mono text-slate-400">
                <span>Cross-Tile Interconnect Latency &amp; Memory Controller Stalls</span>
                <span className="text-blue-300">~535× slower than SilicaCore</span>
              </div>
            </div>
          </div>

          {/* Explanation Box */}
          <div className="p-5 rounded-2xl bg-white/[0.02] border border-white/5 flex items-start space-x-3 text-xs text-slate-300">
            <Sparkles className="w-5 h-5 text-cyan-400 mt-0.5 flex-shrink-0" />
            <div className="space-y-1">
              <span className="font-bold text-white font-mono uppercase tracking-wider text-xs">
                Physical AI Advantage:
              </span>
              <p className="leading-relaxed text-slate-300 font-normal">
                {selectedWorkload.explanation}
              </p>
            </div>
          </div>
        </div>

        {/* 3 Bento Cards: Why LLMs Need Photonic Co-Processors */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="p-6 rounded-3xl bg-white/[0.02] border border-white/10">
            <GitBranch className="w-6 h-6 text-purple-400 mb-3" />
            <h4 className="text-base font-bold text-white mb-1">MCTS Tree-of-Thought</h4>
            <p className="text-xs text-slate-300 leading-relaxed font-normal">
              Reasoning models like DeepSeek-R1 evaluate thousands of candidate thoughts per prompt. SilicaCore solves shortest optimal reasoning routes in 42 ns, slashing inference time from seconds to milliseconds.
            </p>
          </div>

          <div className="p-6 rounded-3xl bg-white/[0.02] border border-white/10">
            <Network className="w-6 h-6 text-cyan-400 mb-3" />
            <h4 className="text-base font-bold text-white mb-1">Zero-Hallucination Graph RAG</h4>
            <p className="text-xs text-slate-300 leading-relaxed font-normal">
              Instead of relying purely on statistical next-token guessing, LLMs verify factual entities across multi-hop knowledge graphs at 23.7 Million path queries per second.
            </p>
          </div>

          <div className="p-6 rounded-3xl bg-white/[0.02] border border-white/10">
            <Cpu className="w-6 h-6 text-emerald-400 mb-3" />
            <h4 className="text-base font-bold text-white mb-1">Sub-Nanosecond MoE Routing</h4>
            <p className="text-xs text-slate-300 leading-relaxed font-normal">
              Mixture-of-Experts architectures switch token embeddings optically via TFLN phase gates in 0.8 ns, eliminating all-to-all electronic network congestion in datacenter GPU clusters.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
