import { createContext, useContext, useMemo, useState } from 'react';
import { api } from '../utils/api';

const AuthContext = createContext(null);
const STORAGE_KEY = 'student-task-hub-user';

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const savedUser = localStorage.getItem(STORAGE_KEY);
    return savedUser ? JSON.parse(savedUser) : null;
  });

  const login = async (username, password) => {
    const data = await api.login(username, password);
    const nextUser = { username: data.username };
    setUser(nextUser);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(nextUser));
    return data;
  };

  const register = async (username, email, password) => {
    const data = await api.register(username, email, password);
    if (data.error) throw new Error(data.error);
    return data;
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem(STORAGE_KEY);
  };

  const value = useMemo(
    () => ({ user, login, register, logout, isAuthenticated: Boolean(user) }),
    [user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}
