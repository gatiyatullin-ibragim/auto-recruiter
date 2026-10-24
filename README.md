# Auto Recruiter

AI-native сервис автоматизации поиска вакансий на Go.

## Запуск

```bash
cp .env.example .env
go test ./...
go run ./cmd/app
```

Frontend:

```bash
cd web
npm install
npm run dev
```

Откройте `http://localhost:5173`. Vite проксирует `/api` на backend.

Проверка:

```bash
curl http://localhost:8080/health
```

## Архитектура

- `internal/domain` — доменные сущности и инварианты.
- `internal/service` — use cases и оркестрация.
- `internal/ai` — контракты LLM, embeddings и structured prompts.
- `internal/repository` — контракты и in-memory адаптеры для разработки.
- `internal/scraper` — клиенты источников вакансий.
- `internal/transport` — HTTP и фоновые worker-ы.
- `pkg` — переиспользуемые технические компоненты.
- `migrations` — начальная схема PostgreSQL и pgvector.
- `web` — место для React/Vite интерфейса на этапе 2.

Интеграции с внешними API, PostgreSQL, Redis и LLM намеренно подключаются через
интерфейсы, чтобы бизнес-логика не зависела от конкретного провайдера.

## Проверка API

Создание вакансии:

```bash
curl -X POST http://localhost:8080/vacancies \
  -H 'Content-Type: application/json' \
  -d '{"id":"v1","source":"demo","title":"Go Developer","company":"Example","url":"https://example.com/jobs/v1","remote":true,"skills":["Go"]}'
```

Получение вакансии:

```bash
curl http://localhost:8080/vacancies/v1
```

По умолчанию приложение использует in-memory хранилище, поэтому оно запускается
без PostgreSQL и Redis. Docker Compose поднимает инфраструктуру для дальнейшего
подключения постоянных адаптеров.