# Go Distributed Job Queue

A production-inspired distributed background job processing system built with **Go**, **Redis**, and **PostgreSQL**. The project demonstrates asynchronous job execution, concurrent worker pools, delayed retries, fault tolerance, graceful shutdown, and REST APIs documented with Swagger.

---

## Features

- REST API for creating and managing jobs
- PostgreSQL for persistent job storage
- Redis-backed job queue
- Concurrent worker pool using goroutines
- Delayed retry mechanism using Redis Sorted Sets
- Configurable retry limits
- Automatic scheduler for delayed jobs
- Job status tracking
- Queue metrics endpoint
- Graceful shutdown
- Swagger/OpenAPI documentation
- Environment-based configuration
- PostgreSQL connection pooling with `pgxpool`

---

## Tech Stack

- **Language:** Go
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **Database Driver:** pgxpool
- **Cache / Queue:** Redis
- **API Documentation:** Swagger (swaggo)
- **Configuration:** Environment Variables

---

## Architecture

```text
                    Client
                      │
                      ▼
                 REST API (Gin)
                      │
                      ▼
               PostgreSQL (Jobs)
                      │
                      ▼
                Redis Queue
          ┌───────────┴───────────┐
          ▼                       ▼
    Active Queue            Delayed Queue
          │                       ▲
          ▼                       │
     Job Fetcher            Scheduler Service
          │
          ▼
    Worker Pool (Goroutines)
          │
          ▼
     Job Processing
```

---

## Job Lifecycle

```text
Pending
   │
   ▼
Processing
   │
   ├──────────────► Completed
   │
   ▼
Failure
   │
   ▼
Retrying
   │
   ▼
Delayed Queue
   │
   ▼
Scheduler
   │
   ▼
Processing
```

If the maximum retry count is exceeded:

```text
Processing
      │
      ▼
Failed
```

---

## Project Structure

```text
.
├── cmd
│   ├── api
│   ├── scheduler
│   └── worker
│
├── internal
│   ├── config
│   ├── jobs
│   ├── scheduler
│   ├── worker
│   ├── database
│   └── redis
│
├── docs
├── go.mod
└── README.md
```

---

## Features Implemented

### Job Management

- Create Job
- Get Job by ID
- List Jobs

### Queue Processing

- Redis-backed queue
- Worker pool
- Concurrent processing
- Graceful shutdown

### Retry System

- Automatic retries
- Configurable retry count
- Delayed retries
- Error tracking

### Monitoring

- Queue metrics endpoint
- Job status tracking

### Documentation

- Swagger/OpenAPI documentation

---

## API Endpoints

### Jobs

| Method | Endpoint     | Description      |
| ------ | ------------ | ---------------- |
| POST   | `/jobs`      | Create a new job |
| GET    | `/jobs`      | List all jobs    |
| GET    | `/jobs/{id}` | Get job by ID    |

### Metrics

| Method | Endpoint   | Description      |
| ------ | ---------- | ---------------- |
| GET    | `/metrics` | Queue statistics |

### Documentation

| Method | Endpoint              |
| ------ | --------------------- |
| GET    | `/swagger/index.html` |

---

## Job Statuses

| Status       | Description               |
| ------------ | ------------------------- |
| `pending`    | Waiting to be processed   |
| `processing` | Currently being processed |
| `retrying`   | Waiting for retry         |
| `completed`  | Successfully completed    |
| `failed`     | Permanently failed        |

---

## Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/go-job-queue.git

cd go-job-queue
```

### 2. Configure environment variables

Create a `.env` file:

```env
PORT=8080

POSTGRES_URL=postgres://postgres:password@localhost:5432/jobqueue

REDIS_ADDR=localhost:6379

MAX_RETRIES=3
```

---

### 3. Start PostgreSQL

Ensure PostgreSQL is running and create the database:

```sql
CREATE DATABASE jobqueue;
```

---

### 4. Start Redis

```bash
redis-server
```

---

### 5. Run the API

```bash
go run ./cmd/api
```

---

### 6. Run the Worker

```bash
go run ./cmd/worker
```

---

### 7. Run the Scheduler

```bash
go run ./cmd/scheduler
```

---

## Swagger

Generate documentation:

```bash
swag init -g cmd/api/main.go
```

Open:

```text
http://localhost:8080/swagger/index.html
```

---

## Metrics Example

```json
{
  "pending": 2,
  "processing": 1,
  "retrying": 1,
  "completed": 45,
  "failed": 3,
  "total": 52
}
```

---

## Key Concepts Demonstrated

- Concurrent programming with goroutines
- Worker pool pattern
- Asynchronous job processing
- Redis queues
- Redis Sorted Sets
- Retry strategies
- Fault tolerance
- Graceful shutdown
- Connection pooling
- REST API design
- OpenAPI documentation

---

## Future Enhancements

- Dead Letter Queue (DLQ)
- Manual retry endpoint
- Job cancellation
- Scheduled jobs
- Job priorities
- Rate limiting
- Prometheus metrics
- Docker Compose support
- Kubernetes deployment
- Authentication & authorization
- Distributed worker scaling

---

## License

This project is licensed under the MIT License.
