import { render, screen } from '@testing-library/react';
import App from './App';

describe('App', () => {
  it('renders the greet button', () => {
    render(<App />);
    expect(screen.getByRole('button', { name: /greet/i })).toBeInTheDocument();
  });
});
