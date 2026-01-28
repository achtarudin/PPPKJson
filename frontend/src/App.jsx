import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import ExamBoard from './pages/ExamBoard';
import Result from './pages/Result';
import Layout from './components/Layout';

const App = () => (
  <Router>
    <Layout>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/exam/:userID" element={<ExamBoard />} />
        <Route path="/results/:userID" element={<Result />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  </Router>
);

export default App;
