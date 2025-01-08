# TinyWiny URL Shortener

TinyWiny is a lightweight URL shortener built with GoLang, Docker, and Redis. It allows you to create short URLs, redirect to original URLs, and persist data using Redis.

---

## Features

- Generate short URLs for given long URLs.
- Redirect users from short URLs to original URLs.
- Persist data using Redis with automatic loading of existing databases.
- Fully Dockerized for easy setup and deployment.

---

## Prerequisites

1. **Docker** and **Docker Compose** installed.
2. (Optional) **GoLang** for local development.

---

## Setup and Usage

### 1. Clone the Repository
```bash
git clone <repository-url>
cd TinyWiny
```

### 2. Build Docker Images
```bash
make build
```

### 3. Start the Application
```bash
make tinywiny
```

### 4. Test the Setup
```bash
curl -X POST -H "Content-Type: application/json" -d '{"long_url":"http://example.com"}' http://localhost:8888/shorten
```
---

## API Endpoints

### Shorten a URL
- **POST /shorten**
  ```json
  {
    "long_url": "http://example.com"
  }
  ```
  Response:
  ```json
  {
    "short_url": "http://localhost:8888/1"
  }
  ```

### Redirect to Original URL
- **GET /:short_url**
  Example:
  ```bash
  curl http://localhost:8888/1
  ```

---

## Directory Structure

```
TinyWiny/
├── cmd/                   # Entry point
├── internal/              # Core logic
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   ├── storage/           # Redis integration
│   └── models/            # Data models
├── Dockerfile             # Dockerfile
├── docker-compose.yml     # Docker Compose setup
├── Makefile               # Build and run commands
└── README.md              # Documentation
```

---

## Configuration

- **Redis Data File:** Stored in `data/url_database.rdb` and loaded automatically if it exists.
- **Environment Variables:**
  - `REDIS_HOST`: Redis hostname (default: `redis`).
  - `REDIS_PORT`: Redis port (default: `6379`).

---

## Commands

- **Build Images:** `make build`
- **Start Services:** `make tinywiny`
- **Stop Services:** `make down`
- **View Logs:** `make logs`
- **Clean Redis Data:** `make redis-clean`

---

## Testing

Run tests:
```bash
go test ./...
```

---

## Future Enhancements

- Add custom aliases for short URLs.
- Support URL expiration with TTL.
- Provide analytics for URL access.

---

## License

This project is licensed under the MIT License.

---

## Contributors

- **Amin Mousavi** (Developer)

