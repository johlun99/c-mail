import { MOD, ENTER } from './platform';

/** Keyboard reference shown in the help overlay and the settings "om" pane. */
export const KEYMAP: [string, [string, string][]][] = [
  [
    'navigation',
    [
      ['j / k', 'nästa / föregående mail'],
      ['g g / G', 'till toppen / botten'],
      ['1–6', 'hoppa till kategori'],
      [`${ENTER} / l`, 'öppna / fokusera'],
      ['h / esc', 'tillbaka / stäng'],
    ],
  ],
  [
    'agent & utkast',
    [
      ['r', 'generera / öppna utkast'],
      ['e', 'redigera utkast'],
      [`${MOD} ${ENTER}`, 'godkänn & skicka'],
      ['d', 'släng utkast'],
      ['c', 'ändra kategori'],
    ],
  ],
  [
    'arkiv & sök',
    [
      ['a / e', 'arkivera mail'],
      ['/', 'sök i inkorgen'],
      ['u', 'markera oläst'],
    ],
  ],
  [
    'överallt',
    [
      [`${MOD} K`, 'kommandopalett'],
      [',', 'inställningar'],
      ['?', 'denna hjälp'],
    ],
  ],
];
