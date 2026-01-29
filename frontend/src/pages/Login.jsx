import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { examAPI } from '../services/api';

const Login = () => {
  const [userID, setUserID] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const response = await examAPI.getOrCreateExam(userID);
      if (response.data.success) {
        navigate(`/exam/${userID}`);
      } else {
        setError(response.data.error || 'Failed to start exam session');
      }
    } catch (err) {
      setError('Connection failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container d-flex justify-content-center align-items-center" style={{ minHeight: '70vh' }}>
      <div className="row w-100 justify-content-center">
        <div className="col-md-6 col-lg-4">
          <div className="card p-4 shadow">
            <div className="text-center mb-4">
              <i className="bi bi-clipboard-check text-primary" style={{fontSize: '3rem'}}></i>
              <h3 className="mb-1">PPPK Exam System</h3>
              <p className="text-muted">Enter your User ID to start</p>
            </div>
            
            <form onSubmit={handleSubmit}>
              <div className="mb-3">
                <label htmlFor="userID" className="form-label">User ID</label>
                <div className="input-group">
                  <span className="input-group-text">
                    <i className="bi bi-person"></i>
                  </span>
                  <input
                    type="text"
                    className="form-control"
                    id="userID"
                    value={userID}
                    onChange={e => setUserID(e.target.value)}
                    placeholder="Enter your User ID"
                    required
                    autoFocus
                  />
                </div>
              </div>
              
              {error && (
                <div className="alert alert-danger py-2">
                  <i className="bi bi-exclamation-triangle"></i> {error}
                </div>
              )}
              
              <button 
                type="submit" 
                className="btn btn-primary w-100 mb-3" 
                disabled={loading || !userID.trim()}
              >
                {loading ? (
                  <>
                    <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                    Loading...
                  </>
                ) : (
                  <>
                    <i className="bi bi-play-circle"></i> Start Exam
                  </>
                )}
              </button>
            </form>
            
            <hr />
            
            <div className="text-center">
              <p className="text-muted mb-2">
                <small>Admin Access</small>
              </p>
              <Link 
                to="/admin" 
                className="btn btn-outline-secondary btn-sm"
              >
                <i className="bi bi-speedometer2"></i> View Dashboard
              </Link>
            </div>
          </div>
          
          <div className="text-center mt-3">
            <small className="text-muted">
              4 questions • 120 minutes • Auto-save enabled
            </small>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
