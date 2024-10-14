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
        url: url
      })

      setEncodedURL(response.data.short_url);
    } catch (error) {
      console.log(error)
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
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://example.com"
              className='input-field'
              required
          />
          <br />
          {encodedURL.length && <p className='text'>Encoded URL: {encodedURL}</p>}
          {encodedURL.length && <button className='button' onClick={() => window.location.href = encodedURL}>Redirect to URL</button>}
          {error.length && <p>Error: {error}</p>}
          <button type="submit" className='submit-button'>
              Submit
          </button>
      </form>
    </div>
  )
}

export default App
