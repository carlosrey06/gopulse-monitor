# GoPulse Monitor

GoPulse Monitor is a small website and service monitoring MVP. Its Go backend checks five configured services concurrently and exposes their current availability, HTTP status, response time, and check timestamp. The Astro dashboard presents that information in a dark, responsive monitoring interface with automatic refreshes and an ECharts response-time comparison.

## Development

Backend:

```bash
go run .
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

URLs:

```text
Backend:  http://localhost:8080
Frontend: http://localhost:4321
Health:   http://localhost:8080/health
Monitors: http://localhost:8080/api/monitors
```

## API

- `GET /health` returns the backend service status.
- `GET /check?url=https://example.com` checks a valid HTTP or HTTPS URL.
- `GET /api/monitors` concurrently checks the five monitors configured in `main.go`.

The development CORS policy permits the Astro dev server at `http://localhost:4321` (and its `127.0.0.1` equivalent). The monitor list is intentionally fixed in Go for this MVP; it does not include persistence, authentication, notifications, or monitor editing.

## Verification

```bash
gofmt -w .
go vet ./...
go test ./...

cd frontend
npm run check
npm run build
```
