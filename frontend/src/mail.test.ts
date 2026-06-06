import { mail } from '../wailsjs/go/models';

// Guards the TS<->Go data contract: the generated models must round-trip the
// shapes the MailService returns (nested agent analysis + optional draft).
describe('mail models', () => {
  it('constructs a Mail with nested agent analysis and draft', () => {
    const m = mail.Mail.createFrom({
      id: 'm1',
      cat: 'svara',
      from: 'Anna Lindqvist',
      confidence: 0.94,
      threadCount: 12,
      body: ['Hej!'],
      agent: { summary: 's', facts: [{ key: 'kategori', value: '[svara]' }], tasks: ['x'] },
      draft: { tone: 'varm', lines: [{ text: 'Hej Anna,', kind: '+' }] },
    });

    expect(m.id).toBe('m1');
    expect(m.cat).toBe('svara');
    expect(m.agent.facts[0].key).toBe('kategori');
    expect(m.draft?.lines[0].kind).toBe('+');
  });

  it('leaves draft undefined when absent', () => {
    const m = mail.Mail.createFrom({ id: 'm3', cat: 'faktura' });
    expect(m.draft).toBeUndefined();
  });
});
