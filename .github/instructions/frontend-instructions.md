PPPKFrontend AI Development Guide

Project Overview

Project inside at folder frontend

This is a React + Bootstrap 5 frontend application designed to interact with the PPPK Exam Hexagonal Go Backend. The goal is to provide a minimalist, responsive interface for users to take exams across 4 categories (MANAJERIAL, SOSIAL_KULTURAL, TEKNIS, WAWANCARA). Each exam session contains 4 questions (1 per category) with a 120-minute time limit.

Tech Stack & Setup

Framework: React 18+ (Vite)

Styling: Bootstrap 5 (via CDN or bootstrap npm) + Standard CSS for minor tweaks.

HTTP Client: Axios

Routing: React Router DOM (v6)

State Management: React useState / useContext (No Redux/Zustand needed for this scope).

Directory Structure (Minimalist)

src/
├── components/       # Reusable UI components
│   ├── Layout.jsx    # Main wrapper with Navbar
│   ├── QuestionCard.jsx # Renders question text & radio options
│   └── Timer.jsx     # Countdown timer based on session
├── pages/
│   ├── Login.jsx     # Simple input for User ID
│   ├── ExamBoard.jsx # Main exam interface (Question navigation)
│   ├── Result.jsx    # Score summary
│   └── AdminDashboard.jsx # Admin dashboard to view all users
├── services/
│   └── api.js        # Axios instance & endpoints
├── App.jsx           # Routes definition
└── main.jsx          # Entry point + Bootstrap import


API Integration Strategy

Base Configuration

Base URL: http://localhost:8080/api/v1

Response Wrapper: Backend returns { success: bool, message: string, data: any, error?: string }.

Rule: Always unwrap response.data.data to get the actual payload and check response.data.success before proceeding.

Error Handling: If response.data.success is false, show response.data.error message to user.

Key Endpoints Mapping

1. Create/Get Exam Session:

GET /exam/{userID}

Description: Creates a new exam session with 4 random questions (1 per category) or returns existing active session.

Response Structure:
```json
{
  "success": true,
  "message": "Exam session ready",
  "data": {
    "session_id": 1,
    "user_id": "1234",
    "session_code": "EXAM_1234_1643356800",
    "status": "NOT_STARTED",
    "expires_at": "2026-01-28T12:00:00Z",
    "duration": 120,
    "questions": [
      {
        "exam_question_id": 1,
        "question_id": 15,
        "category": "MANAJERIAL",
        "order_number": 1,
        "question_text": "Question text...",
        "options": [
          {
            "id": 59,
            "option_text": "Option text..."
          }
        ]
      }
    ],
    "category_stats": [
      {
        "category": "MANAJERIAL",
        "total_questions": 1,
        "answered_count": 0
      }
    ]
  }
}
```

2. Start Exam Session:

POST /exam/{userID}/start

Description: Starts the exam timer and changes status to IN_PROGRESS.

Action: Call this when user clicks "Start Exam" to begin the timer.

3. Submit Answer:

POST /exam/{userID}/answer

Payload: 
```json
{
  "exam_question_id": 1,
  "question_option_id": 59
}
```

Description: Submits user's answer for a specific question. Use `exam_question_id` from the question data and `question_option_id` from the selected option.

Important: Use `exam_question_id` (NOT `question_id`) and `question_option_id` (the option's `id` field).

4. Complete Exam:

POST /exam/{userID}/complete

Description: Completes the exam session and calculates final results. Call this when user finishes or time expires.

5. Get Results:

GET /exam/{userID}/results

Description: Retrieves detailed exam results including summary and category breakdown.

Response Structure:
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_questions": 4,
      "total_answered": 4,
      "total_score": 12,
      "max_score": 16,
      "overall_percentage": 75.0,
      "overall_grade": "B",
      "is_passed": true,
      "completed_at": "2026-01-28T11:30:00Z"
    },
    "results_by_category": [
      {
        "category": "MANAJERIAL",
        "total_questions": 1,
        "total_answered": 1,
        "total_score": 3,
        "max_score": 4,
        "percentage": 75.0,
        "grade": "B",
        "is_passed": true
      }
    ]
  }
}
```

6. Get Individual User Dashboard:

GET /exam/{userID}/dashboard

Description: Gets comprehensive dashboard data for a specific user including exam status, progress, and results.

Response Structure:
```json
{
  "success": true,
  "data": {
    "user_id": "1234",
    "has_exam": true,
    "exam_status": "COMPLETED",
    "exam_session": {
      "session_id": 1,
      "user_id": "1234",
      "session_code": "EXAM_1234_1643356800",
      "status": "COMPLETED",
      "expires_at": "2026-01-28T12:00:00Z",
      "duration": 120
    },
    "exam_results": {
      "summary": { ... },
      "results_by_category": [ ... ]
    },
    "progress_info": {
      "total_questions": 4,
      "answered_questions": 2,
      "remaining_time_minutes": 95
    }
  }
}
```

7. Get All Users Dashboard (Admin):

GET /dashboard/users

Description: Gets dashboard data for all users who have taken exams. Useful for admin interfaces.

Response Structure:
```json
{
  "success": true,
  "data": {
    "total_users": 25,
    "users": [
      {
        "user_id": "1234",
        "exam_status": "COMPLETED",
        "session_code": "EXAM_1234_1643356800",
        "started_at": "2026-01-28T10:00:00Z",
        "completed_at": "2026-01-28T11:30:00Z",
        "total_score": 12,
        "max_score": 16,
        "percentage": 75.0,
        "grade": "B",
        "is_passed": true
      },
      {
        "user_id": "5678",
        "exam_status": "IN_PROGRESS",
        "session_code": "EXAM_5678_1643357200",
        "started_at": "2026-01-28T10:15:00Z"
      },
      {
        "user_id": "9999",
        "exam_status": "EXPIRED",
        "session_code": "EXAM_9999_1643357800"
      }
    ]
  }
}
```

Status Values:
- `NO_EXAM`: User hasn't created any exam session
- `NOT_STARTED`: Exam session created but not started  
- `IN_PROGRESS`: User is currently taking the exam
- `COMPLETED`: Exam finished with results
- `EXPIRED`: Exam session has expired

Component Implementation Rules

1. Login Page (Login.jsx)

Simple centered card.

Input field for UserID (numeric string or text).

"Create/Start Exam" button.

Logic: On submit, call GET /exam/{userID} to create/get exam session, then redirect to /exam/{userID}.

2. Exam Board (ExamBoard.jsx)

State Needed:

examData: Complete exam session data from GET /exam/{userID}.

currentIndex: Integer (0 to 3 for 4 questions).

answers: Object/Map { examQuestionId: optionId }.

loading: Boolean.

examStarted: Boolean (tracks if POST /start has been called).

timeRemaining: Integer in seconds (from duration * 60).

Layout:

Header: Timer (right), User ID and Session Code (left).

Sidebar/Navigation: Question pills numbered 1-4, colored by status:
- Answered: btn-success
- Current: btn-primary  
- Unanswered: btn-outline-secondary
- Group by category with labels

Main Area: QuestionCard component.

Footer: "Previous", "Next", "Submit Answer", "Finish Exam" buttons.

Important Logic:
- Call POST /exam/{userID}/start when user clicks "Start Exam" (only if status is NOT_STARTED).
- Call POST /exam/{userID}/answer immediately when user selects an option.
- Call POST /exam/{userID}/complete when user clicks "Finish" or time expires.
- Timer should count down from duration (120 minutes = 7200 seconds).

3. Question Card (QuestionCard.jsx)

Props: 
- question: Question object with exam_question_id, question_text, options
- selectedOption: Currently selected option ID
- onSelect: Callback function when option is selected

Rendering:

Show question number and category badge.

Render question_text.

Loop through options array.

Use Bootstrap form-check for radio buttons.

Critical: Use unique name attribute per question: `name={question-${exam_question_id}}`.

Call onSelect with the option's id when radio button changes.

4. Result Page (Result.jsx)

Fetch results from GET /exam/{userID}/results on component mount.

Display "Exam Completed" with overall summary.

Show category breakdown in a table format.

Display overall grade, percentage, and pass/fail status.

"Back to Home" button to return to login.

5. Admin Dashboard Page (AdminDashboard.jsx)

Fetch all users data from GET /dashboard/users on component mount.

Display table with columns: User ID, Status, Start Time, End Time, Score, Grade.

Status badges with different colors for each status (COMPLETED, IN_PROGRESS, EXPIRED, NOT_STARTED).

Filter/search functionality by user ID or status.

Click on user row to navigate to individual results page.

Pagination if there are many users.

Auto-refresh data every 30 seconds to show live updates.

Development Workflow Commands

1. Initialization

```bash
npm create vite@latest frontend -- --template react
cd frontend
npm install bootstrap axios react-router-dom
```

2. Main Entry (main.jsx)

Import Bootstrap CSS immediately:

```jsx
import 'bootstrap/dist/css/bootstrap.min.css';
import ReactDOM from 'react-dom/client';
import App from './App';
// ... render
```

3. API Service Example (services/api.js)

```javascript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor to handle API response format
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export const examAPI = {
  // Create or get existing exam session
  getOrCreateExam: (userID) =>
    api.get(`/exam/${userID}`),
    
  // Start exam timer
  startExam: (userID) =>
    api.post(`/exam/${userID}/start`),
    
  // Submit single answer
  submitAnswer: (userID, examQuestionId, optionId) =>
    api.post(`/exam/${userID}/answer`, {
      exam_question_id: examQuestionId,
      question_option_id: optionId
    }),
    
  // Complete exam
  completeExam: (userID) =>
    api.post(`/exam/${userID}/complete`),
    
  // Get results
  getResults: (userID) =>
    api.get(`/exam/${userID}/results`),
    
  // Get individual user dashboard
  getUserDashboard: (userID) =>
    api.get(`/exam/${userID}/dashboard`),
    
  // Get all users dashboard (admin)
  getAllUsersDashboard: () =>
    api.get(`/dashboard/users`)
};

export default api;
```

4. Running

```bash
npm run dev
```


UI/UX Best Practices (Bootstrap)

Use Containers (container, container-fluid) for alignment.

Use Cards (card, card-body) for questions to make them pop.

Use Buttons (btn-primary, btn-outline-secondary) for navigation.

Use Alerts (alert-info) for loading states or instructions.

Ensure Mobile Responsiveness: Stack the question navigation on small screens.

Timer Colors:
- Green (text-success): > 30 minutes remaining
- Yellow (text-warning): 10-30 minutes remaining  
- Red (text-danger): < 10 minutes remaining

Question Status Colors:
- btn-success: Answered questions
- btn-primary: Current question
- btn-outline-secondary: Unanswered questions

Category Badges:
- MANAJERIAL: badge-primary
- SOSIAL_KULTURAL: badge-success  
- TEKNIS: badge-warning
- WAWANCARA: badge-info

Complete Example Implementation

Here's a complete example of how to implement the ExamBoard component with correct API usage:

```jsx
// ExamBoard.jsx
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
  }, [userID]);

  const loadExamSession = async () => {
    try {
      setLoading(true);
      const response = await examAPI.getOrCreateExam(userID);
      
      if (response.data.success) {
        const exam = response.data.data;
        setExamData(exam);
        setExamStarted(exam.status === 'IN_PROGRESS');
        
        // If exam is already completed, redirect to results
        if (exam.status === 'COMPLETED' || exam.status === 'EXPIRED') {
          navigate(`/results/${userID}`);
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

  const startExam = async () => {
    try {
      const response = await examAPI.startExam(userID);
      if (response.data.success) {
        setExamStarted(true);
        // Reload exam data to get updated status
        loadExamSession();
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
      const response = await examAPI.submitAnswer(userID, examQuestionId, optionId);
      if (response.data.success) {
        setAnswers(prev => ({
          ...prev,
          [examQuestionId]: optionId
        }));
      } else {
        console.error('Submit answer error:', response.data.error);
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
      console.error('Complete exam error:', error);
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

        {/* Question Navigation Sidebar */}
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
                        const isAnswered = answers[q.exam_question_id];
                        const isCurrent = globalIdx === currentIndex;
                        
                        return (
                          <button
                            key={q.exam_question_id}
                            className={`btn btn-sm ${
                              isCurrent ? 'btn-primary' : 
                              isAnswered ? 'btn-success' : 'btn-outline-secondary'
                            }`}
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
              
              {/* Navigation Buttons */}
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
```

Error Handling

API Response Format: Backend always returns `{ success, message, data?, error? }`.

Success Check: Always check `response.data.success` before using data.

Error Messages: Show `response.data.error` or `response.data.message` for user feedback.

Network Errors: Handle axios network errors with try-catch and show generic "Connection failed" message.

Session Expiry: If any API returns 404 or indicates session expired, redirect to login page.

Example Error Handling:
```javascript
try {
  const response = await examAPI.getOrCreateExam(userID);
  if (response.data.success) {
    setExamData(response.data.data);
  } else {
    setError(response.data.error || 'Failed to load exam');
  }
} catch (error) {
  setError('Connection failed. Please try again.');
  console.error('API Error:', error);
}
```

UI States:

Loading: Show Bootstrap spinner during API calls.

Success: Use alert-success for successful operations.

Error: Use alert-danger for errors.

Timer Warning: Use alert-warning when < 10 minutes remaining.

Status Indicators:

NOT_STARTED: Show "Start Exam" button.

IN_PROGRESS: Show exam interface with timer.

COMPLETED/EXPIRED: Redirect to results page.