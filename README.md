# URL Shortener

A production-ready URL shortener built with Go, featuring Redis caching, PostgreSQL persistence, rate limiting, graceful shutdown, and URL expiry.

## Features

- **URL Shortening** — generate a random 6-character short code for any URL
- **Custom Alias** — optionally define your own short code (e.g. `/github`)
- **Redis Caching** — cache-aside pattern with 24-hour TTL for fast redirects
- **PostgreSQL Storage** — persistent URL storage with expiry timestamps
- **URL Expiry** — all short URLs automatically expire after 24 hours
- **Rate Limiting** — per-IP rate limiting (10 req/s, burst of 20)
- **Graceful Shutdown** — handles `Ctrl+C` without dropping in-flight requests

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| Database | PostgreSQL (via `pgx`) |
| Cache | Redis 7 |
| HTTP | `net/http` (stdlib) |
| Rate Limiter | `golang.org/x/time/rate` |
| Env Config | `godotenv` |

## System Architecture
<p align="center">
  <img src="https://github.com/user-attachments/assets/3f70c3e5-1f13-4ae2-bd32-a2ae6d33adc8"
       alt="URL Shortener Architecture"
       width="900">
</p>

## Project Structure

```
url-shortener/
├── main.go              # Server setup, routing, graceful shutdown
├── handlers/
│   └── handler.go       # HTTP handlers (ShortenURL, RedirectURL)
├── middleware/
│   └── rate-limiter.go  # Per-IP rate limiting middleware
├── models/
│   └── url.go           # Short URL generation (crypto/rand)
├── storage/
│   ├── store.go         # Redis + PostgreSQL storage layer
│   └── store_test.go    # Unit tests for Save and Get
├── .env                 # Environment variables (not committed)
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.22+
- PostgreSQL 14+
- Redis 7+ (via Docker recommended)

## Setup

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
```

### 2. Start Redis (Docker)

```bash
docker run -d -p 6379:6379 --name redis-shortener redis
```

### 3. Create PostgreSQL database and table

```sql
CREATE DATABASE url_shortener;

\c url_shortener

CREATE TABLE urls (
    short_url  TEXT PRIMARY KEY,
    long_url   TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '24 hours'
);
```

### 4. Configure environment variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/url_shortener?sslmode=disable
REDIS_ADDR=localhost:6379
```

### 5. Install dependencies

```bash
go mod tidy
```

### 6. Run the server

```bash
go run main.go
```

Server starts on `http://localhost:8080`.

## API

### Shorten a URL

```bash
POST /shorten
Content-Type: application/x-www-form-urlencoded

url=https://example.com
```

**Response:**
```
http://shorty.url/xK92mP
```

### Shorten with custom alias

```bash
curl -X POST -d "url=https://github.com&alias=github" http://localhost:8080/shorten
```

**Response:**
```
http://shorty.url/github
```

### Redirect

```bash
GET /:shortCode

curl http://localhost:8080/xK92mP
# → 302 redirect to https://example.com
```

## Running Tests

```bash
go test ./storage/...
```

> Make sure PostgreSQL and Redis are running before running tests. The test suite uses a real database connection and cleans up after itself via `t.Cleanup`.

## How It Works

**Shorten flow:**

1. Request hits the rate limiter middleware
2. Handler extracts `url` (and optional `alias`) from form body
3. Short code is generated randomly or taken from the alias
4. Short URL is saved to Redis (24h TTL) and PostgreSQL (`expires_at = NOW() + 24h`)

**Redirect flow:**

1. Handler extracts the short code from the URL path
2. Redis is checked first — if found and TTL > 0, redirect immediately
3. If cache miss, PostgreSQL is queried with `expires_at > NOW()` check
4. If found in DB, redirect; if not, return 404

## License

MIT
