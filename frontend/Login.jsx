import React, { useState, useEffect, useRef } from 'react';
import './Login.css';

/* ── tiny reusable animated blob ── */
const Blob = ({ className }) => <div className={`blob ${className}`} />;

/* ── floating particle ── */
const Particle = ({ style }) => <span className="particle" style={style} />;

const Login = ({ onLogin }) => {
  const [username, setUsername]       = useState('');
  const [password, setPassword]       = useState('');
  const [remember, setRemember]       = useState(true);
  const [showPass, setShowPass]       = useState(false);
  const [error, setError]             = useState('');
  const [loading, setLoading]         = useState(false);
  const [focused, setFocused]         = useState('');
  const [mounted, setMounted]         = useState(false);
  const cardRef = useRef(null);

  /* mount animation trigger */
  useEffect(() => { setTimeout(() => setMounted(true), 60); }, []);

  /* subtle parallax tilt on card */
  const handleMouseMove = (e) => {
    if (!cardRef.current) return;
    const rect = cardRef.current.getBoundingClientRect();
    const x = (e.clientX - rect.left) / rect.width  - 0.5;
    const y = (e.clientY - rect.top)  / rect.height - 0.5;
    cardRef.current.style.transform =
      `perspective(1200px) rotateY(${x * 4}deg) rotateX(${-y * 4}deg) scale(1.01)`;
  };
  const handleMouseLeave = () => {
    if (cardRef.current)
      cardRef.current.style.transform = 'perspective(1200px) rotateY(0) rotateX(0) scale(1)';
  };

  /* particles data */
  const particles = Array.from({ length: 14 }, (_, i) => ({
    left:  `${(i * 37 + 11) % 100}%`,
    top:   `${(i * 53 + 7)  % 100}%`,
    animationDelay: `${(i * 0.4).toFixed(1)}s`,
    animationDuration: `${3 + (i % 4)}s`,
    width:  `${4 + (i % 5)}px`,
    height: `${4 + (i % 5)}px`,
    opacity: (0.2 + (i % 5) * 0.07).toFixed(2),
  }));

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      setError('Please fill in all fields.');
      return;
    }
    setError('');
    setLoading(true);
    /* simulate async login */
    setTimeout(() => {
      setLoading(false);
      onLogin();
    }, 1200);
  };

  return (
    <div className={`lp-page ${mounted ? 'lp-page--in' : ''}`}>

      {/* ── animated gradient mesh background ── */}
      <div className="lp-bg">
        <div className="lp-bg__mesh" />
        <div className="lp-bg__glow lp-bg__glow--cyan"  />
        <div className="lp-bg__glow lp-bg__glow--pink"  />
        <div className="lp-bg__glow lp-bg__glow--purple"/>
      </div>

      {/* ── card ── */}
      <div
        ref={cardRef}
        className="lp-card"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
      >

        {/* ══ LEFT PANEL ══ */}
        <div className="lp-left">
          <Blob className="blob--a" />
          <Blob className="blob--b" />
          <Blob className="blob--c" />

          {particles.map((p, i) => <Particle key={i} style={p} />)}

          {/* decorative SVG waves */}
          <svg className="lp-waves" viewBox="0 0 420 480" fill="none">
            <path d="M-20 100 Q60 50 120 120 Q180 190 240 130 Q300 70 360 140 Q420 210 480 150"
              stroke="rgba(255,255,255,0.2)" strokeWidth="2.5" fill="none"/>
            <path d="M-20 180 Q70 130 140 200 Q210 270 280 210 Q350 150 420 220"
              stroke="rgba(255,255,255,0.15)" strokeWidth="1.5" fill="none"/>
            <path d="M200 320 Q240 280 270 340 Q300 400 340 370 Q380 340 420 390"
              stroke="rgba(255,255,255,0.18)" strokeWidth="2" fill="none"/>
            <path d="M180 370 Q220 330 250 385 Q280 440 320 415"
              stroke="rgba(255,255,255,0.12)" strokeWidth="1.5" fill="none"/>
          </svg>

          {/* dot grid */}
          <div className="lp-dots" />

          {/* plus signs */}
          <span className="lp-plus lp-plus--1">+</span>
          <span className="lp-plus lp-plus--2">+</span>

          {/* circles */}
          <span className="lp-ring lp-ring--1" />
          <span className="lp-ring lp-ring--2" />

          {/* text */}
          <div className="lp-hero">
            <div className="lp-hero__badge">STUDENT PLATFORM</div>
            <h2 className="lp-hero__title">Welcome back!</h2>
            <p  className="lp-hero__hub">H-HUB</p>
            <p  className="lp-hero__sub">
              Sign in to access your tasks, track progress, and stay on top of your goals.
            </p>
            <div className="lp-hero__pills">
              <span className="pill">📚 Tasks</span>
              <span className="pill">📅 Deadlines</span>
              <span className="pill">⚡ Priorities</span>
            </div>
          </div>
        </div>

        {/* ══ RIGHT PANEL ══ */}
        <div className="lp-right">

          <div className="lp-form-wrap">
            <h1 className="lp-title">Sign In</h1>
            <p className="lp-subtitle">Welcome back — let's get you logged in.</p>

            <form className="lp-form" onSubmit={handleSubmit} noValidate>

              {/* username */}
              <div className={`lp-field ${focused === 'user' ? 'lp-field--focused' : ''} ${username ? 'lp-field--filled' : ''}`}>
                <label className="lp-label">Username or email</label>
                <div className="lp-input-wrap">
                  <span className="lp-icon">
                    <svg width="17" height="17" viewBox="0 0 24 24" fill="none"
                      stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                      <circle cx="12" cy="7" r="4"/>
                    </svg>
                  </span>
                  <input
                    type="text"
                    className="lp-input"
                    placeholder="your@email.com"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                    onFocus={() => setFocused('user')}
                    onBlur={() => setFocused('')}
                    autoComplete="username"
                  />
                </div>
              </div>

              {/* password */}
              <div className={`lp-field ${focused === 'pass' ? 'lp-field--focused' : ''} ${password ? 'lp-field--filled' : ''}`}>
                <label className="lp-label">Password</label>
                <div className="lp-input-wrap">
                  <span className="lp-icon">
                    <svg width="17" height="17" viewBox="0 0 24 24" fill="none"
                      stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                      <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                    </svg>
                  </span>
                  <input
                    type={showPass ? 'text' : 'password'}
                    className="lp-input"
                    placeholder="••••••••"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    onFocus={() => setFocused('pass')}
                    onBlur={() => setFocused('')}
                    autoComplete="current-password"
                  />
                  <button
                    type="button"
                    className="lp-eye"
                    onClick={() => setShowPass(v => !v)}
                    tabIndex={-1}
                    aria-label="Toggle password visibility"
                  >
                    {showPass ? (
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none"
                        stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                        <line x1="1" y1="1" x2="23" y2="23"/>
                      </svg>
                    ) : (
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none"
                        stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                        <circle cx="12" cy="12" r="3"/>
                      </svg>
                    )}
                  </button>
                </div>
              </div>

              {/* options row */}
              <div className="lp-options">
                <label className="lp-check">
                  <input
                    type="checkbox"
                    checked={remember}
                    onChange={e => setRemember(e.target.checked)}
                  />
                  <span className="lp-check__box" />
                  <span className="lp-check__txt">Remember me</span>
                </label>
                <a href="#" className="lp-forgot">Forgot password?</a>
              </div>

              {/* error */}
              {error && (
                <div className="lp-error" role="alert">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="8" x2="12" y2="12"/>
                    <line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                  {error}
                </div>
              )}

              {/* submit */}
              <button type="submit" className={`lp-btn ${loading ? 'lp-btn--loading' : ''}`} disabled={loading}>
                {loading ? (
                  <>
                    <span className="lp-spinner" />
                    Signing in…
                  </>
                ) : (
                  <>
                    Sign In
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none"
                      stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <line x1="5" y1="12" x2="19" y2="12"/>
                      <polyline points="12 5 19 12 12 19"/>
                    </svg>
                  </>
                )}
              </button>
            </form>

            <p className="lp-register">
              New here?{' '}
              <a href="#" className="lp-register__link">Create an Account</a>
            </p>
          </div>
        </div>

      </div>
    </div>
  );
};

export default Login;
