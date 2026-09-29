import React, { useState, useEffect } from 'react';
import { Cpu, Github, FileText, ArrowRight, Sparkles } from 'lucide-react';

interface NavbarProps {
  activeSection: string;
}

export const Navbar: React.FC<NavbarProps> = ({ activeSection }) => {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 20);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const navLinks = [
    { id: 'section-how-it-works', label: 'How It Works' },
    { id: 'section-decision', label: 'Decision' },
    { id: 'section-chip', label: 'The Chip' },
    { id: 'section-compare', label: 'Comparisons' },
    { id: 'section-maxwell', label: 'Waveguide Physics' },
    { id: 'section-energy', label: 'Energy' },
    { id: 'section-status', label: 'Status' },
  ];

  const scrollTo = (id: string) => {
    const el = document.getElementById(id);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex justify-center px-4 pt-4 transition-all duration-300">
      <div
        className={`w-full max-w-5xl rounded-full transition-all duration-300 px-5 py-2.5 flex items-center justify-between border ${
          scrolled
            ? 'bg-black/80 backdrop-blur-2xl border-white/15 shadow-[0_8px_32px_rgba(0,0,0,0.8)]'
            : 'bg-black/50 backdrop-blur-xl border-white/10 shadow-lg'
        }`}
      >
        {/* Brand Logo & Name */}
        <div
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          className="flex items-center space-x-2.5 cursor-pointer group"
        >
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-cyan-500 to-blue-600 p-[1px] shadow-[0_0_12px_rgba(6,182,212,0.4)]">
            <div className="w-full h-full bg-[#050811] rounded-full flex items-center justify-center">
              <Cpu className="w-4 h-4 text-cyan-400 group-hover:scale-110 transition-transform" />
            </div>
          </div>

          <div className="flex items-center space-x-1.5">
            <span className="text-sm font-semibold tracking-tight text-white font-sans">
              Silica<span className="text-cyan-400 font-bold">Core</span>
            </span>
            <span className="hidden sm:inline-block px-1.5 py-0.5 text-[9px] uppercase font-mono tracking-widest bg-cyan-950/80 border border-cyan-500/40 text-cyan-300 rounded-full">
              Photonic ToF
            </span>
          </div>
        </div>

        {/* Center Navigation Links - Apple minimalist style */}
        <nav className="hidden md:flex items-center space-x-1">
          {navLinks.map((item) => {
            const isActive = activeSection === item.id;
            return (
              <button
                key={item.id}
                onClick={() => scrollTo(item.id)}
                className={`px-3.5 py-1.5 rounded-full text-xs font-medium transition-all ${
                  isActive
                    ? 'text-white bg-white/10 shadow-sm'
                    : 'text-slate-400 hover:text-white hover:bg-white/5'
                }`}
              >
                {item.label}
              </button>
            );
          })}
        </nav>

        {/* Right Action Buttons */}
        <div className="flex items-center space-x-2">
          <a
            href="https://github.com/Roliveira96/SilicaCore"
            target="_blank"
            rel="noreferrer"
            className="p-2 rounded-full text-slate-400 hover:text-white hover:bg-white/10 transition-colors"
            title="GitHub Repository"
          >
            <Github className="w-4 h-4" />
          </a>

          <button
            onClick={() => scrollTo('section-chip')}
            className="group px-4 py-1.5 rounded-full text-xs font-medium text-slate-950 bg-gradient-to-r from-cyan-400 to-blue-400 hover:from-cyan-300 hover:to-blue-300 shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all flex items-center space-x-1.5"
          >
            <span>See the chip</span>
            <ArrowRight className="w-3 h-3 group-hover:translate-x-0.5 transition-transform" />
          </button>
        </div>
      </div>
    </header>
  );
};
