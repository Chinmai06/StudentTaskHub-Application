import { useAuth } from './context/AuthContext';
import { useTasks } from '../context/TaskContext';

function TaskCard({ task }) {
  const { user } = useAuth();
  const { claimTask, updateTaskStatus, deleteTask } = useTasks();

  const isCreator = user && user.username === task.created_by;
  const isClaimer = user && user.username === task.claimed_by;
  const canUpdateStatus = isCreator || isClaimer;

  const handleClaim = () => {
    if (user) claimTask(task.id, user.username);
  };

  const handleStatusChange = (event) => {
    if (user) updateTaskStatus(task.id, user.username, event.target.value);
  };

  const handleDelete = () => {
    if (user && window.confirm('Are you sure you want to delete this task?')) {
      deleteTask(task.id, user.username);
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <article className="task-card" data-testid="task-card">
      <div className="task-card-header">
        <div>
          <span className={`status-badge status-${task.status}`}>{task.status}</span>
          <h3>{task.title}</h3>
        </div>
        {task.priority === 'high' ? (
          <span className="priority-badge priority-high">High Priority</span>
        ) : null}
      </div>

      <p className="task-description">{task.description}</p>

      <div className="task-meta-grid">
        <span><strong>Deadline:</strong> {task.deadline}</span>
        <span><strong>Created by:</strong> {task.created_by}</span>
        {task.claimed_by ? (
          <span><strong>Claimed by:</strong> {task.claimed_by}</span>
        ) : null}
        <span><strong>Created:</strong> {formatDate(task.created_at)}</span>
      </div>

      <div style={{ display: 'flex', gap: '0.5rem', marginTop: '0.75rem', flexWrap: 'wrap' }}>
        {task.status === 'open' && user && !isCreator ? (
          <button onClick={handleClaim} className="primary-button" style={{ fontSize: '0.85rem', padding: '0.4rem 0.8rem' }}>
            Claim Task
          </button>
        ) : null}

        {canUpdateStatus && task.status !== 'done' ? (
          <select
            value={task.status}
            onChange={handleStatusChange}
            style={{ padding: '0.4rem', borderRadius: '6px', border: '1px solid #ccc' }}
          >
            <option value="open">Open</option>
            <option value="claimed">Claimed</option>
            <option value="in_progress">In Progress</option>
            <option value="done">Done</option>
          </select>
        ) : null}

        {isCreator ? (
          <button
            onClick={handleDelete}
            style={{
              fontSize: '0.85rem', padding: '0.4rem 0.8rem',
              background: '#ef4444', color: 'white', border: 'none',
              borderRadius: '6px', cursor: 'pointer'
            }}
          >
            Delete
          </button>
        ) : null}
      </div>
    </article>
  );
}

export default TaskCard;
