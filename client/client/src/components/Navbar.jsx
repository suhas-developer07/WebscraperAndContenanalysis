import React from 'react';
import { Database, Search, Upload } from 'lucide-react';
import './navbar.css';

const Navbar = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'search', label: 'Search Content', icon: Search },
    { id: 'submit', label: 'Submit URLs', icon: Upload }
  ];

  return (
    <nav className="navbar">
      <div className="nav-brand">
        <Database size={28} className="nav-icon" />
        <h1>Content Analyzer</h1>
      </div>

      <div className="nav-tabs">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              className={`nav-tab ${activeTab === tab.id ? 'nav-tab-active' : ''}`}
              onClick={() => setActiveTab(tab.id)}
            >
              <Icon size={20} className="tab-icon" />
              {tab.label}
            </button>
          );
        })}
      </div>
    </nav>
  );
};

export default Navbar;