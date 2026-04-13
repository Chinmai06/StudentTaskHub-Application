import { Navigate, Route, Routes } from 'react-router-dom';
import Layout from './src/components/layout';
import ProtectedRoute from './src/components/ProtectedRoute';
import CreateTaskPage from './src/Pages/CreateTaskPage';
import HomePage from './src/Pages/HomePage';
import LoginPage from './src/Pages/LoginPage';
import ProfilePage from './src/Pages/ProfilePage';
import ViewTasksPage from './src/Pages/ViewTasksPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route path="/home" element={<HomePage />} />
        <Route path="/create-task" element={<CreateTaskPage />} />
        <Route path="/view-tasks" element={<ViewTasksPage />} />
        <Route path="/profile" element={<ProfilePage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default App;
