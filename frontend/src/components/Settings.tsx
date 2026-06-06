import { useEffect, useState } from 'react';
import { mail } from '../../wailsjs/go/models';
import { MOD, ENTER } from '../lib/platform';
import { KEYMAP } from '../lib/keymap';
import type { Density, Font, Layout, Tweaks } from '../hooks/useTweaks';

type SetTweak = <K extends keyof Tweaks>(key: K, value: Tweaks[K]) => void;

const SECTIONS = [
  { key: 'konton', label: 'konton', ic: '@' },
  { key: 'utseende', label: 'utseende', ic: '~' },
  { key: 'agent', label: 'agent & automatik', ic: '>' },
  { key: 'kategorier', label: 'kategorier', ic: '#' },
  { key: 'signatur', label: 'röst & signatur', ic: '¶' },
  { key: 'sekretess', label: 'sekretess & data', ic: '*' },
  { key: 'om', label: 'om & kommandon', ic: '?' },
] as const;

const ACC = [
  { v: '#7dd3a8', name: 'grön' },
  { v: '#6aa3f5', name: 'blå' },
  { v: '#c9a4f0', name: 'lila' },
  { v: '#e8995c', name: 'amber' },
];
const TONES = ['professionell · varm', 'kort · direkt', 'formell', 'vänlig · informell'];
const clampBg = (v: number) => Math.max(0.12, Math.min(0.22, Math.round(v * 1000) / 1000));

interface SettingsProps {
  rules: mail.Rule[];
  cats: mail.Category[];
  accounts: mail.Account[];
  connecting: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
  onToggleRule: (i: number) => void;
  onClose: () => void;
  accent: string;
  font: Font;
  density: Density;
  layout: Layout;
  bgl: number;
  setTweak: SetTweak;
}

export function Settings({
  rules,
  cats,
  accounts,
  connecting,
  onConnect,
  onDisconnect,
  onToggleRule,
  onClose,
  accent,
  font,
  density,
  layout,
  bgl,
  setTweak,
}: SettingsProps) {
  const [section, setSection] = useState(0);
  const [focus, setFocus] = useState<'nav' | 'content'>('nav');
  const [itemIdx, setItemIdx] = useState(0);
  const [toneIdx, setToneIdx] = useState(0);
  const [catAuto, setCatAuto] = useState<Record<string, boolean>>(() => {
    const o: Record<string, boolean> = {};
    for (const c of cats) o[c.key] = c.key === 'svara' || c.key === 'vantar';
    return o;
  });
  const [priv, setPriv] = useState({ thread: true, cloud: false, localdraft: true, retain: false });

  const cycle = <T,>(arr: T[], cur: T, set: (v: T) => void) =>
    set(arr[(arr.indexOf(cur) + 1) % arr.length]);
  const cycleAccent = () =>
    cycle(
      ACC.map((a) => a.v),
      accent,
      (v) => setTweak('accent', v),
    );
  const bgPresets = [0.13, 0.165, 0.2];
  const cycleBg = () => {
    const i = bgPresets.reduce(
      (b, p, k) => (Math.abs(p - bgl) < Math.abs(bgPresets[b] - bgl) ? k : b),
      0,
    );
    setTweak('bgl', bgPresets[(i + 1) % 3]);
  };

  const items: { run: () => void }[] = (() => {
    switch (SECTIONS[section].key) {
      case 'konton':
        return [...accounts.map(() => ({ run: onDisconnect })), { run: onConnect }];
      case 'utseende':
        return [
          { run: cycleAccent },
          { run: () => setTweak('font', font === 'mono' ? 'sans' : 'mono') },
          {
            run: () =>
              cycle(['tät', 'normal', 'luftig'] as Density[], density, (v) =>
                setTweak('density', v),
              ),
          },
          { run: () => setTweak('layout', layout === '3 kol' ? '2 kol' : '3 kol') },
          { run: cycleBg },
        ];
      case 'agent':
        return rules.map((r, i) => ({ run: () => !r.locked && onToggleRule(i) }));
      case 'kategorier':
        return cats.map((c) => ({ run: () => setCatAuto((o) => ({ ...o, [c.key]: !o[c.key] })) }));
      case 'signatur':
        return [{ run: () => setToneIdx((i) => (i + 1) % TONES.length) }, { run: () => {} }];
      case 'sekretess':
        return (Object.keys(priv) as (keyof typeof priv)[]).map((k) => ({
          run: () => setPriv((p) => ({ ...p, [k]: !p[k] })),
        }));
      default:
        return [];
    }
  })();

  useEffect(() => {
    const h = (e: KeyboardEvent) => {
      const k = e.key;
      if (k === 'Escape') {
        e.preventDefault();
        if (focus === 'content') setFocus('nav');
        else onClose();
        return;
      }
      if (focus === 'nav') {
        if (k === 'j' || k === 'ArrowDown') {
          e.preventDefault();
          setSection((s) => Math.min(s + 1, SECTIONS.length - 1));
          setItemIdx(0);
        } else if (k === 'k' || k === 'ArrowUp') {
          e.preventDefault();
          setSection((s) => Math.max(s - 1, 0));
          setItemIdx(0);
        } else if (k === 'Enter' || k === 'l' || k === 'ArrowRight') {
          if (items.length) {
            e.preventDefault();
            setFocus('content');
            setItemIdx(0);
          }
        }
      } else {
        if (k === 'j' || k === 'ArrowDown') {
          e.preventDefault();
          setItemIdx((i) => Math.min(i + 1, items.length - 1));
        } else if (k === 'k' || k === 'ArrowUp') {
          e.preventDefault();
          setItemIdx((i) => {
            if (i <= 0) {
              setFocus('nav');
              return 0;
            }
            return i - 1;
          });
        } else if (k === 'h' || k === 'ArrowLeft') {
          e.preventDefault();
          setFocus('nav');
        } else if (k === ' ' || k === 'Enter') {
          e.preventDefault();
          items[itemIdx]?.run();
        }
      }
    };
    window.addEventListener('keydown', h, true);
    return () => window.removeEventListener('keydown', h, true);
  }, [focus, section, itemIdx, items, onClose]);

  const foc = focus === 'content' ? itemIdx : -1;
  const counts: Record<string, number> = {
    konton: accounts.length,
    agent: rules.length,
    kategorier: cats.length,
  };
  const pickSection = (i: number) => {
    setSection(i);
    setItemIdx(0);
    setFocus('nav');
  };
  const key = SECTIONS[section].key;

  return (
    <div className="scrim" onClick={onClose}>
      <div className="settings" onClick={(e) => e.stopPropagation()}>
        <div className="set-top">
          <span className="pr">⚙</span>
          <span className="t">inställningar</span>
          <span className="sp" />
          <button className="btn" onClick={onClose}>
            stäng <kbd>esc</kbd>
          </button>
        </div>
        <div className="set-main">
          <nav className="set-nav">
            <div className="gp">inställningar</div>
            {SECTIONS.map((s, i) => (
              <div
                key={s.key}
                className={
                  'navitem' +
                  (i === section ? ' sel' : '') +
                  (i === section && focus === 'nav' ? ' navfoc' : '')
                }
                onClick={() => pickSection(i)}
              >
                <span className="ic">{s.ic}</span>
                <span>{s.label}</span>
                {counts[s.key] != null && <span className="ct">{counts[s.key]}</span>}
              </div>
            ))}
          </nav>
          <div className="set-content">
            {key === 'konton' && (
              <KontonPane
                accounts={accounts}
                foc={foc}
                connecting={connecting}
                onConnect={onConnect}
                onDisconnect={onDisconnect}
              />
            )}
            {key === 'utseende' && (
              <UtseendePane
                foc={foc}
                accent={accent}
                font={font}
                density={density}
                layout={layout}
                bgl={bgl}
                setTweak={setTweak}
              />
            )}
            {key === 'agent' && <AgentPane rules={rules} foc={foc} onToggle={onToggleRule} />}
            {key === 'kategorier' && (
              <KategoriPane cats={cats} foc={foc} catAuto={catAuto} setCatAuto={setCatAuto} />
            )}
            {key === 'signatur' && (
              <SignaturPane foc={foc} toneIdx={toneIdx} setToneIdx={setToneIdx} />
            )}
            {key === 'sekretess' && <SekretessPane foc={foc} priv={priv} setPriv={setPriv} />}
            {key === 'om' && <OmPane />}
          </div>
        </div>
        <div className="set-foot">
          <span>
            <b>j/k</b> flytta
          </span>
          <span>
            <b>{ENTER}</b> {focus === 'nav' ? 'öppna sektion' : 'växla'}
          </span>
          <span>
            <b>h</b> tillbaka
          </span>
          <span>
            <b>esc</b> stäng
          </span>
          <span className="sp" style={{ flex: 1 }} />
          <span className="muted">eller klicka var som helst</span>
        </div>
      </div>
    </div>
  );
}

function Seg<T extends string>({ opts, val, onSet }: { opts: T[]; val: T; onSet: (v: T) => void }) {
  return (
    <div className="seg2">
      {opts.map((o) => (
        <button key={o} className={val === o ? 'on' : ''} onClick={() => onSet(o)}>
          {o}
        </button>
      ))}
    </div>
  );
}

function KontonPane({
  accounts,
  foc,
  connecting,
  onConnect,
  onDisconnect,
}: {
  accounts: mail.Account[];
  foc: number;
  connecting: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
}) {
  return (
    <div>
      <div className="pane-h">Konton</div>
      <div className="pane-sub">
        Koppla din Gmail. Agenten läser och kategoriserar inkommande mail — men skickar aldrig något
        utan ditt godkännande.
      </div>
      {accounts.map((a, i) => (
        <AccountCard key={a.email} a={a} foc={foc === i} onDisconnect={onDisconnect} />
      ))}
      {accounts.length === 0 && (
        <div className="pane-sub" style={{ marginTop: 0 }}>
          Inget konto anslutet än.
        </div>
      )}
      <div
        className={'addrow' + (foc === accounts.length ? ' foc' : '')}
        onClick={connecting ? undefined : onConnect}
      >
        <span className="plus">+</span> {connecting ? 'ansluter…' : 'Anslut Gmail-konto'}
        <span className="muted">· öppnar Google-inloggning</span>
      </div>
    </div>
  );
}

function AccountCard({
  a,
  foc,
  onDisconnect,
}: {
  a: mail.Account;
  foc: boolean;
  onDisconnect: () => void;
}) {
  const st = a.status;
  const badge: [string, string] =
    st === 'connected'
      ? ['on', 'ansluten']
      : st === 'connecting'
        ? ['warn', 'ansluter…']
        : ['off', 'ej ansluten'];
  return (
    <div className={'acct-card' + (foc ? ' foc' : '')}>
      <div className="acct-row">
        <div className="prov">{a.provider[0]}</div>
        <div className="acct-info">
          <div className="em">{a.email}</div>
          <div className="meta">
            {a.provider}
            {st === 'connected'
              ? ` · ${a.syncedAt}`
              : st === 'connecting'
                ? ' · autentiserar'
                : ' · ej kopplad'}
          </div>
        </div>
        <span className={'status-badge ' + badge[0]}>{badge[1]}</span>
      </div>
      {a.scopes.length > 0 && (
        <div className="scopes">
          {a.scopes.map((s) => (
            <div className="scope-line" key={s}>
              <span className="ck">✓</span> {s}
            </div>
          ))}
        </div>
      )}
      <div className="acct-foot">
        <button className="btn" onClick={onDisconnect}>
          koppla från
        </button>
      </div>
    </div>
  );
}

function UtseendePane({
  foc,
  accent,
  font,
  density,
  layout,
  bgl,
  setTweak,
}: {
  foc: number;
  accent: string;
  font: Font;
  density: Density;
  layout: Layout;
  bgl: number;
  setTweak: SetTweak;
}) {
  const bgLabel = bgl <= 0.14 ? 'mörk' : bgl >= 0.19 ? 'ljus' : 'normal';
  return (
    <div>
      <div className="pane-h">Utseende</div>
      <div className="pane-sub">Tema och layout. Ändras direkt i appen.</div>
      <div className={'ctrl-row' + (foc === 0 ? ' foc' : '')} style={{ paddingBottom: 24 }}>
        <div className="lbl">
          Accentfärg<small>fokus, agentaktivitet, knappar</small>
        </div>
        <div className="swatches">
          {ACC.map((o) => (
            <div
              key={o.v}
              className={'swatch' + (accent === o.v ? ' sel' : '')}
              style={{ background: o.v, ['--sw-c' as string]: o.v }}
              title={o.name}
              onClick={() => setTweak('accent', o.v)}
            >
              <span className="nm">{o.name}</span>
            </div>
          ))}
        </div>
      </div>
      <div className={'ctrl-row' + (foc === 1 ? ' foc' : '')}>
        <div className="lbl">Typsnitt</div>
        <Seg opts={['mono', 'sans']} val={font} onSet={(v) => setTweak('font', v)} />
      </div>
      <div className={'ctrl-row' + (foc === 2 ? ' foc' : '')}>
        <div className="lbl">Densitet</div>
        <Seg
          opts={['tät', 'normal', 'luftig']}
          val={density}
          onSet={(v) => setTweak('density', v)}
        />
      </div>
      <div className={'ctrl-row' + (foc === 3 ? ' foc' : '')}>
        <div className="lbl">Kolumner</div>
        <Seg opts={['2 kol', '3 kol']} val={layout} onSet={(v) => setTweak('layout', v)} />
      </div>
      <div className={'ctrl-row' + (foc === 4 ? ' foc' : '')}>
        <div className="lbl">Bakgrund</div>
        <div className="stepper">
          <button className="sbtn" onClick={() => setTweak('bgl', clampBg(bgl - 0.02))}>
            −
          </button>
          <span className="val">{bgLabel}</span>
          <button className="sbtn" onClick={() => setTweak('bgl', clampBg(bgl + 0.02))}>
            +
          </button>
        </div>
      </div>
    </div>
  );
}

function AgentPane({
  rules,
  foc,
  onToggle,
}: {
  rules: mail.Rule[];
  foc: number;
  onToggle: (i: number) => void;
}) {
  return (
    <div>
      <div className="pane-h">Agent & automatik</div>
      <div className="pane-sub">
        Vad agenten får göra på egen hand. Allt utom dessa kräver din bekräftelse.
      </div>
      {rules.map((r, i) => (
        <div
          key={i}
          className={'rule' + (foc === i ? ' foc' : '')}
          onClick={() => !r.locked && onToggle(i)}
        >
          <div className={'sw' + (r.on ? ' on' : '') + (r.locked ? ' lock' : '')}>
            <span className="kn" />
          </div>
          <span className="txt">{r.text}</span>
          <span className={'scope' + (r.locked ? ' locked' : '')}>
            {r.locked ? 'fast' : r.scope}
          </span>
        </div>
      ))}
    </div>
  );
}

function KategoriPane({
  cats,
  foc,
  catAuto,
  setCatAuto,
}: {
  cats: mail.Category[];
  foc: number;
  catAuto: Record<string, boolean>;
  setCatAuto: React.Dispatch<React.SetStateAction<Record<string, boolean>>>;
}) {
  return (
    <div>
      <div className="pane-h">Kategorier</div>
      <div className="pane-sub">
        Etiketterna agenten sorterar inkommande mail i. Slå på auto-utkast där agenten ska föreslå
        svar direkt.
      </div>
      {cats.map((c, i) => (
        <div
          key={c.key}
          className={'rule' + (foc === i ? ' foc' : '')}
          onClick={() => setCatAuto((o) => ({ ...o, [c.key]: !o[c.key] }))}
        >
          <span
            style={{
              width: 10,
              height: 10,
              borderRadius: 2,
              background: `var(--${c.color})`,
              flex: 'none',
            }}
          />
          <span className="txt">{c.label}</span>
          <span className="muted" style={{ marginRight: 9, fontSize: 'calc(var(--fs) - 2px)' }}>
            auto-utkast
          </span>
          <div className={'sw' + (catAuto[c.key] ? ' on' : '')}>
            <span className="kn" />
          </div>
        </div>
      ))}
      <div className="addrow" style={{ marginTop: 12 }}>
        <span className="plus">+</span> Ny kategori{' '}
        <span className="muted">· beskriv den i ord, agenten lär sig</span>
      </div>
    </div>
  );
}

function SignaturPane({
  foc,
  toneIdx,
  setToneIdx,
}: {
  foc: number;
  toneIdx: number;
  setToneIdx: React.Dispatch<React.SetStateAction<number>>;
}) {
  return (
    <div>
      <div className="pane-h">Röst & signatur</div>
      <div className="pane-sub">
        Hur agentens utkast låter, och vad som signeras. Agenten anpassar tonen efter mottagaren.
      </div>
      <div className={'ctrl-row' + (foc === 0 ? ' foc' : '')}>
        <div className="lbl">
          Standardton<small>utgångspunkt för alla utkast</small>
        </div>
        <div className="stepper">
          <button
            className="sbtn"
            onClick={() => setToneIdx((i) => (i + TONES.length - 1) % TONES.length)}
          >
            ‹
          </button>
          <span className="val" style={{ minWidth: 150 }}>
            {TONES[toneIdx]}
          </span>
          <button className="sbtn" onClick={() => setToneIdx((i) => (i + 1) % TONES.length)}>
            ›
          </button>
        </div>
      </div>
      <div className={'ctrl-row' + (foc === 1 ? ' foc' : '')}>
        <div className="lbl">Signatur</div>
        <span className="scope">Johan · Nordveda AB</span>
        <button className="btn" style={{ marginLeft: 10 }}>
          redigera
        </button>
      </div>
    </div>
  );
}

function SekretessPane({
  foc,
  priv,
  setPriv,
}: {
  foc: number;
  priv: Record<string, boolean>;
  setPriv: React.Dispatch<
    React.SetStateAction<{ thread: boolean; cloud: boolean; localdraft: boolean; retain: boolean }>
  >;
}) {
  const rows: [string, string, string][] = [
    ['thread', 'Agenten läser hela trådhistoriken', 'kontext'],
    ['cloud', 'Använd molnmodell för analys (annars lokalt)', 'modell'],
    ['localdraft', 'Spara utkast endast lokalt tills de skickas', 'utkast'],
    ['retain', 'Radera agentens analys efter 30 dagar', 'lagring'],
  ];
  return (
    <div>
      <div className="pane-h">Sekretess & data</div>
      <div className="pane-sub">Du bestämmer vad agenten ser och var det lagras.</div>
      {rows.map(([k, txt, scope], i) => (
        <div
          key={k}
          className={'rule' + (foc === i ? ' foc' : '')}
          onClick={() => setPriv((p) => ({ ...p, [k]: !p[k as keyof typeof p] }))}
        >
          <div className={'sw' + (priv[k] ? ' on' : '')}>
            <span className="kn" />
          </div>
          <span className="txt">{txt}</span>
          <span className="scope">{scope}</span>
        </div>
      ))}
    </div>
  );
}

function OmPane() {
  const flat: [string, string][] = [];
  for (const [, rows] of KEYMAP) for (const r of rows) flat.push(r);
  return (
    <div>
      <div className="pane-h">Om cmail</div>
      <div className="pane-sub">
        Agentisk e-postklient. Keyboard-first, men allt går att klicka.
      </div>
      <div className="about-row">
        <span>version</span>
        <span className="v">0.1.0 · tidig utveckling</span>
      </div>
      <div className="about-row">
        <span>agentmodell</span>
        <span className="v">lokal · Ollama</span>
      </div>
      <div className="about-row">
        <span>plattform</span>
        <span className="v">
          {MOD === '⌘' ? 'macOS' : 'Linux'} · {MOD}-tangenter
        </span>
      </div>
      <div className="set-sec2">kortkommandon</div>
      {flat.map(([k, d], i) => (
        <div className="about-row" key={i}>
          <span>{d}</span>
          <span className="keys" style={{ display: 'flex', gap: 4 }}>
            {k.split(' / ').map((kk) => (
              <kbd key={kk}>{kk}</kbd>
            ))}
          </span>
        </div>
      ))}
    </div>
  );
}
