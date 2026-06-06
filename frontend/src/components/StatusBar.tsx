export type Mode = 'normal' | 'search' | 'cmd';

interface StatusBarProps {
  mode: Mode;
  selIdx: number;
  total: number;
  current: string;
  hint: string;
  agentMsg: string;
}

export function StatusBar({ mode, selIdx, total, current, hint, agentMsg }: StatusBarProps) {
  const modeCls = mode === 'search' ? 'search' : mode === 'cmd' ? 'cmd' : '';
  const modeTxt = mode === 'search' ? 'SEARCH' : mode === 'cmd' ? 'COMMAND' : 'NORMAL';
  return (
    <div className="status">
      <div className={'st-seg st-mode ' + modeCls}>{modeTxt}</div>
      <div className="st-seg">{current}</div>
      <div className="st-seg">
        <span className="st-key">
          rad <b>{total ? selIdx + 1 : 0}</b>/{total}
        </span>
      </div>
      <div className="st-seg muted">{agentMsg}</div>
      <span className="sp" />
      <div className="st-seg r">
        <span className="st-key">{hint}</span>
      </div>
      <div className="st-seg r">
        <span className="st-key">
          <b>?</b> hjälp
        </span>
      </div>
    </div>
  );
}
