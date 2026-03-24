import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { isUflEmail, isValidUfid } from '../utils/validation';

const initialState = {
  email: '',
  ufid: ''
};

function LoginPage() {
  const [formState, setFormState] = useState(initialState);
  const [errors, setErrors] = useState({});
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormState((currentState) => ({
      ...currentState,
      [name]: value
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const nextErrors = {};

    if (!isUflEmail(formState.email)) {
      nextErrors.email = 'Use a valid @ufl.edu email address.';
    }

    if (!isValidUfid(formState.ufid)) {
      nextErrors.ufid = 'UFID must be exactly 8 digits.';
    }

    setErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    login({
      email: formState.email.trim().toLowerCase(),
      ufid: formState.ufid.trim()
    });

    navigate('/home');
  };

  return (
    <section className="auth-page">
      <div className="auth-card">
        <p className="eyebrow">React Frontend</p>
        <h1>Welcome to StudentTaskHub</h1>
        <p className="auth-copy">
          Sign in using your UFL email and UFID to create and manage student tasks.
        </p>
        <form onSubmit={handleSubmit} className="form-grid" data-testid="login-form">
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

          <label>
            UFID
            <input
              type="text"
              name="ufid"
              placeholder="8 digit UFID"
              value={formState.ufid}
              onChange={handleChange}
              maxLength={8}
            />
          </label>
          {errors.ufid ? <p className="error-text">{errors.ufid}</p> : null}

          <button type="submit" className="primary-button">
            Sign In
          </button>
        </form>
      </div>
    </section>
  );
}

export default LoginPage;
