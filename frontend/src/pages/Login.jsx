import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
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
      <div className="card p-4 shadow" style={{ minWidth: 320 }}>
        <h3 className="mb-3 text-center">PPPK Exam Login</h3>
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label htmlFor="userID" className="form-label">User ID</label>
            <input
              type="text"
              className="form-control"
              id="userID"
              value={userID}
              onChange={e => setUserID(e.target.value)}
              required
              autoFocus
            />
          </div>
          {error && <div className="alert alert-danger py-2">{error}</div>}
          <button type="submit" className="btn btn-primary w-100" disabled={loading}>
            {loading ? 'Loading...' : 'Create/Start Exam'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
