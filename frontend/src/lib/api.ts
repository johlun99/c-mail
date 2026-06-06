// Data layer: fetches from the Go MailService via Wails bindings. Outside a
// Wails runtime (e.g. unit tests, plain browser) the bindings are unavailable,
// so these return empty data instead of throwing.
import { mail } from '../../wailsjs/go/models';
import {
  GetMails,
  GetCategories,
  GetActivity,
  GetRules,
  GetAccounts,
  ConnectGmail,
  DisconnectGmail,
  SendReply,
} from '../../wailsjs/go/main/MailService';
import { EventsOn } from '../../wailsjs/runtime';

function wailsReady(): boolean {
  return typeof window !== 'undefined' && 'go' in window;
}

/** True when running inside the Wails desktop app (vs a plain browser preview). */
export function isApp(): boolean {
  return wailsReady();
}

export function fetchMails(): Promise<mail.Mail[]> {
  return wailsReady() ? GetMails() : Promise.resolve([]);
}

export function fetchCategories(): Promise<mail.Category[]> {
  return wailsReady() ? GetCategories() : Promise.resolve([]);
}

export function fetchActivity(): Promise<mail.Activity[]> {
  return wailsReady() ? GetActivity() : Promise.resolve([]);
}

export function fetchRules(): Promise<mail.Rule[]> {
  return wailsReady() ? GetRules() : Promise.resolve([]);
}

export function fetchAccounts(): Promise<mail.Account[]> {
  return wailsReady() ? GetAccounts() : Promise.resolve([]);
}

/** Runs the Google OAuth flow (opens the system browser). Desktop app only. */
export function connectGmail(): Promise<mail.Account> {
  return wailsReady() ? ConnectGmail() : Promise.reject(new Error('Gmail kräver desktop-appen'));
}

/** Disconnects the current account and clears its local cache + token. */
export function disconnectGmail(): Promise<void> {
  return wailsReady() ? DisconnectGmail() : Promise.resolve();
}

/** Sends a reply — only ever called from an explicit user approval. */
export function sendReply(to: string, subject: string, body: string): Promise<void> {
  return wailsReady() ? SendReply(to, subject, body) : Promise.resolve();
}

/** Subscribes to background mail updates; returns an unsubscribe function. */
export function onMailsUpdated(cb: () => void): () => void {
  if (!wailsReady()) return () => {};
  return EventsOn('mails:updated', cb);
}
