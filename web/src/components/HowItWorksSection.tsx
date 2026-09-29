import React, { useState } from 'react';
import { Radio, Zap, Clock, ShieldCheck, ArrowRight, Sparkles, Cpu, Layers } from 'lucide-react';

export const HowItWorksSection: React.FC = () => {
  const [activeStep, setActiveStep] = useState<number>(1);
  const [qubitState, setQubitState] = useState<'0' | '1' | 'superposition'>('superposition');

  return (
    <section id="section-how-it-works" className="relative py-28 bg-[#020307] border-t border-white/5 overflow-hidden">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-20">
          <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/60 border border-cyan-500/30 text-cyan-300 text-xs font-mono mb-4">
            <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
            <span>HOW LIGHT COMPUTES</span>
          </div>
          <h2 className="text-4xl sm:text-5xl font-black text-white tracking-tight">
            How the Processor Works.
          </h2>
          <p className="mt-4 text-base sm:text-lg text-slate-400 leading-relaxed font-light">
            Instead of shuttling billions of bits back and forth between registers and DRAM, SilicaCore lets physics solve the problem. Here is how a calculation happens in 3 simple steps:
          </p>
        </div>

        {/* 3 Visual Interactive Steps */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
          {/* Step 1 */}
          <div
            onClick={() => setActiveStep(1)}
            className={`p-7 rounded-3xl border transition-all cursor-pointer relative overflow-hidden flex flex-col justify-between ${
              activeStep === 1
                ? 'bg-gradient-to-b from-cyan-950/40 to-black border-cyan-500/60 shadow-[0_0_30px_rgba(6,182,212,0.2)]'
                : 'bg-white/[0.02] hover:bg-white/[0.04] border-white/10'
            }`}
          >
            <div>
              <div className="w-12 h-12 rounded-2xl bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400 font-mono font-bold text-lg mb-6">
                01
              </div>
              <h3 className="text-xl font-bold text-white mb-2">Laser Injection</h3>
              <p className="text-sm text-slate-300 leading-relaxed font-normal">
                A continuous 1550 nm laser injects a clean optical pulse directly into the chip facet through an ultra-low-loss inverse-taper edge coupler (≤ 1.5 dB).
              </p>
            </div>

            <div className="mt-8 pt-4 border-t border-white/10 flex items-center justify-between text-xs font-mono text-cyan-300">
              <span>Wavelength: 1550 nm</span>
              <span>Insertion: ≤ 1.5 dB</span>
            </div>
          </div>

          {/* Step 2 */}
          <div
            onClick={() => setActiveStep(2)}
            className={`p-7 rounded-3xl border transition-all cursor-pointer relative overflow-hidden flex flex-col justify-between ${
              activeStep === 2
                ? 'bg-gradient-to-b from-blue-950/40 to-black border-blue-500/60 shadow-[0_0_30px_rgba(59,130,246,0.2)]'
                : 'bg-white/[0.02] hover:bg-white/[0.04] border-white/10'
            }`}
          >
            <div>
              <div className="w-12 h-12 rounded-2xl bg-blue-500/10 border border-blue-500/30 flex items-center justify-center text-blue-400 font-mono font-bold text-lg mb-6">
                02
              </div>
              <h3 className="text-xl font-bold text-white mb-2">Wavefront Race</h3>
              <p className="text-sm text-slate-300 leading-relaxed font-normal">
                Light branches into all silicon nitride (Si₃N₄) waveguides simultaneously. Spiral lengths encode physical distances. The shortest route wins naturally.
              </p>
            </div>

            <div className="mt-8 pt-4 border-t border-white/10 flex items-center justify-between text-xs font-mono text-blue-300">
              <span>All Paths in Parallel</span>
              <span>0 Joule Heating</span>
            </div>
          </div>

          {/* Step 3 */}
          <div
            onClick={() => setActiveStep(3)}
            className={`p-7 rounded-3xl border transition-all cursor-pointer relative overflow-hidden flex flex-col justify-between ${
              activeStep === 3
                ? 'bg-gradient-to-b from-emerald-950/40 to-black border-emerald-500/60 shadow-[0_0_30px_rgba(16,185,129,0.2)]'
                : 'bg-white/[0.02] hover:bg-white/[0.04] border-white/10'
            }`}
          >
            <div>
              <div className="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center text-emerald-400 font-mono font-bold text-lg mb-6">
                03
              </div>
              <h3 className="text-xl font-bold text-white mb-2">Instant Detection</h3>
              <p className="text-sm text-slate-300 leading-relaxed font-normal">
                The first photon wavefront to reach the destination latches an ultra-fast Time-to-Digital Converter (TDC). The shortest path is locked in 42.2 nanoseconds.
              </p>
            </div>

            <div className="mt-8 pt-4 border-t border-white/10 flex items-center justify-between text-xs font-mono text-emerald-300">
              <span>Latency: 42.2 ns</span>
              <span>Zero Memory Stalls</span>
            </div>
          </div>
        </div>

        {/* Visual Interactive Photonic Qubit Demonstration */}
        <div className="p-8 sm:p-10 rounded-3xl bg-gradient-to-br from-white/[0.04] to-black border border-white/10 backdrop-blur-2xl">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-6 mb-8 pb-6 border-b border-white/10">
            <div>
              <div className="inline-flex items-center space-x-1.5 text-xs font-mono text-cyan-300 uppercase tracking-wider mb-1">
                <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
                <span>Quantum Photonic Superposition</span>
              </div>
              <h3 className="text-2xl font-bold text-white tracking-tight">
                Dual-Rail Photonic Qubit Trajectory
              </h3>
              <p className="text-sm text-slate-300 font-normal mt-1">
                Visualizing how a single photon travels through paths |0⟩, |1⟩, or both paths at once (1 and 0 simultaneously).
              </p>
            </div>

            {/* Qubit State Selectors */}
            <div className="flex items-center p-1 rounded-2xl bg-white/5 border border-white/10">
              <button
                onClick={() => setQubitState('0')}
                className={`px-4 py-2 rounded-xl text-xs font-mono font-medium transition-all ${
                  qubitState === '0'
                    ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/50 shadow-[0_0_12px_rgba(6,182,212,0.3)]'
                    : 'text-slate-300 hover:text-white'
                }`}
              >
                State |0⟩ (Rail 0)
              </button>
              <button
                onClick={() => setQubitState('1')}
                className={`px-4 py-2 rounded-xl text-xs font-mono font-medium transition-all ${
                  qubitState === '1'
                    ? 'bg-purple-500/20 text-purple-300 border border-purple-500/50 shadow-[0_0_12px_rgba(168,85,247,0.3)]'
                    : 'text-slate-300 hover:text-white'
                }`}
              >
                State |1⟩ (Rail 1)
              </button>
              <button
                onClick={() => setQubitState('superposition')}
                className={`px-4 py-2 rounded-xl text-xs font-mono font-medium transition-all ${
                  qubitState === 'superposition'
                    ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/50 shadow-[0_0_12px_rgba(16,185,129,0.3)]'
                    : 'text-slate-300 hover:text-white'
                }`}
              >
                Superposition |0⟩ + |1⟩ (1 &amp; 0)
              </button>
            </div>
          </div>

          {/* Visual Waveguide Diagram */}
          <div className="relative p-6 sm:p-8 rounded-2xl bg-black/70 border border-white/5 overflow-hidden">
            <div className="flex flex-col space-y-8 relative z-10">
              {/* Rail 0 */}
              <div className="flex items-center space-x-4">
                <span className="w-16 font-mono text-xs text-cyan-400 font-bold">Rail |0⟩</span>
                <div className="flex-1 h-3 rounded-full bg-slate-900 border border-slate-800 relative overflow-hidden">
                  {(qubitState === '0' || qubitState === 'superposition') && (
                    <div
                      className={`h-full rounded-full transition-all duration-500 ${
                        qubitState === '0'
                          ? 'w-full bg-cyan-400 shadow-[0_0_15px_rgba(6,182,212,0.8)] animate-pulse'
                          : 'w-full bg-gradient-to-r from-cyan-400 to-emerald-400 opacity-75 shadow-[0_0_15px_rgba(6,182,212,0.5)]'
                      }`}
                    />
                  )}
                </div>
                <span className="font-mono text-xs text-slate-400 w-24 text-right">
                  {qubitState === '0' ? '100% Light' : qubitState === 'superposition' ? '50% Light' : '0% (Dark)'}
                </span>
              </div>

              {/* Central Hadamard 50:50 Beam Splitter */}
              <div className="flex justify-center -my-3">
                <div className="px-4 py-1.5 rounded-full bg-white/5 border border-white/10 text-[11px] font-mono text-slate-300 flex items-center space-x-2">
                  <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
                  <span>50:50 Directional Coupler (Hadamard Gate H = 1/√2)</span>
                </div>
              </div>

              {/* Rail 1 */}
              <div className="flex items-center space-x-4">
                <span className="w-16 font-mono text-xs text-purple-400 font-bold">Rail |1⟩</span>
                <div className="flex-1 h-3 rounded-full bg-slate-900 border border-slate-800 relative overflow-hidden">
                  {(qubitState === '1' || qubitState === 'superposition') && (
                    <div
                      className={`h-full rounded-full transition-all duration-500 ${
                        qubitState === '1'
                          ? 'w-full bg-purple-400 shadow-[0_0_15px_rgba(168,85,247,0.8)] animate-pulse'
                          : 'w-full bg-gradient-to-r from-purple-400 to-emerald-400 opacity-75 shadow-[0_0_15px_rgba(168,85,247,0.5)]'
                      }`}
                    />
                  )}
                </div>
                <span className="font-mono text-xs text-slate-400 w-24 text-right">
                  {qubitState === '1' ? '100% Light' : qubitState === 'superposition' ? '50% Light' : '0% (Dark)'}
                </span>
              </div>
            </div>

            {/* Explanation box */}
            <div className="mt-8 pt-4 border-t border-white/5 flex flex-col sm:flex-row sm:items-center justify-between text-xs text-slate-400 gap-2">
              <div>
                <strong className="text-white">Active Quantum State:</strong>{' '}
                <span className="font-mono text-cyan-300">
                  {qubitState === '0' && '|ψ⟩ = |0⟩'}
                  {qubitState === '1' && '|ψ⟩ = |1⟩'}
                  {qubitState === 'superposition' && '|ψ⟩ = (1/√2)|0⟩ + (1/√2)|1⟩ (Simultaneous 1 and 0)'}
                </span>
              </div>
              <div className="text-slate-500 font-mono">Linear Optical Quantum Computing (LOQC)</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
