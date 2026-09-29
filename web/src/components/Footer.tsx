import React from 'react';
import { Cpu, Github, FileText, ExternalLink } from 'lucide-react';

export const Footer: React.FC = () => {
  return (
    <footer className="bg-black border-t border-white/10 py-16 text-slate-400 text-xs">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-8 pb-10 border-b border-white/10">
          {/* Brand */}
          <div className="space-y-2">
            <div className="flex items-center space-x-2.5">
              <div className="w-8 h-8 rounded-full bg-cyan-500/20 border border-cyan-500/40 flex items-center justify-center">
                <Cpu className="w-4 h-4 text-cyan-400" />
              </div>
              <span className="text-base font-semibold text-white font-sans tracking-tight">
                Silica<span className="text-cyan-400 font-bold">Core</span>
              </span>
              <span className="px-2 py-0.5 rounded-full text-[10px] font-mono bg-white/5 border border-white/10 text-slate-300">
                1550 nm ToF
              </span>
            </div>
            <p className="text-slate-400 text-xs max-w-md font-light leading-relaxed">
              Open research on a photonic co-processor that solves shortest paths with time-of-flight race logic in Si₃N₄ / TFLN. Simulated, not yet fabricated.
            </p>
          </div>

          {/* Quick Links */}
          <div className="flex flex-wrap items-center gap-6 text-xs">
            <a
              href="https://github.com/Roliveira96/SilicaCore/blob/main/docs/papers/artigo-preliminar.md"
              target="_blank"
              rel="noreferrer"
              className="text-slate-300 hover:text-cyan-400 transition-colors flex items-center gap-1.5"
            >
              <FileText className="w-3.5 h-3.5 text-slate-500" />
              <span>Draft paper</span>
            </a>
            <a
              href="https://github.com/Roliveira96/SilicaCore"
              target="_blank"
              rel="noreferrer"
              className="text-slate-300 hover:text-cyan-400 transition-colors flex items-center gap-1.5"
            >
              <Github className="w-3.5 h-3.5 text-slate-500" />
              <span>GitHub Source</span>
            </a>
            <a
              href="https://github.com/Roliveira96/SilicaCore/tree/main/docs/architecture"
              target="_blank"
              rel="noreferrer"
              className="text-slate-300 hover:text-cyan-400 transition-colors flex items-center gap-1.5"
            >
              <ExternalLink className="w-3.5 h-3.5 text-slate-500" />
              <span>13 Architecture Monographs</span>
            </a>
          </div>
        </div>

        {/* Bottom Credits */}
        <div className="pt-8 flex flex-col sm:flex-row items-center justify-end text-[11px] text-slate-500 gap-4">
          <p className="font-mono text-[10px]">
            16×16 results from the open Go simulator; Dijkstra baseline in Go; waveguide physics from 2D Meep FDTD.
          </p>
        </div>
      </div>
    </footer>
  );
};
