import React, { useState, useEffect } from 'react';
import { questionAPI } from '../services/api';

const QuestionManager = () => {
  const [questions, setQuestions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [searchText, setSearchText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Modal states
  const [showModal, setShowModal] = useState(false);
  const [selectedQuestion, setSelectedQuestion] = useState(null);
  const [updatingScore, setUpdatingScore] = useState(false);

  // Server-side pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [pagination, setPagination] = useState({
    currentPage: 1,
    itemsPerPage: 10,
    totalItems: 0,
    totalPages: 0
  });

  useEffect(() => {
    loadCategories();
  }, []);

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      loadQuestions();
    }, 300); // Debounce search by 300ms

    return () => clearTimeout(timeoutId);
  }, [selectedCategory, searchText, currentPage, itemsPerPage]);

  const loadCategories = async () => {
    try {
      const response = await questionAPI.getCategories();
      if (response.data.success) {
        setCategories(response.data.data || []);
      }
    } catch (error) {
      console.error('Failed to load categories:', error);
      setError('Failed to load categories');
    }
  };

  const loadQuestions = async () => {
    try {
      setLoading(true);
      setError('');
      
      const requestParams = {
        category: selectedCategory,
        search: searchText,
        page: currentPage,
        limit: itemsPerPage === 'all' ? 0 : itemsPerPage
      };
      
      console.log('API Request params:', requestParams); // Debug log
      
      const response = await questionAPI.getQuestionsByCategory(requestParams);
      
      if (response.data.success) {
        console.log('API Response:', response.data.data); // Debug log
        setQuestions(response.data.data.questions);
        setPagination(response.data.data.pagination);
      } else {
        setError(response.data.message || 'Failed to load questions');
      }
    } catch (error) {
      console.error('Failed to load questions:', error);
      setError('Failed to load questions. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const updateScore = async (questionId, optionId, newScore) => {
    try {
      setError('');
      setSuccess('');
      setUpdatingScore(true);
      
      const response = await questionAPI.updateOptionScore(questionId, optionId, newScore);
      
      if (response.data.success) {
        // Update local state
        setQuestions(prevQuestions =>
          prevQuestions.map(question =>
            question.id === questionId
              ? {
                  ...question,
                  options: question.options.map(option =>
                    option.id === optionId ? { ...option, score: newScore } : option
                  )
                }
              : question
          )
        );

        // Update modal state if it's the same question
        if (selectedQuestion && selectedQuestion.id === questionId) {
          setSelectedQuestion(prev => ({
            ...prev,
            options: prev.options.map(option =>
              option.id === optionId ? { ...option, score: newScore } : option
            )
          }));
        }
        
        setSuccess('Score updated successfully!');
        setTimeout(() => setSuccess(''), 3000);
      } else {
        setError(response.data.error || 'Failed to update score');
      }
    } catch (error) {
      console.error('Failed to update score:', error);
      setError('Failed to update score. Please try again.');
    } finally {
      setUpdatingScore(false);
    }
  };

  const openEditModal = (question) => {
    setSelectedQuestion(question);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setSelectedQuestion(null);
  };

  const getCategoryBadgeClass = (category) => {
    const categoryClasses = {
      'TEKNIS': 'bg-warning text-dark',
      'MANAJERIAL': 'bg-primary',
      'SOSIAL KULTURAL': 'bg-success',
      'WAWANCARA': 'bg-info'
    };
    return categoryClasses[category] || 'bg-secondary';
  };

  const getScoreColorClass = (score) => {
    if (score >= 8) return 'text-success';     // 8-10: Green (Excellent)
    if (score >= 6) return 'text-primary';    // 6-7: Blue (Good) 
    if (score >= 4) return 'text-info';       // 4-5: Cyan (Fair)
    if (score >= 2) return 'text-warning';    // 2-3: Yellow (Poor)
    return 'text-danger';                     // 0-1: Red (Very Poor)
  };

  const truncateText = (text, maxLength = 80) => {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  const handleItemsPerPageChange = (newItemsPerPage) => {
    setItemsPerPage(newItemsPerPage);
    setCurrentPage(1);
  };

  const getPaginationRange = () => {
    if (!pagination || !pagination.totalPages || pagination.totalPages <= 1) {
      console.log('No pagination range needed:', pagination); // Debug log
      return [];
    }
    
    const totalPages = pagination.totalPages;
    const range = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage - 2);
    let end = Math.min(totalPages, start + maxVisible - 1);
    
    if (end - start < maxVisible - 1) {
      start = Math.max(1, end - maxVisible + 1);
    }
    
    for (let i = start; i <= end; i++) {
      range.push(i);
    }
    
    console.log('Pagination range:', range); // Debug log
    return range;
  };

  // Reset to page 1 when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [selectedCategory, searchText, itemsPerPage]);

  return (
    <div className="container-fluid py-4">
      <div className="row">
        <div className="col-12">
          <div className="d-flex justify-content-between align-items-center mb-4">
            <h2>Question Management</h2>
          </div>

          {/* Category Filter & Search */}
          <div className="card mb-4">
            <div className="card-body">
              <div className="row align-items-center">
                <div className="col-md-2">
                  <label htmlFor="categorySelect" className="form-label fw-bold">
                    Filter by Category:
                  </label>
                </div>
                <div className="col-md-3">
                  <select
                    id="categorySelect"
                    className="form-select"
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                  >
                    <option value="">All Categories</option>
                    {categories.map(category => (
                      <option key={category} value={category}>
                        {category}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="col-md-2">
                  <label htmlFor="searchInput" className="form-label fw-bold">
                    Search Question:
                  </label>
                </div>
                <div className="col-md-4">
                  <input
                    id="searchInput"
                    type="text"
                    className="form-control"
                    placeholder="Search by question text..."
                    value={searchText}
                    onChange={(e) => setSearchText(e.target.value)}
                  />
                </div>
                <div className="col-md-1">
                  <button
                    className="btn btn-primary"
                    onClick={loadQuestions}
                    disabled={loading}
                  >
                    {loading ? (
                      <>
                        <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                        Loading...
                      </>
                    ) : (
                      'Refresh'
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Alert Messages */}
          {error && (
            <div className="alert alert-danger alert-dismissible fade show" role="alert">
              {error}
              <button type="button" className="btn-close" onClick={() => setError('')}></button>
            </div>
          )}

          {success && (
            <div className="alert alert-success alert-dismissible fade show" role="alert">
              {success}
              <button type="button" className="btn-close" onClick={() => setSuccess('')}></button>
            </div>
          )}

          {/* Questions Table */}
          <div className="card">
            <div className="card-header">
              <div className="d-flex justify-content-between align-items-center">
                <h5 className="mb-0">Questions List</h5>
                <span className="badge bg-primary">{pagination.totalItems} Questions</span>
              </div>
            </div>
            <div className="card-body">
              {loading ? (
                <div className="text-center py-5">
                  <div className="spinner-border" role="status">
                    <span className="visually-hidden">Loading questions...</span>
                  </div>
                  <p className="mt-3 text-muted">Loading questions...</p>
                </div>
              ) : (
                <>
                  {questions.length === 0 ? (
                    <div className="alert alert-info" role="alert">
                      {searchText
                        ? `No questions found matching "${searchText}"${selectedCategory ? ` in category "${selectedCategory}"` : ''}`
                        : selectedCategory 
                        ? `No questions found for category "${selectedCategory}"` 
                        : 'No questions found'}
                    </div>
                  ) : (
                    <>
                      <div className="table-responsive">
                        <table className="table table-hover">
                          <thead className="table-light">
                            <tr>
                              <th style={{ width: '80px' }}>No</th>
                              <th style={{ width: '150px' }}>Category</th>
                              <th>Question Text</th>
                              <th style={{ width: '120px' }}>Options Count</th>
                              <th style={{ width: '150px' }}>Actions</th>
                            </tr>
                          </thead>
                          <tbody>
                            {questions.map((question, index) => {
                              const sequentialNumber = (currentPage - 1) * (itemsPerPage === 'all' ? questions.length : itemsPerPage) + index + 1;
                              return (
                                <tr key={question.id}>
                                  <td className="fw-bold">{sequentialNumber}</td>
                              <td>
                                <span className={`badge ${getCategoryBadgeClass(question.category)}`}>
                                  {question.category}
                                </span>
                              </td>
                              <td>
                                <div title={question.question_text}>
                                  {truncateText(question.question_text)}
                                </div>
                              </td>
                              <td>
                                <span className="badge bg-secondary">
                                  {question.options?.length || 0} options
                                </span>
                              </td>
                              <td>
                                <button
                                  className="btn btn-sm btn-outline-primary"
                                  onClick={() => openEditModal(question)}
                                  title="Edit Scores"
                                >
                                  <i className="bi bi-pencil-square"></i> Edit Scores
                                </button>
                              </td>
                            </tr>
                              );
                            })}
                          </tbody>
                        </table>
                      </div>
                    
                      {/* Pagination Navigation with Items Per Page */}
                      <div className="d-flex justify-content-between align-items-center mt-4">
                        <div className="d-flex align-items-center">
                          <span className="text-muted small">
                            {pagination && pagination.totalItems > 0 ? (
                              `Showing ${(currentPage - 1) * (itemsPerPage === 'all' ? pagination.totalItems : itemsPerPage) + 1} to ${Math.min(currentPage * (itemsPerPage === 'all' ? pagination.totalItems : itemsPerPage), pagination.totalItems)} of ${pagination.totalItems} questions`
                            ) : (
                              'Showing 0 to 0 of 0 questions'
                            )}
                          </span>
                        </div>
                        
                        <div className="d-flex align-items-center gap-3">
                          {/* Items per page selector */}
                          <div className="d-flex align-items-center">
                            <label className="form-label fw-bold mb-0 me-2">Items per page:</label>
                            <select
                              className="form-select form-select-sm"
                              style={{ width: 'auto' }}
                              value={itemsPerPage}
                              onChange={(e) => handleItemsPerPageChange(e.target.value === 'all' ? 'all' : parseInt(e.target.value))}
                            >
                              <option value={10}>10</option>
                              <option value={20}>20</option>
                              <option value={25}>25</option>
                              <option value={75}>75</option>
                              <option value="all">All</option>
                            </select>
                          </div>
                          
                          {/* Pagination navigation */}
                          {(() => {
                            console.log('Pagination debug:', { pagination, itemsPerPage, totalPages: pagination?.totalPages }); // Debug log
                            return pagination && pagination.totalPages > 1 && itemsPerPage !== 'all';
                          })() && (
                            <nav aria-label="Questions pagination">
                              <ul className="pagination pagination-sm mb-0">
                                <li className={`page-item ${currentPage <= 1 ? 'disabled' : ''}`}>
                                  <button 
                                    className="page-link"
                                    onClick={() => handlePageChange(currentPage - 1)}
                                    disabled={currentPage <= 1}
                                  >
                                    Previous
                                  </button>
                                </li>
                                
                                {getPaginationRange().map(page => (
                                  <li key={page} className={`page-item ${currentPage === page ? 'active' : ''}`}>
                                    <button 
                                      className="page-link"
                                      onClick={() => handlePageChange(page)}
                                    >
                                      {page}
                                    </button>
                                  </li>
                                ))}
                                
                                <li className={`page-item ${currentPage >= (pagination?.totalPages || 0) ? 'disabled' : ''}`}>
                                  <button 
                                    className="page-link"
                                    onClick={() => handlePageChange(currentPage + 1)}
                                    disabled={currentPage >= (pagination?.totalPages || 0)}
                                  >
                                    Next
                                  </button>
                                </li>
                              </ul>
                            </nav>
                          )}
                        </div>
                      </div>
                    </>
                  )}
                </>
              )}
            </div>
          </div>

          {/* Summary */}
          {questions.length > 0 && (
            <div className="card mt-4">
              <div className="card-body">
                <div className="row text-center">
                  <div className="col-md-3">
                    <h4 className="text-primary">{questions.length}</h4>
                    <p className="text-muted mb-0">Found Questions</p>
                  </div>
                  <div className="col-md-3">
                    <h4 className="text-success">
                      {questions.reduce((sum, q) => sum + (q.options?.length || 0), 0)}
                    </h4>
                    <p className="text-muted mb-0">Total Options</p>
                  </div>
                  <div className="col-md-3">
                    <h4 className="text-info">
                      {selectedCategory || 'All'}
                    </h4>
                    <p className="text-muted mb-0">Current Category</p>
                  </div>
                  <div className="col-md-3">
                    <h4 className="text-warning">{categories.length}</h4>
                    <p className="text-muted mb-0">Total Categories</p>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Edit Score Modal */}
      {showModal && (
        <>
          {/* Modal Backdrop */}
          <div 
            className="modal-backdrop fade show" 
            style={{ zIndex: 1040 }}
            onClick={closeModal}
          ></div>
          
          {/* Modal Content */}
          <div 
            className="modal fade show" 
            style={{ 
              display: 'block',
              zIndex: 1050,
              position: 'fixed',
              top: 0,
              left: 0,
              width: '100%',
              height: '100%',
              pointerEvents: 'none'
            }} 
            tabIndex="-1"
          >
            <div 
              className="modal-dialog modal-lg" 
              style={{ 
                pointerEvents: 'auto',
                margin: '1.75rem auto'
              }}
            >
              <div className="modal-content">
                <div className="modal-header">
                  <h5 className="modal-title">
                    Edit Scores - Question #{selectedQuestion.id}
                    <span className={`badge ${getCategoryBadgeClass(selectedQuestion.category)} ms-2`}>
                      {selectedQuestion.category}
                    </span>
                  </h5>
                  <button 
                    type="button" 
                    className="btn-close" 
                    onClick={closeModal}
                  ></button>
                </div>
                <div className="modal-body">
                  <div className="mb-4">
                    <h6 className="fw-bold">Question:</h6>
                    <p className="text-muted">{selectedQuestion.question_text}</p>
                  </div>

                  <h6 className="fw-bold mb-3">Options & Scores:</h6>
                  <div className="list-group">
                    {selectedQuestion.options?.map((option, index) => (
                      <div key={option.id} className="list-group-item">
                        <div className="d-flex justify-content-between align-items-start mb-2">
                          <span className="badge bg-secondary me-2">
                            Option {String.fromCharCode(65 + index)}
                          </span>
                          <span className="fw-bold">
                            Score: {option.score}
                          </span>
                        </div>
                        
                        <p className="text-muted mb-3" style={{fontSize: '0.9rem'}}>
                          {option.option_text}
                        </p>

                        <div className="d-flex align-items-center justify-content-between">
                          <label className="form-label fw-bold mb-0" style={{fontSize: '0.9rem'}}>
                            Update Score:
                          </label>
                          <select
                            className="form-select form-select-sm"
                            style={{ width: 'auto', minWidth: '80px' }}
                            value={option.score}
                            onChange={(e) => updateScore(selectedQuestion.id, option.id, parseInt(e.target.value))}
                            disabled={updatingScore}
                          >
                            <option value={0}>0</option>
                            <option value={1}>1</option>
                            <option value={2}>2</option>
                            <option value={3}>3</option>
                            <option value={4}>4</option>
                            <option value={5}>5</option>
                            <option value={6}>6</option>
                            <option value={7}>7</option>
                            <option value={8}>8</option>
                            <option value={9}>9</option>
                            <option value={10}>10</option>
                          </select>
                        </div>
                      </div>
                    ))}
                  </div>

                  {updatingScore && (
                    <div className="text-center mt-3">
                      <div className="spinner-border spinner-border-sm text-primary" role="status">
                        <span className="visually-hidden">Updating score...</span>
                      </div>
                      <span className="ms-2 text-muted">Updating score...</span>
                    </div>
                  )}
                </div>
                <div className="modal-footer">
                  <button 
                    type="button" 
                    className="btn btn-secondary" 
                    onClick={closeModal}
                    disabled={updatingScore}
                  >
                    Close
                  </button>
                </div>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default QuestionManager;