import { useState } from 'react';
import TaskCard from '../components/TaskCard';
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

      <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
        <input
          type="text"
          placeholder="Search tasks..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          onKeyDown={(e) => { if (e.key === 'Enter') handleFilter(); }}
          style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc', flex: '1', minWidth: '150px' }}
        />
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}
          style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}>
          <option value="">All Status</option>
          <option value="open">Open</option>
          <option value="claimed">Claimed</option>
          <option value="in_progress">In Progress</option>
          <option value="done">Done</option>
        </select>
        <select value={priorityFilter} onChange={(e) => setPriorityFilter(e.target.value)}
          style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}>
          <option value="">All Priority</option>
          <option value="high">High</option>
          <option value="normal">Normal</option>
        </select>
        <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}
          style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}>
          <option value="">Default Sort</option>
          <option value="deadline">By Deadline</option>
          <option value="priority">By Priority</option>
          <option value="newest">Newest First</option>
          <option value="oldest">Oldest First</option>
        </select>
        <button onClick={handleFilter} className="primary-button" style={{ padding: '0.5rem 1rem' }}>
          Filter
        </button>
        <button onClick={handleClearFilters}
          style={{ padding: '0.5rem 1rem', background: '#6b7280', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }}>
          Clear
        </button>
      </div>

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
