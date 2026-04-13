const API_URL = "http://localhost:8080/api";

export const api = {
  // ============ User Endpoints ============

  async register(username, email, password) {
    const res = await fetch(`${API_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, email, password }),
    });
    return await res.json();
  },

  async login(username, password) {
    const res = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Login failed");
    return data;
  },

  // ============ Task Endpoints ============

  async getTasks(filters = "") {
    const res = await fetch(`${API_URL}/tasks${filters}`);
    return await res.json();
  },

  async getTask(id) {
    const res = await fetch(`${API_URL}/tasks/${id}`);
    return await res.json();
  },

  async createTask(taskData) {
    const res = await fetch(`${API_URL}/tasks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(taskData),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Failed to create task");
    return data;
  },

  async updateTask(id, username, taskData) {
    const res = await fetch(`${API_URL}/tasks/${id}?username=${username}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(taskData),
    });
    return await res.json();
  },

  async deleteTask(id, username) {
    const res = await fetch(`${API_URL}/tasks/${id}?username=${username}`, {
      method: "DELETE",
    });
    return await res.json();
  },

  async claimTask(id, username) {
    const res = await fetch(`${API_URL}/tasks/${id}/claim`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ claimed_by: username }),
    });
    return await res.json();
  },

  async updateTaskStatus(id, username, status) {
    const res = await fetch(`${API_URL}/tasks/${id}/status?username=${username}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ status }),
    });
    return await res.json();
  },

  async searchTasks(keyword) {
    const res = await fetch(`${API_URL}/tasks?search=${keyword}`);
    return await res.json();
  },
};
