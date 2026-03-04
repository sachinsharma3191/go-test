import './Stats.css'

function Stats({ stats }) {
  if (!stats) {
    return (
      <div className="stats-container">
        <div className="empty-state">
          <div className="empty-state-icon">📊</div>
          <h3>No Statistics Available</h3>
          <p>Statistics will appear once users and tasks are created.</p>
        </div>
      </div>
    )
  }

  const totalUsers = stats.users?.total || 0
  const totalTasks = stats.tasks?.total || 0
  const pendingTasks = stats.tasks?.pending || 0
  const inProgressTasks = stats.tasks?.inProgress || 0
  const completedTasks = stats.tasks?.completed || 0

  if (totalUsers === 0 && totalTasks === 0) {
    return (
      <div className="stats-container">
        <div className="empty-state">
          <div className="empty-state-icon">📊</div>
          <h3>No Data Available</h3>
          <p>Create users and tasks to see statistics here.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="stats-container">
      <div className="stat-card">
        <div className="stat-icon">👥</div>
        <div className="stat-content">
          <h3>{totalUsers}</h3>
          <p>Total Users</p>
        </div>
      </div>

      <div className="stat-card">
        <div className="stat-icon">📋</div>
        <div className="stat-content">
          <h3>{totalTasks}</h3>
          <p>Total Tasks</p>
        </div>
      </div>

      <div className="stat-card pending">
        <div className="stat-icon">⏳</div>
        <div className="stat-content">
          <h3>{pendingTasks}</h3>
          <p>Pending</p>
        </div>
      </div>

      <div className="stat-card in-progress">
        <div className="stat-icon">🔄</div>
        <div className="stat-content">
          <h3>{inProgressTasks}</h3>
          <p>In Progress</p>
        </div>
      </div>

      <div className="stat-card completed">
        <div className="stat-icon">✅</div>
        <div className="stat-content">
          <h3>{completedTasks}</h3>
          <p>Completed</p>
        </div>
      </div>
    </div>
  )
}

export default Stats
