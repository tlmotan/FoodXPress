# grabfood-clone — scaffolding

## Quick start
```
docker compose up -d          # start Postgres+PostGIS and Redis
psql -h localhost -U grabfood -d grabfood -f migrations/0001_init.sql
go mod tidy
go run cmd/api/main.go
```

## What's real vs. placeholder
- `docker-compose.yml`, `migrations/0001_init.sql` — usable as-is, adjust as your design evolves.
- `internal/domain/*.go` — basic structs, expand as needed.
- `internal/repository/*.go` — interface signatures are a starting draft, not final — revisit once you know your actual query patterns.
- `internal/service/matching_service.go` — intentionally empty. This is the centerpiece; design and write it yourself.
- `cmd/api/main.go` — bare entrypoint, wire it up as you build.
