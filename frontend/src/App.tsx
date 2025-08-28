import { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsLoading(true);
    setError('');
    setShortUrl('');

    if (!url) {
      setError('Please enter a URL to shorten.');
      setIsLoading(false);
      return;
    }

    try {
      const response = await axios.post('http://localhost:8080', { url });
      setShortUrl(response.data.short_url);
    } catch (err) {
      setError('Failed to shorten URL. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="app-container">
      <h1 className="title">Shortly</h1>
      <p className="subtitle">Your modern URL shortener.</p>
      
      <form className="form-container" onSubmit={handleSubmit}>
        <input
          type="url"
          className="url-input"
          placeholder="Enter a long URL here..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          disabled={isLoading}
        />
        <button type="submit" className="submit-button" disabled={isLoading}>
          {isLoading ? '...' : 'Shorten'}
        </button>
      </form>

      {error && <p className="error-message">{error}</p>}

      {shortUrl && (
        <div className="result-container">
          <p>Your short URL:</p>
          <a href={shortUrl} target="_blank" rel="noopener noreferrer">
            {shortUrl}
          </a>
        </div>
      )}
    </div>
  );
}

export default App;