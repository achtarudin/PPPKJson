import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import ExamBoard from './pages/ExamBoard';
import Result from './pages/Result';
import AdminDashboard from './pages/AdminDashboard';
import DebugExam from './pages/DebugExam';
import Layout from './components/Layout';

const App = () => (
  <Router>
    <Layout>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
        <Route path="/exam/:userID" element={<ExamBoard />} />
        <Route path="/results/:userID" element={<Result />} />
        <Route path="/admin" element={<AdminDashboard />} />
        <Route path="/dashboard" element={<AdminDashboard />} />
        <Route path="/debug" element={<DebugExam />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  </Router>
);

export default App;
