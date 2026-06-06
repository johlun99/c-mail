import type { CSSProperties } from 'react';
import { mail } from '../../wailsjs/go/models';
import { MOD } from '../lib/platform';

interface RailProps {
  cats: mail.Category[];
  counts: Record<string, number>;
  current: string;
  onPick: (key: string) => void;
  onHelp: () => void;
  onPalette: () => void;
}

export function Rail({ cats, counts, current, onPick, onHelp, onPalette }: RailProps) {
  return (
    <div className="rail">
      <div className="rail-sec">vyer</div>
      <RailItem
        label="alla"
        color="c-dim"
        count={counts.alla ?? 0}
        sel={current === 'alla'}
        onPick={() => onPick('alla')}
        forceDot="◆"
      />
      <div className="rail-sec">kategorier</div>
      {cats.map((c) => (
        <RailItem
          key={c.key}
          label={c.label}
          color={c.color}
          count={counts[c.key] ?? 0}
          sel={current === c.key}
          onPick={() => onPick(c.key)}
        />
      ))}
      <div className="spring" />
      <div className="rail-foot">
        <div className="k">agent har sorterat {counts.alla ?? 0} mail</div>
        <div>
          <span className="lk" onClick={onHelp}>
            <kbd>?</kbd> kortkommandon
          </span>{' '}
          ·{' '}
          <span className="lk" onClick={onPalette}>
            <kbd>{MOD}</kbd>
            <kbd>K</kbd> palett
          </span>
        </div>
      </div>
    </div>
  );
}

interface RailItemProps {
  label: string;
  color: string;
  count: number;
  sel: boolean;
  onPick: () => void;
  forceDot?: string;
}

function RailItem({ label, color, count, sel, onPick, forceDot }: RailItemProps) {
  const style = { '--catc': `var(--${color})` } as CSSProperties;
  return (
    <div className={'cat' + (sel ? ' sel' : '')} onClick={onPick} style={style}>
      {forceDot ? (
        <span
          className="dot"
          style={{
            background: 'transparent',
            color: 'var(--tx-3)',
            width: 'auto',
            boxShadow: 'none',
          }}
        >
          {forceDot}
        </span>
      ) : (
        <span className="dot" />
      )}
      <span className="nm">{label}</span>
      <span className="ct">{count}</span>
    </div>
  );
}
