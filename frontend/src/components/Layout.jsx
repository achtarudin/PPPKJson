import React from 'react';
import { Link, useLocation } from 'react-router-dom';

const Layout = ({ children }) => {
  const location = useLocation();
  
  const isActive = (path) => {
    return location.pathname === path || location.pathname.startsWith(path);
  };

  return (
    <div className="min-vh-100 bg-light">
      {/* Navigation */}
      <nav className="navbar navbar-expand-lg navbar-dark bg-primary">
        <div className="container-fluid">
          <Link className="navbar-brand" to="/">
            <i className="bi bi-clipboard-check"></i> PPPK Exam System
          </Link>
          
          <button 
            className="navbar-toggler" 
            type="button" 
            data-bs-toggle="collapse" 
            data-bs-target="#navbarNav"
          >
            <span className="navbar-toggler-icon"></span>
          </button>
          
          <div className="collapse navbar-collapse" id="navbarNav">
            <ul className="navbar-nav me-auto">
              <li className="nav-item">
                <Link 
                  className={`nav-link ${isActive('/') && !isActive('/admin') && !isActive('/questions') ? 'active' : ''}`} 
                  to="/"
                >
                  <i className="bi bi-house"></i> Home
                </Link>
              </li>
              <li className="nav-item">
                <Link 
                  className={`nav-link ${isActive('/admin') ? 'active' : ''}`} 
                  to="/admin"
                >
                  <i className="bi bi-speedometer2"></i> Admin Dashboard
                </Link>
              </li>
              <li className="nav-item">
                <Link 
                  className={`nav-link ${isActive('/questions') ? 'active' : ''}`} 
                  to="/questions"
                >
                  <i className="bi bi-question-circle"></i> Question Manager
                </Link>
              </li>
            </ul>
            
            <ul className="navbar-nav">
              <li className="nav-item dropdown">
                <a 
                  className="nav-link dropdown-toggle" 
                  href="#" 
                  role="button" 
                  data-bs-toggle="dropdown"
                >
                  <i className="bi bi-gear"></i> Options
                </a>
                <ul className="dropdown-menu">
                  <li>
                    <Link className="dropdown-item" to="/login">
                      <i className="bi bi-box-arrow-in-right"></i> Start New Exam
                    </Link>
                  </li>
                  <li><hr className="dropdown-divider" /></li>
                  <li>
                    <Link className="dropdown-item" to="/admin">
                      <i className="bi bi-people"></i> View All Users
                    </Link>
                  </li>
                </ul>
              </li>
            </ul>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="py-3">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-dark text-light text-center py-3 mt-auto">
        <div className="container">
          <p className="mb-0">
            <small>
              © 2026 PPPK Exam System - Built with React & Bootstrap
            </small>
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
