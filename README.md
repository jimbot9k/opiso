# Opiso

Named for the transliteration of a Greek word for 'backwards'. [Source](https://glosbe.com/en/grc/backwards).

This project consists of two applications. A RESTful Golang Backend, and a Vue.js client application to interact with it.

# backend
A minimal Golang 1.23 application. It serves POST requests on /reverse, which takes an array of words and reverses each word. It has basic configuration options via environment variables.

PORT (default:8080) - what port does the server run on.

ROUTINE_LIMIT (default:50000) - how many go routines can run concurrently to reverse words. 

CORS_ORIGIN (default:http://localhost:${port}) - origins of sites accessing this backend allowed. 


# frontend
A basic Vue.js application for sending requests to the backend.

environment variables:

VITE_API_URL(default:http://localhost:8080) - the backend API Base URL

# docker-compose.yaml
```bash
docker-compose up --build
```

Assuming the default docker-compose has been used:

The UI should now be accessible at http://localhost.

The backend is available at http://localhost:4040

# future work
- random messages generator for frontend
- unit tests on backend
- openapi specification for backend
- health endpoint for backend