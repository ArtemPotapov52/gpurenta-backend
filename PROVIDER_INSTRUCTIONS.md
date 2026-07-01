# GPURenta — Provider Setup (шаг за шагом)

Ты — **Provider**. Твой ноутбук будет **зарабатывать**, сдавая GPU в аренду.

---

## Что тебе понадобится

- **Терминал** (macOS/Linux) или **PowerShell** (Windows)
- Права на установку программ
- Доступ в интернет (ngrok выведет твой localhost наружу)

---

## Шаг 1. Установить Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh

# Windows — скачай https://ollama.com/download/windows
```

Скачать модель (лёгкая, 1.2B параметров):

```bash
ollama pull llama3.2:1b
```

Запустить в **фоне** (это окно больше не закрывай):

```bash
ollama serve
```

Проверить:

```bash
curl http://localhost:11434/api/tags
# → {"models":[{"name":"llama3.2:1b",...}]}
```

---

## Шаг 2. Установить ngrok

Как установить:

```bash
# macOS
brew install ngrok

# Windows — скачай https://ngrok.com/download
```

Получить токен (бесплатно, 1 минута):
1. Зайти на https://dashboard.ngrok.com
2. Войти (можно через GitHub)
3. Скопировать **Your Authtoken**
4. В терминале:

```bash
ngrok config add-authtoken СКОПИРОВАННЫЙ_ТОКЕН
```

---

## Шаг 3. Зарегистрировать GPU в бэкенде

Бэкенд доступен по публичному URL:
```
https://backroom-sliceable-yo-yo.ngrok-free.dev
```

Открой **новое окно терминала** и делай по порядку:

**3.1 — Получить токен:**
```bash
JWT=$(curl -s -X POST https://backroom-sliceable-yo-yo.ngrok-free.dev/v1/auth/dev | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
```

**3.2 — Зарегистрировать GPU:**
```bash
curl -s -X POST https://backroom-sliceable-yo-yo.ngrok-free.dev/v1/agents/register \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"gpu_model":"M3 Max","vram_gb":48,"os":"macOS","supported_images":["ollama"],"price_per_hour":25}'
```

> ⚠️ Замени `gpu_model`, `vram_gb`, `os` на свои характеристики (можно узнать через `nvidia-smi` на Windows/Linux или `About This Mac` на macOS).

**Ответ:**
```json
{"agent_id": "c26c5c87-...", "secret": "fd883e67..."}
```

**СОХРАНИ** `agent_id` и `secret` — они вставятся в следующие шаги.

---

## Шаг 4. Открыть туннель (ngrok)

**Новое окно терминала** (не закрывай):

```bash
ngrok http 11434
```

В терминале появится:
```
Forwarding  https://твой-уникальный-адрес.ngrok-free.app → http://localhost:11434
```

**Скопируй URL** — он понадобится на следующем шаге.

---

## Шаг 5. Отправить Heartbeat

Вернись в окно из Шага 3. Выполни (подставь свои значения):

```bash
curl -s -X POST https://backroom-sliceable-yo-yo.ngrok-free.dev/v1/agents/heartbeat \
  -H "X-Agent-ID: ТВОЙ_AGENT_ID" \
  -H "X-Agent-Secret: ТВОЙ_SECRET" \
  -H "Content-Type: application/json" \
  -d '{"frp_url":"https://твой-адрес.ngrok-free.dev"}'
```

Должен прийти ответ:
```json
{"status":"ok"}
```

---

## Шаг 6. Держать heartbeat (автоматически)

Heartbeat нужно отправлять **каждую минуту**, иначе GPU уйдёт в offline через 5 минут.

**macOS/Linux** — одна команда в фоне:
```bash
nohup bash -c 'while true; do curl -s -X POST https://backroom-sliceable-yo-yo.ngrok-free.dev/v1/agents/heartbeat \
  -H "X-Agent-ID: ТВОЙ_AGENT_ID" \
  -H "X-Agent-Secret: ТВОЙ_SECRET" \
  -H "Content-Type: application/json" \
  -d "{\"frp_url\":\"https://твой-адрес.ngrok-free.dev\"}" > /dev/null; sleep 60; done' &
```

**Windows (PowerShell):**
```powershell
while ($true) {
  curl.exe -s -X POST https://backroom-sliceable-yo-yo.ngrok-free.dev/v1/agents/heartbeat `
    -H "X-Agent-ID: ТВОЙ_AGENT_ID" `
    -H "X-Agent-Secret: ТВОЙ_SECRET" `
    -H "Content-Type: application/json" `
    -d '{"frp_url":"https://твой-адрес.ngrok-free.dev"}'
  Start-Sleep -Seconds 60
}
```

---

## Шаг 7. Отдать agent_id арендатору

Скажи ему **agent_id** (токен из Шага 3). Арендатор найдёт твою GPU в каталоге и сможет арендовать.

---

## Как это работает

```
Твой ноутбук                          Бэкенд                           Арендатор
┌──────────────────┐                 ┌────────────────────┐           ┌──────────────────┐
│  Ollama :11434    │                 │  Go API :8081       │           │  Дашборд :8081    │
│       ↕          │                 │  PostgreSQL         │           │                   │
│  ngrok → url     │──heartbeat─────→│  GPU → online       │──→каталог→│  Видит GPU        │
│                  │                 │                     │←──аренда──│  Жмёт "Rent"      │
│       ↕          │←──ngrok tunnel──│  (прямой канал)     │────запрос─│  curl к Ollama    │
│  Ollama отвечает │─────────────────┼─────────────────────│ ←──ответ──│  Получает ответ   │
└──────────────────┘                 └────────────────────┘           └──────────────────┘
```

Всё. Твой ноутбук зарабатывает 🎉
