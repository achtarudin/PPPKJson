import React, { useState } from 'react';
import { examAPI } from '../services/api';

const DebugAnswers = () => {
  const [userID, setUserID] = useState('1234');
  const [answers, setAnswers] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const testGetAnswers = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await examAPI.getUserAnswers(userID);
      console.log('Raw response:', response);
      
      if (response.data.success) {
        setAnswers(response.data.data);
      } else {
        setError(response.data.error || 'Failed to get answers');
      }
    } catch (error) {
      setError('Network error: ' + error.message);
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mt-5">
      <div className="card">
        <div className="card-body">
          <h3>Debug User Answers</h3>
          
          <div className="mb-3">
            <label>User ID:</label>
            <input 
              type="text" 
              className="form-control"
              value={userID}
              onChange={(e) => setUserID(e.target.value)}
            />
          </div>
          
          <button 
            className="btn btn-primary" 
            onClick={testGetAnswers}
            disabled={loading}
          >
            {loading ? 'Loading...' : 'Get User Answers'}
          </button>
          
          {error && (
            <div className="alert alert-danger mt-3">
              {error}
            </div>
          )}
          
          {answers && (
            <div className="mt-3">
              <h5>Answers:</h5>
              <pre>{JSON.stringify(answers, null, 2)}</pre>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default DebugAnswers;