import React, { useState, useEffect } from 'react';
import { examAPI } from '../services/api';

const CategoryAnswerDetails = ({ userID, category, onClose }) => {
  const [detailedAnswers, setDetailedAnswers] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchDetailedAnswers = async () => {
      try {
        const response = await examAPI.getDetailedUserAnswers(userID);
        if (response.data.success) {
          setDetailedAnswers(response.data.data);
        } else {
          setError(response.data.error || 'Failed to fetch detailed answers');
        }
      } catch (err) {
        setError('Connection failed. Please try again.');
        console.error('Fetch detailed answers error:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchDetailedAnswers();
  }, [userID]);

  if (loading) {
    return (
      <div className="modal fade show d-block" tabIndex="-1" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
        <div className="modal-dialog modal-lg">
          <div className="modal-content">
            <div className="modal-body text-center">
              <div className="spinner-border" role="status">
                <span className="visually-hidden">Loading...</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="modal fade show d-block" tabIndex="-1" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
        <div className="modal-dialog modal-lg">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">Error</h5>
              <button type="button" className="btn-close" onClick={onClose}></button>
            </div>
            <div className="modal-body">
              <div className="alert alert-danger" role="alert">
                {error}
              </div>
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" onClick={onClose}>Close</button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const categoryAnswers = detailedAnswers?.[category] || [];

  const getCategoryBadgeClass = (cat) => {
    switch (cat) {
      case 'MANAJERIAL': return 'bg-primary';
      case 'SOSIAL_KULTURAL': return 'bg-success';
      case 'TEKNIS': return 'bg-warning';
      case 'WAWANCARA': return 'bg-info';
      default: return 'bg-secondary';
    }
  };

  const getScoreColor = (score, maxScore) => {
    const percentage = (score / maxScore) * 100;
    if (percentage >= 75) return 'text-success';
    if (percentage >= 50) return 'text-warning';
    return 'text-danger';
  };

  return (
    <div className="modal fade show d-block" tabIndex="-1" style={{ backgroundColor: 'rgba(0,0,0,0.5)' }}>
      <div className="modal-dialog modal-xl">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title">
              <span className={`badge ${getCategoryBadgeClass(category)} me-2`}>
                {category}
              </span>
              Detail Soal dan Jawaban
            </h5>
            <button type="button" className="btn-close" onClick={onClose}></button>
          </div>
          <div className="modal-body" style={{ maxHeight: '70vh', overflowY: 'auto' }}>
            {categoryAnswers.length === 0 ? (
              <div className="alert alert-warning">
                Tidak ada jawaban ditemukan untuk kategori {category}.
              </div>
            ) : (
              <div className="row">
                {categoryAnswers.map((answer, index) => (
                  <div key={answer.exam_question_id} className="col-12 mb-4">
                    <div className="card">
                      <div className="card-header d-flex justify-content-between align-items-center">
                        <div className="d-flex align-items-center">
                          <span className="fw-bold me-2">Soal {index + 1}</span>
                          {answer.is_correct ? (
                            <span className="badge bg-success">
                              <i className="bi bi-check-circle-fill me-1"></i>Benar
                            </span>
                          ) : (
                            <span className="badge bg-danger">
                              <i className="bi bi-x-circle-fill me-1"></i>Salah
                            </span>
                          )}
                        </div>
                        <div>
                          <span className={`badge ${getScoreColor(answer.score, answer.max_score)} fs-6`}>
                            Score: {answer.score}/{answer.max_score}
                          </span>
                          <small className="text-muted ms-2">
                            {new Date(answer.answered_at).toLocaleString('id-ID')}
                          </small>
                        </div>
                      </div>
                      <div className="card-body">
                        <div className="mb-3">
                          <h6 className="text-muted">Pertanyaan:</h6>
                          <p className="mb-0">{answer.question_text}</p>
                        </div>
                        
                        <div className="mb-3">
                          <h6 className="text-muted">Jawaban yang Dipilih:</h6>
                          <div className={`alert border ${answer.is_correct ? 'alert-success' : 'alert-danger'}`}>
                            <div className="d-flex justify-content-between align-items-center">
                              <div className="flex-grow-1">
                                <div className="d-flex align-items-center">
                                  {answer.is_correct ? (
                                    <i className="bi bi-check-circle-fill text-success me-2"></i>
                                  ) : (
                                    <i className="bi bi-x-circle-fill text-danger me-2"></i>
                                  )}
                                  <span>{answer.selected_option}</span>
                                </div>
                              </div>
                              <span className={`badge ${getScoreColor(answer.score, answer.max_score)}`}>
                                +{answer.score} poin
                              </span>
                            </div>
                          </div>
                        </div>

                        {/* Show correct answer if user answered incorrectly */}
                        {!answer.is_correct && (
                          <div className="mb-3">
                            <h6 className="text-muted">
                              <i className="bi bi-lightbulb-fill text-warning me-1"></i>
                              Jawaban Yang Benar:
                            </h6>
                            <div className="alert alert-success border">
                              <div className="d-flex justify-content-between align-items-center">
                                <div className="flex-grow-1">
                                  <div className="d-flex align-items-center">
                                    <i className="bi bi-check-circle-fill text-success me-2"></i>
                                    <span><strong>{answer.correct_option}</strong></span>
                                  </div>
                                </div>
                                <span className="badge bg-success">
                                  +{answer.correct_score} poin
                                </span>
                              </div>
                            </div>
                          </div>
                        )}
                        
                        <div className="row">
                          <div className="col-md-6">
                            <small className="text-muted">
                              <i className="bi bi-clock"></i> Dijawab pada: {' '}
                              {new Date(answer.answered_at).toLocaleString('id-ID', {
                                day: '2-digit',
                                month: 'short',
                                year: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                              })}
                            </small>
                          </div>
                          <div className="col-md-6 text-end">
                            <small className="text-muted">
                              Question ID: {answer.question_id}
                            </small>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CategoryAnswerDetails;