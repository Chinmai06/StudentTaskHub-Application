import { useEffect, useState } from 'react';
import TaskCard from '../components/TaskCard';
<<<<<<< HEAD
import { useAuth } from './context/AuthContext';
=======
import { useAuth } from '../context/AuthContext';
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f
import { api } from '../utils/api';

function HomePage() {
  const { user } = useAuth();
  const [myTasks, setMyTasks] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (user) {
      api.getTasks(`?created_by=${user.username}`)
        .then((data) => setMyTasks(data || []))
        .catch((err) => console.error(err))
        .finally(() => setLoading(false));
    }
  }, [user]);

  return (
    <section className="content-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Home</p>
          <h2>My Created Tasks</h2>
        </div>
        <span className="count-chip">{myTasks.length} task(s)</span>
      </div>

      {loading ? (
        <p>Loading tasks...</p>
      ) : myTasks.length === 0 ? (
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
