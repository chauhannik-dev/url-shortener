# URL Shortener

## Tech Stack
- **Backend**: Go
- **Frontend**: React

## Overview
This project is a URL shortener service that allows users to shorten long URLs, delete shortened URLs, and redirect to the original long URLs based on the shortened version. The application uses hashing techniques to generate unique short URLs while handling potential collisions.

## Features
### 1. Shorten a URL
- **Endpoint**: `POST /`
- **Request Body**: 
  ```json
  {
    "url": "https://www.example.com"
  }
  ```
Description: This endpoint takes a long URL and generates a shortened version using SHA-256 for hashing and Base 62 encoding. It handles collisions by generating a new hash if a collision occurs.

- **Response**
```{
    "key": "unique_hash",
    "short_url": "http://short.ly/unique_hash",
    "long_url": "https://www.example.com"
}
```


### 2. Delete a URL
- **Endpoint**: `DELETE /{short_url}`

Description: Deletes the shortened URL entry from the database.

- **Response**:
```{
    "message": "URL deleted successfully."
}
```

### 3. Redirect to the Original URL
- **Endpoint**: `GET /{short_url}`

Description: When accessed, it looks up the shortened URL and redirects the user to the corresponding long URL.

- **Response**: Redirects to the original long URL associated with the provided shortened URL.
