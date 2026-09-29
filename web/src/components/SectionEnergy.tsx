import React from 'react';
import { BatteryCharging, Flame, Gauge, Info, Zap } from 'lucide-react';

const SINKS = [
  { name: 'Lasers', note: 'mode-locked source or comb feeding every channel; wall-plug efficiency matters most' },
  { name: 'TFLN modulators and drivers', note: '960 re-fire modulators on the 16×16 chip, driven electrically' },
  { name: 'Photodiodes, comparators and TDCs', note: 'one detector per incoming edge, one time-to-digital converter per node' },
  { name: 'Readout link', note: '12-bit times for 256 nodes over a 100 Gb/s link: 30.7 ns per query' },
  { name: 'Sb₂Se₃ heaters', note: 'only when the map changes; the switches hold their state with no power' },
];

export const SectionEnergy: React.FC = () => (
  <section id="section-energy" className="relative overflow-hidden border-t border-white/5 bg-[#030408] py-28">
    <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
      <div className="mx-auto mb-16 max-w-3xl text-center">
        <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-amber-500/30 bg-amber-950/40 px-3 py-1 font-mono text-xs text-amber-300">
          <Gauge className="h-3.5 w-3.5" />
          <span>ENERGY · HONEST ACCOUNTING</span>
        </div>
        <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">Light is not free.</h2>
        <p className="mt-4 text-base font-light leading-relaxed text-slate-400 sm:text-lg">
          Waveguides carry light without resistive heating, but a working chip still spends energy producing, switching,
          detecting and reading out that light. The model adds up every block, with values from published devices.
        </p>
      </div>

      <div className="mb-12 grid grid-cols-1 gap-6 md:grid-cols-3">
        <Card
          icon={<Zap className="h-5 w-5" />}
          value="~7.6 pJ/bit"
          title="Optical transport, from the model"
          text="Laser power shared over the binary aggregate bit rate (64 channels × ~5 GHz ≈ 0.33 Tb/s), plus modulator and detector energy. Same order as off-package electrical links."
          tone="cyan"
        />
        <Card
          icon={<BatteryCharging className="h-5 w-5" />}
          value="0 W"
          title="To hold the map"
          text="Sb₂Se₃ phase-change switches keep the programmed delays without power. Energy is spent only when the map changes."
          tone="emerald"
        />
        <Card
          icon={<Flame className="h-5 w-5" />}
          value="~9.2 W"
          title="Whole chip, bottom-up model"
          text="Lasers 3.2 W, 960 always-on receivers 4.8 W, 256 TDCs 1.05 W, readout 0.15 W at 23.7 M queries/s: 0.39 µJ per query."
          tone="amber"
        />
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
        <div className="rounded-3xl border border-white/10 bg-white/[0.02] p-8">
          <h3 className="mb-5 text-lg font-bold text-white">Where the energy goes</h3>
          <ul className="space-y-3">
            {SINKS.map((s) => (
              <li key={s.name} className="rounded-2xl border border-white/10 bg-black/40 p-4">
                <div className="text-sm font-semibold text-white">{s.name}</div>
                <div className="mt-0.5 text-xs text-slate-400">{s.note}</div>
              </li>
            ))}
          </ul>
        </div>

        <div className="flex flex-col justify-between rounded-3xl border border-white/10 bg-white/[0.02] p-8">
          <div>
            <h3 className="mb-3 text-lg font-bold text-white">A reference point, not a SilicaCore number</h3>
            <p className="text-sm leading-relaxed text-slate-300">
              Photonic cores reach over 100 TOPS/W inside the optics. The best published complete photonic AI system
              (Lightmatter, <em>Nature</em> 2025) delivers 65.5 TOPS for 78 W of electronics and 1.6 W of optical power:
              about <strong className="text-white">0.84 TOPS/W</strong> at system level. The periphery, not the light,
              sets the energy bill.
            </p>
          </div>
          <div className="mt-6 flex items-start gap-2 rounded-2xl border border-white/10 bg-black/40 p-4 text-xs text-slate-400">
            <Info className="mt-0.5 h-4 w-4 flex-shrink-0 text-slate-300" />
            <span>
              SilicaCore has not been fabricated. Every energy figure here comes from the open simulator&apos;s bottom-up model
              and the cited literature; the always-on receivers, not the light, dominate the budget.
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
);

const TONES = {
  cyan: 'text-cyan-300 border-cyan-500/30 bg-cyan-500/10',
  emerald: 'text-emerald-300 border-emerald-500/30 bg-emerald-500/10',
  amber: 'text-amber-300 border-amber-500/30 bg-amber-500/10',
} as const;

const Card: React.FC<{ icon: React.ReactNode; value: string; title: string; text: string; tone: keyof typeof TONES }> = ({
  icon,
  value,
  title,
  text,
  tone,
}) => (
  <div className="rounded-3xl border border-white/10 bg-gradient-to-br from-white/[0.04] to-transparent p-7">
    <div className={`mb-5 flex h-11 w-11 items-center justify-center rounded-2xl border ${TONES[tone]}`}>{icon}</div>
    <div className="font-mono text-3xl font-extrabold tracking-tight text-white">{value}</div>
    <div className="mt-1 text-sm font-semibold text-slate-200">{title}</div>
    <p className="mt-2 text-xs leading-relaxed text-slate-400">{text}</p>
  </div>
);

export default SectionEnergy;
