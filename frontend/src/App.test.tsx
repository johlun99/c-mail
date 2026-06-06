import { render, screen } from '@testing-library/react';
import App from './App';

describe('App', () => {
  it('renders the cmail brand and status bar', async () => {
    render(<App />);
    expect(await screen.findByText('cmail')).toBeInTheDocument();
    expect(screen.getByText('NORMAL')).toBeInTheDocument();
  });
});
