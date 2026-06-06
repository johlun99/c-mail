import { Fragment, useEffect, useRef, useState, type KeyboardEvent } from 'react';

export interface Command {
  id: string;
  name: string;
  group: string;
  ic?: string;
  kbd?: string;
  run: () => void;
}

interface CommandPaletteProps {
  commands: Command[];
  onClose: () => void;
  onRun: (c: Command) => void;
}

export function CommandPalette({ commands, onClose, onRun }: CommandPaletteProps) {
  const [q, setQ] = useState('');
  const [i, setI] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const filtered = commands.filter((c) =>
    `${c.name} ${c.group}`.toLowerCase().includes(q.toLowerCase()),
  );

  useEffect(() => {
    inputRef.current?.focus();
  }, []);
  useEffect(() => {
    setI(0);
  }, [q]);

  const onKey = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'ArrowDown' || (e.ctrlKey && e.key === 'n')) {
      e.preventDefault();
      setI((x) => Math.min(x + 1, filtered.length - 1));
    } else if (e.key === 'ArrowUp' || (e.ctrlKey && e.key === 'p')) {
      e.preventDefault();
      setI((x) => Math.max(x - 1, 0));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      const c = filtered[i];
      if (c) onRun(c);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
    }
  };

  let lastGrp: string | null = null;
  return (
    <div className="scrim" onClick={onClose}>
      <div className="palette" onClick={(e) => e.stopPropagation()}>
        <div className="pal-in">
          <span className="pr">&gt;</span>
          <input
            ref={inputRef}
            value={q}
            placeholder="skriv ett kommando…"
            onChange={(e) => setQ(e.target.value)}
            onKeyDown={onKey}
          />
          <span className="muted">{filtered.length}</span>
        </div>
        <div className="pal-list">
          {filtered.length === 0 && <div className="pal-grp">inga kommandon</div>}
          {filtered.map((c, idx) => {
            const head = c.group && c.group !== lastGrp ? c.group : null;
            lastGrp = c.group;
            return (
              <Fragment key={c.id}>
                {head && <div className="pal-grp">{head}</div>}
                <div
                  className={'pal-item' + (idx === i ? ' sel' : '')}
                  onMouseEnter={() => setI(idx)}
                  onClick={() => onRun(c)}
                >
                  <span className="ic">{c.ic ?? '›'}</span>
                  <span className="nm">{c.name}</span>
                  {c.kbd && (
                    <span className="kb">
                      <kbd>{c.kbd}</kbd>
                    </span>
                  )}
                </div>
              </Fragment>
            );
          })}
        </div>
      </div>
    </div>
  );
}
