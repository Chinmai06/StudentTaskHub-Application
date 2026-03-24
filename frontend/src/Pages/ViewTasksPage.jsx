import TaskCard from '../components/TaskCard';
import { useTasks } from '../context/TaskContext';
import { sortTasksByPriority } from '../utils/taskHelpers';

function ViewTasksPage() {
  const { tasks } = useTasks();
  const sortedTasks = sortTasksByPriority(tasks);

  return (
    <section className="content-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Explore</p>
          <h2>All Shared Tasks</h2>
        </div>
        <span className="count-chip">{sortedTasks.length} total</span>
      </div>

      <div className="task-grid">
        {sortedTasks.map((task) => (
          <TaskCard key={task.id} task={task} />
        ))}
      </div>
    </section>
  );
}

export default ViewTasksPage;
