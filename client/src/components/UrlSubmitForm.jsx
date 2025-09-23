import React, { useState } from 'react';
import { Plus, Trash2, Upload } from 'lucide-react';
import { api } from '../services/api';
import './UrlSubmitForm.css';

const UrlSubmitForm = () => {
  const [urls, setUrls] = useState(['']);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState('');

  const addUrlField = () => {
    setUrls([...urls, '']);
  };

  const removeUrlField = (index) => {
    if (urls.length > 1) {
      setUrls(urls.filter((_, i) => i !== index));
    }
  };

  const updateUrl = (index, value) => {
    const newUrls = [...urls];
    newUrls[index] = value;
    setUrls(newUrls);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setMessage('');

    try {
      const validUrls = urls.filter(url => url.trim() !== '');
      if (validUrls.length === 0) {
        throw new Error('Please enter at least one valid URL');
      }

      const result = await api.submitUrls(validUrls);
      setMessage(`Success! Job ID: ${result.job_id}. Processing ${validUrls.length} URLs...`);
      setUrls(['']);
    } catch (error) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="url-submit-form">
      <h2 className="form-title">Submit URLs for Analysis</h2>
      
      <form onSubmit={handleSubmit} className="url-form">
        <div className="url-fields">
          {urls.map((url, index) => (
            <div key={index} className="url-input-group">
              <input
                type="url"
                placeholder="https://example.com/article"
                value={url}
                onChange={(e) => updateUrl(index, e.target.value)}
                className="url-input"
                required
              />
              {urls.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeUrlField(index)}
                  className="remove-btn"
                  aria-label="Remove URL"
                >
                  <Trash2 size={18} />
                </button>
              )}
            </div>
          ))}
        </div>

        <div className="form-actions">
          <button
            type="button"
            onClick={addUrlField}
            className="secondary-btn"
          >
            <Plus size={18} />
            Add Another URL
          </button>

          <button
            type="submit"
            disabled={isSubmitting}
            className="primary-btn"
          >
            {isSubmitting ? (
              'Processing...'
            ) : (
              <>
                <Upload size={18} />
                Submit URLs
              </>
            )}
          </button>
        </div>

        {message && (
          <div className={`message ${message.includes('Success') ? 'message-success' : 'message-error'}`}>
            {message}
          </div>
        )}
      </form>
    </div>
  );
};

export default UrlSubmitForm;