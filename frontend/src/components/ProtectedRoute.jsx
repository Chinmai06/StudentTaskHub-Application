import { Navigate } from 'react-router-dom';
<<<<<<< HEAD
import { useAuth } from './context/AuthContext';
=======
import { useAuth } from '../context/AuthContext';
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return children;
}

export default ProtectedRoute;
