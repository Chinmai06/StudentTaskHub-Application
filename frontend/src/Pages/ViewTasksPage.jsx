// frontend/src/Pages/ViewTasksPage.jsx
import { useState } from 'react';
import TaskCard from '../components/TaskCard';
import SearchBar from '../components/SearchBar';
import { useTasks } from '../context/TaskContext';

function ViewTasksPage() {
  const { tasks, loading, fetchTasks } = useTasks();
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [priorityFilter, setPriorityFilter] = useState('');
  const [sortBy, setSortBy] = useState('');

  const handleFilter = () => {
    let filters = '?';
    if (searchTerm.trim()) filters += `search=${searchTerm.trim()}&`;
    if (statusFilter) filters += `status=${statusFilter}&`;
    if (priorityFilter) filters += `priority=${priorityFilter}&`;
    if (sortBy) filters += `sort=${sortBy}&`;
    fetchTasks(filters.slice(0, -1));
  };

  const handleClearFilters = () => {
    setSearchTerm('');
    setStatusFilter('');
    setPriorityFilter('');
    setSortBy('');
    fetchTasks();
  };

  return (
    <section className="content-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Explore</p>
          <h2>All Shared Tasks</h2>
        </div>
        <span className="count-chip">{tasks.length} total</span>
      </div>

      <SearchBar
        searchTerm={searchTerm}
        setSearchTerm={setSearchTerm}
        statusFilter={statusFilter}
        setStatusFilter={setStatusFilter}
        priorityFilter={priorityFilter}
        setPriorityFilter={setPriorityFilter}
        sortBy={sortBy}
        setSortBy={setSortBy}
        onFilter={handleFilter}
        onClear={handleClearFilters}
      />

      {loading ? (
        <p>Loading tasks...</p>
      ) : tasks.length === 0 ? (
        <div className="empty-state">
          <h3>No tasks found</h3>
          <p>Try a different search or create a new task.</p>
        </div>
      ) : (
        <div className="task-grid">
          {tasks.map((task) => (
            <TaskCard key={task.id} task={task} />
          ))}
        </div>
      )}
    </section>
  );
}

export default ViewTasksPage;
