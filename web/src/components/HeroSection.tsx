import React from 'react';
import { Play, FileText, ArrowRight, Zap, Clock, ShieldCheck, Flame, Sparkles } from 'lucide-react';

interface HeroSectionProps {
  onLaunchSim: () => void;
}

export const HeroSection: React.FC<HeroSectionProps> = ({ onLaunchSim }) => {
  return (
    <section className="relative min-h-[92vh] flex flex-col justify-center items-center pt-32 pb-24 px-4 sm:px-6 lg:px-8 overflow-hidden bg-[#040814] text-center">
      {/* Background Silicon Photonics Microchip Photograph */}
      <div className="absolute inset-0 z-0 overflow-hidden pointer-events-none">
        <img
          src="/hero-bg.jpg"
          alt="SilicaCore Silicon Photonics Integrated Chip with 1550nm Laser Injection"
          className="w-full h-full object-cover object-center opacity-70"
        />
        {/* Cinematic dark gradients to balance image visibility and text readability */}
        <div className="absolute inset-0 bg-gradient-to-b from-black/60 via-black/40 to-[#040814]" />
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,_transparent_30%,_rgba(0,0,0,0.75)_85%)]" />
      </div>

      {/* Subtle ambient photon glow circle */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[700px] h-[380px] bg-gradient-to-tr from-cyan-500/20 via-blue-600/15 to-transparent blur-[160px] rounded-full pointer-events-none z-0" />

      <div className="max-w-5xl mx-auto flex flex-col items-center relative z-10">
        {/* Apple-style intro tag */}
        <div className="inline-flex items-center space-x-2 px-3.5 py-1.5 rounded-full bg-white/10 border border-white/15 text-cyan-300 text-xs font-mono mb-8 backdrop-blur-xl shadow-lg">
          <span className="w-1.5 h-1.5 rounded-full bg-cyan-400 animate-ping"></span>
          <span>SILICACORE • INTEGRATED PHOTONIC RACE LOGIC</span>
        </div>

        {/* Giant Apple-style Headline */}
        <h1 className="text-5xl sm:text-7xl md:text-8xl font-black tracking-tight text-white max-w-4xl leading-[1.05] drop-shadow-md">
          Light moves.{' '}
          <span className="bg-gradient-to-r from-cyan-400 via-teal-300 to-blue-400 bg-clip-text text-transparent">
            Electrons crawl.
          </span>
        </h1>

        {/* Crisp Sub-headline */}
        <p className="mt-6 text-lg sm:text-xl md:text-2xl text-slate-200 max-w-3xl font-light leading-relaxed drop-shadow">
          The planar photonic co-processor that resolves complex graph pathfinding in{' '}
          <strong className="text-white font-semibold">42.2 nanoseconds</strong>. No memory bottlenecks. No clock cycles. Just pure speed-of-light transit.
        </p>

        {/* Action Buttons */}
        <div className="mt-10 flex flex-wrap items-center justify-center gap-4">
          <button
            onClick={onLaunchSim}
            className="group px-7 py-3.5 rounded-full font-bold text-slate-950 bg-gradient-to-r from-cyan-400 via-teal-300 to-blue-400 hover:from-cyan-300 hover:to-blue-300 shadow-[0_0_30px_rgba(6,182,212,0.5)] transition-all flex items-center space-x-2 text-sm"
          >
            <Play className="w-4 h-4 fill-slate-950 transition-transform group-hover:scale-110" />
            <span>Test Interactive Chip</span>
            <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
          </button>

          <a
            href="https://github.com/Roliveira96/SilicaCore/blob/main/docs/papers/artigo-preliminar.md"
            target="_blank"
            rel="noreferrer"
            className="px-6 py-3.5 rounded-full font-medium text-white bg-white/10 hover:bg-white/15 border border-white/15 hover:border-white/25 transition-all text-sm flex items-center space-x-2 backdrop-blur-xl shadow-md"
          >
            <FileText className="w-4 h-4 text-cyan-400" />
            <span>Read Pre-print Paper</span>
          </a>
        </div>

        {/* 4 Apple-style Key Metric Bento Cards */}
        <div className="mt-16 w-full max-w-4xl grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="p-5 rounded-3xl bg-black/60 border border-white/15 backdrop-blur-2xl text-left hover:border-cyan-500/40 transition-all shadow-xl">
            <div className="flex items-center space-x-1.5 text-cyan-400 text-xs font-mono uppercase mb-2">
              <Clock className="w-3.5 h-3.5" />
              <span>Full-Graph Latency</span>
            </div>
            <div className="text-3xl font-extrabold text-white font-mono tracking-tight">
              42.2 <span className="text-xs font-sans text-cyan-400 font-normal">ns</span>
            </div>
            <p className="text-xs text-slate-300 mt-1">Light wavefront + TDC latch</p>
          </div>

          <div className="p-5 rounded-3xl bg-black/60 border border-white/15 backdrop-blur-2xl text-left hover:border-emerald-500/40 transition-all shadow-xl">
            <div className="flex items-center space-x-1.5 text-emerald-400 text-xs font-mono uppercase mb-2">
              <Zap className="w-3.5 h-3.5" />
              <span>Speedup Factor</span>
            </div>
            <div className="text-3xl font-extrabold text-white font-mono tracking-tight">398× - 535×</div>
            <p className="text-xs text-slate-300 mt-1">Faster than Apple M5 Max &amp; x86</p>
          </div>

          <div className="p-5 rounded-3xl bg-black/60 border border-white/15 backdrop-blur-2xl text-left hover:border-blue-500/40 transition-all shadow-xl">
            <div className="flex items-center space-x-1.5 text-blue-400 text-xs font-mono uppercase mb-2">
              <Zap className="w-3.5 h-3.5" />
              <span>Total System TDP</span>
            </div>
            <div className="text-3xl font-extrabold text-white font-mono tracking-tight">
              18.5 <span className="text-xs font-sans text-blue-400 font-normal">W</span>
            </div>
            <p className="text-xs text-slate-300 mt-1">76% lower power than Apple M5 Max</p>
          </div>

          <div className="p-5 rounded-3xl bg-black/60 border border-white/15 backdrop-blur-2xl text-left hover:border-amber-500/40 transition-all shadow-xl">
            <div className="flex items-center space-x-1.5 text-amber-400 text-xs font-mono uppercase mb-2">
              <Flame className="w-3.5 h-3.5" />
              <span>Operating Temp</span>
            </div>
            <div className="text-3xl font-extrabold text-white font-mono tracking-tight">
              28 <span className="text-xs font-sans text-amber-400 font-normal">°C</span>
            </div>
            <p className="text-xs text-slate-300 mt-1">Ice-cold passive conduction</p>
          </div>
        </div>
      </div>
    </section>
  );
};
