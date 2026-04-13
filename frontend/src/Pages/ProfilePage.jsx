import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { api } from '../utils/api';

function ProfilePage() {
  const { user } = useAuth();
  const [myTaskCount, setMyTaskCount] = useState(0);
  const [claimedCount, setClaimedCount] = useState(0);

  useEffect(() => {
    if (user) {
      api.getTasks(`?created_by=${user.username}`)
        .then((data) => setMyTaskCount((data || []).length))
        .catch(() => {});

      api.getTasks(`?claimed_by=${user.username}`)
        .then((data) => setClaimedCount((data || []).length))
        .catch(() => {});
    }
  }, [user]);

  return (
    <section className="content-card profile-card">
      <p className="eyebrow">Profile</p>
      <h2>{user.username}</h2>
      <div className="profile-grid">
        <div className="stat-card">
          <span>Tasks Created</span>
          <strong>{myTaskCount}</strong>
        </div>
        <div className="stat-card">
          <span>Tasks Claimed</span>
          <strong>{claimedCount}</strong>
        </div>
      </div>
    </section>
  );
}

export default ProfilePage;
