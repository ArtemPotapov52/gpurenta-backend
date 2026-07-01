# GPURenta — Backend Task List

> Статус: **Код написан, компилируется (`go build ./...` ✅)**

---

## Как читать

- `[ ]` — не начато
- `[/]` — в процессе
- `[x]` — завершено
- `[!]` — заблокировано

---

## ✅ 1.1 Структура проекта

- [x] `backend/` — корневая директория
- [x] `cmd/server/main.go` — точка входа (сборка роутов, graceful shutdown)
- [x] `internal/config/config.go` — конфигурация из env
- [x] `internal/types/types.go` — User, Agent, Rental, WorkloadImage
- [x] `go.mod` — модуль + зависимости

## ✅ 1.2 База данных (PostgreSQL)

- [x] `internal/db/db.go` — pgx pool (MaxConns=10, MinConns=2)
- [x] `internal/db/users.go` — FindOrCreateUser, FindUserByGoogleID, CreateUser
- [x] `internal/db/agents.go` — CreateAgent, GetAgentByID, GetAgentBySecret, Heartbeat, ListOnlineGPUs (с фильтрацией min_vram/image), MarkStaleAgentsOffline
- [x] `internal/db/rentals.go` — CreateRental, GetRentalByID, StopRental (с расчётом cost), GetActiveRentalByAgentID
- [x] `internal/db/migrations/001_init.sql` — таблицы users, agents, rentals + индексы

## ✅ 1.3 Auth API

- [x] `internal/auth/google.go` — верификация Google access_token через `oauth2/v3/userinfo`
- [x] `internal/auth/jwt.go` — GenerateToken / ValidateToken (HS256, 24h expiry)
- [x] `POST /v1/auth/google` — принимает `{access_token}`, возвращает `{token, user}`
- [x] `GET /v1/health` — `{"status":"ok"}`

## ✅ 1.4 Agent API

- [x] `POST /v1/agents/register` — создаёт агента, возвращает `{agent_id, secret}`
- [x] `POST /v1/agents/heartbeat` — обновляет last_heartbeat + frp_url
- [x] `GET /v1/images` — список поддерживаемых workload-образов
- [x] **Middleware agent_secret** — проверка X-Agent-ID + X-Agent-Secret

## ✅ 1.5 GPU Catalog API

- [x] `GET /v1/gpus` — список свободных онлайн-GPU с фильтрацией `?min_vram=&image=`

## ✅ 1.6 Rentals API

- [x] `POST /v1/rentals/start` — проверяет online + не занят → создаёт rental
- [x] `POST /v1/rentals/{id}/stop` — завершает, считает cost_cents
- [x] `GET /v1/rentals/{id}` — статус аренды

## 🟡 1.7 Stripe Connect (платежи)

> Отложено до появления реальных арендаторов. На MVP ручной режим.

- [ ] `internal/handler/payments.go` — checkout + webhook + payouts

## 🟡 1.8 FRP Tunnel

> Для MVP используем ngrok. FRP-сервер будет на VPS позже.

- [x] `internal/tunnel/frp.go` — заглушка

## ✅ 1.9 Middleware

- [x] CORS (проверка Origin, разрешён localhost:8080 + FRONTEND_URL)
- [x] Auth (JWT Bearer token, middleware для защищённых роутов)
- [x] Logging (request_id, method, path, status, duration)
- [x] Error handling (Recoverer + JSONError helper)
- [x] RequestID (X-Request-ID / UUID)

---

## Проверка

- [x] `go build ./...` — компилируется
- [x] `go vet ./...` — без ошибок
- [ ] `curl localhost:8080/v1/health` — `{"status":"ok"}`
