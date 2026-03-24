import { useAuth } from './context/AuthContext';
import { useTasks } from './context/TaskContext';
import { filterTasksCreatedBy } from './utils/taskHelpers';

function ProfilePage() {
  const { user } = useAuth();
  const { tasks } = useTasks();
  const myTaskCount = filterTasksCreatedBy(tasks, user.email).length;

  return (
    <section className="content-card profile-card">
      <p className="eyebrow">Profile</p>
      <h2>{user.email}</h2>
      <div className="profile-grid">
        <div className="stat-card">
          <span>UFID</span>
          <strong>{user.ufid}</strong>
        </div>
        <div className="stat-card">
          <span>Tasks created</span>
          <strong>{myTaskCount}</strong>
        </div>
      </div>
    </section>
  );
}

export default ProfilePage;
