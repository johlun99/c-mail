import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { mail } from '../wailsjs/go/models';
import {
  fetchAccounts,
  fetchCategories,
  fetchMails,
  fetchRules,
  connectGmail,
  disconnectGmail,
  sendReply,
  onMailsUpdated,
  isApp,
} from './lib/api';
import { catBy } from './lib/categories';
import { metaPressed, MOD, ENTER } from './lib/platform';
import { useTweaks } from './hooks/useTweaks';
import { TopBar } from './components/TopBar';
import { Rail } from './components/Rail';
import { MessageList } from './components/MessageList';
import { ReadingPane } from './components/ReadingPane';
import { StatusBar, type Mode } from './components/StatusBar';
import { CommandPalette, type Command } from './components/CommandPalette';
import { Help } from './components/Help';
import { Settings } from './components/Settings';

interface Toast {
  msg: string;
  pr: string;
}

type Overlay = 'palette' | 'settings' | 'help' | null;

function App() {
  const [tweaks, setTweak] = useTweaks();

  const [allMails, setAllMails] = useState<mail.Mail[]>([]);
  const [cats, setCats] = useState<mail.Category[]>([]);
  const [accounts, setAccounts] = useState<mail.Account[]>([]);
  const [rules, setRules] = useState<mail.Rule[]>([]);
  const [overlay, setOverlay] = useState<Overlay>(null);

  const [filter, setFilter] = useState('alla');
  const [selId, setSelId] = useState('');
  const [mode, setMode] = useState<Mode>('normal');
  const [query, setQuery] = useState('');
  const [archived, setArchived] = useState<Set<string>>(() => new Set());
  const [sent, setSent] = useState<Set<string>>(() => new Set());
  const [discarded, setDiscarded] = useState<Set<string>>(() => new Set());
  const [overrides, setOverrides] = useState<Record<string, string>>({});
  const [editing, setEditing] = useState(false);
  const [editValue, setEditValue] = useState('');
  const [toast, setToast] = useState<Toast | null>(null);
  const [connecting, setConnecting] = useState(false);

  const gPending = useRef(0);
  const toastTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  // Load all data from the Go MailService.
  const loadData = useCallback(async () => {
    const [m, c, a, r] = await Promise.all([
      fetchMails(),
      fetchCategories(),
      fetchAccounts(),
      fetchRules(),
    ]);
    setAllMails(m);
    setCats(c);
    setAccounts(a);
    setRules(r);
    if (m.length) setSelId((cur) => cur || m[0].id);
  }, []);

  useEffect(() => {
    void loadData();
    // Re-load when the backend signals a background sync / connection change.
    return onMailsUpdated(() => void loadData());
  }, [loadData]);

  const working = useMemo(
    () =>
      allMails
        .filter((m) => !archived.has(m.id))
        .map((m) => {
          const override = overrides[m.id];
          if (!override) return m;
          const cp = mail.Mail.createFrom(m);
          cp.cat = override;
          return cp;
        }),
    [allMails, archived, overrides],
  );

  const counts = useMemo(() => {
    const c: Record<string, number> = { alla: working.length };
    for (const m of working) c[m.cat] = (c[m.cat] ?? 0) + 1;
    return c;
  }, [working]);

  const filtered = useMemo(() => {
    let list = filter === 'alla' ? working : working.filter((m) => m.cat === filter);
    if (query) {
      const q = query.toLowerCase();
      list = list.filter((m) => `${m.from} ${m.subject} ${m.snippet}`.toLowerCase().includes(q));
    }
    return list;
  }, [working, filter, query]);

  const selIdx = filtered.findIndex((m) => m.id === selId);
  const sel = filtered[selIdx] ?? null;
  const filterLabel = filter === 'alla' ? 'alla' : catBy(cats, filter).label;

  // Keep selection valid as the filtered list changes.
  useEffect(() => {
    if (filtered.length === 0) return;
    if (!filtered.some((m) => m.id === selId)) setSelId(filtered[0].id);
  }, [filtered, selId]);

  // Reset the draft editor when the selection changes.
  useEffect(() => {
    setEditing(false);
    const cur = allMails.find((m) => m.id === selId);
    if (cur?.draft) setEditValue(cur.draft.lines.map((l) => l.text).join('\n'));
  }, [selId, allMails]);

  const flash = useCallback((msg: string, pr = '✓') => {
    setToast({ msg, pr });
    clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(() => setToast(null), 1900);
  }, []);

  const toggleRule = useCallback((i: number) => {
    setRules((rs) => rs.map((r, j) => (j === i ? mail.Rule.createFrom({ ...r, on: !r.on }) : r)));
  }, []);

  const onConnect = useCallback(async () => {
    if (!isApp()) {
      flash('Gmail kräver desktop-appen', '⚠');
      return;
    }
    setConnecting(true);
    try {
      await connectGmail();
      await loadData();
      flash('konto anslutet', '✓');
    } catch {
      flash('anslutning misslyckades', '⚠');
    } finally {
      setConnecting(false);
    }
  }, [flash, loadData]);

  const onDisconnect = useCallback(async () => {
    await disconnectGmail();
    await loadData();
    flash('konto frånkopplat', '✓');
  }, [flash, loadData]);

  // ── actions ──
  const move = (delta: number) => {
    if (!filtered.length) return;
    let i = selIdx < 0 ? 0 : selIdx + delta;
    i = Math.max(0, Math.min(filtered.length - 1, i));
    setSelId(filtered[i].id);
  };
  const jump = (where: 'top' | 'bottom') => {
    if (!filtered.length) return;
    setSelId(filtered[where === 'top' ? 0 : filtered.length - 1].id);
  };
  const archive = () => {
    if (!sel) return;
    const next = filtered[selIdx + 1] ?? filtered[selIdx - 1];
    setArchived((s) => new Set(s).add(sel.id));
    if (next) setSelId(next.id);
    flash(`arkiverade · ${sel.from}`, 'e');
  };
  const recategorize = () => {
    if (!sel || !cats.length) return;
    const order = cats.map((c) => c.key);
    const cur = order.indexOf(sel.cat);
    const nextKey = order[(cur + 1) % order.length];
    setOverrides((o) => ({ ...o, [sel.id]: nextKey }));
    flash(`kategori → [${catBy(cats, nextKey).label}]`, 'c');
  };
  const approve = () => {
    if (!sel) return;
    const m = sel;
    const draft = m.draft;
    if (!draft || sent.has(m.id)) return;
    const subject = m.subject.startsWith('Re:') ? m.subject : `Re: ${m.subject}`;
    const body = editing ? editValue : draft.lines.map((l) => l.text).join('\n');
    const markSent = () => {
      setSent((s) => new Set(s).add(m.id));
      setEditing(false);
    };
    if (accounts.length > 0) {
      // Connected: send for real — this is the explicit user approval.
      void sendReply(m.fromAddr, subject, body)
        .then(() => {
          markSent();
          flash(`skickat → ${m.from}`, '↗');
        })
        .catch(() => flash('kunde inte skicka', '⚠'));
    } else {
      markSent();
      flash(`skickat → ${m.from}`, '↗');
    }
  };
  const regenerate = () => {
    if (sel?.draft) {
      setEditing(false);
      flash('agenten genererade ett nytt utkast', 'agent >');
    }
  };
  const discardDraft = () => {
    if (sel?.draft) {
      setDiscarded((s) => new Set(s).add(sel.id));
      flash('utkast slängt', 'd');
    }
  };
  const openDraft = () => {
    if (!sel) return;
    if (sel.draft && !discarded.has(sel.id) && !sent.has(sel.id)) setEditing(true);
    else flash('inget utkast — agenten skapar ett…', 'agent >');
  };

  // ── command palette commands ──
  const commands = useMemo<Command[]>(() => {
    const close = () => setOverlay(null);
    return [
      {
        id: 'settings',
        name: 'Öppna inställningar (agentregler)',
        group: 'navigera',
        ic: '⚙',
        kbd: ',',
        run: () => setOverlay('settings'),
      },
      {
        id: 'help',
        name: 'Visa kortkommandon',
        group: 'navigera',
        ic: '?',
        kbd: '?',
        run: () => setOverlay('help'),
      },
      {
        id: 'search',
        name: 'Sök i inkorgen',
        group: 'navigera',
        ic: '/',
        kbd: '/',
        run: () => {
          close();
          setMode('search');
        },
      },
      ...cats.map((c) => ({
        id: 'cat-' + c.key,
        name: 'Gå till: ' + c.label,
        group: 'kategorier',
        ic: '›',
        run: () => {
          setFilter(c.key);
          close();
        },
      })),
      {
        id: 'cat-alla',
        name: 'Gå till: alla',
        group: 'kategorier',
        ic: '›',
        run: () => {
          setFilter('alla');
          close();
        },
      },
      {
        id: 'approve',
        name: 'Godkänn & skicka utkast',
        group: 'agent',
        ic: '↗',
        kbd: `${MOD} ${ENTER}`,
        run: () => {
          approve();
          close();
        },
      },
      {
        id: 'edit',
        name: 'Redigera utkast',
        group: 'agent',
        ic: '✎',
        kbd: 'e',
        run: () => {
          openDraft();
          close();
        },
      },
      {
        id: 'regen',
        name: 'Generera om utkast',
        group: 'agent',
        ic: '↻',
        kbd: 'r',
        run: () => {
          regenerate();
          close();
        },
      },
      {
        id: 'recat',
        name: 'Ändra kategori på valt mail',
        group: 'agent',
        ic: '⊟',
        kbd: 'c',
        run: () => {
          recategorize();
          close();
        },
      },
      {
        id: 'archive',
        name: 'Arkivera valt mail',
        group: 'åtgärder',
        ic: 'e',
        kbd: 'a',
        run: () => {
          archive();
          close();
        },
      },
      {
        id: 'acc-gron',
        name: 'Tema: grön accent',
        group: 'utseende',
        ic: '●',
        run: () => {
          setTweak('accent', '#7dd3a8');
          close();
        },
      },
      {
        id: 'acc-bla',
        name: 'Tema: blå accent',
        group: 'utseende',
        ic: '●',
        run: () => {
          setTweak('accent', '#6aa3f5');
          close();
        },
      },
      {
        id: 'layout',
        name: 'Växla 2/3 kolumner',
        group: 'utseende',
        ic: '▦',
        run: () => {
          setTweak('layout', tweaks.layout === '3 kol' ? '2 kol' : '3 kol');
          close();
        },
      },
      {
        id: 'font',
        name: 'Växla mono/sans',
        group: 'utseende',
        ic: 'Aa',
        run: () => {
          setTweak('font', tweaks.font === 'mono' ? 'sans' : 'mono');
          close();
        },
      },
    ];
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cats, sel, selIdx, filtered, tweaks.layout, tweaks.font]);

  // ── global keymap ──
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (overlay) return; // overlays own their keys
      const meta = metaPressed(e);
      if (meta && (e.key === 'k' || e.key === 'K')) {
        e.preventDefault();
        setOverlay('palette');
        return;
      }
      if (meta && e.key === 'Enter') {
        e.preventDefault();
        approve();
        return;
      }

      if (mode === 'search') {
        if (e.key === 'Escape') {
          setQuery('');
          setMode('normal');
        } else if (e.key === 'Enter') {
          setMode('normal');
        } else if (e.key === 'Backspace') {
          setQuery((q) => q.slice(0, -1));
        } else if (e.key.length === 1 && !meta) {
          setQuery((q) => q + e.key);
        }
        e.preventDefault();
        return;
      }

      const tag = (e.target as HTMLElement | null)?.tagName?.toLowerCase();
      if (tag === 'textarea' || tag === 'input') {
        if (e.key === 'Escape') {
          (e.target as HTMLElement).blur();
          setEditing(false);
        }
        return;
      }

      switch (e.key) {
        case 'j':
        case 'ArrowDown':
          e.preventDefault();
          move(1);
          break;
        case 'k':
        case 'ArrowUp':
          e.preventDefault();
          move(-1);
          break;
        case 'G':
          e.preventDefault();
          jump('bottom');
          break;
        case 'g': {
          const now = Date.now();
          if (now - gPending.current < 500) {
            jump('top');
            gPending.current = 0;
          } else {
            gPending.current = now;
          }
          break;
        }
        case 'Enter':
        case 'l':
          e.preventDefault();
          openDraft();
          break;
        case 'r':
          e.preventDefault();
          openDraft();
          break;
        case 'e':
          e.preventDefault();
          if (sel?.draft && !sent.has(sel.id) && !discarded.has(sel.id)) setEditing((x) => !x);
          else archive();
          break;
        case 'a':
          e.preventDefault();
          archive();
          break;
        case 'd':
          e.preventDefault();
          discardDraft();
          break;
        case 'c':
          e.preventDefault();
          recategorize();
          break;
        case '/':
          e.preventDefault();
          setQuery('');
          setMode('search');
          break;
        case '?':
          e.preventDefault();
          setOverlay('help');
          break;
        case ',':
          e.preventDefault();
          setOverlay('settings');
          break;
        case 'Escape':
          setFilter('alla');
          break;
        default:
          if (/^[1-6]$/.test(e.key) && cats.length >= +e.key) {
            e.preventDefault();
            setFilter(cats[+e.key - 1].key);
          }
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  });

  const draftVisible = !!sel?.draft && !discarded.has(sel.id);
  const pendingDrafts = working.filter((m) => m.draft && !sent.has(m.id)).length;
  const hint =
    sel?.draft && !sent.has(sel.id) && !discarded.has(sel.id)
      ? `${MOD} ${ENTER} godkänn · e redigera · r om`
      : `j/k flytta · ${ENTER} öppna · / sök`;

  return (
    <div className={'app' + (tweaks.layout === '2 kol' ? ' cols-2' : '')}>
      <TopBar
        filterLabel={filterLabel}
        account={accounts[0]?.email ?? ''}
        onHome={() => setFilter('alla')}
        onSearch={() => {
          setQuery('');
          setMode('search');
        }}
        onPalette={() => setOverlay('palette')}
        onSettings={() => setOverlay('settings')}
        onHelp={() => setOverlay('help')}
      />
      <div className="body">
        <Rail
          cats={cats}
          counts={counts}
          current={filter}
          onPick={setFilter}
          onHelp={() => setOverlay('help')}
          onPalette={() => setOverlay('palette')}
        />
        <MessageList
          mails={filtered}
          cats={cats}
          selId={selId}
          onSelect={setSelId}
          filterLabel={filterLabel}
          searchMode={mode === 'search'}
          query={query}
        />
        <ReadingPane
          m={sel}
          cats={cats}
          draftVisible={draftVisible}
          sent={!!sel && sent.has(sel.id)}
          editing={editing}
          editValue={editValue}
          onEditChange={setEditValue}
          onStartEdit={setEditing}
          onApprove={approve}
          onRegenerate={regenerate}
          onDiscard={discardDraft}
        />
      </div>
      <StatusBar
        mode={mode}
        selIdx={selIdx}
        total={filtered.length}
        current={filterLabel}
        hint={hint}
        agentMsg={`agent: sorterade ${counts.alla} · ${pendingDrafts} utkast`}
      />

      {overlay === 'palette' && (
        <CommandPalette
          commands={commands}
          onClose={() => setOverlay(null)}
          onRun={(c) => c.run()}
        />
      )}
      {overlay === 'settings' && (
        <Settings
          rules={rules}
          cats={cats}
          accounts={accounts}
          connecting={connecting}
          onConnect={onConnect}
          onDisconnect={onDisconnect}
          onToggleRule={toggleRule}
          onClose={() => setOverlay(null)}
          accent={tweaks.accent}
          font={tweaks.font}
          density={tweaks.density}
          layout={tweaks.layout}
          bgl={tweaks.bgl}
          setTweak={setTweak}
        />
      )}
      {overlay === 'help' && <Help onClose={() => setOverlay(null)} />}

      {toast && (
        <div className="toast">
          <span className="pr">{toast.pr}</span>
          {toast.msg}
        </div>
      )}
    </div>
  );
}

export default App;
