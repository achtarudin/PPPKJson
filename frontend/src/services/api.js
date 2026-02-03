import axios from 'axios';

// Use API base URL from Vite config, with fallback to auto-detection
const getApiBaseUrl = () => {
  // Use global variable from vite.config.mjs
  if (typeof __API_BASE_URL__ !== 'undefined') {
    return __API_BASE_URL__;
  }
  
  // Fallback: Auto-detect based on current host
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  const port = window.location.port;
  
  // If running on localhost, assume development
  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return 'http://localhost:8080/api/v1';
  }
  
  // For production, use same host as frontend
  const baseUrl = port ? `${protocol}//${hostname}:${port}` : `${protocol}//${hostname}`;
  return `${baseUrl}/api/v1`;
};

const API_BASE_URL = getApiBaseUrl();

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export const examAPI = {
  getOrCreateExam: (userID) => api.get(`/exam/${userID}`),
  startExam: (userID) => api.post(`/exam/${userID}/start`),
  submitAnswer: (userID, examQuestionId, optionId) =>
    api.post(`/exam/${userID}/answer`, {
      exam_question_id: examQuestionId,
      question_option_id: optionId
    }),
  completeExam: (userID) => api.post(`/exam/${userID}/complete`),
  getResults: (userID) => api.get(`/exam/${userID}/results`),
  
  // Get existing user answers for repopulation
  getUserAnswers: (userID) => api.get(`/exam/${userID}/answers`),
  
  // Get detailed answers with questions and scores for completed exam
  getDetailedUserAnswers: (userID) => api.get(`/exam/${userID}/detailed-answers`),
  
  // Dashboard endpoints
  getUserDashboard: (userID) => api.get(`/exam/${userID}/dashboard`),
  getAllUsersDashboard: () => api.get(`/dashboard/users`)
};

export const questionAPI = {
  // Get all questions with optional category filter, search text, and pagination
  getQuestionsByCategory: ({ category = '', search = '', page = 1, limit = 10 } = {}) => {
    const params = {
      page: page,
      limit: limit
    };
    if (category) params.category = category;
    if (search) params.search = search;
    return api.get('/questions', { params });
  },
  
  // Update score for a specific question option
  updateOptionScore: (questionId, optionId, score) =>
    api.put(`/questions/${questionId}/option/${optionId}/score`, { score }),
  
  // Get all available categories
  getCategories: () => api.get('/questions/categories')
};

export default api;
