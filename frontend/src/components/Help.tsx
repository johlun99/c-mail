import { useEffect } from 'react';
import { KEYMAP } from '../lib/keymap';

interface HelpProps {
  onClose: () => void;
}

export function Help({ onClose }: HelpProps) {
  useEffect(() => {
    const h = (e: KeyboardEvent) => {
      if (e.key === 'Escape' || e.key === '?') onClose();
    };
    window.addEventListener('keydown', h, true);
    return () => window.removeEventListener('keydown', h, true);
  }, [onClose]);

  const cols = [KEYMAP.slice(0, 2), KEYMAP.slice(2)];
  return (
    <div className="scrim" onClick={onClose}>
      <div
        className="sheet"
        onClick={(e) => e.stopPropagation()}
        style={{ width: 'min(860px, 95vw)' }}
      >
        <div className="sheet-hd">
          <span className="pr">?</span>
          <span className="t">kortkommandon</span>
          <span className="sp" />
          <span className="muted">vim-first</span>
        </div>
        <div className="sheet-body">
          <div className="help-grid">
            {cols.map((col, ci) => (
              <div className="help-col" key={ci}>
                {col.map(([h, rows]) => (
                  <div key={h}>
                    <h4>{h}</h4>
                    {rows.map(([k, d]) => (
                      <div className="help-row" key={k}>
                        <span className="desc">{d}</span>
                        <span className="keys">
                          {k.split(' / ').map((kk) => (
                            <kbd key={kk}>{kk}</kbd>
                          ))}
                        </span>
                      </div>
                    ))}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
