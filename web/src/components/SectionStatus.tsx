import React from 'react';
import { CheckCircle2, CircleDashed, Compass, FlaskConical, Map, TriangleAlert } from 'lucide-react';

const VALIDATED = [
  '0 decoding errors in 25.5 M distances at 100 ps delay unit (16×16, 10 chips × 10⁴ queries)',
  'Measured error per hop count follows the Gaussian noise prediction',
  '42.2 ns per query: 11.5 ns light race + 30.7 ns readout; 648 mm² die',
  'Si₃N₄ + TFLN routing: 1.3 dB per gate, vs. 46.6 dB for mirrors in bulk glass',
  '2D FDTD (Meep): 0.012 dB per 90° bend at R = 30 µm; Si₃N₄ → TFLN taper < 0.003 dB from 25 µm',
  'Multi-chip 64×64: 0 errors with 150 ps unit, 20–31 ns origin → target',
  'Bottom-up energy model: ~9.2 W for the 16×16 chip, 0.39 µJ per query, ~5 °C above ambient',
];

const ASSUMED = [
  'Node regeneration latency (20 ps) and jitter (1.5 ps rms)',
  'Static delay error per edge after calibration (0.5 ps rms)',
  'Chip-to-chip coupling ≤ 1.5 dB per facet',
  'Sb₂Se₃ reprogramming time (1 µs) and 100 Gb/s readout link',
  'Energy model inputs: 5 mW per receiver, 4.1 mW per TDC, 20% laser efficiency',
];

const LIMITS = [
  'Not O(1): race time grows with the longest path, readout grows with the number of nodes.',
  'Area caps one chip at a 16×16 map at 100 ps; larger maps need several chips.',
  'Speed against current CPUs is estimated from one measured CPU and Geekbench scores; A* would narrow the gap.',
  'It is an accelerator for graph problems, not a CPU replacement.',
];

const NEXT = [
  { icon: FlaskConical, title: 'Fiber bench', text: 'ToF gate, one race-logic node and a 3×3 graph with commercial 1550 nm parts.' },
  { icon: Compass, title: '3D simulation', text: 'Full 3D FDTD or eigenmode expansion of bends and tapers.' },
  { icon: Map, title: 'Paper submission', text: 'Preliminary paper targeting WSCAD / SBESC, then an international venue.' },
];

export const SectionStatus: React.FC = () => (
  <section id="section-status" className="relative overflow-hidden border-t border-white/5 bg-[#020307] py-28">
    <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
      <div className="mx-auto mb-16 max-w-3xl text-center">
        <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-cyan-500/30 bg-cyan-950/60 px-3 py-1 font-mono text-xs text-cyan-300">
          <CheckCircle2 className="h-3.5 w-3.5" />
          <span>RESEARCH STATUS</span>
        </div>
        <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">What is proven. What is not yet.</h2>
        <p className="mt-4 text-base font-light leading-relaxed text-slate-400 sm:text-lg">
          SilicaCore is open research. Every number on this page is reproducible from the repository, and every assumption is
          labeled as one.
        </p>
      </div>

      <div className="mb-8 grid grid-cols-1 gap-8 lg:grid-cols-2">
        <List title="Validated by simulation" icon={<CheckCircle2 className="h-5 w-5 text-emerald-400" />} items={VALIDATED} />
        <List title="Assumptions still to be measured" icon={<CircleDashed className="h-5 w-5 text-amber-400" />} items={ASSUMED} />
      </div>

      <div className="mb-8">
        <List title="Known limits" icon={<TriangleAlert className="h-5 w-5 text-rose-400" />} items={LIMITS} />
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        {NEXT.map(({ icon: Icon, title, text }) => (
          <div key={title} className="rounded-3xl border border-white/10 bg-white/[0.02] p-7">
            <Icon className="mb-4 h-6 w-6 text-cyan-400" />
            <div className="text-lg font-bold text-white">{title}</div>
            <p className="mt-1 text-sm text-slate-400">{text}</p>
          </div>
        ))}
      </div>
    </div>
  </section>
);

const List: React.FC<{ title: string; icon: React.ReactNode; items: string[] }> = ({ title, icon, items }) => (
  <div className="h-full rounded-3xl border border-white/10 bg-white/[0.02] p-8">
    <div className="mb-5 flex items-center gap-2">
      {icon}
      <h3 className="text-lg font-bold text-white">{title}</h3>
    </div>
    <ul className="space-y-2.5">
      {items.map((it) => (
        <li key={it} className="flex gap-2.5 text-sm leading-relaxed text-slate-300">
          <span className="mt-2 h-1.5 w-1.5 flex-shrink-0 rounded-full bg-slate-500" />
          <span>{it}</span>
        </li>
      ))}
    </ul>
  </div>
);

export default SectionStatus;
