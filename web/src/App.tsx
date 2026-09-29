import React, { useState, useEffect } from 'react';
import { Navbar } from './components/Navbar';
import { HeroSection } from './components/HeroSection';
import { HowItWorksSection } from './components/HowItWorksSection';
import { InteractiveChipSection } from './components/InteractiveChipSection';
import { SectionLLMBenchmark } from './components/SectionLLMBenchmark';
import { SectionEnergyChart } from './components/SectionEnergyChart';
import { SectionAdvantages } from './components/SectionAdvantages';
import { Footer } from './components/Footer';

export const App: React.FC = () => {
  const [activeSection, setActiveSection] = useState<string>('section-how-it-works');

  // Intersection observer / scroll spy for active navbar tab
  useEffect(() => {
    const sectionIds = [
      'section-how-it-works',
      'section-chip',
      'section-llm',
      'section-energy',
      'section-advantages',
    ];

    const handleScroll = () => {
      const scrollPosition = window.scrollY + 250;

      for (const id of sectionIds) {
        const el = document.getElementById(id);
        if (el) {
          const top = el.offsetTop;
          const height = el.offsetHeight;
          if (scrollPosition >= top && scrollPosition < top + height) {
            setActiveSection(id);
            break;
          }
        }
      }
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const handleLaunchSim = () => {
    const el = document.getElementById('section-chip');
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div className="min-h-screen bg-black text-slate-100 flex flex-col font-sans selection:bg-cyan-500 selection:text-slate-950">
      {/* Apple-Style Floating Pill Navbar */}
      <Navbar activeSection={activeSection} />

      {/* Main Product Presentation Flow */}
      <main className="flex-1">
        {/* 1. Cinematic Apple-Style Hero */}
        <HeroSection onLaunchSim={handleLaunchSim} />

        {/* 2. Visual 3-Step "How It Works" + Photonic Qubit */}
        <HowItWorksSection />

        {/* 3. Live Processor Benchmark Arena (SilicaCore vs Apple M5 Max) + 3D Die */}
        <InteractiveChipSection />

        {/* 4. AI & Large Language Models (LLM) Acceleration Arena */}
        <SectionLLMBenchmark />

        {/* 5. The 3 Comparative Graphs: Speed, TDP Power & Heat (Apple M5 Max, Intel, AMD, NVIDIA) */}
        <SectionEnergyChart />

        {/* 6. The Investor & Commercial Advantages (Bento Grid) */}
        <SectionAdvantages />
      </main>

      {/* Apple-Style Minimal Footer */}
      <Footer />
    </div>
  );
};

export default App;
