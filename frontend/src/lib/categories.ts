import { mail } from '../../wailsjs/go/models';

// Fallback used before categories have loaded or for unknown keys.
const FALLBACK: mail.Category = mail.Category.createFrom({
  key: 'alla',
  label: 'alla',
  short: 'alla',
  color: 'c-dim',
  desc: '',
});

/** Look up a category by key, falling back gracefully. */
export function catBy(cats: mail.Category[], key: string): mail.Category {
  return cats.find((c) => c.key === key) ?? cats[0] ?? FALLBACK;
}

/** The CSS custom-property reference for a category's color token. */
export function colorVar(c: mail.Category): string {
  return `var(--${c.color})`;
}
