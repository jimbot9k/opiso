# Opiso

Named for the transliteration of a Greek word for 'backwards'. [Source](https://glosbe.com/en/grc/backwards).

This project consists of two applications. A RESTful Golang Backend, and a Vue.js client application to interact with it.

# backend
A minimal Golang 1.23 application. It serves POST requests on /reverse, which takes an array of words and reverses each word. It has basic configuration options via environment variables. A health endpoint is exposed on /health, and a Prometheus metric endpoint is served on /metrics. An OpenAPI endpoint is available on /openapi, this can be used via a Swagger UI deployment.

PORT (default:8080) - what port does the server run on.

ROUTINE_LIMIT (default:5000) - how many go routines can run concurrently to reverse words. 

CORS_ORIGIN (default:http://localhost:${port}) - origins of sites accessing this backend allowed. 

CACHE_SIZE (default:1000) - how many messsages can be cached.

CACHE_MINIMUM (default:10) - how long a message has to be to be cached.

## backend tests

Tests have been writted for the the "github.com/jimbot9k/opiso/internal/reverse" package. This package has 100% coverage. It is the only one tested right now due to it containing all the business logic of the application.

```bash
cd backend
go test -cover ./...
```

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