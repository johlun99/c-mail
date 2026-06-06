import type { CSSProperties } from 'react';
import { mail } from '../../wailsjs/go/models';
import { catBy, colorVar } from '../lib/categories';
import { MOD, ENTER } from '../lib/platform';

interface ReadingPaneProps {
  m: mail.Mail | null;
  cats: mail.Category[];
  draftVisible: boolean;
  sent: boolean;
  editing: boolean;
  editValue: string;
  onEditChange: (v: string) => void;
  onStartEdit: (v: boolean) => void;
  onApprove: () => void;
  onRegenerate: () => void;
  onDiscard: () => void;
}

export function ReadingPane(props: ReadingPaneProps) {
  const { m, cats, draftVisible, sent } = props;
  if (!m) {
    return (
      <div className="read">
        <div className="read-empty">
          <div className="acc" style={{ fontSize: '1.5em' }}>
            ~/cmail
          </div>
          <div>
            välj ett mail · <kbd>j</kbd>/<kbd>k</kbd> för att navigera
          </div>
        </div>
      </div>
    );
  }
  const c = catBy(cats, m.cat);
  const style = { '--catc': colorVar(c) } as CSSProperties;
  return (
    <div className="read" style={style}>
      <div className="read-scroll">
        <div className="read-hd">
          <div className="read-subj">{m.subject}</div>
          <div className="read-meta">
            <div className="av">{m.avatar}</div>
            <div className="read-from">
              <span className="nm">{m.from}</span>
              <span className="ad">{m.fromAddr}</span>
            </div>
            <span className="sp" />
            <div className="read-meta-r">
              <span className="tag">{c.label}</span>
            </div>
            <div className="tm">
              {m.day} {m.time}
              <br />
              {m.threadCount > 1 && <span className="muted">{m.threadCount} i tråden</span>}
            </div>
          </div>
        </div>

        <AgentBlock m={m} />

        {m.draft &&
          (sent ? <SentBanner m={m} /> : draftVisible ? <DraftBlock {...props} m={m} /> : null)}

        <div className="read-body">
          {m.body.map((p, i) => (
            <p key={i}>{p}</p>
          ))}
        </div>
      </div>
    </div>
  );
}

function AgentBlock({ m }: { m: mail.Mail }) {
  const a = m.agent;
  const maxK = a.facts.reduce((n, f) => Math.max(n, f.key.length), 0);
  return (
    <div className="agent">
      <div className="agent-hd">
        <span className="pr">agent &gt;</span>
        <span className="lbl">analyserade mailet</span>
        <span className="sp" />
        <span className="conf">konfidens {m.confidence.toFixed(2)}</span>
      </div>
      <div className="agent-body">
        <div className="sum">{a.summary}</div>
        <div className="tree">
          {a.facts.map((f, i) => {
            const last = i === a.facts.length - 1;
            const warn = /⚠/.test(f.value);
            return (
              <div key={i}>
                <span className="br">{last ? ' └ ' : ' ├ '}</span>
                <span className="key">{f.key.padEnd(maxK, ' ')}</span>
                <span className="br">{'  ·  '}</span>
                <span className={warn ? 'warn' : 'val'}>{f.value}</span>
              </div>
            );
          })}
        </div>
        {a.tasks.length > 0 && (
          <div className="agent-acts">
            {a.tasks.map((t, i) => (
              <span key={i} className="chip task">
                <span className="ck">+</span>
                {t}
              </span>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

interface DiffRow {
  text: string;
  kind: string;
  mod: boolean;
}

function DraftBlock({
  m,
  editing,
  editValue,
  onEditChange,
  onStartEdit,
  onApprove,
  onRegenerate,
  onDiscard,
}: ReadingPaneProps & { m: mail.Mail }) {
  const draft = m.draft!;
  const orig = draft.lines.map((l) => l.text);
  const toName = m.from.split(' ')[0];

  let rows: DiffRow[];
  if (editing) {
    rows = editValue.split('\n').map((text, i) => {
      const isMod = i >= orig.length || text !== orig[i];
      return { text, kind: isMod ? '~' : '+', mod: isMod };
    });
  } else {
    rows = orig.map((text) => ({ text, kind: '+', mod: false }));
  }

  return (
    <div className="draft">
      <div className="draft-hd">
        <span className="pr">utkast &gt;</span>
        <span className="to">svar till {toName}</span>
        <span className="tone muted">· {draft.tone}</span>
        <span className="sp" />
        <span className="badge">{editing ? 'redigerar' : 'väntar godkännande'}</span>
      </div>

      {editing ? (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 0 }}>
          <textarea
            className="draft-edit"
            autoFocus
            value={editValue}
            onChange={(e) => onEditChange(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                e.stopPropagation();
                onStartEdit(false);
              }
            }}
          />
          <div className="diff" style={{ borderLeft: '1px solid var(--border)' }}>
            {rows.map((r, i) => (
              <DiffLine key={i} r={r} />
            ))}
          </div>
        </div>
      ) : (
        <div className="diff">
          {rows.map((r, i) => (
            <DiffLine key={i} r={r} />
          ))}
        </div>
      )}

      <div className="draft-foot">
        <button className="btn primary" onClick={onApprove}>
          <span>godkänn & skicka</span>
          <kbd>{MOD}</kbd>
          <kbd>{ENTER}</kbd>
        </button>
        <button className="btn" onClick={() => onStartEdit(!editing)}>
          {editing ? 'klar' : 'redigera'} <kbd>e</kbd>
        </button>
        <button className="btn" onClick={onRegenerate}>
          generera om <kbd>r</kbd>
        </button>
        <span className="sp" />
        <button className="btn" onClick={onDiscard}>
          släng <kbd>d</kbd>
        </button>
      </div>
    </div>
  );
}

function DiffLine({ r }: { r: DiffRow }) {
  const empty = r.text.length === 0;
  return (
    <div className={'dl ' + (r.mod ? 'mod' : 'add') + (empty ? ' empty' : '')}>
      <span className="sign">{r.kind}</span>
      <span className="txt">{r.text || ' '}</span>
    </div>
  );
}

function SentBanner({ m }: { m: mail.Mail }) {
  const draft = m.draft!;
  const toName = m.from.split(' ')[0];
  return (
    <div className="draft" style={{ boxShadow: 'none' }}>
      <div className="draft-hd">
        <span className="pr" style={{ color: 'var(--accent)' }}>
          ↗
        </span>
        <span className="to">svar skickat till {toName}</span>
        <span className="sp" />
        <span className="badge">✓ skickat</span>
      </div>
      <div className="diff">
        {draft.lines.map((l, i) => (
          <div className="dl add empty" key={i} style={{ opacity: 0.55 }}>
            <span className="sign">✓</span>
            <span className="txt">{l.text || ' '}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
