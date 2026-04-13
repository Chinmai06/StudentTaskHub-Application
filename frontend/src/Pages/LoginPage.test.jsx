import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LoginPage from './LoginPage.jsx';

const mockNavigate = vi.fn();
const mockLogin = vi.fn();
const mockRegister = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock('../context/AuthContext', () => ({
  useAuth: () => ({
    login: mockLogin,
    register: mockRegister,
  }),
}));

describe('LoginPage', () => {
  beforeEach(() => {
    mockNavigate.mockReset();
    mockLogin.mockReset();
    mockRegister.mockReset();
  });

  it('shows validation errors when username and password are empty', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(screen.getByText(/username is required/i)).toBeInTheDocument();
    expect(screen.getByText(/password is required/i)).toBeInTheDocument();
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

    await user.type(screen.getByPlaceholderText(/enter username/i), 'student1');
    await user.type(screen.getByPlaceholderText(/enter password/i), 'pass123');

    const form = screen.getByRole('button', { name: /^sign in$/i }).closest('form');
    const submitButton = within(form).getByRole('button', { name: /^sign in$/i });

    await user.click(submitButton);

    expect(mockLogin).toHaveBeenCalledWith('student1', 'pass123');
    expect(mockNavigate).toHaveBeenCalledWith('/home');
  });

  it('switches to register mode and shows email field', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    const toggleButtons = screen.getAllByRole('button', { name: /register/i });
    await user.click(toggleButtons[toggleButtons.length - 1]);

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
  });

  it('shows validation errors in register mode when fields are invalid', async () => {
    const user = userEvent.setup();

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    const toggleButtons = screen.getAllByRole('button', { name: /register/i });
    await user.click(toggleButtons[toggleButtons.length - 1]);

    const form = screen.getByRole('button', { name: /^register$/i }).closest('form');
    const submitButton = within(form).getByRole('button', { name: /^register$/i });

    await user.click(submitButton);

    expect(screen.getByText(/username is required/i)).toBeInTheDocument();
    expect(screen.getByText(/email is required/i)).toBeInTheDocument();
    expect(screen.getByText(/password is required/i)).toBeInTheDocument();
    expect(mockRegister).not.toHaveBeenCalled();
  });

  it('registers successfully with valid data', async () => {
    const user = userEvent.setup();
    mockRegister.mockResolvedValueOnce({ message: 'registered' });

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    const toggleButtons = screen.getAllByRole('button', { name: /register/i });
    await user.click(toggleButtons[toggleButtons.length - 1]);

    await user.type(screen.getByPlaceholderText(/enter username/i), 'student1');
    await user.type(screen.getByLabelText(/email/i), 'student@ufl.edu');
    await user.type(screen.getByPlaceholderText(/enter password/i), 'pass123');

    const form = screen.getByRole('button', { name: /^register$/i }).closest('form');
    const submitButton = within(form).getByRole('button', { name: /^register$/i });

    await user.click(submitButton);

    expect(mockRegister).toHaveBeenCalledWith('student1', 'student@ufl.edu', 'pass123');
    expect(
      screen.getByText(/registration successful! you can now sign in./i)
    ).toBeInTheDocument();
  });
});

// import { MemoryRouter } from 'react-router-dom';
// import { beforeEach, describe, expect, it, vi } from 'vitest';
// import { render, screen, within } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';
// import LoginPage from './LoginPage.jsx';

// const mockNavigate = vi.fn();
// const mockLogin = vi.fn();
// const mockRegister = vi.fn();

// vi.mock('react-router-dom', async () => {
//   const actual = await vi.importActual('react-router-dom');
//   return {
//     ...actual,
//     useNavigate: () => mockNavigate,
//   };
// });

// vi.mock('../context/AuthContext', () => ({
//   useAuth: () => ({
//     login: mockLogin,
//     register: mockRegister,
//   }),
// }));

// describe('LoginPage', () => {
//   beforeEach(() => {
//     mockNavigate.mockReset();
//     mockLogin.mockReset();
//     mockRegister.mockReset();
//   });

//   it('shows validation errors when username and password are empty', async () => {
//     const user = userEvent.setup();

//     render(
//       <MemoryRouter>
//         <LoginPage />
//       </MemoryRouter>
//     );

//     await user.click(screen.getByRole('button', { name: /sign in/i }));

//     expect(screen.getByText(/username is required/i)).toBeInTheDocument();
//     expect(screen.getByText(/password is required/i)).toBeInTheDocument();
//     expect(mockLogin).not.toHaveBeenCalled();
//     expect(mockNavigate).not.toHaveBeenCalled();
//   });

//   it('logs in a valid user and navigates to /home', async () => {
//     const user = userEvent.setup();

//     render(
//       <MemoryRouter>
//         <LoginPage />
//       </MemoryRouter>
//     );

//     await user.type(screen.getByPlaceholderText(/enter username/i), 'student1');
//     await user.type(screen.getByPlaceholderText(/enter password/i), 'pass123');

//     const form = screen.getByRole('button', { name: /^sign in$/i }).closest('form');
//     const submitButton = within(form).getByRole('button', { name: /^sign in$/i });

//     await user.click(submitButton);

//     expect(mockLogin).toHaveBeenCalledWith('student1', 'pass123');
//     expect(mockNavigate).toHaveBeenCalledWith('/home');
//   });

//   it('switches to register mode and shows email field', async () => {
//     const user = userEvent.setup();

//     render(
//       <MemoryRouter>
//         <LoginPage />
//       </MemoryRouter>
//     );

//     const toggleButtons = screen.getAllByRole('button', { name: /register/i });
//     await user.click(toggleButtons[toggleButtons.length - 1]);

//     expect(screen.getByPlaceholderText(/name@ufl.edu/i)).toBeInTheDocument();
//   });

//   it('shows validation errors in register mode when fields are invalid', async () => {
//     const user = userEvent.setup();

//     render(
//       <MemoryRouter>
//         <LoginPage />
//       </MemoryRouter>
//     );

//     const toggleButtons = screen.getAllByRole('button', { name: /register/i });
//     await user.click(toggleButtons[toggleButtons.length - 1]);

//     const registerSubmitButton = screen
//       .getByRole('button', { name: /^register$/i })
//       .closest('form');

//     const submitButton = within(registerSubmitButton).getByRole('button', {
//       name: /^register$/i,
//     });

//     await user.click(submitButton);

//     expect(screen.getByText(/username is required/i)).toBeInTheDocument();
//     expect(screen.getByText(/email is required/i)).toBeInTheDocument();
//     expect(screen.getByText(/password is required/i)).toBeInTheDocument();
//     expect(mockRegister).not.toHaveBeenCalled();
//   });

//   it('registers successfully with valid data', async () => {
//     const user = userEvent.setup();
//     mockRegister.mockResolvedValueOnce({ message: 'registered' });

//     render(
//       <MemoryRouter>
//         <LoginPage />
//       </MemoryRouter>
//     );

//     const toggleButtons = screen.getAllByRole('button', { name: /register/i });
//     await user.click(toggleButtons[toggleButtons.length - 1]);

//     await user.type(screen.getByPlaceholderText(/enter username/i), 'student1');
//     await user.type(screen.getByPlaceholderText(/name@ufl.edu/i), 'student@ufl.edu');
//     await user.type(screen.getByPlaceholderText(/enter password/i), 'pass123');

//     const form = screen.getByRole('button', { name: /^register$/i }).closest('form');
//     const submitButton = within(form).getByRole('button', { name: /^register$/i });

//     await user.click(submitButton);

//     expect(mockRegister).toHaveBeenCalledWith('student1', 'student@ufl.edu', 'pass123');
//     expect(
//       screen.getByText(/registration successful! you can now sign in./i)
//     ).toBeInTheDocument();
//   });
// });
