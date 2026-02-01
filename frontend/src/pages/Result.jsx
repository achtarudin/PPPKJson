import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { examAPI } from '../services/api';
import CategoryAnswerDetails from '../components/CategoryAnswerDetails';

const Result = () => {
  const { userID } = useParams();
  const navigate = useNavigate();
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState(null);

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const response = await examAPI.getResults(userID);
        if (response.data.success) {
          setResult(response.data.data);
        } else {
          setError(response.data.error || 'Failed to fetch results');
        }
      } catch (err) {
        setError('Connection failed. Please try again.');
      } finally {
        setLoading(false);
      }
    };
    fetchResults();
  }, [userID]);

  if (loading) {
    return (
      <div className="container mt-5 text-center">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-5">
        <div className="alert alert-danger" role="alert">
          {error}
        </div>
      </div>
    );
  }

  if (!result) return null;

  const { summary, results_by_category } = result;

  const handleShowDetail = (category) => {
    setSelectedCategory(category);
    setShowDetailModal(true);
  };

  const handleCloseDetail = () => {
    setShowDetailModal(false);
    setSelectedCategory(null);
  };

  return (
    <div className="container mt-5">
      <div className="card mx-auto" style={{ maxWidth: 600 }}>
        <div className="card-body">
          <h3 className="mb-3 text-center">Exam Completed</h3>
          <div className="mb-3">
            <strong>Overall Grade:</strong> {summary.overall_grade} <br />
            <strong>Score:</strong> {summary.total_score} / {summary.max_score} <br />
            <strong>Percentage:</strong> {summary.overall_percentage}% <br />
            <strong>Status:</strong> {summary.is_passed ? (
              <span className="badge bg-success">Passed</span>
            ) : (
              <span className="badge bg-danger">Failed</span>
            )}
          </div>
          <h5>Category Breakdown</h5>
          <table className="table table-bordered">
            <thead>
              <tr>
                <th>Category</th>
                <th>Score</th>
                <th>Grade</th>
                <th>Passed</th>
                <th>Detail</th>
              </tr>
            </thead>
            <tbody>
              {results_by_category.map(cat => (
                <tr key={cat.category}>
                  <td>{cat.category}</td>
                  <td>{cat.total_score} / {cat.max_score}</td>
                  <td>{cat.grade}</td>
                  <td>{cat.is_passed ? <span className="text-success">Yes</span> : <span className="text-danger">No</span>}</td>
                  <td>
                    <button
                      className="btn btn-sm btn-outline-primary"
                      onClick={() => handleShowDetail(cat.category)}
                      title="Lihat detail soal dan jawaban"
                    >
                      <i className="bi bi-eye"></i> Lihat Detail
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <button className="btn btn-secondary w-100 mt-3" onClick={() => navigate('/')}>Back to Home</button>
        </div>
      </div>
      
      {/* Detail Modal */}
      {showDetailModal && selectedCategory && (
        <CategoryAnswerDetails
          userID={userID}
          category={selectedCategory}
          onClose={handleCloseDetail}
        />
      )}
    </div>
  );
};

export default Result;
