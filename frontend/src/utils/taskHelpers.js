export function createTaskPayload(formData, user) {
  return {
    id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    title: formData.title.trim(),
    category: formData.category.trim(),
    description: formData.description.trim(),
    location: formData.location.trim(),
    priority: formData.priority,
    dueDate: formData.dueDate,
    createdBy: user.email,
    creatorUfid: user.ufid,
    createdAt: new Date().toISOString()
  };
}

export function filterTasksCreatedBy(tasks, email) {
  return tasks.filter((task) => task.createdBy === email);
}

export function sortTasksByPriority(tasks) {
  const priorityOrder = { High: 1, Medium: 2, Low: 3 };
  return [...tasks].sort(
    (left, right) => priorityOrder[left.priority] - priorityOrder[right.priority]
  );
}

export function formatPostedDate(dateString) {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}
