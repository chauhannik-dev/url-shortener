import { useState } from 'react'
import './App.css'

const App: React.FC = () => {
  const [url, setUrl] = useState<string>("")

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    console.log('Submitted value:', url);
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
          <button type="submit" className='submit-button'>
              Submit
          </button>
      </form>
    </div>
  )
}

export default App
