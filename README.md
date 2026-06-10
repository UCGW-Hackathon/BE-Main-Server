---
title: SiTukang API
emoji: 🛠️
colorFrom: blue
colorTo: indigo
sdk: docker
app_port: 7860
---

# SiTukang Backend

Production-oriented Gin/GORM backend for the SiTukang user and worker API.

## Run Locally

1. Create PostgreSQL database:

```sql
CREATE DATABASE situkang;
```

2. Copy environment values:

```bash
cp .env.example .env
```

3. Update `DB_PASSWORD` and `JWT_SECRET`, then start:

```bash
go run .
```

The API listens on `http://localhost:8080/v1` by default. On startup, `AUTO_MIGRATE=true` creates the SiTukang tables/enums and `SEED_DATA=true` inserts the default service catalog, FAQ, article, and promo rows.

## Key Environment

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `JWT_SECRET`, `JWT_ACCESS_TOKEN_TTL_SECONDS`, `JWT_REFRESH_TOKEN_TTL_HOURS`
- `AUTO_MIGRATE`, `SEED_DATA`
- `HOST_ADDRESS`, `HOST_PORT`

## Verification

```bash
go test ./...
go vet ./...
```
