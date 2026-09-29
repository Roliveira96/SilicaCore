import React, { useState } from 'react';
import { BarChart3, Cpu, Gauge, Play, Sparkles } from 'lucide-react';
import { measureDijkstraUs } from '../lib/racelogic';

// Log-scale axis for the speed chart, in nanoseconds.
const AXIS_MIN_NS = 10;
const AXIS_MAX_NS = 1_000_000;
const logPos = (ns: number) =>
  (100 * (Math.log10(ns) - Math.log10(AXIS_MIN_NS))) / (Math.log10(AXIS_MAX_NS) - Math.log10(AXIS_MIN_NS));

interface SpeedRow {
  label: string;
  detail: string;
  minNs: number;
  maxNs: number;
  kind: 'simulated' | 'measured' | 'estimated' | 'live';
}

const KIND_STYLE: Record<SpeedRow['kind'], { bar: string; tag: string; text: string }> = {
  simulated: { bar: 'bg-cyan-400', tag: 'border-cyan-500/40 text-cyan-300', text: 'simulated' },
  measured: { bar: 'bg-slate-300', tag: 'border-slate-500/40 text-slate-300', text: 'measured' },
  estimated: { bar: 'bg-slate-500/70', tag: 'border-slate-600/40 text-slate-400', text: 'estimated' },
  live: { bar: 'bg-amber-400', tag: 'border-amber-500/40 text-amber-300', text: 'measured now' },
};

const fmt = (ns: number) => (ns >= 1000 ? `${(ns / 1000).toFixed(ns >= 10_000 ? 0 : 1)} µs` : `${ns.toFixed(1)} ns`);

const ENERGY = [
  { label: 'Moving a bit across a CMOS die', value: '~1 pJ/bit', note: 'on-chip copper interconnect (literature)' },
  { label: 'Fetching a bit from DRAM / HBM', value: '10–100 pJ/bit', note: 'off-chip memory access (literature)' },
  { label: 'SilicaCore optical transport', value: '~1.9 pJ/bit', note: 'laser share + modulator + detector (model)' },
];

interface Flagship {
  chip: string;
  maker: string;
  kind: string;
  compute: string;
  memory: string;
  power: string;
  source: string;
  url: string;
}

const FLAGSHIPS: Flagship[] = [
  { chip: 'M5 Max', maker: 'Apple', kind: 'Laptop SoC (2026)', compute: '18-core CPU (6 super + 12 perf.), 40-core GPU', memory: 'up to 128 GB, 614 GB/s', power: 'not published', source: 'Apple', url: 'https://www.apple.com/newsroom/2026/03/apple-debuts-m5-pro-and-m5-max-to-supercharge-the-most-demanding-pro-workflows/' },
  { chip: 'Rubin', maker: 'NVIDIA', kind: 'AI GPU (shipping H2 2026)', compute: '50 PFLOPS NVFP4 inference, TSMC 3 nm', memory: '288 GB HBM4, 22 TB/s', power: 'not in cited source', source: 'NVIDIA', url: 'https://developer.nvidia.com/blog/inside-the-nvidia-rubin-platform-six-new-chips-one-ai-supercomputer/' },
  { chip: 'B300 (Blackwell Ultra)', maker: 'NVIDIA', kind: 'AI GPU', compute: '15 PFLOPS dense FP4', memory: '288 GB HBM3e, 8 TB/s', power: '1,400 W', source: "Tom's Hardware", url: 'https://www.tomshardware.com/pc-components/gpus/nvidia-announces-blackwell-ultra-b300-1-5x-faster-than-b200-with-288gb-hbm3e-and-15-pflops-dense-fp4' },
  { chip: 'EPYC Venice (Zen 6)', maker: 'AMD', kind: 'Server CPU (Q3 2026)', compute: 'up to 256 cores, TSMC 2 nm', memory: '—', power: 'not in cited source', source: 'Hardware Busters', url: 'https://hwbusters.com/news/amd-launches-zen-6-with-epyc-venice-256-cores-on-tsmc-2nm/' },
  { chip: 'EPYC 9965', maker: 'AMD', kind: 'Server CPU', compute: '192 cores, up to 3.7 GHz', memory: '576 GB/s per socket', power: '500 W', source: 'HPE', url: 'https://buy.hpe.com/us/en/Options/Processors/Third-Party-Processors/Third-Party-Processor-Options/AMD-EPYC-9965-2-25GHz-192%E2%80%91core-500W-Processor-for-HPE/p/P75019-B21' },
  { chip: 'Xeon 6+ (Clearwater Forest)', maker: 'Intel', kind: 'Server CPU (2026)', compute: 'up to 288 E-cores, Intel 18A', memory: '12-ch DDR5-8000', power: '330–450 W', source: 'TechPowerUp', url: 'https://www.techpowerup.com/346941/intel-launches-xeon-6-clearwater-forest-xeon-with-288-e-cores-on-18a-process' },
  { chip: 'Xeon 6980P', maker: 'Intel', kind: 'Server CPU', compute: '128 P-cores, up to 3.9 GHz', memory: '12-ch DDR5-6400 / MRDIMM-8800', power: '500 W', source: 'VideoCardz', url: 'https://videocardz.com/press-release/intel-launches-xeon-6-granite-rapids-6980p-with-128-cores-and-500w-tdp' },
];

const PHOTONIC = [
  { system: 'Lightmatter', source: 'Nature, 2025', task: 'Neural networks (ResNet, BERT)', result: '65.5 TOPS at 78 W + 1.6 W optical', status: 'Fabricated, measured' },
  { system: 'PACE (Lightelligence)', source: 'Nature, 2025', task: 'Ising optimisation, 64×64 matrix', result: '~5 ns latency, >16,000 components', status: 'Fabricated, measured' },
  { system: 'Taichi', source: 'Science, 2024', task: 'AI chiplet', result: '160 TOPS/W inside the chiplet', status: 'Fabricated, measured' },
  { system: 'SilicaCore', source: 'this project', task: 'Shortest path (race logic), 16×16 map', result: '42.2 ns per query, 0 errors in 25.5 M', status: 'Simulated' },
];

export const SectionComparison: React.FC = () => {
  const [liveUs, setLiveUs] = useState<number | null>(null);
  const [measuring, setMeasuring] = useState(false);

  const rows: SpeedRow[] = [
    { label: 'SilicaCore', detail: '16×16 race + TDC readout (Go simulator)', minNs: 42.2, maxNs: 42.2, kind: 'simulated' },
    { label: 'Current desktop CPU', detail: 'Dijkstra, from the paper’s 150–800× estimate', minNs: 42.2 * 150, maxNs: 42.2 * 800, kind: 'estimated' },
    { label: 'Intel Core i3-3217U (2012)', detail: 'Dijkstra in Go, single thread, 20,000 queries', minNs: 69_570, maxNs: 69_570, kind: 'measured' },
  ];
  if (liveUs !== null) {
    rows.push({ label: 'Your computer', detail: 'Dijkstra in JavaScript, in this browser', minNs: liveUs * 1000, maxNs: liveUs * 1000, kind: 'live' });
  }

  const measure = () => {
    setMeasuring(true);
    setTimeout(() => {
      setLiveUs(measureDijkstraUs(16, 400));
      setMeasuring(false);
    }, 30);
  };

  return (
    <section id="section-compare" className="relative overflow-hidden border-t border-white/5 bg-[#03050b] py-28">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto mb-16 max-w-3xl text-center">
          <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-cyan-500/30 bg-cyan-950/60 px-3 py-1 font-mono text-xs text-cyan-300">
            <BarChart3 className="h-3.5 w-3.5" />
            <span>COMPARISONS · EVERY NUMBER SOURCED</span>
          </div>
          <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">How it stacks up.</h2>
          <p className="mt-4 text-base font-light leading-relaxed text-slate-400 sm:text-lg">
            Same problem, different hardware: all distances from one origin on a 16×16 map. Each value says whether it was
            simulated, measured or estimated.
          </p>
        </div>

        {/* Speed */}
        <div className="mb-10 rounded-3xl border border-white/10 bg-white/[0.02] p-6 sm:p-8">
          <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <Gauge className="h-5 w-5 text-cyan-400" />
              <h3 className="text-lg font-bold text-white">Time per 16×16 shortest-path query</h3>
              <span className="text-xs text-slate-500">(log scale, lower is better)</span>
            </div>
            <button
              onClick={measure}
              disabled={measuring}
              className="flex items-center gap-2 rounded-full border border-amber-500/40 bg-amber-500/10 px-4 py-2 text-xs font-semibold text-amber-200 hover:bg-amber-500/20 disabled:opacity-50"
            >
              <Play className="h-3.5 w-3.5" />
              <span>{measuring ? 'Measuring…' : liveUs === null ? 'Measure my computer' : 'Measure again'}</span>
            </button>
          </div>

          <div className="space-y-5">
            {rows.map((r) => {
              const style = KIND_STYLE[r.kind];
              const left = logPos(r.minNs);
              const width = Math.max(1.2, logPos(r.maxNs) - left);
              const value = r.minNs === r.maxNs ? fmt(r.minNs) : `${fmt(r.minNs)} – ${fmt(r.maxNs)}`;
              return (
                <div key={r.label} className="grid grid-cols-1 gap-2 md:grid-cols-[260px_1fr] md:items-center md:gap-6">
                  <div>
                    <div className="flex items-center gap-2">
                      <span className={`text-sm font-semibold ${r.kind === 'simulated' ? 'text-white' : 'text-slate-200'}`}>{r.label}</span>
                      <span className={`rounded-full border px-1.5 py-px font-mono text-[10px] ${style.tag}`}>{style.text}</span>
                    </div>
                    <div className="text-xs text-slate-500">{r.detail}</div>
                  </div>
                  <div className="relative h-7 rounded-md bg-white/[0.03]">
                    <div
                      className={`absolute top-1 h-5 rounded ${style.bar}`}
                      style={{ left: `${left}%`, width: `${width}%` }}
                      title={value}
                    />
                    <span
                      className="absolute top-1/2 -translate-y-1/2 whitespace-nowrap pl-2 font-mono text-xs text-slate-200"
                      style={{ left: `${Math.min(88, left + width)}%` }}
                    >
                      {value}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>

          <div className="mt-4 hidden grid-cols-[260px_1fr] gap-6 md:grid">
            <div />
            <div className="flex justify-between font-mono text-[10px] text-slate-500">
              {['10 ns', '100 ns', '1 µs', '10 µs', '100 µs', '1 ms'].map((t) => (
                <span key={t}>{t}</span>
              ))}
            </div>
          </div>

          <p className="mt-6 text-xs leading-relaxed text-slate-500">
            SilicaCore numbers come from simulation; no chip exists yet. The i3-3217U figure is measured with{' '}
            <code className="rounded bg-white/10 px-1 font-mono text-slate-300">cmd/dijkstrabench</code>; the current-CPU band is an estimate, not a benchmark. The browser runs JavaScript, which is slower than
            native code, so your measurement is an upper bound for your machine.
          </p>
        </div>

        {/* Flagship silicon */}
        <div className="mb-10 rounded-3xl border border-white/10 bg-white/[0.02] p-6 sm:p-8">
          <div className="mb-2 flex items-center gap-2">
            <Cpu className="h-5 w-5 text-cyan-400" />
            <h3 className="text-lg font-bold text-white">The strongest silicon of 2026, for scale</h3>
          </div>
          <p className="mb-6 text-xs leading-relaxed text-slate-500">
            Published specifications from each vendor or trade press. These chips are general-purpose giants; SilicaCore is a
            narrow accelerator, so this is a sense of scale, not a benchmark. None of them has a published 16×16 shortest-path
            time yet: run <code className="rounded bg-white/10 px-1 font-mono text-slate-300">go run ./cmd/dijkstrabench</code>{' '}
            on one to add a real number to the speed chart.
          </p>
          <div className="overflow-x-auto">
            <table className="w-full min-w-[760px] border-separate border-spacing-y-1.5 text-left text-xs">
              <thead className="font-mono text-[10px] uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-3 py-1 font-medium">Chip</th>
                  <th className="px-3 py-1 font-medium">Type</th>
                  <th className="px-3 py-1 font-medium">Compute</th>
                  <th className="px-3 py-1 font-medium">Memory</th>
                  <th className="px-3 py-1 font-medium">Power</th>
                  <th className="px-3 py-1 font-medium">Source</th>
                </tr>
              </thead>
              <tbody>
                {FLAGSHIPS.map((f) => (
                  <tr key={f.chip} className="bg-black/40 text-slate-300">
                    <td className="rounded-l-xl px-3 py-2.5">
                      <div className="font-semibold text-white">{f.chip}</div>
                      <div className="text-[10px] text-slate-500">{f.maker}</div>
                    </td>
                    <td className="px-3 py-2.5">{f.kind}</td>
                    <td className="px-3 py-2.5">{f.compute}</td>
                    <td className="px-3 py-2.5">{f.memory}</td>
                    <td className="px-3 py-2.5 font-mono">{f.power}</td>
                    <td className="rounded-r-xl px-3 py-2.5">
                      <a href={f.url} target="_blank" rel="noreferrer" className="text-cyan-400 hover:underline">
                        {f.source}
                      </a>
                    </td>
                  </tr>
                ))}
                <tr className="bg-cyan-950/40 text-slate-200">
                  <td className="rounded-l-xl px-3 py-2.5">
                    <div className="font-semibold text-cyan-300">SilicaCore RL-16</div>
                    <div className="text-[10px] text-slate-500">this project</div>
                  </td>
                  <td className="px-3 py-2.5">Photonic accelerator (simulated)</td>
                  <td className="px-3 py-2.5">256 race nodes, 960 delay lines; 42.2 ns per 16×16 query</td>
                  <td className="px-3 py-2.5">graph weights stored as optical delays</td>
                  <td className="px-3 py-2.5 font-mono">not yet derived</td>
                  <td className="rounded-r-xl px-3 py-2.5">
                    <a
                      href="https://github.com/Roliveira96/SilicaCore/blob/main/docs/papers/artigo-preliminar.md"
                      target="_blank"
                      rel="noreferrer"
                      className="text-cyan-400 hover:underline"
                    >
                      Draft paper
                    </a>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
          {/* Energy */}
          <div className="rounded-3xl border border-white/10 bg-white/[0.02] p-8">
            <div className="mb-5 flex items-center gap-2">
              <Cpu className="h-5 w-5 text-cyan-400" />
              <h3 className="text-lg font-bold text-white">Energy per bit moved</h3>
            </div>
            <div className="space-y-3">
              {ENERGY.map((e) => (
                <div key={e.label} className="flex items-center justify-between gap-4 rounded-2xl border border-white/10 bg-black/40 p-4">
                  <div>
                    <div className="text-sm text-slate-200">{e.label}</div>
                    <div className="text-xs text-slate-500">{e.note}</div>
                  </div>
                  <div className="whitespace-nowrap font-mono text-sm font-bold text-white">{e.value}</div>
                </div>
              ))}
            </div>
            <p className="mt-4 text-xs leading-relaxed text-slate-500">
              Per bit, light is in the same range as on-chip copper. The gain for graph search is elsewhere: the edge weights
              live in the chip as optical delays, so a query never fetches them from DRAM.
            </p>
          </div>

          {/* Photonic landscape */}
          <div className="rounded-3xl border border-white/10 bg-white/[0.02] p-8">
            <div className="mb-5 flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-cyan-400" />
              <h3 className="text-lg font-bold text-white">Other photonic processors</h3>
            </div>
            <div className="space-y-3">
              {PHOTONIC.map((p) => (
                <div
                  key={p.system}
                  className={`rounded-2xl border p-4 ${p.system === 'SilicaCore' ? 'border-cyan-500/40 bg-cyan-950/30' : 'border-white/10 bg-black/40'}`}
                >
                  <div className="flex flex-wrap items-baseline justify-between gap-2">
                    <span className="text-sm font-semibold text-white">{p.system}</span>
                    <span className="font-mono text-[10px] text-slate-500">
                      {p.source} · {p.status}
                    </span>
                  </div>
                  <div className="mt-1 text-xs text-slate-400">{p.task}</div>
                  <div className="mt-1 font-mono text-xs text-slate-200">{p.result}</div>
                </div>
              ))}
            </div>
            <p className="mt-4 text-xs leading-relaxed text-slate-500">
              These systems solve different problems, so the table shows the landscape rather than a head-to-head score.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default SectionComparison;
