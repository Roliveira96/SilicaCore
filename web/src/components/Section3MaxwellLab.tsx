import React, { useState, useEffect, useRef } from 'react';
import { Activity, Zap, CheckCircle2, ChevronRight, Sliders, Waves, Layers } from 'lucide-react';

interface BendData {
  radiusUm: number;
  lossDb: number;
  transmission: number;
  leakagePct: number;
  description: string;
}

const BEND_SWEEPS: BendData[] = [
  {
    radiusUm: 10,
    lossDb: 0.2134,
    transmission: 0.952,
    leakagePct: 4.8,
    description: 'Tight bend: the mode shifts outwards and radiates into the cladding; mismatch at the straight-bend junctions adds loss.',
  },
  {
    radiusUm: 20,
    lossDb: 0.035,
    transmission: 0.992,
    leakagePct: 0.8,
    description: 'Moderate bend: radiation tails drastically attenuated; practical for loose interconnects.',
  },
  {
    radiusUm: 30,
    lossDb: 0.0121,
    transmission: 0.9972,
    leakagePct: 0.28,
    description: 'Very low loss: tight modal overlap with negligible scattering.',
  },
  {
    radiusUm: 50,
    lossDb: 0.003,
    transmission: 0.9993,
    leakagePct: 0.07,
    description: 'Near-unity transmission (0.003 dB per 90°), comfortable for compact delay spirals.',
  },
];

const TAPER_SWEEPS = [
  { lengthUm: 5, lossDb: 0.526, transmission: 88.6, status: 'Too short: not adiabatic' },
  { lengthUm: 10, lossDb: 0.058, transmission: 98.7, status: 'Residual mode mismatch' },
  { lengthUm: 25, lossDb: 0.0024, transmission: 99.94, status: 'Adiabatic (< 0.003 dB)' },
  { lengthUm: 50, lossDb: 0.0018, transmission: 99.96, status: 'Saturated' },
];

export const Section3MaxwellLab: React.FC = () => {
  const [selectedBendRadius, setSelectedBendRadius] = useState<number>(50);
  const [selectedTaperLen, setSelectedTaperLen] = useState<number>(25);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const currentBend = BEND_SWEEPS.find((b) => b.radiusUm === selectedBendRadius) || BEND_SWEEPS[3];

  // Animated field canvas (illustrative mode profile)
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let animId: number;
    let t = 0;

    const render = () => {
      t += 0.04;
      const width = canvas.width;
      const height = canvas.height;

      // Dark background
      ctx.fillStyle = '#060913';
      ctx.fillRect(0, 0, width, height);

      // Draw bending waveguide path
      const centerX = width * 0.35;
      const centerY = height * 0.75;
      const r = (selectedBendRadius / 50) * 120 + 70;
      const coreWidth = 26;

      // Draw Cladding field radiation leakage
      const leakage = currentBend.leakagePct / 3.13; // 0 to 1
      if (leakage > 0.1) {
        const leakGrad = ctx.createRadialGradient(centerX, centerY, r + coreWidth / 2, centerX, centerY, r + coreWidth * 3);
        leakGrad.addColorStop(0, `rgba(244, 63, 94, ${0.45 * leakage})`);
        leakGrad.addColorStop(1, 'rgba(244, 63, 94, 0)');
        ctx.fillStyle = leakGrad;
        ctx.beginPath();
        ctx.arc(centerX, centerY, r + coreWidth * 3, -Math.PI / 2, 0);
        ctx.lineTo(centerX + r + coreWidth * 3, centerY);
        ctx.stroke();
        ctx.fill();
      }

      // Draw Waveguide Core
      ctx.strokeStyle = '#1e293b';
      ctx.lineWidth = coreWidth;
      ctx.beginPath();
      ctx.arc(centerX, centerY, r, -Math.PI / 2, 0);
      ctx.stroke();

      // Electromagnetic wave oscillating along the bend
      const segments = 90;
      for (let i = 0; i < segments; i++) {
        const angle = -Math.PI / 2 + (i / segments) * (Math.PI / 2);
        const wave = Math.sin(i * 0.6 - t * 4);
        const amp = (1 - (i / segments) * (currentBend.lossDb / 0.15) * 0.35);

        // Core intensity
        const wx = centerX + Math.cos(angle) * r;
        const wy = centerY + Math.sin(angle) * r;

        ctx.fillStyle = wave > 0 ? `rgba(0, 242, 254, ${Math.abs(wave) * amp * 0.9})` : `rgba(79, 172, 254, ${Math.abs(wave) * amp * 0.5})`;
        ctx.beginPath();
        ctx.arc(wx, wy, 8 * amp, 0, Math.PI * 2);
        ctx.fill();

        // Radiative tail outside the bend if low radius
        if (selectedBendRadius <= 20 && wave > 0.3) {
          const leakX = wx + Math.cos(angle) * (18 * leakage);
          const leakY = wy + Math.sin(angle) * (18 * leakage);
          ctx.fillStyle = `rgba(244, 63, 94, ${wave * 0.6 * leakage})`;
          ctx.beginPath();
          ctx.arc(leakX, leakY, 4, 0, Math.PI * 2);
          ctx.fill();
        }
      }

      animId = requestAnimationFrame(render);
    };

    render();
    return () => cancelAnimationFrame(animId);
  }, [selectedBendRadius, currentBend]);

  return (
    <section id="section-maxwell" className="relative py-24 bg-[#050811] border-t border-slate-800/80">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="flex flex-col md:flex-row md:items-end justify-between mb-12">
          <div>
            <div className="flex items-center space-x-2 text-cyan-400 font-mono text-xs uppercase tracking-wider mb-2">
              <Activity className="w-4 h-4" />
              <span>Section 3 • Maxwell Equations &amp; FDTD Meep Validation</span>
            </div>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
              The Maxwell Lab: Sub-dB Waveguide Bends &amp; Tapers
            </h2>
            <p className="mt-3 text-base text-slate-300 max-w-2xl">
              Waveguide bend and taper losses are checked with finite-difference time-domain (FDTD) simulations in Meep 1.34 using the 2D Effective Index Method (EIM). 3D full-wave validation is slated for future work.
            </p>
          </div>

          <div className="mt-4 md:mt-0 flex items-center space-x-2 text-xs font-mono text-cyan-400 bg-cyan-950/60 px-3 py-1.5 rounded-lg border border-cyan-500/30">
            <span>Meep 2D EIM • 20–70 px/µm</span>
          </div>
        </div>

        {/* Experiment 1: Bending Loss Inspector */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start mb-16">
          {/* Canvas Wave Field Display */}
          <div className="lg:col-span-7 flex flex-col space-y-3">
            <div className="relative rounded-2xl bg-gradient-to-b from-[#090f1d] to-[#04060b] border border-cyan-500/30 p-4 shadow-2xl overflow-hidden">
              <div className="flex items-center justify-between mb-3 text-xs font-mono">
                <span className="text-cyan-300 flex items-center gap-1.5 font-semibold">
                  <Waves className="w-4 h-4 text-cyan-400" />
                  Electric field mode profile (1550 nm, quasi-TE₀, illustrative)
                </span>
                <span className="text-slate-400">Si₃N₄ (800 nm × 0.7 µm)</span>
              </div>

              <div className="relative flex justify-center items-center bg-[#030509] rounded-xl border border-slate-800/80 overflow-hidden">
                <canvas
                  ref={canvasRef}
                  width={520}
                  height={320}
                  className="w-full h-auto max-h-[340px]"
                />

                {/* Floating telemetry */}
                <div className="absolute bottom-3 left-3 glass-panel px-3 py-2 rounded-lg border border-slate-700 text-xs font-mono">
                  <div className="text-slate-400 text-[10px]">Transmission per 90°:</div>
                  <div className="text-emerald-400 font-bold text-sm">
                    {(currentBend.transmission * 100).toFixed(2)}%
                  </div>
                </div>

                <div className="absolute bottom-3 right-3 glass-panel px-3 py-2 rounded-lg border border-slate-700 text-xs font-mono text-right">
                  <div className="text-slate-400 text-[10px]">Bend Loss:</div>
                  <div
                    className={`font-bold text-sm ${
                      selectedBendRadius >= 30 ? 'text-emerald-400' : 'text-rose-400'
                    }`}
                  >
                    {currentBend.lossDb.toFixed(4)} dB
                  </div>
                </div>
              </div>

              {/* Radius Selector Pills */}
              <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
                <span className="text-xs font-mono text-slate-300">Bend Radius (R):</span>
                <div className="flex items-center space-x-2">
                  {BEND_SWEEPS.map((bend) => (
                    <button
                      key={bend.radiusUm}
                      onClick={() => setSelectedBendRadius(bend.radiusUm)}
                      className={`px-3 py-1.5 rounded-lg text-xs font-mono font-medium transition-all ${
                        selectedBendRadius === bend.radiusUm
                          ? 'bg-cyan-500 text-slate-950 font-bold shadow-[0_0_12px_rgba(6,182,212,0.4)]'
                          : 'bg-slate-900 text-slate-400 hover:text-white border border-slate-800'
                      }`}
                    >
                      {bend.radiusUm} µm
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Theoretical vs Measured Inspector Details */}
          <div className="lg:col-span-5 flex flex-col space-y-4">
            <div className="glass-panel p-5 rounded-xl border border-slate-800">
              <h3 className="text-sm font-mono uppercase tracking-wider text-slate-400 mb-3 flex items-center justify-between">
                <span>Bending Physics Telemetry</span>
                <span className="text-cyan-400 font-bold">R = {currentBend.radiusUm} µm</span>
              </h3>

              <p className="text-xs text-slate-300 leading-relaxed mb-4">
                {currentBend.description}
              </p>

              <div className="space-y-3 font-mono text-xs">
                <div className="p-3 rounded-lg bg-slate-950/80 border border-slate-800 flex justify-between items-center">
                  <span className="text-slate-400">Core Dimensions:</span>
                  <span className="text-white font-medium">800 nm thick × 0.7 µm width</span>
                </div>
                <div className="p-3 rounded-lg bg-slate-950/80 border border-slate-800 flex justify-between items-center">
                  <span className="text-slate-400">Effective Index (n_eff):</span>
                  <span className="text-cyan-300 font-medium">1.7137 (slab: 1.8836)</span>
                </div>
                <div className="p-3 rounded-lg bg-slate-950/80 border border-slate-800 flex justify-between items-center">
                  <span className="text-slate-400">Power lost from the fundamental mode:</span>
                  <span
                    className={`font-semibold ${
                      currentBend.leakagePct > 1 ? 'text-rose-400' : 'text-emerald-400'
                    }`}
                  >
                    {currentBend.leakagePct.toFixed(2)}% of guided power
                  </span>
                </div>
              </div>

              <div className="mt-4 p-3 rounded-lg bg-slate-900/60 border border-slate-800 text-[11px] text-slate-400">
                <strong className="text-slate-300">Design consequence:</strong> from R ≈ 30 µm up, each 90° bend costs about
                0.01 dB or less, the value assumed by the link model (10 dB power margin per edge). Delay spirals are therefore
                drawn with R ≥ 30 µm. These are 2D effective-index results; 3D simulation and measurement are still pending.
              </div>
            </div>
          </div>
        </div>

        {/* Experiment 2: Transition Taper Visualizer */}
        <div className="glass-panel p-6 rounded-2xl border border-slate-800">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6 pb-4 border-b border-slate-800">
            <div>
              <div className="flex items-center space-x-2 text-indigo-400 font-mono text-xs uppercase tracking-wider mb-1">
                <Layers className="w-3.5 h-3.5" />
                <span>Adiabatic Taper Coupling (Si₃N₄ → TFLN)</span>
              </div>
              <h3 className="text-xl font-bold text-white">
                Sub-0.002 dB Adiabatic Transition Saturation
              </h3>
            </div>
            <div className="mt-2 sm:mt-0 text-xs font-mono text-slate-400">
              Optimal length: <strong className="text-cyan-400">25 µm</strong>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-5 gap-3">
            {TAPER_SWEEPS.map((taper) => {
              const isSelected = selectedTaperLen === taper.lengthUm;
              return (
                <div
                  key={taper.lengthUm}
                  onClick={() => setSelectedTaperLen(taper.lengthUm)}
                  className={`p-4 rounded-xl cursor-pointer transition-all border ${
                    isSelected
                      ? 'bg-cyan-950/40 border-cyan-400 shadow-[0_0_15px_rgba(6,182,212,0.2)] ring-1 ring-cyan-400/30'
                      : 'bg-slate-900/50 hover:bg-slate-900 border-slate-800'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-xs font-mono font-bold text-white">
                      L = {taper.lengthUm} µm
                    </span>
                    {taper.lengthUm === 25 && (
                      <span className="px-1.5 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-300 border border-emerald-500/30">
                        Sweet Spot
                      </span>
                    )}
                  </div>

                  <div className="text-lg font-bold font-mono text-white mb-1">
                    {taper.lossDb.toFixed(4)} <span className="text-xs text-slate-400 font-normal">dB</span>
                  </div>

                  <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mb-2">
                    <div
                      className="bg-cyan-400 h-full rounded-full"
                      style={{ width: `${taper.transmission}%` }}
                    />
                  </div>

                  <div className="text-[11px] text-slate-400 font-mono">
                    {taper.transmission.toFixed(2)}% Transmission
                  </div>

                  <div className="text-[10px] text-slate-500 mt-2 line-clamp-1">
                    {taper.status}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </section>
  );
};
