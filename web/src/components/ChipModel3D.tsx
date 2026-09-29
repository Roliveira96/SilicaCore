import React, { useEffect, useMemo, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { ContactShadows, Environment, Html, Lightformer, OrbitControls } from '@react-three/drei';
import { Bloom, EffectComposer, Vignette } from '@react-three/postprocessing';
import * as THREE from 'three';
import { mergeGeometries } from 'three/examples/jsm/utils/BufferGeometryUtils.js';
import { createGridGraph, simulatePhysicalRace, DEFAULT_PARAMS, RaceHardwareResult } from '../lib/racelogic';

// ============================================================================
// Physically scaled model of the 16x16 race-logic chip (1 scene unit = 1 mm).
// The die is 25.5 x 25.5 mm (the 648 mm^2 block from the simulator), wire
// bonded on an organic package, with a fiber array unit (FAU) edge-coupled to
// the left facet. The glowing race is the real simulated 16x16 race rendered
// in false color: 1550 nm light is infrared and invisible to the eye.
// ============================================================================

const DIE = 25.5;
const DIE_T = 0.775;
const PKG = 44;
const PKG_T = 1.4;
const GRID = 16;
const MARGIN = 2.2;
const PITCH = (DIE - 2 * MARGIN) / (GRID - 1);
const TEX = 2048;
const PX = TEX / DIE;
const SOURCE_NODE = 8 * GRID; // left edge, row 8: fed by the input fiber
const TARGET_NODE = 7 * GRID + (GRID - 1); // right edge, row 7
const RACE_ANIMATION_S = 4.0; // the ~11.5 ns race is slowed down to 4 s
const INJECT_ANIMATION_S = 0.6;

const PKG_TOP = PKG_T / 2;
const DIE_TOP = PKG_TOP + DIE_T;

function nodeXZ(i: number): [number, number] {
  const x = i % GRID;
  const y = Math.floor(i / GRID);
  return [-DIE / 2 + MARGIN + x * PITCH, -DIE / 2 + MARGIN + y * PITCH];
}

const toPx = (mm: number) => (mm + DIE / 2) * PX;

// ---------------------------------------------------------------------------
// Procedural textures
// ---------------------------------------------------------------------------

type Layer = 'composite' | 'nitride' | 'tfln' | 'top';

function drawSpiral(ctx: CanvasRenderingContext2D, cx: number, cy: number, turns: number, rMax: number) {
  const steps = turns * 30;
  ctx.beginPath();
  for (let k = 0; k <= steps; k++) {
    const t = k / steps;
    const a = t * turns * Math.PI * 2;
    const r = rMax * (0.15 + 0.85 * t);
    const x = cx + r * Math.cos(a);
    const y = cy + r * Math.sin(a);
    if (k === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.stroke();
}

function forEachEdgePair(fn: (u: number, v: number, horizontal: boolean) => void) {
  for (let y = 0; y < GRID; y++) {
    for (let x = 0; x < GRID; x++) {
      const u = y * GRID + x;
      if (x + 1 < GRID) fn(u, u + 1, true);
      if (y + 1 < GRID) fn(u, u + GRID, false);
    }
  }
}

function drawWaveguides(ctx: CanvasRenderingContext2D, color: string, width: number) {
  ctx.strokeStyle = color;
  ctx.lineWidth = width;
  const off = 0.3 * PX;
  forEachEdgePair((u, v, horizontal) => {
    const [ux, uz] = nodeXZ(u);
    const [vx, vz] = nodeXZ(v);
    const mx = toPx((ux + vx) / 2);
    const mz = toPx((uz + vz) / 2);
    // One delay spiral per direction, side by side.
    for (const s of [-1, 1]) {
      const sx = horizontal ? mx : mx + s * off;
      const sz = horizontal ? mz + s * off : mz;
      ctx.beginPath();
      ctx.moveTo(toPx(ux), toPx(uz));
      ctx.lineTo(sx, sz);
      ctx.lineTo(toPx(vx), toPx(vz));
      ctx.stroke();
      drawSpiral(ctx, sx, sz, 4, 0.22 * PX);
    }
  });
  // Edge coupler from the left facet to the source node.
  const [sx, sz] = nodeXZ(SOURCE_NODE);
  ctx.lineWidth = width * 1.6;
  ctx.beginPath();
  ctx.moveTo(0, toPx(sz));
  ctx.lineTo(toPx(sx), toPx(sz));
  ctx.stroke();
}

function drawNodes(ctx: CanvasRenderingContext2D, withDetectors: boolean) {
  for (let i = 0; i < GRID * GRID; i++) {
    const [x, z] = nodeXZ(i);
    const cx = toPx(x);
    const cz = toPx(z);
    // TFLN modulator: ground-signal-ground travelling-wave electrodes.
    ctx.fillStyle = '#c9a24a';
    for (const k of [-1, 0, 1]) {
      ctx.fillRect(cx - 0.34 * PX, cz - 0.03 * PX + k * 0.1 * PX, 0.68 * PX, 0.06 * PX);
    }
    if (withDetectors) {
      // Photodiode with its contact frame.
      ctx.fillStyle = '#0b0d12';
      ctx.fillRect(cx + 0.18 * PX, cz + 0.18 * PX, 0.2 * PX, 0.2 * PX);
      ctx.strokeStyle = '#b8923e';
      ctx.lineWidth = 2;
      ctx.strokeRect(cx + 0.18 * PX, cz + 0.18 * PX, 0.2 * PX, 0.2 * PX);
    }
  }
}

function drawPhaseChangeSwitches(ctx: CanvasRenderingContext2D) {
  ctx.fillStyle = '#4b4f58';
  const off = 0.3 * PX;
  forEachEdgePair((u, v, horizontal) => {
    const [ux, uz] = nodeXZ(u);
    const [vx, vz] = nodeXZ(v);
    const mx = toPx((ux + vx) / 2);
    const mz = toPx((uz + vz) / 2);
    for (const s of [-1, 1]) {
      const sx = horizontal ? mx : mx + s * off;
      const sz = horizontal ? mz + s * off : mz;
      for (const d of [-1, 1]) {
        const px = horizontal ? sx + d * 0.3 * PX : sx;
        const pz = horizontal ? sz : sz + d * 0.3 * PX;
        ctx.fillRect(px - 5, pz - 5, 10, 10);
      }
    }
  });
}

function drawMetalRouting(ctx: CanvasRenderingContext2D, alpha: number) {
  ctx.strokeStyle = `rgba(196, 150, 80, ${alpha})`;
  ctx.lineWidth = 3;
  for (let y = 0; y < GRID; y++) {
    const [, z] = nodeXZ(y * GRID);
    ctx.beginPath();
    ctx.moveTo(toPx(-DIE / 2 + MARGIN), toPx(z) - 0.45 * PX);
    ctx.lineTo(TEX - 0.8 * PX, toPx(z) - 0.45 * PX);
    ctx.stroke();
  }
}

function drawDieFrame(ctx: CanvasRenderingContext2D) {
  // Seal ring, alignment marks and die marking.
  ctx.strokeStyle = 'rgba(200, 170, 110, 0.9)';
  ctx.lineWidth = 10;
  ctx.strokeRect(0.25 * PX, 0.25 * PX, TEX - 0.5 * PX, TEX - 0.5 * PX);
  ctx.lineWidth = 4;
  for (const [ax, az] of [
    [1.0, 1.0],
    [DIE - 1.0, 1.0],
    [1.0, DIE - 1.0],
    [DIE - 1.0, DIE - 1.0],
  ]) {
    const x = ax * PX;
    const z = az * PX;
    ctx.beginPath();
    ctx.moveTo(x - 20, z);
    ctx.lineTo(x + 20, z);
    ctx.moveTo(x, z - 20);
    ctx.lineTo(x, z + 20);
    ctx.stroke();
  }
  ctx.fillStyle = 'rgba(210, 215, 225, 0.55)';
  ctx.font = '600 30px "JetBrains Mono", "Fira Code", monospace';
  ctx.fillText('SILICACORE RL-16   Si3N4 / TFLN   ToF RACE LOGIC   REV A', 1.3 * PX, DIE * PX - 0.75 * PX);
}

function makeLayerTexture(layer: Layer): THREE.CanvasTexture {
  const canvas = document.createElement('canvas');
  canvas.width = TEX;
  canvas.height = TEX;
  const ctx = canvas.getContext('2d') as CanvasRenderingContext2D;

  if (layer === 'composite') {
    const g = ctx.createLinearGradient(0, 0, TEX, TEX);
    g.addColorStop(0, '#1b2330');
    g.addColorStop(0.5, '#141b26');
    g.addColorStop(1, '#1d2533');
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, TEX, TEX);
    drawWaveguides(ctx, 'rgba(110, 190, 220, 0.45)', 1.6);
    drawPhaseChangeSwitches(ctx);
    drawMetalRouting(ctx, 0.35);
    drawNodes(ctx, true);
    drawDieFrame(ctx);
  } else if (layer === 'nitride') {
    ctx.clearRect(0, 0, TEX, TEX);
    drawWaveguides(ctx, 'rgba(120, 220, 245, 0.95)', 2.2);
  } else if (layer === 'tfln') {
    ctx.fillStyle = 'rgba(150, 110, 210, 0.35)';
    ctx.fillRect(0, 0, TEX, TEX);
    drawNodes(ctx, false);
  } else {
    ctx.fillStyle = 'rgba(200, 220, 240, 0.08)';
    ctx.fillRect(0, 0, TEX, TEX);
    drawPhaseChangeSwitches(ctx);
    drawMetalRouting(ctx, 0.9);
    drawNodes(ctx, true);
    drawDieFrame(ctx);
  }

  const tex = new THREE.CanvasTexture(canvas);
  tex.colorSpace = THREE.SRGBColorSpace;
  tex.anisotropy = 8;
  return tex;
}

function makePackageTexture(): THREE.CanvasTexture {
  const size = 1024;
  const px = size / PKG;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;
  const ctx = canvas.getContext('2d') as CanvasRenderingContext2D;
  const g = ctx.createRadialGradient(size / 2, size / 2, size * 0.1, size / 2, size / 2, size * 0.75);
  g.addColorStop(0, '#12301f');
  g.addColorStop(1, '#0a1f14');
  ctx.fillStyle = g;
  ctx.fillRect(0, 0, size, size);

  // Traces from the bond fingers towards the package edge (vias).
  ctx.strokeStyle = 'rgba(190, 150, 70, 0.55)';
  ctx.lineWidth = 2;
  const ring = DIE / 2 + 2.4;
  for (let k = 0; k < 22; k++) {
    const t = -DIE / 2 + 1.5 + (k * (DIE - 3)) / 21;
    for (const [x0, z0, x1, z1] of [
      [t, -ring, t * 1.5, -PKG / 2 + 2],
      [t, ring, t * 1.5, PKG / 2 - 2],
      [ring, t, PKG / 2 - 2, t * 1.5],
    ]) {
      ctx.beginPath();
      ctx.moveTo((x0 + PKG / 2) * px, (z0 + PKG / 2) * px);
      ctx.lineTo((x1 + PKG / 2) * px, (z1 + PKG / 2) * px);
      ctx.stroke();
      ctx.fillStyle = '#c9a24a';
      ctx.beginPath();
      ctx.arc((x1 + PKG / 2) * px, (z1 + PKG / 2) * px, 4, 0, Math.PI * 2);
      ctx.fill();
    }
  }
  ctx.fillStyle = 'rgba(235, 240, 235, 0.85)';
  ctx.font = '600 22px "Inter", sans-serif';
  ctx.fillText('SILICACORE  RL-16  ENGINEERING CONCEPT', 2.2 * px, PKG * px - 1.6 * px);
  ctx.beginPath();
  ctx.moveTo(1.5 * px, 1.5 * px);
  ctx.lineTo(3.5 * px, 1.5 * px);
  ctx.lineTo(1.5 * px, 3.5 * px);
  ctx.closePath();
  ctx.fill();

  const tex = new THREE.CanvasTexture(canvas);
  tex.colorSpace = THREE.SRGBColorSpace;
  tex.anisotropy = 8;
  return tex;
}

// ---------------------------------------------------------------------------
// Static hardware
// ---------------------------------------------------------------------------

const Package: React.FC = () => {
  const topTex = useMemo(() => makePackageTexture(), []);
  const balls = useRef<THREE.InstancedMesh>(null);
  const BALLS = 22;

  useEffect(() => {
    if (!balls.current) return;
    const m = new THREE.Matrix4();
    let k = 0;
    for (let i = 0; i < BALLS; i++) {
      for (let j = 0; j < BALLS; j++) {
        const x = -PKG / 2 + 2 + (i * (PKG - 4)) / (BALLS - 1);
        const z = -PKG / 2 + 2 + (j * (PKG - 4)) / (BALLS - 1);
        m.makeTranslation(x, -PKG_T / 2 - 0.25, z);
        balls.current.setMatrixAt(k++, m);
      }
    }
    balls.current.instanceMatrix.needsUpdate = true;
  }, []);

  return (
    <group>
      <mesh castShadow receiveShadow>
        <boxGeometry args={[PKG, PKG_T, PKG]} />
        <meshStandardMaterial color="#0b2014" roughness={0.6} metalness={0.05} />
      </mesh>
      <mesh position={[0, PKG_TOP + 0.002, 0]} rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
        <planeGeometry args={[PKG, PKG]} />
        <meshStandardMaterial map={topTex} roughness={0.55} metalness={0.1} />
      </mesh>
      <instancedMesh ref={balls} args={[undefined, undefined, BALLS * BALLS]} castShadow>
        <sphereGeometry args={[0.32, 12, 12]} />
        <meshStandardMaterial color="#c8ccd2" metalness={0.9} roughness={0.25} />
      </instancedMesh>
      {/* Decoupling capacitors (0402) around the die */}
      {[
        [-17, -17],
        [-13, -17],
        [13, -17],
        [17, -17],
        [17, 13],
        [17, 17],
        [-13, 17],
        [-17, 17],
      ].map(([x, z], i) => (
        <group key={i} position={[x, PKG_TOP + 0.25, z]}>
          <mesh castShadow>
            <boxGeometry args={[1.0, 0.5, 0.5]} />
            <meshStandardMaterial color="#b89a6e" roughness={0.6} />
          </mesh>
          {[-0.45, 0.45].map((dx) => (
            <mesh key={dx} position={[dx, 0, 0]}>
              <boxGeometry args={[0.14, 0.52, 0.52]} />
              <meshStandardMaterial color="#d6d9de" metalness={0.95} roughness={0.2} />
            </mesh>
          ))}
        </group>
      ))}
    </group>
  );
};

const PAD_COUNT = 22;

function padPositions(): { pad: THREE.Vector3; finger: THREE.Vector3 }[] {
  const out: { pad: THREE.Vector3; finger: THREE.Vector3 }[] = [];
  const ring = DIE / 2 + 2.4;
  for (let k = 0; k < PAD_COUNT; k++) {
    const t = -DIE / 2 + 1.5 + (k * (DIE - 3)) / (PAD_COUNT - 1);
    out.push({ pad: new THREE.Vector3(t, DIE_TOP, -DIE / 2 + 0.45), finger: new THREE.Vector3(t, PKG_TOP, -ring) });
    out.push({ pad: new THREE.Vector3(t, DIE_TOP, DIE / 2 - 0.45), finger: new THREE.Vector3(t, PKG_TOP, ring) });
    out.push({ pad: new THREE.Vector3(DIE / 2 - 0.45, DIE_TOP, t), finger: new THREE.Vector3(ring, PKG_TOP, t) });
  }
  return out;
}

const WireBonds: React.FC<{ lift: number }> = ({ lift }) => {
  const { wires, pads, fingers } = useMemo(() => {
    const list = padPositions();
    const tubes: THREE.BufferGeometry[] = [];
    for (const { pad, finger } of list) {
      const top = pad.clone().add(new THREE.Vector3(0, 1.1, 0));
      const mid = pad.clone().lerp(finger, 0.55).add(new THREE.Vector3(0, 0.9, 0));
      const curve = new THREE.CatmullRomCurve3([pad, top, mid, finger]);
      tubes.push(new THREE.TubeGeometry(curve, 24, 0.035, 6, false));
    }
    const padGeoms = list.map(({ pad }) => new THREE.BoxGeometry(0.55, 0.02, 0.55).translate(pad.x, pad.y + 0.01, pad.z));
    const fingerGeoms = list.map(({ finger }) =>
      new THREE.BoxGeometry(0.6, 0.03, 0.9).translate(finger.x, finger.y + 0.015, finger.z),
    );
    return {
      wires: mergeGeometries(tubes) as THREE.BufferGeometry,
      pads: mergeGeometries(padGeoms) as THREE.BufferGeometry,
      fingers: mergeGeometries(fingerGeoms) as THREE.BufferGeometry,
    };
  }, []);

  const gold = <meshStandardMaterial color="#e2b65a" metalness={1} roughness={0.22} />;
  return (
    <group>
      <mesh geometry={wires} castShadow visible={lift < 0.05}>
        {gold}
      </mesh>
      <mesh geometry={pads} position={[0, lift * 4.2, 0]}>
        {gold}
      </mesh>
      <mesh geometry={fingers}>{gold}</mesh>
    </group>
  );
};

/** Fiber array unit butt-coupled to the left die facet, with a yellow SMF ribbon. */
const FiberArray: React.FC<{ glow: number }> = ({ glow }) => {
  const [, sz] = nodeXZ(SOURCE_NODE);
  const faceX = -DIE / 2;
  const y = DIE_TOP - 0.35;

  const { jacket, fibers } = useMemo(() => {
    const fiberCurves: THREE.CatmullRomCurve3[] = [];
    for (let k = 0; k < 8; k++) {
      const z = sz + (k - 3.5) * 0.25;
      fiberCurves.push(
        new THREE.CatmullRomCurve3([
          new THREE.Vector3(faceX - 3.0, y, z),
          new THREE.Vector3(faceX - 7.0, y + 0.2, z * 0.98),
          new THREE.Vector3(faceX - 10.0, y + 0.8, sz),
        ]),
      );
    }
    const jacketCurve = new THREE.CatmullRomCurve3([
      new THREE.Vector3(faceX - 10.0, y + 0.8, sz),
      new THREE.Vector3(faceX - 16.0, y + 3.0, sz - 2),
      new THREE.Vector3(faceX - 24.0, y + 2.0, sz - 10),
      new THREE.Vector3(faceX - 34.0, y - 1.0, sz - 24),
    ]);
    return {
      jacket: new THREE.TubeGeometry(jacketCurve, 64, 0.9, 16, false),
      fibers: mergeGeometries(fiberCurves.map((c) => new THREE.TubeGeometry(c, 24, 0.07, 8, false))) as THREE.BufferGeometry,
    };
  }, [faceX, sz, y]);

  return (
    <group>
      {/* Glass V-groove block and lid */}
      <mesh position={[faceX - 1.5, y - 0.35, sz]} castShadow>
        <boxGeometry args={[3.0, 1.4, 4.0]} />
        <meshPhysicalMaterial color="#dff3ff" transmission={0.92} roughness={0.08} thickness={1.4} ior={1.46} />
      </mesh>
      <mesh position={[faceX - 1.5, y + 0.55, sz]} castShadow>
        <boxGeometry args={[3.0, 0.4, 4.0]} />
        <meshPhysicalMaterial color="#e8f6ff" transmission={0.9} roughness={0.12} thickness={0.4} ior={1.46} />
      </mesh>
      <mesh geometry={fibers}>
        <meshStandardMaterial color="#f1f3f5" roughness={0.35} />
      </mesh>
      <mesh geometry={jacket} castShadow>
        <meshStandardMaterial color="#e3c01c" roughness={0.45} />
      </mesh>
      {/* False-color guided light in the input fiber */}
      <mesh position={[faceX - 5.0, y + 0.1, sz]} rotation={[0, 0, Math.PI / 2]} visible={glow > 0.01}>
        <cylinderGeometry args={[0.09, 0.09, 10, 8]} />
        <meshBasicMaterial color={new THREE.Color(1.5, 4, 5).multiplyScalar(glow)} toneMapped={false} />
      </mesh>
    </group>
  );
};

// ---------------------------------------------------------------------------
// Die, exploded layers and the simulated race
// ---------------------------------------------------------------------------

interface DieProps {
  explode: number;
  race: RaceHardwareResult | null;
  raceTime: React.MutableRefObject<number>; // simulated race time in ps, -1 when idle
}

const LAYERS: { key: Exclude<Layer, 'composite'>; label: string; detail: string; gap: number }[] = [
  { key: 'nitride', label: 'Si₃N₄ waveguides', detail: '800 nm × 0.7 µm, delay spirals (R ≥ 30 µm)', gap: 1 },
  { key: 'tfln', label: 'Thin-film LiNbO₃', detail: 'electro-optic switches and re-fire modulators', gap: 2 },
  { key: 'top', label: 'Cladding, metal and Sb₂Se₃', detail: 'non-volatile delay selection, routing, pads', gap: 3 },
];

const Die: React.FC<DieProps> = ({ explode, race, raceTime }) => {
  const composite = useMemo(() => makeLayerTexture('composite'), []);
  const layerTex = useMemo(
    () => ({
      nitride: makeLayerTexture('nitride'),
      tfln: makeLayerTexture('tfln'),
      top: makeLayerTexture('top'),
    }),
    [],
  );
  const nodes = useRef<THREE.InstancedMesh>(null);
  const tree = useRef<THREE.LineSegments>(null);
  const spread = explode * 4.2;

  const treeGeom = useMemo(() => {
    const g = new THREE.BufferGeometry();
    g.setAttribute('position', new THREE.BufferAttribute(new Float32Array(GRID * GRID * 2 * 3), 3));
    g.setAttribute('color', new THREE.BufferAttribute(new Float32Array(GRID * GRID * 2 * 3), 3));
    return g;
  }, []);

  // Winning path from source to target, highlighted brighter.
  const winning = useMemo(() => {
    const set = new Set<number>();
    if (!race) return set;
    let v = TARGET_NODE;
    while (v !== -1 && v !== SOURCE_NODE) {
      set.add(v);
      v = race.predecessors[v];
    }
    return set;
  }, [race]);

  useEffect(() => {
    if (!nodes.current) return;
    const m = new THREE.Matrix4();
    for (let i = 0; i < GRID * GRID; i++) {
      const [x, z] = nodeXZ(i);
      m.makeTranslation(x, 0, z);
      nodes.current.setMatrixAt(i, m);
      nodes.current.setColorAt(i, new THREE.Color(0, 0, 0));
    }
    nodes.current.instanceMatrix.needsUpdate = true;
  }, []);

  useFrame(() => {
    const t = raceTime.current;
    const lift = spread * 1; // nitride layer height offset when exploded
    if (nodes.current) nodes.current.position.y = DIE_TOP + 0.03 + lift;
    if (tree.current) tree.current.position.y = DIE_TOP + 0.02 + lift;
    if (!race || !nodes.current) return;

    const color = new THREE.Color();
    const pos = treeGeom.getAttribute('position') as THREE.BufferAttribute;
    const col = treeGeom.getAttribute('color') as THREE.BufferAttribute;
    for (let i = 0; i < GRID * GRID; i++) {
      const fire = race.fireTimesPS[i];
      const lit = t >= 0 && fire <= t;
      const age = lit ? (t - fire) / 1500 : 0;
      const hot = winning.has(i) || i === SOURCE_NODE;
      const k = lit ? (hot ? 5 : 1.6 + 2.5 * Math.exp(-age)) : 0;
      color.setRGB(0.25 * k, 0.8 * k, 1.0 * k);
      nodes.current.setColorAt(i, color);

      // Shortest-path tree: segment from predecessor, drawn up to the pulse front.
      const p = race.predecessors[i];
      const base = i * 2;
      if (p >= 0 && t >= 0) {
        const [px, pz] = nodeXZ(p);
        const [x, z] = nodeXZ(i);
        const start = race.fireTimesPS[p];
        const f = Math.max(0, Math.min(1, (t - start) / Math.max(1, fire - start)));
        pos.setXYZ(base, px, 0, pz);
        pos.setXYZ(base + 1, px + (x - px) * f, 0, pz + (z - pz) * f);
        const c = hot && f >= 1 ? 4 : 1.2;
        col.setXYZ(base, 0.2 * c, 0.75 * c, c);
        col.setXYZ(base + 1, 0.2 * c, 0.75 * c, c);
      } else {
        pos.setXYZ(base, 0, 0, 0);
        pos.setXYZ(base + 1, 0, 0, 0);
      }
    }
    if (nodes.current.instanceColor) nodes.current.instanceColor.needsUpdate = true;
    pos.needsUpdate = true;
    col.needsUpdate = true;
  });

  const assembled = Math.max(0, 1 - explode * 5);

  return (
    <group>
      {/* Silicon handle wafer */}
      <mesh position={[0, PKG_TOP + DIE_T / 2, 0]} castShadow receiveShadow>
        <boxGeometry args={[DIE, DIE_T, DIE]} />
        <meshPhysicalMaterial color="#2a2f37" metalness={0.55} roughness={0.28} clearcoat={0.6} />
      </mesh>

      {/* Assembled die surface: patterned, clear-coated, with thin-film iridescence */}
      <mesh position={[0, DIE_TOP + 0.004, 0]} rotation={[-Math.PI / 2, 0, 0]} visible={assembled > 0.01}>
        <planeGeometry args={[DIE, DIE]} />
        <meshPhysicalMaterial
          map={composite}
          transparent
          opacity={assembled}
          metalness={0.35}
          roughness={0.32}
          clearcoat={1}
          clearcoatRoughness={0.08}
          iridescence={0.55}
          iridescenceIOR={1.35}
          iridescenceThicknessRange={[180, 520]}
        />
      </mesh>

      {/* Exploded process stack */}
      {explode > 0.02 && (
        <group>
          <mesh position={[0, DIE_TOP + 0.03, 0]}>
            <boxGeometry args={[DIE, 0.06, DIE]} />
            <meshPhysicalMaterial color="#cfe6ff" transparent opacity={0.35} roughness={0.1} />
          </mesh>
          <Html position={[DIE / 2 + 1, DIE_TOP, 0]} className="pointer-events-none select-none">
            <LayerLabel title="Buried oxide on Si" detail="thermal SiO₂ isolates the modes from the substrate" />
          </Html>
          {LAYERS.map((l) => (
            <group key={l.key} position={[0, DIE_TOP + 0.08 + spread * l.gap, 0]}>
              <mesh rotation={[-Math.PI / 2, 0, 0]}>
                <planeGeometry args={[DIE, DIE]} />
                <meshStandardMaterial
                  map={layerTex[l.key]}
                  transparent
                  opacity={Math.min(1, explode * 3)}
                  side={THREE.DoubleSide}
                  depthWrite={false}
                  metalness={l.key === 'top' ? 0.4 : 0.1}
                  roughness={0.35}
                />
              </mesh>
              <Html position={[DIE / 2 + 1, 0, 0]} className="pointer-events-none select-none">
                <LayerLabel title={l.label} detail={l.detail} />
              </Html>
            </group>
          ))}
        </group>
      )}

      {/* Simulated race: nodes firing and the shortest-path tree (false color) */}
      <instancedMesh ref={nodes} args={[undefined, undefined, GRID * GRID]}>
        <boxGeometry args={[0.34, 0.02, 0.34]} />
        <meshBasicMaterial toneMapped={false} />
      </instancedMesh>
      <lineSegments ref={tree} geometry={treeGeom}>
        <lineBasicMaterial vertexColors toneMapped={false} />
      </lineSegments>
    </group>
  );
};

const LayerLabel: React.FC<{ title: string; detail: string }> = ({ title, detail }) => (
  <div className="whitespace-nowrap rounded-lg border border-white/15 bg-black/75 px-2.5 py-1.5 text-[11px] leading-tight text-slate-200 backdrop-blur">
    <div className="font-semibold text-white">{title}</div>
    <div className="text-slate-400">{detail}</div>
  </div>
);

// ---------------------------------------------------------------------------
// Scene
// ---------------------------------------------------------------------------

export interface ChipModel3DProps {
  explode: number;
  fireToken: number; // increment to start a new race
  onRaceComplete?: (race: RaceHardwareResult) => void;
}

const RaceDriver: React.FC<{
  fireToken: number;
  race: RaceHardwareResult | null;
  raceTime: React.MutableRefObject<number>;
  fiberGlow: React.MutableRefObject<number>;
  onDone: () => void;
}> = ({ fireToken, race, raceTime, fiberGlow, onDone }) => {
  const start = useRef<number | null>(null);
  const finished = useRef(true);

  useEffect(() => {
    if (fireToken > 0) {
      start.current = null;
      finished.current = false;
    }
  }, [fireToken]);

  useFrame(({ clock }) => {
    if (finished.current || !race) return;
    if (start.current === null) start.current = clock.getElapsedTime();
    const e = clock.getElapsedTime() - start.current;
    if (e < INJECT_ANIMATION_S) {
      fiberGlow.current = Math.sin((e / INJECT_ANIMATION_S) * Math.PI) + 0.2;
      raceTime.current = -1;
      return;
    }
    fiberGlow.current = Math.max(0, fiberGlow.current - 0.05);
    const f = (e - INJECT_ANIMATION_S) / RACE_ANIMATION_S;
    raceTime.current = f * race.maxToF_PS;
    if (f >= 1.15) {
      finished.current = true;
      onDone();
    }
  });
  return null;
};

const FiberGlow: React.FC<{ level: React.MutableRefObject<number> }> = ({ level }) => {
  const [glow, setGlow] = React.useState(0);
  useFrame(() => {
    const v = Math.round(level.current * 20) / 20;
    if (v !== glow) setGlow(v);
  });
  return <FiberArray glow={glow} />;
};

const AutoRotate: React.FC<{ enabled: boolean }> = ({ enabled }) => {
  const ref = useRef<React.ElementRef<typeof OrbitControls>>(null);
  useFrame(() => {
    if (ref.current) ref.current.autoRotate = enabled;
  });
  return (
    <OrbitControls
      ref={ref}
      enablePan={false}
      enableDamping
      autoRotateSpeed={0.35}
      minDistance={24}
      maxDistance={110}
      maxPolarAngle={Math.PI / 2.1}
      target={[0, 1, 0]}
    />
  );
};

export const ChipModel3D: React.FC<ChipModel3DProps> = ({ explode, fireToken, onRaceComplete }) => {
  const graph = useMemo(() => createGridGraph(GRID, GRID, 15, 42), []);
  const [race, setRace] = React.useState<RaceHardwareResult | null>(null);
  const raceTime = useRef(-1);
  const fiberGlow = useRef(0);

  useEffect(() => {
    if (fireToken === 0) return;
    raceTime.current = -1;
    setRace(simulatePhysicalRace(graph, SOURCE_NODE, DEFAULT_PARAMS));
  }, [fireToken, graph]);

  return (
    <Canvas shadows dpr={[1, 2]} camera={{ position: [34, 40, 50], fov: 32 }} gl={{ antialias: true }}>
      <color attach="background" args={['#05070c']} />
      <fog attach="fog" args={['#05070c', 90, 190]} />

      <ambientLight intensity={0.15} />
      <directionalLight position={[20, 40, 15]} intensity={1.2} castShadow shadow-mapSize={[2048, 2048]} />
      <Environment resolution={256}>
        <Lightformer form="rect" intensity={3} position={[0, 30, 0]} rotation-x={Math.PI / 2} scale={[40, 40, 1]} />
        <Lightformer form="rect" intensity={1.6} position={[-35, 12, 10]} rotation-y={Math.PI / 2} scale={[30, 8, 1]} />
        <Lightformer form="rect" intensity={1.2} position={[35, 10, -10]} rotation-y={-Math.PI / 2} scale={[30, 6, 1]} />
        <Lightformer form="ring" color="#5ad8ff" intensity={1.2} position={[0, 12, -40]} scale={12} />
      </Environment>

      <group>
        <Package />
        <Die explode={explode} race={race} raceTime={raceTime} />
        <WireBonds lift={explode} />
        <FiberGlow level={fiberGlow} />
      </group>

      <ContactShadows position={[0, -PKG_T / 2 - 0.6, 0]} opacity={0.65} scale={120} blur={2.4} far={30} />
      <RaceDriver
        fireToken={fireToken}
        race={race}
        raceTime={raceTime}
        fiberGlow={fiberGlow}
        onDone={() => race && onRaceComplete?.(race)}
      />
      <AutoRotate enabled={explode < 0.02} />

      <EffectComposer multisampling={4}>
        <Bloom mipmapBlur intensity={1.1} luminanceThreshold={1} luminanceSmoothing={0.25} />
        <Vignette offset={0.3} darkness={0.55} />
      </EffectComposer>
    </Canvas>
  );
};

export default ChipModel3D;
