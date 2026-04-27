// frontend/src/components/SearchBar.jsx
// Drop-in search + filter bar for ViewTasksPage (or any task list)

function SearchBar({ searchTerm, setSearchTerm, statusFilter, setStatusFilter,
  priorityFilter, setPriorityFilter, sortBy, setSortBy,
  onFilter, onClear }) {

  return (
    <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
      <input
        type="text"
        placeholder="Search tasks..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        onKeyDown={(e) => { if (e.key === 'Enter') onFilter(); }}
        style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc', flex: '1', minWidth: '150px' }}
        aria-label="Search tasks"
      />
      <select
        value={statusFilter}
        onChange={(e) => setStatusFilter(e.target.value)}
        style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}
        aria-label="Filter by status"
      >
        <option value="">All Status</option>
        <option value="open">Open</option>
        <option value="claimed">Claimed</option>
        <option value="in_progress">In Progress</option>
        <option value="done">Done</option>
      </select>
      <select
        value={priorityFilter}
        onChange={(e) => setPriorityFilter(e.target.value)}
        style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}
        aria-label="Filter by priority"
      >
        <option value="">All Priority</option>
        <option value="High">High</option>
        <option value="Medium">Medium</option>
        <option value="Low">Low</option>
      </select>
      <select
        value={sortBy}
        onChange={(e) => setSortBy(e.target.value)}
        style={{ padding: '0.5rem', borderRadius: '6px', border: '1px solid #ccc' }}
        aria-label="Sort tasks"
      >
        <option value="">Default Sort</option>
        <option value="deadline">By Deadline</option>
        <option value="priority">By Priority</option>
        <option value="newest">Newest First</option>
        <option value="oldest">Oldest First</option>
      </select>
      <button onClick={onFilter} className="primary-button" style={{ padding: '0.5rem 1rem' }}>
        Filter
      </button>
      <button
        onClick={onClear}
        style={{ padding: '0.5rem 1rem', background: '#6b7280', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }}
      >
        Clear
      </button>
    </div>
  );
}

export default SearchBar;
