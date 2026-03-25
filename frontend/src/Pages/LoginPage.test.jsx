import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, beforeEach, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LoginPage from './LoginPage.jsx';

const mockNavigate = vi.fn();
const mockLogin = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate
  };
});

vi.mock('../context/AuthContext', () => ({
  useAuth: () => ({ login: mockLogin })
}));

describe('Issue #22 - Sign In Functionality', () => {
  beforeEach(() => {
    mockNavigate.mockReset();
    mockLogin.mockReset();
  });

  it('shows validation errors for invalid email and UFID', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    await user.type(screen.getByPlaceholderText('name@ufl.edu'), 'wrong@gmail.com');
    await user.type(screen.getByPlaceholderText('8 digit UFID'), '1234');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(screen.getByText(/use a valid @ufl.edu email address/i)).toBeInTheDocument();
    expect(screen.getByText(/ufid must be exactly 8 digits/i)).toBeInTheDocument();
    expect(mockLogin).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('logs in a valid user and navigates to /home', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    await user.type(screen.getByPlaceholderText('name@ufl.edu'), 'Student@UFL.edu');
    await user.type(screen.getByPlaceholderText('8 digit UFID'), '12345678');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(mockLogin).toHaveBeenCalledWith({
      email: 'student@ufl.edu',
      ufid: '12345678'
    });
    expect(mockNavigate).toHaveBeenCalledWith('/home');
  });
});
