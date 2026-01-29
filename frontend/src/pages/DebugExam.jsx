import React, { useState, useEffect } from 'react';
import { examAPI } from '../services/api';

const DebugExam = () => {
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const testAPI = async () => {
    setLoading(true);
    setError(null);
    
    try {
      console.log('Testing API call...');
      const response = await examAPI.getOrCreateExam('debuguser');
      console.log('API Response:', response);
      
      setResult(JSON.stringify(response.data, null, 2));
    } catch (err) {
      console.error('API Error:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    testAPI();
  }, []);

  return (
    <div className="container mt-5">
      <h1>API Debug Page</h1>
      
      <button onClick={testAPI} disabled={loading} className="btn btn-primary">
        {loading ? 'Testing...' : 'Test API'}
      </button>
      
      <div className="mt-3">
        {error && (
          <div className="alert alert-danger">
            <strong>Error:</strong> {error}
          </div>
        )}
        
        {result && (
          <div className="alert alert-success">
            <h5>API Response:</h5>
            <pre>{result}</pre>
          </div>
        )}
        
        {loading && (
          <div className="alert alert-info">
            Loading...
          </div>
        )}
      </div>
    </div>
  );
};

export default DebugExam;