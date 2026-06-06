import { useCallback, useEffect, useState } from 'react';

export type Font = 'mono' | 'sans';
export type Density = 'tät' | 'normal' | 'luftig';
export type Layout = '2 kol' | '3 kol';

export interface Tweaks {
  accent: string;
  font: Font;
  density: Density;
  layout: Layout;
  bgl: number;
}

export const DEFAULT_TWEAKS: Tweaks = {
  accent: '#7dd3a8',
  font: 'mono',
  density: 'normal',
  layout: '3 kol',
  bgl: 0.165,
};

// Ink (text-on-accent) color per accent swatch.
const ACCENT_INK: Record<string, string> = {
  '#7dd3a8': '#06150e', // grön
  '#6aa3f5': '#04101f', // blå
  '#c9a4f0': '#1a0d2b', // lila
  '#e8995c': '#1d0f03', // amber
};

const DENSITY: Record<Density, { fs: string; row: string; pad: string; lh: number }> = {
  tät: { fs: '12px', row: '26px', pad: '9px', lh: 1.45 },
  normal: { fs: '13px', row: '30px', pad: '11px', lh: 1.55 },
  luftig: { fs: '14px', row: '38px', pad: '15px', lh: 1.7 },
};

const STORAGE_KEY = 'cmail.tweaks';

function load(): Tweaks {
  if (typeof localStorage === 'undefined') return DEFAULT_TWEAKS;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? { ...DEFAULT_TWEAKS, ...(JSON.parse(raw) as Partial<Tweaks>) } : DEFAULT_TWEAKS;
  } catch {
    return DEFAULT_TWEAKS;
  }
}

/**
 * Appearance preferences applied to :root as CSS custom properties — the same
 * single-knob model as the design prototype (accent, font, density, --bgl).
 */
export function useTweaks(): [Tweaks, <K extends keyof Tweaks>(key: K, value: Tweaks[K]) => void] {
  const [tweaks, setTweaks] = useState<Tweaks>(load);

  useEffect(() => {
    const root = document.documentElement;
    root.style.setProperty('--accent', tweaks.accent);
    root.style.setProperty('--accent-ink', ACCENT_INK[tweaks.accent] ?? '#06150e');
    root.style.setProperty('--font', tweaks.font === 'sans' ? 'var(--sans)' : 'var(--mono)');
    root.style.setProperty('--bgl', String(tweaks.bgl));
    const d = DENSITY[tweaks.density];
    root.style.setProperty('--fs', d.fs);
    root.style.setProperty('--row', d.row);
    root.style.setProperty('--pad', d.pad);
    root.style.setProperty('--lh', String(d.lh));
    try {
      localStorage?.setItem(STORAGE_KEY, JSON.stringify(tweaks));
    } catch {
      // ignore persistence errors (e.g. private mode)
    }
  }, [tweaks]);

  const setTweak = useCallback(<K extends keyof Tweaks>(key: K, value: Tweaks[K]) => {
    setTweaks((t) => ({ ...t, [key]: value }));
  }, []);

  return [tweaks, setTweak];
}
