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
  const [downloading, setDownloading] = useState(false);

  // Modal states
  const [showModal, setShowModal] = useState(false);
  const [selectedQuestion, setSelectedQuestion] = useState(null);
  const [selectedTextNo, setSelectedTextNo] = useState(null);

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

        setError('Failed to load categories');
    }
  };

  const downloadQuestions = async () => {
    try {
      setDownloading(true);
      setError('');
      
      const response = await questionAPI.downloadQuestionsJSON({
        category: selectedCategory,
        search: searchText
      });
      
      // Create filename based on filters
      let filename = 'questions';
      if (selectedCategory) {
        filename += `_${selectedCategory}`;
      }
      if (searchText) {
        filename += '_search';
      }
      filename += '.json';
      
      // Create blob and download
      const blob = new Blob([JSON.stringify(response.data, null, 2)], {
        type: 'application/json'
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      setSuccess(`Questions downloaded successfully as ${filename}`);
    } catch (error) {
      setError('Failed to download questions');
      console.error('Download error:', error);
    } finally {
      setDownloading(false);
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
      
      
      const response = await questionAPI.getQuestionsByCategory(requestParams);
      
      if (response.data.success) {
        const responseData = response.data.data;
        const questionsList = responseData.questions || responseData || [];
        setQuestions(questionsList);
        
        // Ensure pagination data is properly set
        let newPagination;
        if (responseData.pagination) {
            newPagination = {
							currentPage: responseData.pagination.current_page,
							itemsPerPage: itemsPerPage === 'all' ? 0 : itemsPerPage,
							totalItems: responseData.pagination.total_items,
							totalPages: responseData.pagination.total_pages
          };
        } else {
          // Fallback pagination - assume there might be more data
          const totalItems = questionsList.length;
          // If we got exactly itemsPerPage questions, assume there might be more
          const estimatedTotal = (questionsList.length === itemsPerPage && itemsPerPage !== 'all') ? totalItems * 2 : totalItems;
          const totalPages = itemsPerPage === 'all' || itemsPerPage === 0 ? 1 : Math.ceil(estimatedTotal / itemsPerPage);

          newPagination = {
            currentPage: currentPage,
            itemsPerPage: itemsPerPage === 'all' ? 0 : itemsPerPage,
            totalItems: estimatedTotal,
            totalPages: totalPages
          };

        }

        setPagination(newPagination);
      } else {
        setError(response.data.message || 'Failed to load questions');
      }
    } catch (error) {
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
      setError('Failed to update score. Please try again.');
    } finally {
      setUpdatingScore(false);
    }
  };

  const openEditModal = (question, textNo) => {
    setSelectedQuestion(question);
    setSelectedTextNo(textNo);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setSelectedQuestion(null);
    setSelectedTextNo(null);
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
    const totalPages = pagination?.totalPages || 0;
    
    if (totalPages <= 1) {
      return [];
    }
    
    const range = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage - 2);
    let end = Math.min(totalPages, start + maxVisible - 1);
    
    // Adjust start if we don't have enough pages at the end
    if (end - start < maxVisible - 1) {
      start = Math.max(1, end - maxVisible + 1);
    }
    
    for (let i = start; i <= end; i++) {
      range.push(i);
    }
    
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

              <div className='row justify-content-between align-items-end '>
                <div className='col-5'>
                  <label htmlFor="categorySelect" className="form-label fw-bold">
                    Filter by Category:
                  </label>
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

                <div className='col-5'>
                  <label htmlFor="searchInput" className="form-label fw-bold">
                    Search Question:
                  </label>
                  <input
                    id="searchInput"
                    type="text"
                    className="form-control"
                    placeholder="Search by question text..."
                    value={searchText}
                    onChange={(e) => setSearchText(e.target.value)}
                  />
                </div>

                 <div className="col-2">
                  <div className='d-flex justify-content-between'>
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

                    <button
                      className="btn btn-success"
                      onClick={downloadQuestions}
                      disabled={downloading || loading}
                    >
                      {downloading ? (
                        <>
                          <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                          Downloading...
                        </>
                      ) : (
                        <>
                          <i className="bi bi-download me-1"></i>
                          Download JSON
                        </>
                      )}
                    </button>
                  </div>
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
		
          <div className='d-flex justify-content-end mb-4'>
						<div className="d-flex align-items-end">
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
									<option value={0}>All</option>
							</select>
					</div>
					</div>
					

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
                              <th style={{ width: '50px' }}>No</th>
                              <th style={{ width: '150px' }}>Category</th>
                              <th>Question Text</th>
                              <th style={{ width: '120px' }}>Options Count</th>
                              <th style={{ width: '250px' }}>Actions</th>
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
                                      {truncateText(question.question_text, 100)}
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
                                      onClick={() => openEditModal(question, index + 1)}
                                      title="Edit Scores"
                                    >
                                      <i className="bi bi-pencil-square"></i> Edit Scores (Question {index + 1})  
                                    </button>
                                  </td>
                            </tr>
                              );
                            })}
                          </tbody>
                        </table>
                      </div>
                    
                    
                    
                    </>
                  )}
                </>
              )}
            </div>
          </div>
					
					<div className="row align-items-center mt-4 gap-3">
						<div className="col"></div>

						<div className="col-auto">
							<div className="d-flex align-items-center gap-3">
								{/* Pagination navigation - Show when we have multiple pages */}
								{((pagination?.totalPages || 0) > 1 && itemsPerPage !== 'all') || 
								(questions.length >= itemsPerPage && itemsPerPage !== 'all' && itemsPerPage < 50) ? ( // Show pagination if we have full page of results
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
										{/* Generate page numbers 1, 2, 3, etc. */}
										{Array.from({ length: Math.max(pagination?.totalPages || 2, 2) }, (_, i) => i + 1).map(page => (
												<li key={page} className={`page-item ${currentPage === page ? 'active' : ''}`}>
												<button 
														className="page-link"
														onClick={() => handlePageChange(page)}
												>
														{page}
												</button>
												</li>
										))}
										
										<li className={`page-item ${currentPage >= Math.max(pagination?.totalPages || 2, 2) ? 'disabled' : ''}`}>
												<button 
												className="page-link"
												onClick={() => handlePageChange(currentPage + 1)}
												disabled={currentPage >= Math.max(pagination?.totalPages || 2, 2)}
												>
												Next
												</button>
										</li>
										</ul>
								</nav>
								
								) : null}
							</div>
						</div>

						<div className="col text-end">
							<div id="div2" className=" d-inline-block">
								<div className="d-flex align-items-end">
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
											<option value={0}>All</option>
									</select>
								</div>
							</div>
						</div>
					</div>
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
                    Edit Scores - Text No #{selectedTextNo}  
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