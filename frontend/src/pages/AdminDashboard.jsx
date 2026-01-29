import React, { useState, useEffect } from 'react';
import { examAPI } from '../services/api';
import UserDetailModal from '../components/UserDetailModal';

const AdminDashboard = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedUser, setSelectedUser] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('ALL');

  useEffect(() => {
    loadUsers();
    // Auto-refresh every 30 seconds
    const interval = setInterval(loadUsers, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadUsers = async () => {
    try {
      const response = await examAPI.getAllUsersDashboard();
      if (response.data.success) {
        setUsers(response.data.data.users);
      } else {
        setError(response.data.error || 'Failed to load users');
      }
    } catch (error) {
      setError('Connection failed. Please try again.');
      console.error('Load users error:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleViewDetails = (userID) => {
    setSelectedUser(userID);
    setShowModal(true);
  };

  const getStatusBadge = (status) => {
    const badgeClasses = {
      'COMPLETED': 'badge bg-success',
      'IN_PROGRESS': 'badge bg-primary', 
      'EXPIRED': 'badge bg-danger',
      'NOT_STARTED': 'badge bg-warning text-dark'
    };
    return badgeClasses[status] || 'badge bg-secondary';
  };

  const formatDateTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('id-ID', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const filteredUsers = users.filter(user => {
    const matchesSearch = user.user_id.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'ALL' || user.exam_status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  if (loading) {
    return (
      <div className="container mt-5 text-center">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
        <p className="mt-2">Loading users data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-5">
        <div className="alert alert-danger" role="alert">
          {error}
          <button className="btn btn-outline-danger ms-3" onClick={loadUsers}>
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container-fluid mt-4">
      <div className="row">
        <div className="col-12">
          <div className="d-flex justify-content-between align-items-center mb-4">
            <h2>PPPK Exam Admin Dashboard</h2>
            <button 
              className="btn btn-outline-primary" 
              onClick={loadUsers}
              disabled={loading}
            >
              <i className="bi bi-arrow-clockwise"></i> Refresh
            </button>
          </div>

          {/* Filters */}
          <div className="card mb-4">
            <div className="card-body">
              <div className="row">
                <div className="col-md-6">
                  <label htmlFor="search" className="form-label">Search User ID</label>
                  <input
                    id="search"
                    type="text"
                    className="form-control"
                    placeholder="Enter user ID..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                  />
                </div>
                <div className="col-md-4">
                  <label htmlFor="status" className="form-label">Filter by Status</label>
                  <select
                    id="status"
                    className="form-select"
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                  >
                    <option value="ALL">All Status</option>
                    <option value="COMPLETED">Completed</option>
                    <option value="IN_PROGRESS">In Progress</option>
                    <option value="NOT_STARTED">Not Started</option>
                    <option value="EXPIRED">Expired</option>
                  </select>
                </div>
                <div className="col-md-2 d-flex align-items-end">
                  <div className="text-muted">
                    <small>{filteredUsers.length} of {users.length} users</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Users Table */}
          <div className="card">
            <div className="card-header">
              <h5 className="mb-0">Exam Participants ({users.length} total)</h5>
            </div>
            <div className="card-body p-0">
              <div className="table-responsive">
                <table className="table table-striped table-hover mb-0">
                  <thead className="table-dark">
                    <tr>
                      <th>User ID</th>
                      <th>Status</th>
                      <th>Session Code</th>
                      <th>Started At</th>
                      <th>Completed At</th>
                      <th>Score</th>
                      <th>Grade</th>
                      <th>Result</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredUsers.length === 0 ? (
                      <tr>
                        <td colSpan="9" className="text-center py-4 text-muted">
                          No users found
                        </td>
                      </tr>
                    ) : (
                      filteredUsers.map((user) => (
                        <tr key={user.user_id}>
                          <td>
                            <strong>{user.user_id}</strong>
                          </td>
                          <td>
                            <span className={getStatusBadge(user.exam_status)}>
                              {user.exam_status.replace('_', ' ')}
                            </span>
                          </td>
                          <td>
                            <small className="text-muted">{user.session_code}</small>
                          </td>
                          <td>{formatDateTime(user.started_at)}</td>
                          <td>{formatDateTime(user.completed_at)}</td>
                          <td>
                            {user.total_score !== null ? (
                              <span>
                                {user.total_score}/{user.max_score}
                                <small className="text-muted"> ({user.percentage?.toFixed(1)}%)</small>
                              </span>
                            ) : (
                              '-'
                            )}
                          </td>
                          <td>
                            {user.grade && (
                              <span className={`badge ${
                                user.grade === 'A' ? 'bg-success' :
                                user.grade === 'B' ? 'bg-info' :
                                user.grade === 'C' ? 'bg-warning text-dark' :
                                user.grade === 'D' ? 'bg-warning text-dark' :
                                'bg-danger'
                              }`}>
                                {user.grade}
                              </span>
                            )}
                          </td>
                          <td>
                            {user.is_passed !== null && (
                              <span className={`badge ${user.is_passed ? 'bg-success' : 'bg-danger'}`}>
                                {user.is_passed ? 'PASSED' : 'FAILED'}
                              </span>
                            )}
                          </td>
                          <td>
                            <button
                              className="btn btn-sm btn-outline-primary"
                              onClick={() => handleViewDetails(user.user_id)}
                            >
                              <i className="bi bi-eye"></i> Details
                            </button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          {/* Summary Stats */}
          <div className="row mt-4">
            <div className="col-md-3">
              <div className="card bg-primary text-white">
                <div className="card-body">
                  <h5>{users.filter(u => u.exam_status === 'COMPLETED').length}</h5>
                  <p className="mb-0">Completed</p>
                </div>
              </div>
            </div>
            <div className="col-md-3">
              <div className="card bg-info text-white">
                <div className="card-body">
                  <h5>{users.filter(u => u.exam_status === 'IN_PROGRESS').length}</h5>
                  <p className="mb-0">In Progress</p>
                </div>
              </div>
            </div>
            <div className="col-md-3">
              <div className="card bg-warning text-dark">
                <div className="card-body">
                  <h5>{users.filter(u => u.exam_status === 'NOT_STARTED').length}</h5>
                  <p className="mb-0">Not Started</p>
                </div>
              </div>
            </div>
            <div className="col-md-3">
              <div className="card bg-danger text-white">
                <div className="card-body">
                  <h5>{users.filter(u => u.exam_status === 'EXPIRED').length}</h5>
                  <p className="mb-0">Expired</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* User Detail Modal */}
      {showModal && (
        <UserDetailModal
          userID={selectedUser}
          show={showModal}
          onHide={() => setShowModal(false)}
        />
      )}
    </div>
  );
};

export default AdminDashboard;