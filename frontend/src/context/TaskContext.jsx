import { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { api } from '../utils/api';

const TaskContext = createContext(null);

export function TaskProvider({ children }) {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchTasks = async (filters = "") => {
    try {
      setLoading(true);
      const data = await api.getTasks(filters);
      setTasks(data || []);
    } catch (err) {
      console.error("Failed to fetch tasks:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  const addTask = async (taskData) => {
    try {
      const newTask = await api.createTask(taskData);
      if (newTask.error) throw new Error(newTask.error);
      await fetchTasks();
      return newTask;
    } catch (err) {
      throw err;
    }
  };

  const deleteTask = async (taskId, username) => {
    try {
      await api.deleteTask(taskId, username);
      await fetchTasks();
    } catch (err) {
      console.error("Failed to delete task:", err);
    }
  };

  const claimTask = async (taskId, username) => {
    try {
      await api.claimTask(taskId, username);
      await fetchTasks();
    } catch (err) {
      console.error("Failed to claim task:", err);
    }
  };

  const updateTaskStatus = async (taskId, username, status) => {
    try {
      await api.updateTaskStatus(taskId, username, status);
      await fetchTasks();
    } catch (err) {
      console.error("Failed to update status:", err);
    }
  };

  const value = useMemo(
    () => ({ tasks, loading, fetchTasks, addTask, deleteTask, claimTask, updateTaskStatus }),
    [tasks, loading]
  );

  return <TaskContext.Provider value={value}>{children}</TaskContext.Provider>;
}

export function useTasks() {
  return useContext(TaskContext);
}
