# FoodXPress — GrabFood clone

## Quick start
```
docker compose up -d          # start Postgres+PostGIS and Redis
psql -h localhost -U grabfood -d grabfood -f migrations/0001_init.sql
go mod tidy
go run cmd/api/main.go
```

Note: `docker compose up` does not run migrations — apply `migrations/0001_init.sql` manually on a fresh database.

### Running tests
Repository tests are integration tests against the Docker Postgres and are skipped if `DATABASE_URL` is unset:
```
DATABASE_URL=postgres://grabfood:grabfood@localhost:5432/grabfood go test ./...
```

# 1. Data/Request Flow
The sequence diagram illustrates these three core features:
1. Order placement
2. Driver matching
3. Live tracking

![alt text](data-request-sequence-diagram.png)

# 2. Core DB and Backend Groundwork
1. Docker Compose as the setup. Then, Postgres + PostGIS and Redis running in containers in a `docker-compose.yml` file.
- so that we don't have to run PostGIS and Redis manually.

## What's real vs. placeholder
- `docker-compose.yml`, `migrations/0001_init.sql` — usable as-is, adjust as your design evolves.
- `internal/domain/*.go` — basic structs, expand as needed.
- `internal/repository/*.go` — interface signatures are a starting draft, not final — revisit once you know your actual query patterns.
- `internal/service/matching_service.go` — intentionally empty. This is the centerpiece; design and write it yourself.
- `cmd/api/main.go` — bare entrypoint, wire it up as you build.
