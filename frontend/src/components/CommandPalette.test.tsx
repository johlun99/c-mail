import { render, screen, fireEvent } from '@testing-library/react';
import { CommandPalette, type Command } from './CommandPalette';

const commands: Command[] = [
  { id: 'a', name: 'Alpha', group: 'navigera', run: () => {} },
  { id: 'b', name: 'Beta', group: 'navigera', run: () => {} },
];

describe('CommandPalette', () => {
  it('filters commands by the typed query', () => {
    render(<CommandPalette commands={commands} onClose={() => {}} onRun={() => {}} />);
    fireEvent.change(screen.getByPlaceholderText('skriv ett kommando…'), {
      target: { value: 'bet' },
    });
    expect(screen.getByText('Beta')).toBeInTheDocument();
    expect(screen.queryByText('Alpha')).toBeNull();
  });

  it('runs the selected command on Enter', () => {
    let ran: string | null = null;
    const cmds: Command[] = [{ id: 'x', name: 'Kör mig', group: 'g', run: () => (ran = 'x') }];
    render(<CommandPalette commands={cmds} onClose={() => {}} onRun={(c) => c.run()} />);
    fireEvent.keyDown(screen.getByPlaceholderText('skriv ett kommando…'), { key: 'Enter' });
    expect(ran).toBe('x');
  });
});
