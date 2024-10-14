import { useState } from 'react'
import axios from 'axios'
import './App.css'

const App: React.FC = () => {
  const [url, setUrl] = useState<string>("")
  const [encodedURL, setEncodedURL] = useState<string>("")
  const [error, setError] = useState<string>("")

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    try {
      const response = await axios.post('/api', {
        url
      })

      setEncodedURL(response.data.short_url);
    } catch (error) {
      setError("Failed to fetch data: " + (error as Error).message);
    }
    
    setUrl('');
  }

  return (
    <div className='card'>
      <h2 className="card-title">Enter your URL</h2>
      <form onSubmit={handleSubmit}>
          <input
              type="text"
              value={url}
              onChange={(e) => {
                setUrl(e.target.value); setError('')
              }}
              placeholder="https://example.com"
              className='input-field'
              required
          />
          <br />
          { encodedURL.length > 0 &&
            <div className='paragraph'>
              <p>Encoded URL: {encodedURL}</p>
            </div>
          }
          
          {
            encodedURL.length && <button className='button' onClick={() => window.location.href = encodedURL}>Redirect to URL</button>
          }
          
          {error && <div><p className='error'>{error}</p></div>}
          
          <button type="submit" className='submit-button'>
              Submit
          </button>
      </form>
    </div>
  )
}

export default App
