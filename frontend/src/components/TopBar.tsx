import { MOD } from '../lib/platform';

interface TopBarProps {
  filterLabel: string;
  account: string;
  onHome: () => void;
  onSearch: () => void;
  onPalette: () => void;
  onSettings: () => void;
  onHelp: () => void;
}

export function TopBar({
  filterLabel,
  account,
  onHome,
  onSearch,
  onPalette,
  onSettings,
  onHelp,
}: TopBarProps) {
  return (
    <div className="topbar">
      <div className="brand" onClick={onHome} title="till inkorgen">
        <span className="sigil">~/</span>
        <span>cmail</span>
        <span className="cur" />
      </div>
      <span className="path">agent@inbox · {filterLabel}</span>
      <span className="spacer" />
      <span className="agent-pill">
        <span className="agent-dot" />
        agent: vilande
      </span>
      <div className="tbar">
        <button className="tbtn" onClick={onSearch} title="Sök (/)">
          <span className="ic">/</span>sök
        </button>
        <button className="tbtn" onClick={onPalette} title="Kommandopalett">
          <kbd>{MOD}</kbd>
          <kbd>K</kbd>
        </button>
        <button className="tbtn" onClick={onSettings} title="Inställningar">
          <span className="ic">≡</span>inställningar
        </button>
        <button className="tbtn" onClick={onHelp} title="Kortkommandon (?)">
          <span className="ic">?</span>
        </button>
      </div>
      {account && (
        <span className="acct" onClick={onSettings} title="Konton & inställningar">
          {account}
        </span>
      )}
    </div>
  );
}
