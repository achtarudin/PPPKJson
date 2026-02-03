import React, { useState, useEffect } from 'react';
import QuestionManager from '../pages/QuestionManager';

const ProtectedQuestions = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  // Check if already authenticated in this session
  useEffect(() => {
    const authenticated = sessionStorage.getItem('questionsAuth');
    if (authenticated === 'true') {
      setIsAuthenticated(true);
    }
  }, []);

  const handlePasswordSubmit = (e) => {
    e.preventDefault();
    const correctPassword = 'moncos214';
    
    if (password === correctPassword) {
      setIsAuthenticated(true);
      setError('');
      // Store authentication in session storage
      sessionStorage.setItem('questionsAuth', 'true');
    } else {
      setError('Incorrect password. Please try again.');
      setPassword('');
    }
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    sessionStorage.removeItem('questionsAuth');
    setPassword('');
    setError('');
  };

  if (isAuthenticated) {
    return (
      <div>
        <div className="d-flex justify-content-end mb-3">
          <button 
            onClick={handleLogout}
            className="btn btn-outline-secondary btn-sm"
          >
            Logout
          </button>
        </div>
        <QuestionManager />
      </div>
    );
  }

  return (
    <div className="container mt-5">
      <div className="row justify-content-center">
        <div className="col-md-6">
          <div className="card">
            <div className="card-header">
              <h4 className="mb-0">Access Protected Area</h4>
            </div>
            <div className="card-body">
              <p className="text-muted mb-4">
                This area is password protected. Please enter the password to continue.
              </p>
              
              <form onSubmit={handlePasswordSubmit}>
                <div className="mb-3">
                  <label htmlFor="password" className="form-label">
                    Password
                  </label>
                  <input
                    type="password"
                    className={`form-control ${error ? 'is-invalid' : ''}`}
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Enter password"
                    required
                    autoFocus
                  />
                  {error && (
                    <div className="invalid-feedback">
                      {error}
                    </div>
                  )}
                </div>
                
                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">
                    Access Questions
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProtectedQuestions;