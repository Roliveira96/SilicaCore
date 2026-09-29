import React from 'react';
import { Zap, ShieldCheck, Flame, Server, Cpu, TrendingUp, Sparkles, Building2, CheckCircle2 } from 'lucide-react';

export const SectionAdvantages: React.FC = () => {
  return (
    <section id="section-advantages" className="relative py-28 bg-[#030408] border-t border-white/5 overflow-hidden">
      {/* Background glow */}
      <div className="absolute top-1/3 left-1/2 -translate-x-1/2 w-[700px] h-[350px] bg-gradient-to-tr from-cyan-600/10 via-blue-600/10 to-transparent blur-[160px] rounded-full pointer-events-none -z-10" />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-20">
          <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/60 border border-cyan-500/30 text-cyan-300 text-xs font-mono mb-4">
            <TrendingUp className="w-3.5 h-3.5 text-cyan-400" />
            <span>INVESTOR &amp; COMMERCIAL VALUE PROPOSITION</span>
          </div>
          <h2 className="text-4xl sm:text-5xl font-black text-white tracking-tight">
            The SilicaCore Advantage.
          </h2>
          <p className="mt-4 text-base sm:text-lg text-slate-400 font-light leading-relaxed">
            Why traditional silicon has reached its physical limits, and how photonic computing creates an insurmountable economic and performance moat.
          </p>
        </div>

        {/* 4 Apple-Style Bento Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
          {/* Bento Card 1: Memory Wall Solved */}
          <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-br from-white/[0.03] to-white/[0.01] border border-white/10 hover:border-cyan-500/30 transition-all flex flex-col justify-between">
            <div>
              <div className="w-12 h-12 rounded-2xl bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400 mb-6">
                <Cpu className="w-6 h-6" />
              </div>
              <span className="text-xs font-mono text-cyan-400 uppercase tracking-wider font-semibold">
                Architecture Breakthrough
              </span>
              <h3 className="text-2xl font-bold text-white mt-1 mb-3">
                Bypassing the Von Neumann Memory Wall
              </h3>
              <p className="text-slate-400 text-sm leading-relaxed font-light">
                Modern CPUs and GPUs spend over <strong>80% of their time and energy</strong> waiting for data to travel from DRAM to compute cores. SilicaCore eliminates registers, ALU cycles, and cache misses entirely. The calculation occurs continuously <strong>in-flight</strong> as light propagates across dielectric waveguides.
              </p>
            </div>

            <div className="mt-8 pt-5 border-t border-white/5 flex items-center justify-between text-xs font-mono text-slate-400">
              <span>Time-to-Solution:</span>
              <span className="text-cyan-400 font-bold">42.2 ns Continuous</span>
            </div>
          </div>

          {/* Bento Card 2: Zero Joule Heating */}
          <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-br from-white/[0.03] to-white/[0.01] border border-white/10 hover:border-amber-500/30 transition-all flex flex-col justify-between">
            <div>
              <div className="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400 mb-6">
                <Flame className="w-6 h-6" />
              </div>
              <span className="text-xs font-mono text-amber-400 uppercase tracking-wider font-semibold">
                Thermal &amp; Power Superiority
              </span>
              <h3 className="text-2xl font-bold text-white mt-1 mb-3">
                0 W Dynamic Joule Heating in Waveguides
              </h3>
              <p className="text-slate-400 text-sm leading-relaxed font-light">
                Photons have zero rest mass and zero electrical charge. Propagating light produces <strong>zero resistive heat</strong> ($P = I^2 R = 0$). While modern GPUs demand thousands of watts and noisy liquid chillers, SilicaCore operates near ambient temperature (<strong>28°C</strong>) using simple passive conduction pads.
              </p>
            </div>

            <div className="mt-8 pt-5 border-t border-white/5 flex items-center justify-between text-xs font-mono text-slate-400">
              <span>System TDP:</span>
              <span className="text-amber-400 font-bold">18.5 W (54× Lower than GPUs)</span>
            </div>
          </div>

          {/* Bento Card 3: Standard 300mm Foundry Ready */}
          <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-br from-white/[0.03] to-white/[0.01] border border-white/10 hover:border-blue-500/30 transition-all flex flex-col justify-between">
            <div>
              <div className="w-12 h-12 rounded-2xl bg-blue-500/10 border border-blue-500/30 flex items-center justify-center text-blue-400 mb-6">
                <Building2 className="w-6 h-6" />
              </div>
              <span className="text-xs font-mono text-blue-400 uppercase tracking-wider font-semibold">
                Commercial Manufacturability
              </span>
              <h3 className="text-2xl font-bold text-white mt-1 mb-3">
                Fabricated on Commercial Silicon Lines
              </h3>
              <p className="text-slate-400 text-sm leading-relaxed font-light">
                SilicaCore is engineered for high-volume commercial production. It uses industry-standard <strong>Silicon Nitride (Si₃N₄)</strong> and <strong>Thin-Film Lithium Niobate (TFLN)</strong> on thermal oxide, fully compatible with existing 300mm CMOS foundry processes (TSMC, GlobalFoundries). No exotic or unproven manufacturing plants required.
              </p>
            </div>

            <div className="mt-8 pt-5 border-t border-white/5 flex items-center justify-between text-xs font-mono text-slate-400">
              <span>Die Footprint:</span>
              <span className="text-blue-400 font-bold">648 mm² (Fits Standard Reticle)</span>
            </div>
          </div>

          {/* Bento Card 4: Massive $85B Addressable Market */}
          <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-br from-white/[0.03] to-white/[0.01] border border-white/10 hover:border-emerald-500/30 transition-all flex flex-col justify-between">
            <div>
              <div className="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center text-emerald-400 mb-6">
                <Server className="w-6 h-6" />
              </div>
              <span className="text-xs font-mono text-emerald-400 uppercase tracking-wider font-semibold">
                Market Opportunity
              </span>
              <h3 className="text-2xl font-bold text-white mt-1 mb-3">
                $85B+ High-Performance Computing Market
              </h3>
              <p className="text-slate-400 text-sm leading-relaxed font-light">
                From <strong>hyperscale datacenter packet routing</strong> to <strong>real-time AI reasoning trees</strong> (Monte Carlo Tree Search, Graph RAG), <strong>autonomous vehicle path planning</strong>, and <strong>ultra-low-latency financial trading</strong>, shortest path calculations represent one of computing’s most critical and compute-hungry bottlenecks.
              </p>
            </div>

            <div className="mt-8 pt-5 border-t border-white/5 flex items-center justify-between text-xs font-mono text-slate-400">
              <span>Target Sectors:</span>
              <span className="text-emerald-400 font-bold">AI Data Centers • Autonomous Systems • Finance</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
