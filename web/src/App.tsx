import React, { useState, useEffect } from 'react';
import { Navbar } from './components/Navbar';
import { HeroSection } from './components/HeroSection';
import { HowItWorksSection } from './components/HowItWorksSection';
import { InteractiveChipSection } from './components/InteractiveChipSection';
import { SectionComparison } from './components/SectionComparison';
import { Section3MaxwellLab } from './components/Section3MaxwellLab';
import { SectionEnergy } from './components/SectionEnergy';
import { SectionStatus } from './components/SectionStatus';
import { Footer } from './components/Footer';

export const App: React.FC = () => {
  const [activeSection, setActiveSection] = useState<string>('section-how-it-works');

  // Intersection observer / scroll spy for active navbar tab
  useEffect(() => {
    const sectionIds = [
      'section-how-it-works',
      'section-chip',
      'section-compare',
      'section-maxwell',
      'section-energy',
      'section-status',
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

        {/* 2. How one query works, plus the long-term photonic qubit concept */}
        <HowItWorksSection />

        {/* 3. Physically scaled 3D chip with the simulated race, and race vs. measured Dijkstra */}
        <InteractiveChipSection />

        {/* 4. Sourced comparisons: speed, energy per bit and other photonic processors */}
        <SectionComparison />

        {/* 5. 2D FDTD results for bends and the Si3N4 -> TFLN taper */}
        <Section3MaxwellLab />

        {/* 6. Energy: where the power goes and what is still undetermined */}
        <SectionEnergy />

        {/* 7. Research status: validated results, assumptions, limits and next steps */}
        <SectionStatus />
      </main>

      {/* Apple-Style Minimal Footer */}
      <Footer />
    </div>
  );
};

export default App;
