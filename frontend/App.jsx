import { Navigate, Route, Routes } from 'react-router-dom';
<<<<<<< HEAD
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import CreateTaskPage from './pages/CreateTaskPage';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import ProfilePage from './pages/ProfilePage';
import ViewTasksPage from './pages/ViewTasksPage';
=======
import Layout from './src/components/layout';
import ProtectedRoute from './src/components/ProtectedRoute';
import CreateTaskPage from './src/Pages/CreateTaskPage';
import HomePage from './src/Pages/HomePage';
import LoginPage from './src/Pages/LoginPage';
import ProfilePage from './src/Pages/ProfilePage';
import ViewTasksPage from './src/Pages/ViewTasksPage';
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f

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
