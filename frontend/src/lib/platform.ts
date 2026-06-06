// Platform-aware key labels — the app runs on macOS + Linux.
const isMac =
  typeof navigator !== 'undefined' &&
  /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent || '');

export const MOD = isMac ? '⌘' : 'Ctrl';
export const ENTER = isMac ? '⏎' : 'Enter';

/** True when the platform's command/meta key is held. */
export function metaPressed(e: KeyboardEvent): boolean {
  return e.metaKey || e.ctrlKey;
}
