import React, { useEffect, useMemo, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { ContactShadows, Environment, Html, Lightformer, OrbitControls } from '@react-three/drei';
import { Bloom, EffectComposer } from '@react-three/postprocessing';
import * as THREE from 'three';

// ============================================================================
// 3D model of how a bit is decided.
// ToF mode: a TFLN switch sends the pulse through a short or a long waveguide;
// both run at the same speed, so the long one arrives later. The detector's
// TDC stamps the arrival on a time ruler with two windows and a threshold.
// Quantum mode: one photon meets a 50:50 coupler, travels as an amplitude in
// both rails, and exactly one cryogenic detector clicks.
// Timing values come from the project model (t1 = 96.7 ps, t0 = 196.7 ps) on
// Si3N4 spirals with group index 2.0 (14.5 mm and 29.5 mm).
// ============================================================================

export const T1_PS = 96.73;
export const T0_PS = 196.73;
export const HALF_WINDOW_PS = 47.46;
export const THRESHOLD_PS = (T1_PS + T0_PS) / 2;
export const HIST_BIN_PS = 2;
export const HIST_MAX_PS = 250;
export const HIST_BINS = HIST_MAX_PS / HIST_BIN_PS;

const RULER_X0 = 8.5;
const RULER_LEN = 11;
const psToX = (ps: number) => RULER_X0 + (Math.min(Math.max(ps, 0), HIST_MAX_PS) / HIST_MAX_PS) * RULER_LEN;

const PULSE_SPEED = 7.5; // scene units per second, identical on both paths

export type QuantumState = 'zero' | 'one' | 'super';

export interface Decision3DProps {
  mode: 'tof' | 'quantum';
  bit: 0 | 1; // ToF: 1 = pass through the short path, 0 = divert to the long path
  shotToken: number; // increment to fire one animated pulse / photon
  lastArrivalPS: number | null; // measured arrival of the last ToF pulse
  histogram: number[]; // ToF arrival histogram, HIST_BINS bins
  quantumState: QuantumState;
  quantumOutcome: 0 | 1 | null; // which detector clicked for the last photon
  onArrive: () => void; // called when the animated pulse reaches the detector
}

// ---------------------------------------------------------------------------
// Shared pieces
// ---------------------------------------------------------------------------

const Waveguide: React.FC<{ curve: THREE.Curve<THREE.Vector3>; lit: number; color?: string }> = ({ curve, lit, color = '#5fd4f0' }) => {
  const geom = useMemo(() => new THREE.TubeGeometry(curve, 240, 0.07, 10, false), [curve]);
  const c = new THREE.Color(color);
  return (
    <mesh geometry={geom}>
      <meshStandardMaterial
        color="#2a5563"
        emissive={c}
        emissiveIntensity={0.15 + lit * 2.2}
        roughness={0.25}
        metalness={0.1}
        toneMapped={lit < 0.05}
      />
    </mesh>
  );
};

const Label: React.FC<{ position: [number, number, number]; children: React.ReactNode; tone?: string }> = ({ position, children, tone = 'text-slate-200' }) => (
  <Html position={position} center className="pointer-events-none select-none">
    <div className={`whitespace-nowrap rounded-md border border-white/10 bg-black/70 px-2 py-1 font-mono text-[11px] ${tone}`}>{children}</div>
  </Html>
);

const Base: React.FC<{ width: number; depth: number; x: number; z?: number }> = ({ width, depth, x, z = 0 }) => (
  <mesh position={[x, -0.18, z]} receiveShadow>
    <boxGeometry args={[width, 0.2, depth]} />
    <meshPhysicalMaterial color="#141a24" roughness={0.35} metalness={0.4} clearcoat={0.8} clearcoatRoughness={0.15} />
  </mesh>
);

const Pulse: React.FC<{ curve: THREE.Curve<THREE.Vector3> | null; token: number; onDone: () => void; color?: string; scale?: number }> = ({
  curve,
  token,
  onDone,
  color = '#9ff3ff',
  scale = 1,
}) => {
  const ref = useRef<THREE.Group>(null);
  const start = useRef<number | null>(null);
  const active = useRef(false);

  useEffect(() => {
    if (token > 0 && curve) {
      start.current = null;
      active.current = true;
    }
  }, [token, curve]);

  useFrame(({ clock }) => {
    if (!ref.current) return;
    if (!active.current || !curve) {
      ref.current.visible = false;
      return;
    }
    if (start.current === null) start.current = clock.getElapsedTime();
    const duration = curve.getLength() / PULSE_SPEED;
    const f = (clock.getElapsedTime() - start.current) / duration;
    if (f >= 1) {
      active.current = false;
      ref.current.visible = false;
      onDone();
      return;
    }
    ref.current.visible = true;
    ref.current.position.copy(curve.getPointAt(f));
  });

  const glow = new THREE.Color(color).multiplyScalar(6);
  return (
    <group ref={ref} visible={false} scale={scale}>
      <mesh>
        <sphereGeometry args={[0.24, 20, 20]} />
        <meshBasicMaterial color={glow} toneMapped={false} />
      </mesh>
      <pointLight color={color} intensity={4} distance={3} />
    </group>
  );
};

// ---------------------------------------------------------------------------
// ToF scene
// ---------------------------------------------------------------------------

function tofCurves() {
  const v = (x: number, z: number) => new THREE.Vector3(x, 0, z);
  const input = [v(-9, 0), v(-6.5, 0), v(-4.9, 0)];
  const shortPts = [...input, v(-3.1, -0.9), v(-1.5, -1.2), v(3.6, -1.2), v(4.6, 0), v(6.0, 0)];

  // Serpentine delay line: 7 meander legs sized so the whole path is 2.03x the
  // short one, the same ratio as the 29.5 mm and 14.5 mm Si3N4 paths of the model.
  const legs = 7;
  const legHeight = 1.68;
  const zBase = 1.2;
  const x0 = -2.2;
  const pitch = (3.4 - x0) / legs;
  const longPts: THREE.Vector3[] = [...input, v(-3.4, 0.9), v(-2.6, zBase)];
  let z = zBase;
  for (let k = 0; k < legs; k++) {
    const x = x0 + k * pitch;
    const top = k % 2 === 0;
    const zEnd = top ? zBase + legHeight : zBase;
    for (let s = 0; s <= 6; s++) longPts.push(v(x, z + ((zEnd - z) * s) / 6));
    const cx = x + pitch / 2;
    for (let s = 1; s < 12; s++) {
      const a = (Math.PI * s) / 12;
      longPts.push(v(cx - (pitch / 2) * Math.cos(a), zEnd + (top ? 1 : -1) * (pitch / 2) * Math.sin(a)));
    }
    z = zEnd;
  }
  const xLast = x0 + legs * pitch;
  for (let s = 0; s <= 6; s++) longPts.push(v(xLast, z + ((zBase - z) * s) / 6));
  longPts.push(v(4.0, 0.9), v(4.6, 0.2), v(6.0, 0));

  return {
    short: new THREE.CatmullRomCurve3(shortPts, false, 'centripetal', 0.3),
    long: new THREE.CatmullRomCurve3(longPts, false, 'centripetal', 0.3),
  };
}

const Histogram: React.FC<{ histogram: number[] }> = ({ histogram }) => {
  const mesh = useRef<THREE.InstancedMesh>(null);
  useEffect(() => {
    if (!mesh.current) return;
    const max = Math.max(1, ...histogram);
    const m = new THREE.Matrix4();
    const col = new THREE.Color();
    const binW = (RULER_LEN / HIST_BINS) * 0.85;
    for (let i = 0; i < HIST_BINS; i++) {
      const h = 0.05 + (2.6 * histogram[i]) / max;
      const empty = histogram[i] === 0;
      const center = (i + 0.5) * HIST_BIN_PS;
      m.compose(new THREE.Vector3(psToX(center), 0.05 + h / 2, -0.9), new THREE.Quaternion(), empty ? new THREE.Vector3(0, 0, 0) : new THREE.Vector3(binW, h, 0.35));
      mesh.current.setMatrixAt(i, m);
      const in1 = Math.abs(center - T1_PS) <= HALF_WINDOW_PS;
      const in0 = Math.abs(center - T0_PS) <= HALF_WINDOW_PS;
      col.set(in1 ? '#22d3ee' : in0 ? '#a78bfa' : '#f43f5e');
      mesh.current.setColorAt(i, col);
    }
    mesh.current.instanceMatrix.needsUpdate = true;
    if (mesh.current.instanceColor) mesh.current.instanceColor.needsUpdate = true;
  }, [histogram]);
  return (
    <instancedMesh ref={mesh} args={[undefined, undefined, HIST_BINS]} castShadow>
      <boxGeometry args={[1, 1, 1]} />
      <meshStandardMaterial roughness={0.4} metalness={0.1} />
    </instancedMesh>
  );
};

const TofScene: React.FC<Decision3DProps> = ({ bit, shotToken, lastArrivalPS, histogram, onArrive }) => {
  const curves = useMemo(() => tofCurves(), []);
  const [lit, setLit] = React.useState(0);
  const activeCurve = bit === 1 ? curves.short : curves.long;

  useEffect(() => {
    if (shotToken > 0) setLit(1);
  }, [shotToken]);

  const decoded =
    lastArrivalPS === null
      ? null
      : Math.abs(lastArrivalPS - T1_PS) <= HALF_WINDOW_PS
        ? 1
        : Math.abs(lastArrivalPS - T0_PS) <= HALF_WINDOW_PS
          ? 0
          : -1;

  return (
    <group>
      <Base width={17} depth={7.4} x={-1.5} z={0.6} />
      <Base width={12.5} depth={4.2} x={RULER_X0 + RULER_LEN / 2} />

      {/* Source: fiber ferrule */}
      <mesh position={[-9.4, 0, 0]} rotation={[0, 0, Math.PI / 2]} castShadow>
        <cylinderGeometry args={[0.28, 0.28, 0.9, 24]} />
        <meshStandardMaterial color="#d8dde3" metalness={0.9} roughness={0.2} />
      </mesh>
      <Label position={[-8.6, 0.3, 1.3]}>1550 nm pulse in</Label>

      {/* TFLN switch with gold electrodes */}
      <mesh position={[-4.2, 0.02, 0]} castShadow>
        <boxGeometry args={[1.6, 0.18, 1.3]} />
        <meshPhysicalMaterial color="#7c5cc9" transparent opacity={0.55} roughness={0.1} clearcoat={1} />
      </mesh>
      {[-0.35, 0, 0.35].map((z) => (
        <mesh key={z} position={[-4.2, 0.14, z]}>
          <boxGeometry args={[1.4, 0.04, 0.12]} />
          <meshStandardMaterial
            color="#e2b65a"
            metalness={1}
            roughness={0.25}
            emissive="#ffb347"
            emissiveIntensity={bit === 0 ? 1.4 : 0}
            toneMapped={bit !== 0}
          />
        </mesh>
      ))}
      <Label position={[-6.0, 0.4, -1.5]} tone={bit === 0 ? 'text-amber-300' : 'text-slate-200'}>
        TFLN switch: {bit === 1 ? 'pass (short path)' : 'divert (long path)'}
      </Label>

      <Waveguide curve={curves.short} lit={bit === 1 ? lit : 0} />
      <Waveguide curve={curves.long} lit={bit === 0 ? lit : 0} color="#c4a5ff" />
      <Label position={[2.6, 0.4, -2.0]} tone="text-cyan-300">
        short path · 14.5 mm · ~96.7 ps · reads 1
      </Label>
      <Label position={[0.6, 0.4, 3.75]} tone="text-violet-300">
        long path · 29.5 mm · ~196.7 ps · reads 0
      </Label>

      {/* Photodiode and TDC */}
      <mesh position={[6.25, 0.05, 0]} castShadow>
        <boxGeometry args={[0.5, 0.3, 0.5]} />
        <meshStandardMaterial color="#0c0f14" roughness={0.4} />
      </mesh>
      <mesh position={[6.25, 0.21, 0]}>
        <torusGeometry args={[0.2, 0.03, 8, 24]} />
        <meshStandardMaterial color="#e2b65a" metalness={1} roughness={0.2} />
      </mesh>
      <mesh position={[7.3, 0.05, 0]} castShadow>
        <boxGeometry args={[1.1, 0.25, 0.9]} />
        <meshStandardMaterial color="#111318" roughness={0.5} metalness={0.3} />
      </mesh>
      <Label position={[6.8, 0.75, 0]}>photodiode → TDC</Label>

      {/* Time ruler, windows and threshold */}
      <mesh position={[RULER_X0 + RULER_LEN / 2, 0.01, 0.6]}>
        <boxGeometry args={[RULER_LEN, 0.02, 0.08]} />
        <meshStandardMaterial color="#8a93a3" />
      </mesh>
      {[0, 50, 100, 150, 200, 250].map((ps) => (
        <group key={ps}>
          <mesh position={[psToX(ps), 0.02, 0.6]}>
            <boxGeometry args={[0.03, 0.04, 0.3]} />
            <meshStandardMaterial color="#8a93a3" />
          </mesh>
          <Label position={[psToX(ps), 0.05, 1.15]} tone="text-slate-400">
            {ps} ps
          </Label>
        </group>
      ))}
      <mesh position={[psToX(T1_PS), 1.4, 0]}>
        <boxGeometry args={[(2 * HALF_WINDOW_PS * RULER_LEN) / HIST_MAX_PS, 2.8, 3.2]} />
        <meshPhysicalMaterial color="#22d3ee" transparent opacity={0.08} roughness={0.2} depthWrite={false} />
      </mesh>
      <mesh position={[psToX(T0_PS), 1.4, 0]}>
        <boxGeometry args={[(2 * HALF_WINDOW_PS * RULER_LEN) / HIST_MAX_PS, 2.8, 3.2]} />
        <meshPhysicalMaterial color="#a78bfa" transparent opacity={0.08} roughness={0.2} depthWrite={false} />
      </mesh>
      <mesh position={[psToX(THRESHOLD_PS), 1.4, 0]}>
        <boxGeometry args={[0.03, 2.8, 3.2]} />
        <meshBasicMaterial color="#f43f5e" transparent opacity={0.45} depthWrite={false} />
      </mesh>
      <Label position={[psToX(T1_PS), 2.9, -1.7]} tone="text-cyan-300">
        window “1”
      </Label>
      <Label position={[psToX(T0_PS), 2.9, -1.7]} tone="text-violet-300">
        window “0”
      </Label>
      <Label position={[psToX(THRESHOLD_PS), 4.4, -1.7]} tone="text-rose-300">
        threshold 146.7 ps
      </Label>

      <Histogram histogram={histogram} />

      {lastArrivalPS !== null && (
        <group position={[psToX(lastArrivalPS), 0.35, 0.6]}>
          <mesh>
            <coneGeometry args={[0.16, 0.4, 16]} />
            <meshBasicMaterial
              color={new THREE.Color(decoded === 1 ? '#22d3ee' : decoded === 0 ? '#a78bfa' : '#f43f5e').multiplyScalar(4)}
              toneMapped={false}
            />
          </mesh>
          <Label
            position={[0, 0.75, 0]}
            tone={decoded === 1 ? 'text-cyan-300' : decoded === 0 ? 'text-violet-300' : 'text-rose-300'}
          >
            {lastArrivalPS.toFixed(1)} ps → {decoded === -1 ? 'rejected' : `bit ${decoded}`}
          </Label>
        </group>
      )}

      <Pulse
        curve={activeCurve}
        token={shotToken}
        onDone={() => {
          setLit(0);
          onArrive();
        }}
      />
    </group>
  );
};

// ---------------------------------------------------------------------------
// Quantum scene
// ---------------------------------------------------------------------------

function quantumCurves() {
  const v = (x: number, z: number) => new THREE.Vector3(x, 0, z);
  // Rails come close in the directional coupler (x = -2..0) and separate towards the detectors.
  const rail0 = new THREE.CatmullRomCurve3([v(-8, -1.4), v(-4, -1.4), v(-2.4, -0.28), v(0.4, -0.28), v(2, -1.4), v(6.2, -1.4)]);
  const rail1 = new THREE.CatmullRomCurve3([v(-8, 1.4), v(-4, 1.4), v(-2.4, 0.28), v(0.4, 0.28), v(2, 1.4), v(6.2, 1.4)]);
  return { rail0, rail1 };
}

const QuantumScene: React.FC<Decision3DProps> = ({ shotToken, quantumState, quantumOutcome, onArrive }) => {
  const { rail0, rail1 } = useMemo(() => quantumCurves(), []);
  const [flying, setFlying] = React.useState(false);
  useEffect(() => {
    if (shotToken > 0) setFlying(true);
  }, [shotToken]);

  // In superposition the photon is drawn as a half-bright amplitude on each rail.
  const inRail0 = quantumState !== 'one';
  const inRail1 = quantumState !== 'zero';
  const amp = quantumState === 'super' ? 0.55 : 1;

  return (
    <group>
      <Base width={17} depth={5.5} x={-0.8} />
      <Waveguide curve={rail0} lit={flying && inRail0 ? 0.6 : 0} />
      <Waveguide curve={rail1} lit={flying && inRail1 ? 0.6 : 0} color="#c4a5ff" />
      <Label position={[-7.4, 0.55, -1.4]} tone="text-cyan-300">
        rail |0⟩
      </Label>
      <Label position={[-7.4, 0.55, 1.4]} tone="text-violet-300">
        rail |1⟩
      </Label>

      {/* 50:50 directional coupler region */}
      <mesh position={[-1, 0.02, 0]}>
        <boxGeometry args={[2.8, 0.05, 1.1]} />
        <meshPhysicalMaterial
          color="#34d399"
          transparent
          opacity={quantumState === 'super' ? 0.28 : 0.06}
          roughness={0.1}
          depthWrite={false}
        />
      </mesh>
      <Label position={[-1, 0.75, 0]} tone={quantumState === 'super' ? 'text-emerald-300' : 'text-slate-500'}>
        50:50 coupler {quantumState === 'super' ? '(Hadamard)' : '(bypassed)'}
      </Label>

      {/* Cryostat with two SNSPD detectors */}
      <mesh position={[7.4, 0.8, 0]}>
        <cylinderGeometry args={[2.1, 2.1, 2, 48, 1, true]} />
        <meshStandardMaterial color="#bfe3ff" transparent opacity={0.1} roughness={0.1} side={THREE.DoubleSide} depthWrite={false} />
      </mesh>
      <mesh position={[7.4, 1.8, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[2.1, 0.06, 12, 64]} />
        <meshStandardMaterial color="#d8dde3" metalness={0.9} roughness={0.25} />
      </mesh>
      <mesh position={[7.4, -0.05, 0]}>
        <cylinderGeometry args={[2.1, 2.1, 0.1, 48]} />
        <meshStandardMaterial color="#2c3845" metalness={0.6} roughness={0.5} />
      </mesh>
      {[0, 1].map((k) => {
        const clicked = quantumOutcome === k;
        return (
          <mesh key={k} position={[6.9, 0.1, k === 0 ? -1.4 : 1.4]}>
            <boxGeometry args={[0.6, 0.25, 0.6]} />
            <meshStandardMaterial
              color="#10151d"
              emissive={k === 0 ? '#22d3ee' : '#a78bfa'}
              emissiveIntensity={clicked ? 4 : 0}
              toneMapped={!clicked}
            />
          </mesh>
        );
      })}
      <Label position={[7.4, 2.4, 0]} tone="text-sky-200">
        SNSPD detectors · ~2 K
      </Label>
      {quantumOutcome !== null && (
        <Label position={[6.9, 0.8, quantumOutcome === 0 ? -1.4 : 1.4]} tone={quantumOutcome === 0 ? 'text-cyan-300' : 'text-violet-300'}>
          click → read {quantumOutcome}
        </Label>
      )}

      {inRail0 && (
        <Pulse curve={rail0} token={shotToken} onDone={() => { setFlying(false); onArrive(); }} scale={amp} />
      )}
      {inRail1 && (
        <Pulse
          curve={rail1}
          token={shotToken}
          onDone={() => {
            setFlying(false);
            if (!inRail0) onArrive();
          }}
          color="#d7c6ff"
          scale={amp}
        />
      )}
    </group>
  );
};

// ---------------------------------------------------------------------------
// Canvas
// ---------------------------------------------------------------------------

export const DecisionMechanism3D: React.FC<Decision3DProps> = (props) => (
  <Canvas shadows dpr={[1, 2]} camera={{ position: props.mode === 'tof' ? [0, 25, 16] : [0, 15, 12], fov: 42 }}>
    <color attach="background" args={['#05070c']} />
    <ambientLight intensity={0.25} />
    <directionalLight position={[8, 16, 10]} intensity={1.1} castShadow shadow-mapSize={[2048, 2048]} />
    <Environment resolution={256}>
      <Lightformer form="rect" intensity={2.5} position={[0, 20, 0]} rotation-x={Math.PI / 2} scale={[30, 30, 1]} />
      <Lightformer form="rect" intensity={1.2} position={[-20, 8, 6]} rotation-y={Math.PI / 2} scale={[20, 6, 1]} />
    </Environment>
    <group position={props.mode === 'tof' ? [-4.4, 0, -0.6] : [-1.2, 0, 0]}>{props.mode === 'tof' ? <TofScene {...props} /> : <QuantumScene {...props} />}</group>
    <ContactShadows position={[0, -0.29, 0]} opacity={0.55} scale={50} blur={2.2} far={10} />
    <OrbitControls enablePan={false} enableDamping minDistance={8} maxDistance={40} maxPolarAngle={Math.PI / 2.15} />
    <EffectComposer multisampling={4}>
      <Bloom mipmapBlur intensity={0.9} luminanceThreshold={1} luminanceSmoothing={0.25} />
    </EffectComposer>
  </Canvas>
);

export default DecisionMechanism3D;
