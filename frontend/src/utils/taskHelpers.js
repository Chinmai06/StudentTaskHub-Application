export function filterTasksCreatedBy(tasks, username) {
  return tasks.filter((task) => task.created_by === username);
}

export function sortTasksByPriority(tasks) {
  const priorityOrder = { high: 1, normal: 2 };
  return [...tasks].sort(
    (left, right) => (priorityOrder[left.priority] || 3) - (priorityOrder[right.priority] || 3)
  );
}

export function formatPostedDate(dateString) {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}
