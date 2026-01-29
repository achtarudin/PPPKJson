import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { examAPI } from '../services/api';
import QuestionCard from '../components/QuestionCard';
import Timer from '../components/Timer';

const ExamBoard = () => {
  const { userID } = useParams();
  const navigate = useNavigate();
  
  const [examData, setExamData] = useState(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState({}); // examQuestionId -> optionId
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [examStarted, setExamStarted] = useState(false);

  useEffect(() => {
    loadExamSession();
    // eslint-disable-next-line
  }, [userID]);

  // Load existing answers when exam data is available and exam is in progress
  useEffect(() => {
    if (examData && examData.status === 'IN_PROGRESS') {
      loadExistingAnswers();
    }
    // eslint-disable-next-line
  }, [examData]);

  const loadExamSession = async () => {
    try {
      setLoading(true);
      const response = await examAPI.getOrCreateExam(userID);
      
      if (response.data.success) {
        const exam = response.data.data;
        setExamData(exam);
        setExamStarted(exam.status === 'IN_PROGRESS');
        
        // If exam is already completed or expired, redirect to results
        if (exam.status === 'COMPLETED' || exam.status === 'EXPIRED') {
          navigate(`/results/${userID}`);
          return;
        }
      } else {
        setError(response.data.error || 'Failed to load exam');
      }
    } catch (error) {
      setError('Connection failed. Please try again.');
      console.error('Load exam error:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadExistingAnswers = async () => {
    try {
      console.log('Loading existing answers for user:', userID);
      const response = await examAPI.getUserAnswers(userID);
      console.log('getUserAnswers response:', response);
      
      if (response.data.success) {
        const userAnswers = response.data.data || {};
        console.log('userAnswers from backend:', userAnswers);
        
        // userAnswers format from backend: { "examQuestionId": optionId }
        // Convert string keys to numbers and set directly
        const newAnswers = {};
        Object.entries(userAnswers).forEach(([examQuestionIdStr, optionId]) => {
          const examQuestionIdInt = parseInt(examQuestionIdStr);
          const optionIdInt = parseInt(optionId);
          newAnswers[examQuestionIdInt] = optionIdInt;
        });
        
        console.log('Converted answers:', newAnswers);
        setAnswers(newAnswers);
      } else {
        console.log('getUserAnswers failed:', response.data.error);
      }
    } catch (error) {
      console.error('Failed to load existing answers:', error);
      // Don't show error to user, just log it
    }
  };

  const startExam = async () => {
    try {
      const response = await examAPI.startExam(userID);
      if (response.data.success) {
        setExamStarted(true);
        // Reload exam data to get updated status
        await loadExamSession();
      } else {
        setError(response.data.error || 'Failed to start exam');
      }
    } catch (error) {
      setError('Failed to start exam');
      console.error('Start exam error:', error);
    }
  };

  const submitAnswer = async (examQuestionId, optionId) => {
    try {
      console.log('Submitting answer:', { examQuestionId, optionId });
      const response = await examAPI.submitAnswer(userID, examQuestionId, optionId);
      console.log('submitAnswer response:', response);
      
      if (response.data.success) {
        console.log('Answer submitted successfully, updating state');
        setAnswers(prev => ({ ...prev, [examQuestionId]: optionId }));
      } else {
        console.error('Submit answer failed:', response.data.error);
      }
    } catch (error) {
      console.error('Submit answer error:', error);
    }
  };

  const completeExam = async () => {
    try {
      const response = await examAPI.completeExam(userID);
      if (response.data.success) {
        navigate(`/results/${userID}`);
      } else {
        setError(response.data.error || 'Failed to complete exam');
      }
    } catch (error) {
      setError('Failed to complete exam');
    }
  };

  const handleOptionSelect = (examQuestionId, optionId) => {
    submitAnswer(examQuestionId, optionId);
  };

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

  if (!examData) return null;

  const currentQuestion = examData.questions[currentIndex];

  return (
    <div className="container-fluid">
      <div className="row">
        {/* Header */}
        <div className="col-12 bg-light py-3 mb-3">
          <div className="d-flex justify-content-between align-items-center">
            <div>
              <span className="h5">User: {userID}</span>
              <small className="text-muted ms-3">Session: {examData.session_code}</small>
            </div>
            <div>
              {examStarted && <Timer expiresAt={examData.expires_at} onExpire={completeExam} />}
            </div>
          </div>
        </div>
        {/* Sidebar */}
        <div className="col-md-3">
          <div className="card">
            <div className="card-header">
              <h6>Questions ({Object.keys(answers).length} / {examData.questions.length})</h6>
            </div>
            <div className="card-body">
              {examData.category_stats.map(stat => (
                <div key={stat.category} className="mb-3">
                  <h6 className="text-muted">{stat.category}</h6>
                  <div className="d-flex flex-wrap gap-1">
                    {examData.questions
                      .filter(q => q.category === stat.category)
                      .map((q, idx) => {
                        const globalIdx = examData.questions.findIndex(quest => quest.exam_question_id === q.exam_question_id);
                        const isAnswered = answers[q.exam_question_id] !== undefined;
                        const isCurrent = globalIdx === currentIndex;
                        return (
                          <button
                            key={q.exam_question_id}
                            className={`btn btn-sm ${isCurrent ? 'btn-primary' : isAnswered ? 'btn-success' : 'btn-outline-secondary'}`}
                            onClick={() => setCurrentIndex(globalIdx)}
                            disabled={!examStarted}
                          >
                            {globalIdx + 1}
                          </button>
                        );
                      })}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
        {/* Main Content */}
        <div className="col-md-9">
          {!examStarted && examData.status === 'NOT_STARTED' ? (
            <div className="text-center">
              <div className="card">
                <div className="card-body">
                  <h3>Ready to Start Your PPPK Exam</h3>
                  <p>You have {examData.duration} minutes to complete {examData.questions.length} questions.</p>
                  <button className="btn btn-primary btn-lg" onClick={startExam}>
                    Start Exam
                  </button>
                </div>
              </div>
            </div>
          ) : (
            <>
              <QuestionCard
                question={currentQuestion}
                selectedOption={answers[currentQuestion.exam_question_id]}
                onSelect={handleOptionSelect}
                questionNumber={currentIndex + 1}
              />
              <div className="d-flex justify-content-between mt-4">
                <button
                  className="btn btn-outline-secondary"
                  onClick={() => setCurrentIndex(Math.max(0, currentIndex - 1))}
                  disabled={currentIndex === 0}
                >
                  Previous
                </button>
                <div className="d-flex gap-2">
                  {currentIndex < examData.questions.length - 1 ? (
                    <button
                      className="btn btn-primary"
                      onClick={() => setCurrentIndex(currentIndex + 1)}
                    >
                      Next
                    </button>
                  ) : (
                    <button
                      className="btn btn-success"
                      onClick={completeExam}
                      disabled={Object.keys(answers).length === 0}
                    >
                      Finish Exam
                    </button>
                  )}
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ExamBoard;
