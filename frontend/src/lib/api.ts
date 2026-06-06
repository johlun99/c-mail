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
} from '../../wailsjs/go/main/MailService';

function wailsReady(): boolean {
  return typeof window !== 'undefined' && 'go' in window;
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
