import React, { useState } from 'react';
import { BarChart3, BrainCircuit, Gauge, Thermometer, Trophy, Zap } from 'lucide-react';

// ============================================================================
// One comparison table, four tabs. Rows are always the same processors; each
// tab changes the metric. The best value is on top, rows without public data
// go to the bottom with the reason. Integers are written without thousands
// separators so they read the same in English and Portuguese.
// ============================================================================

type TabId = 'performance' | 'power' | 'temperature' | 'ai';
type Basis = 'simulated' | 'model' | 'estimated' | 'vendor' | 'measured' | 'n/a';

interface Cell {
  value: number | null; // used for sorting and the bar; null = no public data
  min?: number; // optional range start
  display: string;
  extra: string; // secondary column
  basis: Basis;
}

interface Processor {
  name: string;
  kind: string;
  silica?: boolean;
  cells: Record<TabId, Cell>;
}

// Shortest-path estimate: Dial's bucket-queue Dijkstra (the best algorithm for integer
// weights 1..15) measured at 25.38 us on an i3-3217U (Geekbench 6 single-core 307), scaled
// by each chip's single-core score. A 16x16 map runs on one core from cache.
const I3_US = 25.38;
const I3_GB6_SC = 307;
const estUs = (gb6sc: number) => (I3_US * I3_GB6_SC) / gb6sc;
const fmtUs = (us: number) => `${us.toFixed(1)} µs`;
// Energy per query estimate: package power x estimated time (upper bound).
const fmtEnergy = (uj: number) => (uj >= 1000 ? `${(uj / 1000).toFixed(1)} mJ` : `${Math.round(uj)} µJ`);

const na = (why: string, display = 'no public data'): Cell => ({ value: null, display, extra: why, basis: 'n/a' });
const GPU_SERIAL = 'one small shortest-path query is sequential; GPUs gain only on large batches';

const PROCESSORS: Processor[] = [
  {
    name: 'SilicaCore RL-16',
    kind: 'Photonic accelerator (simulated)',
    silica: true,
    cells: {
      performance: { value: 0.0422, display: '42.2 ns', extra: 'light race + readout, 16×16 map', basis: 'simulated' },
      power: { value: 9.2, display: '~9.2 W', extra: '0.39 µJ per query', basis: 'model' },
      temperature: { value: 5, display: '~5 °C above ambient', extra: '1.4 W/cm² with a small fan heatsink', basis: 'model' },
      ai: na('solves graphs, not neural networks', 'not an AI chip'),
    },
  },
  {
    name: 'Apple M5 Max',
    kind: 'Laptop SoC',
    cells: {
      performance: { value: estUs(4349), display: fmtUs(estUs(4349)), extra: 'Geekbench 6: 4349 single · 29363 multi', basis: 'estimated' },
      power: { value: 62, display: '~62 W', extra: `${fmtEnergy(62 * estUs(4349))} per query (est.)`, basis: 'measured' },
      temperature: na('Apple does not publish thermal limits', 'not published'),
      ai: na('Neural Engine TOPS not stated for M5 Max', 'not published'),
    },
  },
  {
    name: 'Intel Core Ultra 9 285K',
    kind: 'Desktop CPU',
    cells: {
      performance: { value: estUs(3450), display: fmtUs(estUs(3450)), extra: 'Geekbench 6: 3450 single · 23024 multi', basis: 'estimated' },
      power: { value: 250, min: 125, display: '125–250 W', extra: `${fmtEnergy(125 * estUs(3450))} per query (est.)`, basis: 'vendor' },
      temperature: { value: 105, display: 'up to 105 °C', extra: 'maximum junction temperature', basis: 'vendor' },
      ai: { value: 13, display: '13 TOPS', extra: 'on-chip NPU', basis: 'vendor' },
    },
  },
  {
    name: 'AMD Ryzen 9 9950X3D',
    kind: 'Desktop CPU',
    cells: {
      performance: { value: estUs(3440), display: fmtUs(estUs(3440)), extra: 'Geekbench 6: ~3440 single · ~22100 multi', basis: 'estimated' },
      power: { value: 230, min: 170, display: '170–230 W', extra: `${fmtEnergy(170 * estUs(3440))} per query (est.)`, basis: 'vendor' },
      temperature: { value: 95, display: 'up to 95 °C', extra: 'maximum junction temperature', basis: 'vendor' },
      ai: na('no NPU; AI runs on the CPU cores'),
    },
  },
  {
    name: 'AMD EPYC 9965',
    kind: 'Server CPU, 192 cores',
    cells: {
      performance: na('no public Geekbench 6 single-core score'),
      power: { value: 500, display: '500 W', extra: 'TDP', basis: 'vendor' },
      temperature: na('not in the cited sources'),
      ai: na('no dedicated AI TOPS figure'),
    },
  },
  {
    name: 'Intel Xeon 6980P',
    kind: 'Server CPU, 128 cores',
    cells: {
      performance: { value: estUs(2131), display: fmtUs(estUs(2131)), extra: 'Geekbench 6: ~2131 single', basis: 'estimated' },
      power: { value: 500, display: '500 W', extra: `${fmtEnergy(500 * estUs(2131))} per query (est.)`, basis: 'vendor' },
      temperature: { value: 80, display: 'up to 80 °C', extra: 'maximum package temperature', basis: 'vendor' },
      ai: na('AMX matrix units; no TOPS headline'),
    },
  },
  {
    name: 'NVIDIA GeForce RTX 5090',
    kind: 'Gaming / AI GPU',
    cells: {
      performance: na(GPU_SERIAL, 'no single-query gain'),
      power: { value: 575, display: '575 W', extra: 'total graphics power', basis: 'vendor' },
      temperature: { value: 90, display: 'up to 90 °C', extra: 'maximum GPU temperature', basis: 'vendor' },
      ai: { value: 3352, display: '3352 TOPS', extra: 'FP4, sparse', basis: 'vendor' },
    },
  },
  {
    name: 'NVIDIA B300',
    kind: 'Datacenter AI GPU',
    cells: {
      performance: na(GPU_SERIAL, 'no single-query gain'),
      power: { value: 1400, display: '1400 W', extra: 'Blackwell Ultra', basis: 'vendor' },
      temperature: na('not in the cited sources'),
      ai: { value: 15000, display: '15 PFLOPS', extra: 'FP4, dense', basis: 'vendor' },
    },
  },
  {
    name: 'NVIDIA Rubin',
    kind: 'Datacenter AI GPU (2026)',
    cells: {
      performance: na(GPU_SERIAL, 'no single-query gain'),
      power: { value: 2300, min: 1800, display: '1800–2300 W', extra: 'Max-Q to Max-P profiles', basis: 'vendor' },
      temperature: na('not in the cited sources'),
      ai: { value: 50000, display: '50 PFLOPS', extra: 'NVFP4 inference', basis: 'vendor' },
    },
  },
];

const TABS: {
  id: TabId;
  label: string;
  icon: typeof Gauge;
  metric: string;
  extraHeader: string;
  better: 'lower' | 'higher';
  scale: { min: number; max: number };
  headline: { value: string; text: string };
}[] = [
  {
    id: 'performance',
    label: 'Performance',
    icon: Gauge,
    metric: 'Time per shortest-path query',
    extraHeader: 'Details',
    better: 'lower',
    scale: { min: 0.01, max: 100 },
    headline: { value: '~42×', text: 'faster than the fastest CPU in the table (Apple M5 Max, estimated, best algorithm) on one 16×16 shortest-path query when the map changes. On a fixed map a precomputed table answers faster.' },
  },
  {
    id: 'power',
    label: 'Power',
    icon: Zap,
    metric: 'Power under load',
    extraHeader: 'Energy per shortest-path query',
    better: 'lower',
    scale: { min: 1, max: 10000 },
    headline: { value: '~290×', text: 'less energy per shortest-path query than the most efficient CPU in the table (Apple M5 Max, estimated), with the whole chip drawing ~9 W.' },
  },
  {
    id: 'temperature',
    label: 'Temperature',
    icon: Thermometer,
    metric: 'Temperature',
    extraHeader: 'Details',
    better: 'lower',
    scale: { min: 1, max: 120 },
    headline: { value: '~5 °C', text: 'above ambient for the whole chip: ~9 W spread over 6.5 cm² of die. The other chips are rated to run at 80–105 °C.' },
  },
  {
    id: 'ai',
    label: 'AI Usage',
    icon: BrainCircuit,
    metric: 'Peak AI throughput',
    extraHeader: 'Format',
    better: 'higher',
    scale: { min: 1, max: 100000 },
    headline: {
      value: 'Not its job',
      text: 'The race-logic chip accelerates graph search, not neural networks. The AI leaders are shown for context; photonic AI chips such as Lightmatter (Nature, 2025) already exist.',
    },
  },
];

const BASIS_STYLE: Record<Basis, string> = {
  simulated: 'border-cyan-500/40 text-cyan-300',
  model: 'border-cyan-500/40 text-cyan-300',
  estimated: 'border-amber-500/40 text-amber-300',
  measured: 'border-slate-400/40 text-slate-300',
  vendor: 'border-slate-500/40 text-slate-400',
  'n/a': 'border-slate-700 text-slate-500',
};

const logPos = (v: number, min: number, max: number) =>
  (100 * (Math.log10(Math.min(Math.max(v, min), max)) - Math.log10(min))) / (Math.log10(max) - Math.log10(min));

export const SectionComparison: React.FC = () => {
  const [tabId, setTabId] = useState<TabId>('performance');
  const tab = TABS.find((t) => t.id === tabId)!;

  const withData = PROCESSORS.filter((p) => p.cells[tabId].value !== null).sort((a, b) => {
    const va = a.cells[tabId].value as number;
    const vb = b.cells[tabId].value as number;
    return tab.better === 'lower' ? va - vb : vb - va;
  });
  const without = PROCESSORS.filter((p) => p.cells[tabId].value === null);

  return (
    <section id="section-compare" className="relative overflow-hidden border-t border-white/5 bg-[#03050b] py-28">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto mb-10 max-w-3xl text-center">
          <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-cyan-500/30 bg-cyan-950/60 px-3 py-1 font-mono text-xs text-cyan-300">
            <BarChart3 className="h-3.5 w-3.5" />
            <span>COMPARISONS · EVERY NUMBER SOURCED</span>
          </div>
          <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">How it stacks up.</h2>
          <p className="mt-4 text-base font-light leading-relaxed text-slate-400 sm:text-lg">
            SilicaCore next to the top processors of 2026. Pick a tab; the best value is always on top.
          </p>
        </div>

        <div className="overflow-hidden rounded-3xl border border-white/10 bg-white/[0.02]">
          <div role="tablist" className="flex flex-wrap gap-1 border-b border-white/10 bg-black/30 p-2">
            {TABS.map(({ id, label, icon: Icon }) => (
              <button
                key={id}
                role="tab"
                aria-selected={tabId === id}
                onClick={() => setTabId(id)}
                className={`flex flex-1 items-center justify-center gap-2 rounded-2xl px-4 py-3 text-sm font-bold transition-all ${
                  tabId === id ? 'bg-gradient-to-r from-cyan-400 to-blue-500 text-slate-950' : 'text-slate-300 hover:bg-white/5 hover:text-white'
                }`}
              >
                <Icon className="h-4 w-4" />
                <span>{label}</span>
              </button>
            ))}
          </div>

          <div className="flex flex-col gap-2 border-b border-white/10 px-6 py-6 sm:flex-row sm:items-center sm:gap-6 sm:px-8">
            <div className="whitespace-nowrap font-mono text-4xl font-black text-cyan-300">{tab.headline.value}</div>
            <p className="text-sm leading-relaxed text-slate-300">{tab.headline.text}</p>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full min-w-[820px] text-left text-sm">
              <thead className="font-mono text-[11px] uppercase tracking-wider text-slate-500">
                <tr className="border-b border-white/10">
                  <th className="w-10 px-6 py-3 font-medium">#</th>
                  <th className="px-3 py-3 font-medium">Processor</th>
                  <th className="px-3 py-3 font-medium">
                    {tab.metric} <span className="normal-case tracking-normal">({tab.better} is better)</span>
                  </th>
                  <th className="px-3 py-3 font-medium">{tab.extraHeader}</th>
                  <th className="px-6 py-3 font-medium">Basis</th>
                </tr>
              </thead>
              <tbody>
                {withData.map((p, i) => {
                  const c = p.cells[tabId];
                  const end = Math.max(2, logPos(c.value as number, tab.scale.min, tab.scale.max));
                  const start = c.min !== undefined ? logPos(c.min, tab.scale.min, tab.scale.max) : end;
                  return (
                    <tr key={p.name} className={`border-b border-white/5 ${p.silica ? 'bg-cyan-950/30' : ''}`}>
                      <td className="px-6 py-3.5 font-mono text-slate-500">
                        {i === 0 ? <Trophy className="h-4 w-4 text-amber-300" aria-label="best" /> : i + 1}
                      </td>
                      <td className="px-3 py-3.5">
                        <div className={`font-semibold ${p.silica ? 'text-cyan-300' : 'text-white'}`}>{p.name}</div>
                        <div className="text-xs text-slate-500">{p.kind}</div>
                      </td>
                      <td className="px-3 py-3.5">
                        <div className="flex items-center gap-3">
                          <div className="relative h-2.5 w-40 flex-shrink-0 rounded-full bg-white/5">
                            <div
                              className={`absolute left-0 top-0 h-full rounded-full ${p.silica ? 'bg-cyan-400' : 'bg-slate-300'}`}
                              style={{ width: `${start}%` }}
                            />
                            {c.min !== undefined && (
                              <div
                                className="absolute top-0 h-full rounded-r-full bg-slate-300/40"
                                style={{ left: `${start}%`, width: `${Math.max(1, end - start)}%` }}
                              />
                            )}
                          </div>
                          <span className="whitespace-nowrap font-mono font-semibold text-white">{c.display}</span>
                        </div>
                      </td>
                      <td className="px-3 py-3.5 text-xs text-slate-400">{c.extra}</td>
                      <td className="px-6 py-3.5">
                        <span className={`rounded-full border px-2 py-0.5 font-mono text-[10px] ${BASIS_STYLE[c.basis]}`}>{c.basis}</span>
                      </td>
                    </tr>
                  );
                })}
                {without.map((p) => {
                  const c = p.cells[tabId];
                  return (
                    <tr key={p.name} className={`border-b border-white/5 text-slate-500 ${p.silica ? 'bg-cyan-950/20' : ''}`}>
                      <td className="px-6 py-3 font-mono">–</td>
                      <td className="px-3 py-3">
                        <div className={`font-semibold ${p.silica ? 'text-cyan-300/80' : 'text-slate-400'}`}>{p.name}</div>
                        <div className="text-xs">{p.kind}</div>
                      </td>
                      <td className="px-3 py-3 font-mono text-xs">{c.display}</td>
                      <td className="px-3 py-3 text-xs">{c.extra}</td>
                      <td className="px-6 py-3">
                        <span className={`rounded-full border px-2 py-0.5 font-mono text-[10px] ${BASIS_STYLE[c.basis]}`}>{c.basis}</span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>

          <p className="px-6 py-5 text-xs leading-relaxed text-slate-500 sm:px-8">
            <strong className="text-slate-400">Method.</strong> SilicaCore values come from the project&apos;s simulator and
            energy model; no chip has been fabricated yet. CPU shortest-path times are estimated from the fastest algorithm for
            these integer weights (Dijkstra with a bucket queue, Dial) measured on an Intel i3-3217U (25.4 µs, Geekbench 6
            single-core 307), scaled by each chip&apos;s single-core score; energy per query is package power × that time. On
            a map that never changes, a precomputed table of all distances returns a query in ~0.1 µs on the same i3 (a few ns
            on current CPUs), so the photonic advantage applies when the map changes between queries. A dedicated synchronous CMOS race-logic circuit at 3 GHz would reach similar latency (~40 ns) with far less area and energy; the table compares against general-purpose processors. Other figures are vendor or press specifications:{' '}
            <Src href="https://www.macrumors.com/2026/03/05/m5-max-geekbench-benchmarks/">M5 Max</Src>,{' '}
            <Src href="https://www.notebookcheck.net/Apple-M5-Pro-M5-Max-CPU-Analysis-M5-Max-is-not-much-faster-than-the-M4-Max.1246054.0.html">M5 Max power</Src>,{' '}
            <Src href="https://www.intel.com/content/www/us/en/products/sku/241060/intel-core-ultra-9-processor-285k-36m-cache-up-to-5-70-ghz/specifications.html">285K</Src>,{' '}
            <Src href="https://en.wikipedia.org/wiki/AMD_Ryzen_9_9950X3D">9950X3D</Src>,{' '}
            <Src href="https://buy.hpe.com/us/en/Options/Processors/Third-Party-Processors/Third-Party-Processor-Options/AMD-EPYC-9965-2-25GHz-192%E2%80%91core-500W-Processor-for-HPE/p/P75019-B21">EPYC 9965</Src>,{' '}
            <Src href="https://videocardz.com/press-release/intel-launches-xeon-6-granite-rapids-6980p-with-128-cores-and-500w-tdp">Xeon 6980P</Src>,{' '}
            <Src href="https://www.nvidia.com/en-us/geforce/graphics-cards/50-series/rtx-5090/">RTX 5090</Src>,{' '}
            <Src href="https://www.tomshardware.com/pc-components/gpus/nvidia-announces-blackwell-ultra-b300-1-5x-faster-than-b200-with-288gb-hbm3e-and-15-pflops-dense-fp4">B300</Src>,{' '}
            <Src href="https://developer.nvidia.com/blog/inside-the-nvidia-rubin-platform-six-new-chips-one-ai-supercomputer/">Rubin</Src>,{' '}
            <Src href="https://cputronic.com/index.php/cpu/intel-core-i3-3217u">i3-3217U</Src>,{' '}
            <Src href="https://technical.city/en/cpu/Xeon-6980P">Xeon 6980P temperature</Src>,{' '}
            <Src href="https://gamersnexus.net/gpus/nvidia-geforce-rtx-5090-founders-edition-review-benchmarks-gaming-thermals-power">RTX 5090 temperature</Src>.
          </p>
        </div>
      </div>
    </section>
  );
};

const Src: React.FC<{ href: string; children: React.ReactNode }> = ({ href, children }) => (
  <a href={href} target="_blank" rel="noreferrer" className="text-cyan-400 hover:underline">
    {children}
  </a>
);

export default SectionComparison;
