import { render, screen } from '@testing-library/react';
import { MessageList } from './MessageList';
import { mail } from '../../wailsjs/go/models';

const cats = [
  mail.Category.createFrom({
    key: 'svara',
    label: 'att svara på',
    short: 'svara',
    color: 'c-accent',
  }),
];
const mails = [
  mail.Mail.createFrom({
    id: 'm1',
    cat: 'svara',
    from: 'Anna',
    subject: 'Hej',
    snippet: 'snippet',
    time: '09:12',
    day: 'idag',
  }),
];

describe('MessageList', () => {
  it('renders rows, the day group and the category tag', () => {
    render(
      <MessageList
        mails={mails}
        cats={cats}
        selId="m1"
        onSelect={() => {}}
        filterLabel="att svara på"
        searchMode={false}
        query=""
      />,
    );
    expect(screen.getByText('Anna')).toBeInTheDocument();
    expect(screen.getByText('Hej')).toBeInTheDocument();
    expect(screen.getByText('idag')).toBeInTheDocument();
    expect(screen.getByText('svara')).toBeInTheDocument();
  });

  it('shows an empty state when there are no mails', () => {
    render(
      <MessageList
        mails={[]}
        cats={cats}
        selId=""
        onSelect={() => {}}
        filterLabel="alla"
        searchMode={false}
        query=""
      />,
    );
    expect(screen.getByText('inga träffar')).toBeInTheDocument();
  });
});
