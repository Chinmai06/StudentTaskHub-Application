import { createContext, useContext, useMemo, useState } from 'react';
import { mockTasks } from './data/mockTasks';
import { createTaskPayload } from './utils/taskHelpers';

const TaskContext = createContext(null);
const STORAGE_KEY = 'student-task-hub-tasks';

export function TaskProvider({ children }) {
  const [tasks, setTasks] = useState(() => {
    const savedTasks = localStorage.getItem(STORAGE_KEY);
    return savedTasks ? JSON.parse(savedTasks) : mockTasks;
  });

  const addTask = (formData, user) => {
    const nextTask = createTaskPayload(formData, user);
    setTasks((currentTasks) => {
      const updatedTasks = [nextTask, ...currentTasks];
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedTasks));
      return updatedTasks;
    });
    return nextTask;
  };

  const value = useMemo(() => ({ tasks, addTask }), [tasks]);

  return <TaskContext.Provider value={value}>{children}</TaskContext.Provider>;
}

export function useTasks() {
  return useContext(TaskContext);
}
