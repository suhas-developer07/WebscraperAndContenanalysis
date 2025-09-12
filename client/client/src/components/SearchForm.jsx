import React, { useState } from 'react';
import { Search, Filter } from 'lucide-react';
import { CATEGORIES } from '../services/constant';
import './SearchForm.css';

const SearchForm = ({ onSearch, isLoading }) => {
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSearch({ query: query.trim(), category });
  };

  const clearFilters = () => {
    setQuery('');
    setCategory('');
    onSearch({ query: '', category: '' });
  };

  const hasFilters = query.trim() || category;

  return (
    <div className="search-form">
      <h2 className="search-title">Search Analyzed Content</h2>
      
      <form onSubmit={handleSubmit} className="search-form-container">
        <div className="search-input-group">
          <div className="search-input-wrapper">
            <Search size={22} className="search-icon" />
            <input
              type="text"
              placeholder="Search for anything... (AI, sports, technology, etc.)"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="search-input"
            />
          </div>

          <div className="filter-group">
            <Filter size={18} className="filter-icon" />
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="category-select"
            >
              <option value="">All Categories</option>
              {CATEGORIES.map(cat => (
                <option key={cat} value={cat.toLowerCase()}>
                  {cat}
                </option>
              ))}
            </select>
          </div>

          <button 
            type="submit" 
            disabled={isLoading}
            className="search-btn"
          >
            {isLoading ? 'Searching...' : 'Search'}
          </button>

          {hasFilters && (
            <button 
              type="button" 
              onClick={clearFilters}
              className="clear-btn"
            >
              Clear
            </button>
          )}
        </div>
      </form>
    </div>
  );
};

export default SearchForm;