import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { isUflEmail, isValidUfid } from '../utils/validation';

function LoginPage() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const [formData, setFormData] = useState({
    email: '',
    ufid: ''
  });

  const [errors, setErrors] = useState({});

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData((current) => ({
      ...current,
      [name]: value
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const nextErrors = {};

    if (!isUflEmail(formData.email)) {
      nextErrors.email = 'Use a valid @ufl.edu email address.';
    }

    if (!isValidUfid(formData.ufid)) {
      nextErrors.ufid = 'UFID must be exactly 8 digits.';
    }

    setErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    login({
      email: formData.email.trim().toLowerCase(),
      ufid: formData.ufid.trim()
    });

    navigate('/home');
  };

  return (
    <section className="auth-page">
      <div className="auth-card">
        <p className="eyebrow">Student Task Hub</p>
        <h1>Sign In</h1>
        <p className="auth-subtitle">
          Use your UFL email and UFID to continue.
        </p>

        <form onSubmit={handleSubmit} className="auth-form">
          <label>
            Email
            <input
              type="email"
              name="email"
              placeholder="name@ufl.edu"
              value={formData.email}
              onChange={handleChange}
            />
          </label>
          {errors.email && <p>{errors.email}</p>}

          <label>
            UFID
            <input
              type="text"
              name="ufid"
              placeholder="8 digit UFID"
              value={formData.ufid}
              onChange={handleChange}
            />
          </label>
          {errors.ufid && <p>{errors.ufid}</p>}

          <button type="submit">Sign In</button>
        </form>
      </div>
    </section>
  );
}

export default LoginPage;