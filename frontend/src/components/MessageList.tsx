import { useEffect, useRef, type CSSProperties } from 'react';
import { mail } from '../../wailsjs/go/models';
import { catBy, colorVar } from '../lib/categories';
import { ENTER } from '../lib/platform';

interface MessageListProps {
  mails: mail.Mail[];
  cats: mail.Category[];
  selId: string;
  onSelect: (id: string) => void;
  filterLabel: string;
  searchMode: boolean;
  query: string;
}

export function MessageList({
  mails,
  cats,
  selId,
  onSelect,
  filterLabel,
  searchMode,
  query,
}: MessageListProps) {
  const listRef = useRef<HTMLDivElement>(null);

  // Keep the selected row in view — manual scroll math, never scrollIntoView.
  useEffect(() => {
    const c = listRef.current;
    const el = c?.querySelector<HTMLElement>(`[data-id="${selId}"]`);
    if (!el || !c) return;
    const r = el.getBoundingClientRect();
    const cr = c.getBoundingClientRect();
    if (r.top < cr.top + 40) c.scrollTop -= cr.top + 40 - r.top;
    else if (r.bottom > cr.bottom - 8) c.scrollTop += r.bottom - cr.bottom + 8;
  }, [selId, mails]);

  // Group by day label, preserving order.
  const groups: { day: string; items: mail.Mail[] }[] = [];
  for (const m of mails) {
    let g = groups[groups.length - 1];
    if (!g || g.day !== m.day) {
      g = { day: m.day, items: [] };
      groups.push(g);
    }
    g.items.push(m);
  }

  return (
    <div className="list-wrap">
      <div className="list-head">
        <span className="title">{filterLabel}</span>
        <span className="count">{mails.length}</span>
        <span className="sp" />
        {searchMode ? (
          <span className="srch">
            /{query}
            <span className="cur" style={{ height: 11, width: 6 }} />
          </span>
        ) : (
          <span className="muted">↑↓/jk · {ENTER} öppna</span>
        )}
      </div>
      <div className="list" ref={listRef}>
        {mails.length === 0 && <div className="daygrp">inga träffar</div>}
        {groups.map((g) => (
          <div key={g.day}>
            <div className="daygrp">{g.day}</div>
            {g.items.map((m) => (
              <Row
                key={m.id}
                m={m}
                cat={catBy(cats, m.cat)}
                sel={m.id === selId}
                onSelect={onSelect}
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

interface RowProps {
  m: mail.Mail;
  cat: mail.Category;
  sel: boolean;
  onSelect: (id: string) => void;
}

function Row({ m, cat, sel, onSelect }: RowProps) {
  const style = { '--catc': colorVar(cat) } as CSSProperties;
  return (
    <div
      data-id={m.id}
      className={'row' + (sel ? ' sel' : '') + (m.unread ? ' unreadrow' : ' read')}
      style={style}
      onClick={() => onSelect(m.id)}
    >
      <div className="gut">
        <span className="unread" />
      </div>
      <div className="mid">
        <div className="l1">
          <span className="from">{m.from}</span>
          <span className="subj">{m.subject}</span>
        </div>
        <div className="snip">{m.snippet}</div>
      </div>
      <div className="r">
        <span className="time">{m.time}</span>
        <span className="marks">
          {m.deadline && (
            <span className="mk dead" title={'deadline ' + m.deadline}>
              ◷
            </span>
          )}
          {m.draft && (
            <span className="mk draft" title="utkast klart">
              ◆
            </span>
          )}
          <span className="tag">{cat.short}</span>
        </span>
      </div>
    </div>
  );
}
