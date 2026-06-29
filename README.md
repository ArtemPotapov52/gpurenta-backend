# GPURenta Backend

Go API for GPU marketplace — register, discover, and rent GPU compute.

## Deploy to Railway

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/new?template=https://github.com/ArtemPotapov52/gpurenta-backend)

1. Click the button above or create a new Railway project
2. Connect your GitHub repo
3. Railway auto-detects Go, builds and starts the server
4. Add a PostgreSQL plugin — Railway sets `DATABASE_URL` automatically
5. Set environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | ✅ | — | Secret key for JWT tokens (change from default!) |
| `GOOGLE_CLIENT_ID` | ❌ | — | Google OAuth client ID |
| `FRONTEND_URL` | ❌ | `*` | Allowed CORS origin |

All other defaults work out of the box. Tables are created automatically on startup.

## Local dev

```bash
cp .env.example .env
go run ./cmd/server
```

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/v1/health` | — | Health check |
| GET | `/v1/images` | — | Supported workload images |
| POST | `/v1/auth/dev` | — | Dev login (no Google) |
| POST | `/v1/google` | — | Google OAuth login |
| POST | `/v1/agents/register` | JWT | Register GPU agent |
| POST | `/v1/agents/heartbeat` | Agent secret | Update agent status |
| GET | `/v1/gpus` | JWT | List free GPUs |
| POST | `/v1/rentals/start` | JWT | Rent a GPU |
| POST | `/v1/rentals/{id}/stop` | JWT | Stop rental |
| GET | `/v1/rentals/{id}` | JWT | Rental status |
