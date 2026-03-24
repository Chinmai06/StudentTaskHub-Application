import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';

function NavBar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <header className="topbar">
      <div>
        <p className="brand-mark">H HUB</p>
        <h1 className="brand-title">StudentTaskHub</h1>
      </div>
      <nav className="nav-links" aria-label="Primary navigation">
        <NavLink to="/home">Home</NavLink>
        <NavLink to="/create-task">Create Task</NavLink>
        <NavLink to="/view-tasks">View Tasks</NavLink>
        <NavLink to="/profile">Profile</NavLink>
      </nav>
      <div className="user-chip-wrap">
        <span className="user-chip">Signed in: {user?.email}</span>
        <button className="secondary-button" onClick={handleLogout}>
          Logout
        </button>
      </div>
    </header>
  );
}

export default NavBar;
