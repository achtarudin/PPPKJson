import React from 'react';

const Layout = ({ children }) => (
  <div>
    <nav className="navbar navbar-expand-lg navbar-dark bg-primary mb-4">
      <div className="container-fluid">
        <span className="navbar-brand">PPPK Exam System</span>
      </div>
    </nav>
    <main>{children}</main>
  </div>
);

export default Layout;
