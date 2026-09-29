import React, { useCallback, useMemo, useRef, useState } from 'react';
import { Atom, Info, Play, RotateCcw, Timer, Zap } from 'lucide-react';
import DecisionMechanism3D, {
  HALF_WINDOW_PS,
  HIST_BIN_PS,
  HIST_BINS,
  T0_PS,
  T1_PS,
  type QuantumState,
} from './DecisionMechanism3D';

const PLATFORM_SIGMA_PS = 2.125; // Si3N4 + InGaAs photodiode + 5 ps TDC (cmd/tofplatform)
const LEGACY_SIGMA_PS = 11.24; // original silica concept with a 25 ps SPAD
const BATCH_SIZE = 500;

// Standard normal sample (Box-Muller).
function gaussian(): number {
  const u = 1 - Math.random();
  const v = Math.random();
  return Math.sqrt(-2 * Math.log(u)) * Math.cos(2 * Math.PI * v);
}

// Complementary error function (Abramowitz-Stegun 7.1.26, |error| < 1.5e-7).
function erfc(x: number): number {
  const t = 1 / (1 + 0.3275911 * Math.abs(x));
  const y =
    t * (0.254829592 + t * (-0.284496736 + t * (1.421413741 + t * (-1.453152027 + t * 1.061405429)))) * Math.exp(-x * x);
  return x >= 0 ? y : 2 - y;
}

function decode(arrivalPS: number): 0 | 1 | -1 {
  if (Math.abs(arrivalPS - T1_PS) <= HALF_WINDOW_PS) return 1;
  if (Math.abs(arrivalPS - T0_PS) <= HALF_WINDOW_PS) return 0;
  return -1;
}

function formatRate(r: number): string {
  if (r === 0) return '0';
  if (r >= 0.01) return `${(r * 100).toFixed(1)}%`;
  return r.toExponential(1);
}

export const SectionDecision: React.FC = () => {
  const [mode, setMode] = useState<'tof' | 'quantum'>('tof');

  // ToF state
  const [bit, setBit] = useState<0 | 1>(1);
  const [sigma, setSigma] = useState(PLATFORM_SIGMA_PS);
  const [histogram, setHistogram] = useState<number[]>(() => new Array(HIST_BINS).fill(0));
  const [lastArrival, setLastArrival] = useState<number | null>(null);
  const [pulses, setPulses] = useState(0);
  const [errors, setErrors] = useState(0);
  const pendingArrival = useRef<number | null>(null);

  // Quantum state
  const [qState, setQState] = useState<QuantumState>('super');
  const [qOutcome, setQOutcome] = useState<0 | 1 | null>(null);
  const [qCounts, setQCounts] = useState<[number, number]>([0, 0]);
  const pendingOutcome = useRef<0 | 1 | null>(null);

  const [shotToken, setShotToken] = useState(0);
  const [busy, setBusy] = useState(false);

  const q = 100 / (2 * sigma);
  // A pulse is wrong or rejected when it leaves its own window: both Gaussian tails.
  const theoryError = erfc(HALF_WINDOW_PS / (sigma * Math.SQRT2));

  const sampleArrival = useCallback((b: 0 | 1) => (b === 1 ? T1_PS : T0_PS) + sigma * gaussian(), [sigma]);

  const record = useCallback(
    (arrivals: number[], b: 0 | 1) => {
      setHistogram((h) => {
        const next = h.slice();
        for (const a of arrivals) {
          const i = Math.floor(a / HIST_BIN_PS);
          if (i >= 0 && i < HIST_BINS) next[i] += 1;
        }
        return next;
      });
      setPulses((p) => p + arrivals.length);
      setErrors((e) => e + arrivals.filter((a) => decode(a) !== b).length);
    },
    [],
  );

  const resetTof = () => {
    setHistogram(new Array(HIST_BINS).fill(0));
    setLastArrival(null);
    setPulses(0);
    setErrors(0);
  };

  const fireOne = () => {
    if (busy) return;
    if (mode === 'tof') {
      pendingArrival.current = sampleArrival(bit);
      setLastArrival(null);
    } else {
      pendingOutcome.current = qState === 'zero' ? 0 : qState === 'one' ? 1 : Math.random() < 0.5 ? 0 : 1;
      setQOutcome(null);
    }
    setBusy(true);
    setShotToken((t) => t + 1);
  };

  const fireBatch = () => {
    if (mode === 'tof') {
      const arrivals = Array.from({ length: BATCH_SIZE }, () => sampleArrival(bit));
      record(arrivals, bit);
      setLastArrival(arrivals[arrivals.length - 1]);
    } else {
      let zeros = 0;
      for (let i = 0; i < BATCH_SIZE; i++) {
        const o = qState === 'zero' ? 0 : qState === 'one' ? 1 : Math.random() < 0.5 ? 0 : 1;
        if (o === 0) zeros++;
      }
      setQCounts(([c0, c1]) => [c0 + zeros, c1 + BATCH_SIZE - zeros]);
    }
  };

  const onArrive = useCallback(() => {
    setBusy(false);
    if (mode === 'tof' && pendingArrival.current !== null) {
      const a = pendingArrival.current;
      pendingArrival.current = null;
      setLastArrival(a);
      record([a], bit);
    } else if (mode === 'quantum' && pendingOutcome.current !== null) {
      const o = pendingOutcome.current;
      pendingOutcome.current = null;
      setQOutcome(o);
      setQCounts(([c0, c1]) => (o === 0 ? [c0 + 1, c1] : [c0, c1 + 1]));
    }
  }, [mode, bit, record]);

  const switchMode = (m: 'tof' | 'quantum') => {
    setMode(m);
    setBusy(false);
    pendingArrival.current = null;
    pendingOutcome.current = null;
  };

  const qTotal = qCounts[0] + qCounts[1];
  const measuredRate = pulses > 0 ? errors / pulses : null;
  const histogramMemo = useMemo(() => histogram, [histogram]);

  return (
    <section id="section-decision" className="relative overflow-hidden border-t border-white/5 bg-[#020308] py-28">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto mb-12 max-w-3xl text-center">
          <div className="mb-4 inline-flex items-center space-x-2 rounded-full border border-cyan-500/30 bg-cyan-950/40 px-3 py-1 font-mono text-xs text-cyan-300">
            <Timer className="h-3.5 w-3.5" />
            <span>HOW A BIT IS DECIDED</span>
          </div>
          <h2 className="text-4xl font-black tracking-tight text-white sm:text-5xl">Time is the bit.</h2>
          <p className="mt-4 text-base font-light leading-relaxed text-slate-400 sm:text-lg">
            The chip never measures how bright the light is. It measures <em>when</em> the pulse arrives. A short path
            means 1, a path twice as long means 0, and a detector with a time-to-digital converter reads which window the
            pulse fell into. The qubit tab shows the long-term concept, where a single photon is in both rails until it
            is detected.
          </p>
        </div>

        <div className="mb-6 flex justify-center">
          <div className="inline-flex rounded-full border border-white/10 bg-white/[0.03] p-1">
            {(
              [
                { id: 'tof', label: 'Classical bit (time of flight)', icon: <Zap className="h-4 w-4" /> },
                { id: 'quantum', label: 'Qubit (0 and 1)', icon: <Atom className="h-4 w-4" /> },
              ] as const
            ).map((t) => (
              <button
                key={t.id}
                onClick={() => switchMode(t.id)}
                className={`flex items-center gap-2 rounded-full px-4 py-2 text-sm font-semibold transition ${
                  mode === t.id ? 'bg-white text-black' : 'text-slate-300 hover:text-white'
                }`}
              >
                {t.icon}
                {t.label}
              </button>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
          <div className="relative h-[460px] overflow-hidden rounded-3xl border border-white/10 bg-black sm:h-[560px]">
              <DecisionMechanism3D
                key={mode}
                mode={mode}
                bit={bit}
                shotToken={shotToken}
                lastArrivalPS={lastArrival}
                histogram={histogramMemo}
                quantumState={qState}
                quantumOutcome={qOutcome}
                onArrive={onArrive}
              />
            <div className="pointer-events-none absolute bottom-3 left-4 font-mono text-[11px] text-slate-500">
              drag to orbit · scroll to zoom · schematic, not to scale
            </div>
          </div>

          <div className="flex flex-col gap-4">
            {mode === 'tof' ? (
              <>
                <Panel title="Bit to send">
                  <div className="grid grid-cols-2 gap-2">
                    {([1, 0] as const).map((b) => (
                      <button
                        key={b}
                        onClick={() => setBit(b)}
                        className={`rounded-xl border px-3 py-2 text-sm font-semibold ${
                          bit === b
                            ? b === 1
                              ? 'border-cyan-400/60 bg-cyan-500/15 text-cyan-200'
                              : 'border-violet-400/60 bg-violet-500/15 text-violet-200'
                            : 'border-white/10 text-slate-400 hover:text-white'
                        }`}
                      >
                        bit {b} · {b === 1 ? 'short path' : 'long path'}
                      </button>
                    ))}
                  </div>
                </Panel>

                <Panel title={`Timing jitter σ = ${sigma.toFixed(1)} ps`}>
                  <input
                    type="range"
                    min={1}
                    max={40}
                    step={0.1}
                    value={sigma}
                    onChange={(e) => setSigma(Number(e.target.value))}
                    className="w-full accent-cyan-400"
                  />
                  <div className="mt-1 flex justify-between font-mono text-[10px] text-slate-500">
                    <button onClick={() => setSigma(PLATFORM_SIGMA_PS)} className="text-cyan-400 hover:underline">
                      chip model 2.12 ps
                    </button>
                    <button onClick={() => setSigma(LEGACY_SIGMA_PS)} className="text-amber-400 hover:underline">
                      original SPAD 11.24 ps
                    </button>
                  </div>
                </Panel>

                <Controls busy={busy} onOne={fireOne} onBatch={fireBatch} onReset={resetTof} />

                <Panel title="Decision statistics">
                  <Stat label="Pulses" value={String(pulses)} />
                  <Stat label="Wrong or rejected" value={String(errors)} />
                  <Stat label="Measured error rate" value={measuredRate === null ? '–' : formatRate(measuredRate)} />
                  <Stat label={`Theory (Q = ${q.toFixed(1)})`} value={theoryError < 1e-9 ? '< 1e-9' : formatRate(theoryError)} />
                  <p className="mt-3 text-[11px] leading-relaxed text-slate-500">
                    Windows are ±{HALF_WINDOW_PS} ps around {T1_PS} ps (bit 1) and {T0_PS} ps (bit 0); a pulse outside its
                    window counts as an error, so the theory is erfc(W/2σ√2), with Q = Δt/2σ and Δt = 100 ps. With the chip&apos;s
                    photodiode (σ ≈ 2.1 ps) errors are negligible. The SPAD of the original concept (σ = 11.24 ps) errs about
                    1 in 40000 pulses; raise σ further to watch the two peaks merge.
                  </p>
                </Panel>
              </>
            ) : (
              <>
                <Panel title="Photon state">
                  <div className="grid grid-cols-3 gap-2">
                    {(
                      [
                        { id: 'zero', label: '|0⟩' },
                        { id: 'one', label: '|1⟩' },
                        { id: 'super', label: '|0⟩+|1⟩' },
                      ] as const
                    ).map((s) => (
                      <button
                        key={s.id}
                        onClick={() => {
                          setQState(s.id);
                          setQOutcome(null);
                          setQCounts([0, 0]);
                        }}
                        className={`rounded-xl border px-2 py-2 font-mono text-sm ${
                          qState === s.id ? 'border-emerald-400/60 bg-emerald-500/15 text-emerald-200' : 'border-white/10 text-slate-400 hover:text-white'
                        }`}
                      >
                        {s.label}
                      </button>
                    ))}
                  </div>
                </Panel>

                <Controls
                  busy={busy}
                  onOne={fireOne}
                  onBatch={fireBatch}
                  onReset={() => {
                    setQCounts([0, 0]);
                    setQOutcome(null);
                  }}
                  oneLabel="Send 1 photon"
                />

                <Panel title="Detector clicks">
                  {([0, 1] as const).map((k) => {
                    const share = qTotal > 0 ? qCounts[k] / qTotal : 0;
                    return (
                      <div key={k} className="mb-3">
                        <div className="mb-1 flex justify-between font-mono text-xs text-slate-300">
                          <span>detector {k}</span>
                          <span>
                            {qCounts[k]} · {(share * 100).toFixed(1)}%
                          </span>
                        </div>
                        <div className="h-2 overflow-hidden rounded-full bg-white/5">
                          <div
                            className={`h-full rounded-full ${k === 0 ? 'bg-cyan-400' : 'bg-violet-400'}`}
                            style={{ width: `${share * 100}%` }}
                          />
                        </div>
                      </div>
                    );
                  })}
                  <p className="mt-2 text-[11px] leading-relaxed text-slate-500">
                    Each photon produces exactly one click, never both. In superposition the choice is random with 50%
                    each, and only the statistics over many photons show the “0 and 1” state. This needs single-photon
                    detectors at about 2 K; it is a concept, not part of the room-temperature chip.
                  </p>
                </Panel>
              </>
            )}

            <div className="flex items-start gap-2 rounded-2xl border border-white/10 bg-white/[0.02] p-4 text-xs text-slate-400">
              <Info className="mt-0.5 h-4 w-4 flex-shrink-0 text-slate-300" />
              <span>
                Arrival times follow the simulator&apos;s timing model: 14.5 mm vs 29.5 mm of Si₃N₄ spiral (group index
                2.0), laser, photodiode and TDC jitter σ ≈ 2.1 ps. The animation slows light down by tens of billions of times
                so the difference can be seen.
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

const Panel: React.FC<{ title: string; children: React.ReactNode }> = ({ title, children }) => (
  <div className="rounded-2xl border border-white/10 bg-white/[0.02] p-4">
    <div className="mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400">{title}</div>
    {children}
  </div>
);

const Stat: React.FC<{ label: string; value: string }> = ({ label, value }) => (
  <div className="flex justify-between border-b border-white/5 py-1.5 text-sm last:border-0">
    <span className="text-slate-400">{label}</span>
    <span className="font-mono text-white">{value}</span>
  </div>
);

const Controls: React.FC<{ busy: boolean; onOne: () => void; onBatch: () => void; onReset: () => void; oneLabel?: string }> = ({
  busy,
  onOne,
  onBatch,
  onReset,
  oneLabel = 'Fire 1 pulse',
}) => (
  <div className="grid grid-cols-3 gap-2">
    <button
      onClick={onOne}
      disabled={busy}
      className="flex items-center justify-center gap-1.5 rounded-xl bg-white px-2 py-2 text-xs font-bold text-black disabled:opacity-40"
    >
      <Play className="h-3.5 w-3.5" />
      {oneLabel}
    </button>
    <button onClick={onBatch} className="rounded-xl border border-white/15 px-2 py-2 text-xs font-semibold text-white hover:bg-white/5">
      Fire {BATCH_SIZE}
    </button>
    <button
      onClick={onReset}
      className="flex items-center justify-center gap-1.5 rounded-xl border border-white/15 px-2 py-2 text-xs font-semibold text-slate-300 hover:bg-white/5"
    >
      <RotateCcw className="h-3.5 w-3.5" />
      Reset
    </button>
  </div>
);

export default SectionDecision;
