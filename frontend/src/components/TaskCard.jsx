import { formatPostedDate } from '../utils/taskHelpers';

function TaskCard({ task }) {
  return (
    <article className="task-card" data-testid="task-card">
      <div className="task-card-header">
        <div>
          <p className="task-category">{task.category}</p>
          <h3>{task.title}</h3>
        </div>
        <span className={`priority-badge priority-${task.priority.toLowerCase()}`}>
          {task.priority}
        </span>
      </div>
      <p className="task-description">{task.description}</p>
      <div className="task-meta-grid">
        <span><strong>Location:</strong> {task.location}</span>
        <span><strong>Due date:</strong> {task.dueDate}</span>
        <span><strong>Posted:</strong> {formatPostedDate(task.createdAt)}</span>
        <span><strong>Created by:</strong> {task.createdBy}</span>
      </div>
    </article>
  );
}

export default TaskCard;
