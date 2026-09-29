import React, { useState, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Text, Float, Billboard } from '@react-three/drei';
import * as THREE from 'three';
import { Layers, Eye, Sliders, CheckCircle2, Info, Sparkles, Box, Zap, Radio, Atom } from 'lucide-react';

interface LayerInfo {
  id: number;
  name: string;
  tier: string;
  material: string;
  description: string;
  spec: string;
  color: string;
  wireframeColor: string;
  baseY: number;
  explodeMultiplier: number;
}

const LAYERS: LayerInfo[] = [
  {
    id: 1,
    tier: 'Top Tier',
    name: 'Electrical Control & Non-Volatile Phase Change',
    material: 'Sb₂Se₃ (Antimony Selenide) + ITO Micro-heaters',
    description:
      'Zero-static-power optical memory cells. Sb₂Se₃ switches switch between amorphous and crystalline states to latch binary delay lines without consuming continuous electrical power.',
    spec: '0.25 dB insertion loss per switch • 1 µs parallel reprogramming pulse',
    color: '#f59e0b', // Amber/gold
    wireframeColor: '#fbbf24',
    baseY: 1.5,
    explodeMultiplier: 1.8,
  },
  {
    id: 2,
    tier: 'Active Layer',
    name: 'Electro-Optic Modulation & Fast Gating',
    material: 'Thin-Film Lithium Niobate (TFLN) on Insulator',
    description:
      'Ultra-fast Pockels-effect electro-optic switches and Mach-Zehnder re-fire modulators. When a pulse arrives at a node, the UTC detector triggers the TFLN modulator to release a fresh 1550 nm pulse into outgoing delay spirals.',
    spec: '> 100 GHz modulation bandwidth • 20 ps node regeneration latency',
    color: '#8b5cf6', // Violet
    wireframeColor: '#a78bfa',
    baseY: 0.5,
    explodeMultiplier: 0.6,
  },
  {
    id: 3,
    tier: 'Waveguide Routing Tier',
    name: 'Thick Silicon Nitride Low-Loss Delay Spirals',
    material: 'Si₃N₄ (800 nm × 0.7 µm) in SiO₂ Cladding',
    description:
      'Planar spiral delay lines encoding graph edge weights. Thick 800 nm core eliminates bend radiation losses down to 50 µm bend radius, enabling dense binary delay stages.',
    spec: '0.0026 dB/90° at 50 µm bend radius • 0.05 dB/cm propagation loss • 3 µm spiral pitch',
    color: '#06b6d4', // Cyan
    wireframeColor: '#22d3ee',
    baseY: -0.5,
    explodeMultiplier: -0.6,
  },
  {
    id: 4,
    tier: 'Base Substrate',
    name: 'Thermal Silicon & Edge-Coupling Interface',
    material: 'Thermal SiO₂ (3 µm) on 775 µm Bulk Silicon Wafer',
    description:
      'High-grade thermal oxide isolating photonic modes from the silicon substrate. Equipped with inverse-taper edge couplers for sub-1.5 dB fiber-to-chip packaging.',
    spec: '≤ 1.5 dB edge-coupling insertion loss • 648 mm² monolithic block • fits 858 mm² reticle',
    color: '#334155', // Slate metal
    wireframeColor: '#64748b',
    baseY: -1.5,
    explodeMultiplier: -1.8,
  },
];

// Individual 3D Layer Mesh in the Exploded Stack
const LayerMesh: React.FC<{
  layer: LayerInfo;
  explodeFactor: number;
  isWireframe: boolean;
  isSelected: boolean;
  onSelect: () => void;
}> = ({ layer, explodeFactor, isWireframe, isSelected, onSelect }) => {
  const meshRef = useRef<THREE.Mesh>(null);
  const targetY = layer.baseY + layer.explodeMultiplier * explodeFactor * 1.5;

  useFrame((_, delta) => {
    if (meshRef.current) {
      meshRef.current.position.y = THREE.MathUtils.damp(meshRef.current.position.y, targetY, 8, delta);
    }
  });

  return (
    <group>
      <mesh
        ref={meshRef}
        position={[0, targetY, 0]}
        onClick={(e) => {
          e.stopPropagation();
          onSelect();
        }}
        castShadow
        receiveShadow
      >
        <boxGeometry args={[5.5, 0.28, 4.2]} />
        <meshStandardMaterial
          color={layer.color}
          wireframe={isWireframe}
          transparent={true}
          opacity={isSelected ? 0.95 : isWireframe ? 0.8 : 0.75}
          roughness={0.2}
          metalness={0.7}
          emissive={isSelected ? layer.color : '#000000'}
          emissiveIntensity={isSelected ? 0.4 : 0.05}
        />

        {/* Decorative waveguide / circuit traces on layer surface */}
        <mesh position={[0, 0.15, 0]} rotation={[-Math.PI / 2, 0, 0]}>
          <planeGeometry args={[5.2, 3.9]} />
          <meshBasicMaterial
            color={isSelected ? '#ffffff' : layer.wireframeColor}
            wireframe={true}
            transparent={true}
            opacity={isSelected ? 0.6 : 0.25}
          />
        </mesh>
      </mesh>
    </group>
  );
};

// 3D Laser Cannon Assembly injecting 1550 nm optical carrier at chip edge
const LaserCannonFacet: React.FC = () => {
  const beamRef = useRef<THREE.Mesh>(null);

  useFrame(({ clock }) => {
    if (beamRef.current) {
      const s = 1 + Math.sin(clock.getElapsedTime() * 8) * 0.15;
      beamRef.current.scale.set(1, s, s);
    }
  });

  return (
    <group position={[-3.8, -0.5, 0]}>
      {/* Metallic Laser Cannon Barrel Housing */}
      <mesh rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[0.22, 0.32, 1.3, 16]} />
        <meshStandardMaterial color="#0f172a" metalness={0.9} roughness={0.15} />
      </mesh>

      {/* Front Optical Lens Ring */}
      <mesh position={[0.66, 0, 0]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[0.18, 0.24, 0.12, 16]} />
        <meshStandardMaterial color="#00f2fe" emissive="#00f2fe" emissiveIntensity={0.8} />
      </mesh>

      {/* Glowing 1550 nm Laser Beam entering the chip edge facet */}
      <mesh ref={beamRef} position={[1.35, 0, 0]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[0.07, 0.07, 1.25, 12]} />
        <meshBasicMaterial color="#00f2fe" transparent opacity={0.85} />
      </mesh>

      {/* Spot Glow at edge coupling interface */}
      <mesh position={[2.0, 0, 0]}>
        <sphereGeometry args={[0.18, 16, 16]} />
        <meshBasicMaterial color="#38bdf8" transparent opacity={0.9} />
      </mesh>

      {/* Floating 3D Label indicating entry point */}
      <Billboard position={[0.7, 0.75, 0]}>
        <Text fontSize={0.22} color="#00f2fe" anchorX="center" anchorY="bottom" outlineWidth={0.02} outlineColor="#04060a">
          LASER CANNON INJECTION PORT
        </Text>
        <Text position={[0, -0.22, 0]} fontSize={0.15} color="#94a3b8" anchorX="center" anchorY="bottom" outlineWidth={0.01} outlineColor="#04060a">
          1550 nm Carrier Entry (≤ 1.5 dB Edge Facet)
        </Text>
      </Billboard>
    </group>
  );
};

// 3D Dual-Rail Photonic Qubit Trajectory & Animation (|0⟩, |1⟩, and |0⟩+|1⟩ Superposition)
const PhotonicQubitTrajectory: React.FC<{
  layer3Y: number;
  qubitState: '0' | '1' | 'superposition';
}> = ({ layer3Y, qubitState }) => {
  const pulseRef0 = useRef<THREE.Mesh>(null);
  const pulseRef1 = useRef<THREE.Mesh>(null);
  const meshGroupRef = useRef<THREE.Group>(null);

  useFrame((state, delta) => {
    if (meshGroupRef.current) {
      meshGroupRef.current.position.y = THREE.MathUtils.damp(
        meshGroupRef.current.position.y,
        layer3Y + 0.16,
        8,
        delta
      );
    }

    // Animation progress (0 to 1 looping)
    const time = (state.clock.getElapsedTime() * 0.9) % 1.0;

    // Pulse 0 position along Rail 0 (from X = -2.6 to X = 2.4, Z = -0.7)
    if (pulseRef0.current) {
      const active0 = qubitState === '0' || qubitState === 'superposition';
      pulseRef0.current.visible = active0;
      if (active0) {
        let x = -2.6 + time * 5.0;
        let z = 0;
        if (x > -1.2) {
          // Diverges to Rail 0 (Z = -0.75)
          const blend = Math.min(1, Math.max(0, (x + 1.2) / 0.8));
          z = -0.75 * blend;
        }
        pulseRef0.current.position.set(x, 0.05, z);
      }
    }

    // Pulse 1 position along Rail 1 (from X = -2.6 to X = 2.4, Z = +0.7)
    if (pulseRef1.current) {
      const active1 = qubitState === '1' || qubitState === 'superposition';
      pulseRef1.current.visible = active1;
      if (active1) {
        let x = -2.6 + time * 5.0;
        let z = 0;
        if (x > -1.2) {
          // Diverges to Rail 1 (Z = +0.75)
          const blend = Math.min(1, Math.max(0, (x + 1.2) / 0.8));
          z = 0.75 * blend;
        }
        pulseRef1.current.position.set(x, 0.05, z);
      }
    }
  });

  const rail0Color = qubitState === '0' || qubitState === 'superposition' ? '#00f2fe' : '#334155';
  const rail1Color = qubitState === '1' || qubitState === 'superposition' ? '#c084fc' : '#334155';

  return (
    <group ref={meshGroupRef} position={[0, layer3Y + 0.16, 0]}>
      {/* Input Waveguide from Edge Facet (X = -2.7 to -1.3, Z = 0) */}
      <mesh position={[-2.0, 0, 0]}>
        <boxGeometry args={[1.4, 0.04, 0.08]} />
        <meshStandardMaterial color="#00f2fe" emissive="#00f2fe" emissiveIntensity={0.6} />
      </mesh>

      {/* 50:50 Directional Coupler / Hadamard Splitter Region (X = -1.3 to -0.5) */}
      <group position={[-0.9, 0, 0]}>
        <mesh position={[0, 0, -0.375]} rotation={[0, -0.55, 0]}>
          <boxGeometry args={[1.0, 0.04, 0.07]} />
          <meshStandardMaterial color={rail0Color} emissive={rail0Color} emissiveIntensity={0.5} />
        </mesh>
        <mesh position={[0, 0, 0.375]} rotation={[0, 0.55, 0]}>
          <boxGeometry args={[1.0, 0.04, 0.07]} />
          <meshStandardMaterial color={rail1Color} emissive={rail1Color} emissiveIntensity={0.5} />
        </mesh>
        {/* Splitter Label */}
        <Billboard position={[0, 0.45, 0]}>
          <Text fontSize={0.16} color="#38bdf8" anchorX="center" anchorY="bottom" outlineWidth={0.015} outlineColor="#04060a">
            50:50 HADAMARD BEAM SPLITTER
          </Text>
        </Billboard>
      </group>

      {/* Rail 0 Waveguide (Path of |0⟩, Z = -0.75, X = -0.4 to 2.4) */}
      <mesh position={[1.0, 0, -0.75]}>
        <boxGeometry args={[2.8, 0.04, 0.07]} />
        <meshStandardMaterial
          color={rail0Color}
          emissive={rail0Color}
          emissiveIntensity={qubitState === '0' || qubitState === 'superposition' ? 0.7 : 0.0}
        />
      </mesh>

      {/* Rail 1 Waveguide (Path of |1⟩, Z = +0.75, X = -0.4 to 2.4) */}
      <mesh position={[1.0, 0, 0.75]}>
        <boxGeometry args={[2.8, 0.04, 0.07]} />
        <meshStandardMaterial
          color={rail1Color}
          emissive={rail1Color}
          emissiveIntensity={qubitState === '1' || qubitState === 'superposition' ? 0.7 : 0.0}
        />
      </mesh>

      {/* Evanescent Coupling Bridge between Rail 0 and Rail 1 */}
      <mesh position={[0.8, 0, 0]}>
        <boxGeometry args={[0.5, 0.03, 1.45]} />
        <meshStandardMaterial
          color="#1e293b"
          transparent
          opacity={0.3}
          wireframe
        />
      </mesh>

      {/* Animated Traveling Photon Packet for Rail 0 */}
      <mesh ref={pulseRef0} position={[-2.6, 0.05, 0]}>
        <sphereGeometry args={[0.13, 16, 16]} />
        <meshBasicMaterial color="#00f2fe" />
      </mesh>

      {/* Animated Traveling Photon Packet for Rail 1 */}
      <mesh ref={pulseRef1} position={[-2.6, 0.05, 0]}>
        <sphereGeometry args={[0.13, 16, 16]} />
        <meshBasicMaterial color="#c084fc" />
      </mesh>

      {/* Output Detectors at Chip Edge */}
      {/* Detector 0 (State |0⟩) */}
      <group position={[2.45, 0, -0.75]}>
        <mesh>
          <cylinderGeometry args={[0.18, 0.18, 0.1, 16]} />
          <meshStandardMaterial
            color={qubitState === '0' || qubitState === 'superposition' ? '#00f2fe' : '#1e293b'}
            emissive={qubitState === '0' || qubitState === 'superposition' ? '#00f2fe' : '#000000'}
            emissiveIntensity={0.8}
          />
        </mesh>
        <Billboard position={[0, 0.35, 0]}>
          <Text fontSize={0.16} color="#00f2fe" anchorX="center" anchorY="bottom" outlineWidth={0.015} outlineColor="#04060a">
            Rail 0: State |0⟩
          </Text>
        </Billboard>
      </group>

      {/* Detector 1 (State |1⟩) */}
      <group position={[2.45, 0, 0.75]}>
        <mesh>
          <cylinderGeometry args={[0.18, 0.18, 0.1, 16]} />
          <meshStandardMaterial
            color={qubitState === '1' || qubitState === 'superposition' ? '#c084fc' : '#1e293b'}
            emissive={qubitState === '1' || qubitState === 'superposition' ? '#c084fc' : '#000000'}
            emissiveIntensity={0.8}
          />
        </mesh>
        <Billboard position={[0, 0.35, 0]}>
          <Text fontSize={0.16} color="#c084fc" anchorX="center" anchorY="bottom" outlineWidth={0.015} outlineColor="#04060a">
            Rail 1: State |1⟩
          </Text>
        </Billboard>
      </group>

      {/* Dynamic Superposition Banner */}
      <Billboard position={[0.5, 0.9, 0]}>
        <Text fontSize={0.2} color="#ffffff" anchorX="center" anchorY="bottom" outlineWidth={0.02} outlineColor="#04060a">
          {qubitState === '0' && '|ψ⟩ = |0⟩  (Photon in Waveguide 0: |1, 0⟩)'}
          {qubitState === '1' && '|ψ⟩ = |1⟩  (Photon in Waveguide 1: |0, 1⟩)'}
          {qubitState === 'superposition' && '|ψ⟩ = (|0⟩ + |1⟩)/√2  (Dual-Rail Superposition: 1 & 0 simultaneously)'}
        </Text>
      </Billboard>
    </group>
  );
};

// 3D Scene Container
const DieScene: React.FC<{
  explodeFactor: number;
  isWireframe: boolean;
  selectedId: number;
  qubitState: '0' | '1' | 'superposition';
  onSelectLayer: (id: number) => void;
}> = ({ explodeFactor, isWireframe, selectedId, qubitState, onSelectLayer }) => {
  // Dynamic Y position of Layer 3 (Waveguide Routing Tier)
  const layer3Y = LAYERS[2].baseY + LAYERS[2].explodeMultiplier * explodeFactor * 1.5;

  return (
    <Canvas
      camera={{ position: [6.5, 5.0, 7.0], fov: 42 }}
      className="w-full h-full cursor-grab active:cursor-grabbing"
    >
      <ambientLight intensity={0.7} />
      <directionalLight position={[10, 15, 10]} intensity={1.5} castShadow />
      <pointLight position={[-10, -5, -5]} intensity={0.8} color="#00f2fe" />
      <pointLight position={[5, 10, -5]} intensity={0.6} color="#f59e0b" />

      <Float speed={1.2} rotationIntensity={0.15} floatIntensity={0.2}>
        <group rotation={[0.1, -0.4, 0]}>
          {/* Laser Cannon Edge-Coupler Injection Facet */}
          <LaserCannonFacet />

          {/* Photonic Qubit Trajectory & Dual-Rail Waveguide Circuit */}
          <PhotonicQubitTrajectory layer3Y={layer3Y} qubitState={qubitState} />

          {LAYERS.map((layer) => (
            <LayerMesh
              key={layer.id}
              layer={layer}
              explodeFactor={explodeFactor}
              isWireframe={isWireframe}
              isSelected={selectedId === layer.id}
              onSelect={() => onSelectLayer(layer.id)}
            />
          ))}

          {/* Central optical optical interconnect vias linking the layers */}
          {explodeFactor > 0.1 && (
            <group>
              {[-1.8, 0, 1.8].map((x) =>
                [-1.2, 0, 1.2].map((z) => (
                  <mesh key={`via-${x}-${z}`} position={[x, 0, z]}>
                    <cylinderGeometry args={[0.04, 0.04, (LAYERS[0].baseY - LAYERS[3].baseY) + explodeFactor * 5.4, 8]} />
                    <meshBasicMaterial color="#00f2fe" transparent opacity={0.4} />
                  </mesh>
                ))
              )}
            </group>
          )}
        </group>
      </Float>

      <OrbitControls enablePan={false} minDistance={5} maxDistance={15} maxPolarAngle={Math.PI / 1.7} />
    </Canvas>
  );
};

export const Section1DieArchitecture: React.FC = () => {
  const [explodeValue, setExplodeValue] = useState<number>(45);
  const [isWireframe, setIsWireframe] = useState<boolean>(false);
  const [selectedId, setSelectedId] = useState<number>(3); // Si3N4 active by default
  const [qubitState, setQubitState] = useState<'0' | '1' | 'superposition'>('superposition');

  const selectedLayer = LAYERS.find((l) => l.id === selectedId) || LAYERS[2];

  return (
    <section id="section-die" className="relative py-24 border-t border-slate-800/80 bg-[#060911]/90">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="flex flex-col md:flex-row md:items-end justify-between mb-8">
          <div>
            <div className="flex items-center space-x-2 text-cyan-400 font-mono text-xs uppercase tracking-wider mb-2">
              <Layers className="w-4 h-4" />
              <span>Section 1 • Physical Monolithic Integration</span>
            </div>
            <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
              Die Architecture &amp; Photonic Qubit Trajectory
            </h2>
            <p className="mt-3 text-base text-slate-300 max-w-2xl">
              SilicaCore integrates classic race-logic spirals with linear optical quantum computing (LOQC). Observe the laser cannon injection and the dual-rail photonic qubit trajectories for states |0⟩, |1⟩, and quantum superposition |0⟩+|1⟩.
            </p>
          </div>

          <div className="mt-6 md:mt-0 flex items-center space-x-3 text-xs font-mono text-slate-400">
            <span className="w-2.5 h-2.5 rounded-full bg-cyan-400 animate-pulse"></span>
            <span>Commercial 858 mm² Reticle Compliant</span>
          </div>
        </div>

        {/* Photonic Qubit State Selector Toolbar */}
        <div className="mb-6 p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center space-x-2 text-xs font-mono text-slate-300">
            <Atom className="w-4 h-4 text-cyan-400" />
            <span className="font-semibold text-white">Photonic Qubit Trajectory:</span>
            <span className="text-slate-500">Select quantum state to animate beam path</span>
          </div>

          <div className="flex items-center space-x-2">
            <button
              onClick={() => setQubitState('0')}
              className={`px-3 py-1.5 rounded-lg text-xs font-mono font-medium transition-all ${
                qubitState === '0'
                  ? 'bg-cyan-500 text-slate-950 font-bold shadow-[0_0_12px_rgba(0,242,254,0.4)]'
                  : 'bg-slate-950 text-slate-400 hover:text-white border border-slate-800'
              }`}
            >
              State |0⟩ (Rail 0)
            </button>

            <button
              onClick={() => setQubitState('1')}
              className={`px-3 py-1.5 rounded-lg text-xs font-mono font-medium transition-all ${
                qubitState === '1'
                  ? 'bg-purple-500 text-white font-bold shadow-[0_0_12px_rgba(192,132,252,0.4)]'
                  : 'bg-slate-950 text-slate-400 hover:text-white border border-slate-800'
              }`}
            >
              State |1⟩ (Rail 1)
            </button>

            <button
              onClick={() => setQubitState('superposition')}
              className={`px-3.5 py-1.5 rounded-lg text-xs font-mono font-bold transition-all ${
                qubitState === 'superposition'
                  ? 'bg-gradient-to-r from-cyan-400 to-purple-400 text-slate-950 shadow-[0_0_15px_rgba(34,211,238,0.5)]'
                  : 'bg-slate-950 text-slate-300 hover:text-white border border-slate-800'
              }`}
            >
              Superposition |0⟩ + |1⟩ (1 &amp; 0 simultaneously)
            </button>
          </div>
        </div>

        {/* Main Content Grid: 3D Canvas Left/Center, Interactive Layer Inspector Right */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* 3D Interactive Viewport */}
          <div className="lg:col-span-7 flex flex-col space-y-4">
            <div className="relative h-[480px] sm:h-[560px] rounded-2xl bg-gradient-to-b from-[#090e1a] to-[#04060a] border border-cyan-500/20 shadow-2xl overflow-hidden">
              {/* Overlay Badges */}
              <div className="absolute top-4 left-4 z-10 flex flex-wrap gap-2">
                <span className="px-2.5 py-1 rounded-md bg-slate-900/80 backdrop-blur-md border border-slate-700/60 text-[11px] font-mono text-cyan-300 flex items-center gap-1.5">
                  <Box className="w-3.5 h-3.5 text-cyan-400" />
                  3D Exploded Viewport &amp; Laser Trajectory
                </span>
                <span className="px-2.5 py-1 rounded-md bg-slate-900/80 backdrop-blur-md border border-slate-700/60 text-[11px] font-mono text-slate-400">
                  Drag to rotate • Scroll to zoom
                </span>
              </div>

              {/* 3D Canvas */}
              <DieScene
                explodeFactor={explodeValue / 100}
                isWireframe={isWireframe}
                selectedId={selectedId}
                qubitState={qubitState}
                onSelectLayer={setSelectedId}
              />

              {/* In-viewport Controls Bar */}
              <div className="absolute bottom-4 left-4 right-4 z-10 glass-panel p-3.5 rounded-xl flex flex-wrap items-center justify-between gap-4 border border-slate-700/60">
                {/* Explode Slider */}
                <div className="flex items-center space-x-3 flex-1 min-w-[200px]">
                  <Sliders className="w-4 h-4 text-cyan-400" />
                  <span className="text-xs font-mono text-slate-300 whitespace-nowrap">
                    Explode: <span className="text-cyan-400 font-bold">{explodeValue}%</span>
                  </span>
                  <input
                    type="range"
                    min="0"
                    max="100"
                    value={explodeValue}
                    onChange={(e) => setExplodeValue(Number(e.target.value))}
                    className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-400"
                  />
                </div>

                {/* Wireframe Toggle */}
                <button
                  onClick={() => setIsWireframe(!isWireframe)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-mono font-medium transition-all flex items-center space-x-1.5 ${
                    isWireframe
                      ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/50'
                      : 'bg-slate-800 text-slate-400 hover:text-white border border-slate-700'
                  }`}
                >
                  <Eye className="w-3.5 h-3.5" />
                  <span>{isWireframe ? 'Wireframe ON' : 'Physical Render'}</span>
                </button>
              </div>
            </div>

            {/* LOQC Dual-Rail Quantum Telemetry Card */}
            <div className="glass-panel p-4 rounded-xl border border-cyan-500/30 text-xs font-mono space-y-2">
              <div className="flex items-center justify-between text-cyan-300 font-semibold pb-2 border-b border-slate-800">
                <span className="flex items-center gap-1.5">
                  <Atom className="w-4 h-4 text-cyan-400" />
                  <span>Dual-Rail Encoding Trajectory Physics</span>
                </span>
                <span className="text-[10px] text-slate-400">LOQC in Si₃N₄ / TFLN</span>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-2 pt-1 text-[11px]">
                <div className={`p-2 rounded-lg border ${qubitState === '0' ? 'bg-cyan-950/60 border-cyan-400 text-cyan-200' : 'bg-slate-950/60 border-slate-800 text-slate-400'}`}>
                  <strong>State |0⟩ = |1, 0⟩:</strong> Single photon traverses Rail 0 exclusively. Rail 1 in vacuum.
                </div>
                <div className={`p-2 rounded-lg border ${qubitState === '1' ? 'bg-purple-950/60 border-purple-400 text-purple-200' : 'bg-slate-950/60 border-slate-800 text-slate-400'}`}>
                  <strong>State |1⟩ = |0, 1⟩:</strong> TFLN electro-optic gate modulates photon into Rail 1.
                </div>
                <div className={`p-2 rounded-lg border ${qubitState === 'superposition' ? 'bg-gradient-to-br from-cyan-950 to-purple-950 border-cyan-400 text-cyan-200' : 'bg-slate-950/60 border-slate-800 text-slate-400'}`}>
                  <strong>Superposition (1 &amp; 0):</strong> (|0⟩+|1⟩)/√2 via 50:50 Hadamard beam splitter.
                </div>
              </div>
            </div>
          </div>

          {/* Right Column: Layer Breakdown Details */}
          <div className="lg:col-span-5 flex flex-col space-y-3">
            <h3 className="text-sm font-mono uppercase tracking-wider text-slate-400 mb-1 flex items-center gap-2">
              <span>Stack Architecture Hierarchy</span>
              <span className="text-slate-600">/</span>
              <span className="text-cyan-400">4 Monolithic Tiers</span>
            </h3>

            {LAYERS.map((layer) => {
              const isSelected = selectedId === layer.id;
              return (
                <div
                  key={layer.id}
                  onClick={() => setSelectedId(layer.id)}
                  className={`cursor-pointer p-4 rounded-xl transition-all border ${
                    isSelected
                      ? 'glass-panel border-cyan-500 shadow-[0_0_20px_rgba(6,182,212,0.15)] ring-1 ring-cyan-500/30'
                      : 'bg-slate-900/40 hover:bg-slate-900/80 border-slate-800/80'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <div className="flex items-center space-x-2 mb-1">
                        <span
                          className="w-2.5 h-2.5 rounded-full"
                          style={{ backgroundColor: layer.color }}
                        />
                        <span className="text-[11px] font-mono font-semibold uppercase tracking-wider text-slate-400">
                          {layer.tier}
                        </span>
                        {isSelected && (
                          <span className="px-1.5 py-0.2 rounded bg-cyan-950 text-cyan-300 text-[10px] font-mono border border-cyan-500/30">
                            Inspecting
                          </span>
                        )}
                      </div>
                      <h4 className="text-base font-semibold text-white">{layer.name}</h4>
                      <p className="text-xs font-mono text-cyan-300/90 mt-0.5">{layer.material}</p>
                    </div>

                    <span className="text-xs font-mono text-slate-500">#{layer.id}</span>
                  </div>

                  {isSelected && (
                    <div className="mt-3 pt-3 border-t border-slate-800 text-xs text-slate-300 space-y-2">
                      <p className="leading-relaxed">{layer.description}</p>
                      <div className="p-2.5 rounded-lg bg-slate-950/70 border border-cyan-500/20 text-cyan-200 font-mono text-[11px] flex items-start gap-2">
                        <CheckCircle2 className="w-3.5 h-3.5 text-cyan-400 mt-0.5 flex-shrink-0" />
                        <span>{layer.spec}</span>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </section>
  );
};
