import { Outlet } from 'react-router-dom';
import NavBar from './NavBar';

function Layout() {
  return (
    <div className="app-shell">
      <NavBar />
      <main className="page-shell">
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;
