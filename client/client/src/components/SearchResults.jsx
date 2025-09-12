import React from "react";
import { ExternalLink, Calendar, Tag, FileText } from "lucide-react";
import "./SearchResults.css";

const SearchResults = ({ results, isLoading, hasMore, onLoadMore }) => {
  if (isLoading) {
    return (
      <div className="search-results loading">
        <div className="loading-spinner">🔍</div>
        <p>Searching through analyzed content...</p>
      </div>
    );
  }

  if (!results || !results.results || results.results.length === 0) {
    return (
      <div className="search-results empty">
        <div className="empty-icon">📝</div>
        <h3>No results found</h3>
        <p>Try a different search term or category</p>
      </div>
    );
  }

  return (
    <div className="search-results">
      <div className="results-header">
        <h3>Found {results.total} results</h3>
      </div>

      <div className="results-grid">
        {results.results.map((result, index) => {
          const source = result._source; // Extract _source safely

          return (
            <div key={result._id || index} className="result-card">
              <div className="card-header">
                <span className="content-type">{source.content_type}</span>
                <span className="score">
                  Score: {result._score?.toFixed(2) || "N/A"}
                </span>
              </div>

              <h4 className="result-title">
                <a
                  href={source.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="result-link"
                >
                  {source.url}
                  <ExternalLink size={16} className="external-icon" />
                </a>
              </h4>

              <p className="summary">{source.summary}</p>

              <div className="card-footer">
                <div className="tags">
                  <span className="tag category-tag">
                    <Tag size={14} />
                    {source.domain_category}
                  </span>
                  <span className="tag sentiment-tag">{source.sentiment_tone}</span>
                </div>

                <div className="entities">
                  {source.key_entities?.slice(0, 3).map((entity, i) => (
                    <span key={i} className="entity-tag">
                      {entity}
                    </span>
                  ))}
                  {source.key_entities?.length > 3 && (
                    <span className="more-tags">
                      +{source.key_entities.length - 3} more
                    </span>
                  )}
                </div>

                <div className="meta-info">
                  <span className="meta-item">
                    <Calendar size={14} />
                    {new Date(source.processed_at).toLocaleDateString()}
                  </span>
                  <span className="meta-item">
                    <FileText size={14} />
                    Job #{source.job_id}
                  </span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {hasMore && (
        <div className="load-more">
          <button onClick={onLoadMore} className="load-more-btn">
            Load More Results
          </button>
        </div>
      )}
    </div>
  );
};

export default SearchResults;
