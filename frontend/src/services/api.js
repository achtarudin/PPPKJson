import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

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
  
  // Dashboard endpoints
  getUserDashboard: (userID) => api.get(`/exam/${userID}/dashboard`),
  getAllUsersDashboard: () => api.get(`/dashboard/users`)
};

export default api;
