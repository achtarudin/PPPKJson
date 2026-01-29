import React, { useState, useEffect } from 'react';
import { examAPI } from '../services/api';

const UserDetailModal = ({ userID, show, onHide }) => {
  const [userData, setUserData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (show && userID) {
      loadUserDetail();
    }
  }, [show, userID]);

  const loadUserDetail = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await examAPI.getUserDashboard(userID);
      
      if (response.data.success) {
        setUserData(response.data.data);
      } else {
        setError(response.data.error || 'Failed to load user details');
      }
    } catch (error) {
      setError('Connection failed. Please try again.');
      console.error('Load user detail error:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatDateTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('id-ID', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  const getStatusBadge = (status) => {
    const badges = {
      'NO_EXAM': 'badge bg-secondary',
      'NOT_STARTED': 'badge bg-warning text-dark',
      'IN_PROGRESS': 'badge bg-primary',
      'COMPLETED': 'badge bg-success',
      'EXPIRED': 'badge bg-danger'
    };
    return badges[status] || 'badge bg-secondary';
  };

  const getGradeBadge = (grade) => {
    const badges = {
      'A': 'badge bg-success',
      'B': 'badge bg-info',
      'C': 'badge bg-warning text-dark',
      'D': 'badge bg-warning text-dark',
      'E': 'badge bg-danger'
    };
    return badges[grade] || 'badge bg-secondary';
  };

  if (!show) return null;

  return (
    <div className="modal fade show d-block" tabIndex="-1" style={{backgroundColor: 'rgba(0,0,0,0.5)'}}>
      <div className="modal-dialog modal-lg">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title">
              <i className="bi bi-person-circle"></i> User Dashboard - {userID}
            </h5>
            <button type="button" className="btn-close" onClick={onHide}></button>
          </div>
          
          <div className="modal-body">
            {loading ? (
              <div className="text-center py-4">
                <div className="spinner-border" role="status">
                  <span className="visually-hidden">Loading...</span>
                </div>
                <p className="mt-2">Loading user details...</p>
              </div>
            ) : error ? (
              <div className="alert alert-danger">
                <i className="bi bi-exclamation-triangle"></i> {error}
                <button className="btn btn-outline-danger ms-3 btn-sm" onClick={loadUserDetail}>
                  Retry
                </button>
              </div>
            ) : userData ? (
              <div>
                {/* Basic Info */}
                <div className="card mb-3">
                  <div className="card-header">
                    <h6 className="mb-0">
                      <i className="bi bi-info-circle"></i> Basic Information
                    </h6>
                  </div>
                  <div className="card-body">
                    <div className="row">
                      <div className="col-md-6">
                        <p><strong>User ID:</strong> {userData.user_id}</p>
                        <p><strong>Has Exam:</strong> {userData.has_exam ? 'Yes' : 'No'}</p>
                        <p>
                          <strong>Status:</strong> 
                          <span className={getStatusBadge(userData.exam_status)} style={{marginLeft: '8px'}}>
                            {userData.exam_status.replace('_', ' ')}
                          </span>
                        </p>
                      </div>
                      {userData.exam_session && (
                        <div className="col-md-6">
                          <p><strong>Session Code:</strong> {userData.exam_session.session_code}</p>
                          <p><strong>Duration:</strong> {userData.exam_session.duration} minutes</p>
                          <p><strong>Expires At:</strong> {formatDateTime(userData.exam_session.expires_at)}</p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* Progress Info - for IN_PROGRESS or NOT_STARTED */}
                {userData.progress_info && (
                  <div className="card mb-3">
                    <div className="card-header">
                      <h6 className="mb-0">
                        <i className="bi bi-clock-history"></i> Progress Information
                      </h6>
                    </div>
                    <div className="card-body">
                      <div className="row">
                        <div className="col-md-4">
                          <div className="text-center">
                            <h4 className="text-primary">{userData.progress_info.answered_questions}</h4>
                            <small className="text-muted">Answered</small>
                          </div>
                        </div>
                        <div className="col-md-4">
                          <div className="text-center">
                            <h4 className="text-info">{userData.progress_info.total_questions}</h4>
                            <small className="text-muted">Total Questions</small>
                          </div>
                        </div>
                        <div className="col-md-4">
                          <div className="text-center">
                            <h4 className={userData.progress_info.remaining_time_minutes > 10 ? 'text-success' : 'text-danger'}>
                              {userData.progress_info.remaining_time_minutes}
                            </h4>
                            <small className="text-muted">Minutes Left</small>
                          </div>
                        </div>
                      </div>
                      
                      <div className="mt-3">
                        <div className="progress">
                          <div 
                            className="progress-bar" 
                            role="progressbar" 
                            style={{
                              width: `${(userData.progress_info.answered_questions / userData.progress_info.total_questions) * 100}%`
                            }}
                          >
                            {Math.round((userData.progress_info.answered_questions / userData.progress_info.total_questions) * 100)}%
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Results - for COMPLETED */}
                {userData.exam_results && (
                  <div className="card mb-3">
                    <div className="card-header">
                      <h6 className="mb-0">
                        <i className="bi bi-trophy"></i> Exam Results
                      </h6>
                    </div>
                    <div className="card-body">
                      {/* Overall Summary */}
                      <div className="row mb-4">
                        <div className="col-md-12">
                          <h6 className="text-primary">Overall Summary</h6>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <h5 className="text-primary">{userData.exam_results.summary.total_answered}</h5>
                            <small>Answered</small>
                          </div>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <h5>{userData.exam_results.summary.total_questions}</h5>
                            <small>Total</small>
                          </div>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <h5 className="text-success">{userData.exam_results.summary.total_score}</h5>
                            <small>Score</small>
                          </div>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <h5>{userData.exam_results.summary.overall_percentage.toFixed(1)}%</h5>
                            <small>Percentage</small>
                          </div>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <span className={getGradeBadge(userData.exam_results.summary.overall_grade)} style={{fontSize: '16px', padding: '8px 12px'}}>
                              {userData.exam_results.summary.overall_grade}
                            </span>
                            <br />
                            <small>Grade</small>
                          </div>
                        </div>
                        <div className="col-md-2">
                          <div className="text-center">
                            <span className={`badge ${userData.exam_results.summary.is_passed ? 'bg-success' : 'bg-danger'}`} style={{fontSize: '14px', padding: '8px'}}>
                              {userData.exam_results.summary.is_passed ? 'PASSED' : 'FAILED'}
                            </span>
                            <br />
                            <small>Result</small>
                          </div>
                        </div>
                      </div>

                      {/* Category Results */}
                      <div>
                        <h6 className="text-primary mb-3">Results by Category</h6>
                        <div className="table-responsive">
                          <table className="table table-sm table-striped">
                            <thead className="table-light">
                              <tr>
                                <th>Category</th>
                                <th>Score</th>
                                <th>Percentage</th>
                                <th>Grade</th>
                                <th>Status</th>
                              </tr>
                            </thead>
                            <tbody>
                              {userData.exam_results.results_by_category.map((result) => (
                                <tr key={result.category}>
                                  <td>
                                    <span className="badge bg-info">{result.category}</span>
                                  </td>
                                  <td>{result.total_score}/{result.max_score}</td>
                                  <td>{result.percentage.toFixed(1)}%</td>
                                  <td>
                                    <span className={getGradeBadge(result.grade)}>
                                      {result.grade}
                                    </span>
                                  </td>
                                  <td>
                                    <span className={`badge ${result.is_passed ? 'bg-success' : 'bg-danger'}`}>
                                      {result.is_passed ? 'PASS' : 'FAIL'}
                                    </span>
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </div>
                      </div>

                      <div className="mt-3 text-muted">
                        <small>
                          <i className="bi bi-calendar-check"></i> 
                          Completed: {formatDateTime(userData.exam_results.summary.completed_at)}
                        </small>
                      </div>
                    </div>
                  </div>
                )}

                {/* No Exam State */}
                {userData.exam_status === 'NO_EXAM' && (
                  <div className="alert alert-info text-center">
                    <i className="bi bi-info-circle"></i>
                    <h6>No Exam Session</h6>
                    <p className="mb-0">This user hasn't created any exam session yet.</p>
                  </div>
                )}
              </div>
            ) : (
              <div className="text-center py-4">
                <p className="text-muted">No data available</p>
              </div>
            )}
          </div>
          
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onHide}>
              Close
            </button>
            {userData && userData.exam_status === 'COMPLETED' && (
              <button 
                type="button" 
                className="btn btn-outline-primary"
                onClick={() => window.open(`/results/${userID}`, '_blank')}
              >
                <i className="bi bi-file-earmark-text"></i> View Full Results
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserDetailModal;