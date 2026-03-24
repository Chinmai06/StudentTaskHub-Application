import TaskCard from '../components/TaskCard';
import { useAuth } from './context/AuthContext';
import { useTasks } from './context/TaskContext';
import { filterTasksCreatedBy, sortTasksByPriority } from '../utils/taskHelpers';

function HomePage() {
  const { user } = useAuth();
  const { tasks } = useTasks();
  const myTasks = sortTasksByPriority(filterTasksCreatedBy(tasks, user.email));

  return (
    <section className="content-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Home</p>
          <h2>My Created Tasks</h2>
        </div>
        <span className="count-chip">{myTasks.length} task(s)</span>
      </div>

      {myTasks.length === 0 ? (
        <div className="empty-state">
          <h3>No tasks yet</h3>
          <p>Create your first task to see it here.</p>
        </div>
      ) : (
        <div className="task-grid">
          {myTasks.map((task) => (
            <TaskCard key={task.id} task={task} />
          ))}
        </div>
      )}
    </section>
  );
}

export default HomePage;
