import React from 'react';
import ReactDOM from 'react-dom/client';
<<<<<<< HEAD
=======
import { BrowserRouter } from 'react-router-dom';
import { AuthProvider } from './src/context/AuthContext';
import { TaskProvider } from './src/context/TaskContext';
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
<<<<<<< HEAD
    <App />
  </React.StrictMode>
);
=======
    <BrowserRouter>
      <AuthProvider>
        <TaskProvider>
          <App />
        </TaskProvider>
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
);
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f
