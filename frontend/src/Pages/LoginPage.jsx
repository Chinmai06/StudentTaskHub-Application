import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function LoginPage() {
  const [isRegistering, setIsRegistering] = useState(false);
  const [formState, setFormState] = useState({ username: '', email: '', password: '' });
  const [errors, setErrors] = useState({});
  const [message, setMessage] = useState('');
  const { login, register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormState((current) => ({ ...current, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setErrors({});
    setMessage('');

    const nextErrors = {};

    if (!formState.username.trim()) {
      nextErrors.username = 'Username is required.';
    }
    if (!formState.password.trim()) {
      nextErrors.password = 'Password is required.';
    }
    if (isRegistering && !formState.email.trim()) {
      nextErrors.email = 'Email is required.';
    }
    if (isRegistering && formState.password.length < 6) {
      nextErrors.password = 'Password must be at least 6 characters.';
    }

    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) return;

    try {
      if (isRegistering) {
        await register(formState.username.trim(), formState.email.trim(), formState.password);
        setMessage('Registration successful! You can now sign in.');
        setIsRegistering(false);
        setFormState({ username: formState.username, email: '', password: '' });
      } else {
        await login(formState.username.trim(), formState.password);
        navigate('/home');
      }
    } catch (err) {
      setErrors({ form: err.message });
    }
  };

  return (
    <section className="auth-page">
      <div className="auth-card">
        <p className="eyebrow">React Frontend</p>
        <h1>Welcome to StudentTaskHub</h1>
        <p className="auth-copy">
          {isRegistering
            ? 'Create an account to start managing tasks.'
            : 'Sign in to create and manage student tasks.'}
        </p>
        <form onSubmit={handleSubmit} className="form-grid" data-testid="login-form">
          <label>
            Username
            <input
              type="text"
              name="username"
              placeholder="Enter username"
              value={formState.username}
              onChange={handleChange}
            />
          </label>
          {errors.username ? <p className="error-text">{errors.username}</p> : null}

          {isRegistering ? (
            <>
              <label>
                Email
                <input
                  type="email"
                  name="email"
                  placeholder="name@ufl.edu"
                  value={formState.email}
                  onChange={handleChange}
                />
              </label>
              {errors.email ? <p className="error-text">{errors.email}</p> : null}
            </>
          ) : null}

          <label>
            Password
            <input
              type="password"
              name="password"
              placeholder="Enter password"
              value={formState.password}
              onChange={handleChange}
            />
          </label>
          {errors.password ? <p className="error-text">{errors.password}</p> : null}

          {errors.form ? <p className="error-text">{errors.form}</p> : null}
          {message ? <p className="success-text">{message}</p> : null}

          <button type="submit" className="primary-button">
            {isRegistering ? 'Register' : 'Sign In'}
          </button>
        </form>

        <p style={{ textAlign: 'center', marginTop: '1rem' }}>
          {isRegistering ? (
            <>
              Already have an account?{' '}
              <button
                type="button"
                onClick={() => { setIsRegistering(false); setErrors({}); setMessage(''); }}
                style={{ background: 'none', border: 'none', color: '#2563eb', cursor: 'pointer', textDecoration: 'underline' }}
              >
                Sign In
              </button>
            </>
          ) : (
            <>
              Don't have an account?{' '}
              <button
                type="button"
                onClick={() => { setIsRegistering(true); setErrors({}); setMessage(''); }}
                style={{ background: 'none', border: 'none', color: '#2563eb', cursor: 'pointer', textDecoration: 'underline' }}
              >
                Register
              </button>
            </>
          )}
        </p>
      </div>
    </section>
  );
}

export default LoginPage;
