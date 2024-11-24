# Opiso

Named for the transliteration of a Greek word for 'backwards'. [Source](https://glosbe.com/en/grc/backwards).

This project consists of two applications. A RESTful Golang Backend, and a Vue.js client application to interact with it.

# backend
A minimal Golang 1.23 application. It serves POST requests on /reverse, which takes an array of words and reverses each word. It has basic configuration options via environment variables. A health endpoint is exposed on /health, and a Prometheus metric endpoint is served on /metrics. An OpenAPI endpoint is available on /openapi, this can be used via a Swagger UI deployment.

PORT (default:4040) - what port does the server run on.

ROUTINE_LIMIT (default:20000) - how many go routines can run concurrently to reverse words. 

CORS_ORIGIN (default:*) - origins of sites accessing this backend allowed. 

CACHE_SIZE (default:1000) - how many messsages can be cached.

CACHE_MINIMUM (default:10) - how long a message has to be to be cached.

## backend tests

Tests have been writted for the the "github.com/jimbot9k/opiso/internal/reverse" package. This package has 100% coverage. It is the only one tested right now due to it containing all the business logic of the application.

```bash
cd backend
go test -cover ./...
```

## backend design overview

The backend is a RESTful API service with prometheus metrics, health endpoint, open-api documentation, that provides it's primary function of reversing messages on POST /reverse. The backend was initially built to reverse a single messages, then concurrency was added to handle multiple messages in parallel. To prevent overusage of resources, a semaphore was used to limit the amount of routines that can be used at once. Each reversal had a space and time complexity of O(N). To improve this, a LRU (least-recently-used) cache was written to store a given limit of messages above a certain size threshold. This improves time complexity for cached words to O(1).

The project was built using the standard library as much as possible. One of the beautiful parts of Golang is how functional its standard lib is. The only external dependencies are the prometheus client, and its dependencies. 

There are multiple custom metrics the prometheus endpoint presents. These track the number of messages and requests received, the number of messages cached, and the number of routines used from the limit. These are all labled with "opiso" so they can be found easily.

# frontend
A basic Vue.js application for sending requests to the backend.

environment variables:

VITE_API_URL(default:http://localhost:4040) - the backend API Base URL

## frontend design overview

The frontend is a basic Vue.js SPA application. It has a form with a text area and submission button, result area for reversed messages, and a snackbar to tell users abour errors and such. The app uses minimal Vue.js dependencies alongside fetch for HTTP requests. 

Like the backend, the frontend tried to minimise dependencies. State libraries, HTTP request libraries and other tools could be used for more complicated applications. 

# docker-compose.yaml
```bash
docker-compose up --build
```

Assuming the default docker-compose has been used:

The UI should now be accessible at http://localhost.

The backend is available at http://localhost:4040

# future work
- more unit tests on backend
- snackbar on UI sometimes bugs out. need to investigate this.

# lessons learned
If I were to build this project again, I would use an OpenAPI code generator to create my HTTP boilerplate. I'm mostly used to tooling where OpenAPI specs get generated from code, but it appears the better tools in Golang do spec->code. Spec->code is probably better too as it puts more emphasis on upfront design. 
