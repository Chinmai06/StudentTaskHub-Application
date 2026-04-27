import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, register } = useAuth();

  const [isRegister, setIsRegister] = useState(false);
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState({});
  const [message, setMessage] = useState('');

  const validate = () => {
    const newErrors = {};

    if (!username.trim()) {
      newErrors.username = 'Username is required';
    }

    if (isRegister && !email.trim()) {
      newErrors.email = 'Email is required';
    }

    if (!password.trim()) {
      newErrors.password = 'Password is required';
    }

    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const validationErrors = validate();
    setErrors(validationErrors);
    setMessage('');

    if (Object.keys(validationErrors).length > 0) {
      return;
    }

    try {
      if (isRegister) {
        await register(username, email, password);
        setMessage('Registration successful! You can now sign in.');
        setIsRegister(false);
        setEmail('');
        setPassword('');
      } else {
        await login(username, password);
        navigate('/home');
      }
    } catch (err) {
      setMessage(err?.message || 'Something went wrong');
    }
  };

  return (
    <div className="login-page">
      <h2>{isRegister ? 'Register' : 'Sign In'}</h2>

      <form onSubmit={handleSubmit}>
        <div>
          <input
            type="text"
            placeholder="Enter username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          {errors.username && <p>{errors.username}</p>}
        </div>

        {isRegister && (
          <div>
            <input
              type="email"
              placeholder="name@ufl.edu"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            {errors.email && <p>{errors.email}</p>}
          </div>
        )}

        <div>
          <input
            type="password"
            placeholder="Enter password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {errors.password && <p>{errors.password}</p>}
        </div>

        <button type="submit">
          {isRegister ? 'Register' : 'Sign In'}
        </button>
      </form>

      {message && <p>{message}</p>}

      <div style={{ marginTop: '12px' }}>
        <button type="button" onClick={() => setIsRegister(false)}>
          Sign In
        </button>
        <button type="button" onClick={() => setIsRegister(true)}>
          Register
        </button>
      </div>
    </div>
  );
}
