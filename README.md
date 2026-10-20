# Leaderboard System

Real-time leaderboard backend for games, built with **Go**, **Redis** and **PostgreSQL**.

Players register, submit scores, and see rankings update immediately. Redis sorted sets serve the live leaderboards, while PostgreSQL keeps the full score history.

## Features

- JWT authentication with `user` and `admin` roles
- Global and per-game leaderboards backed by Redis sorted sets
- Score history and per-player stats stored in PostgreSQL
- Game management (create / delete) for admins
- One-command local setup with Docker Compose, schema and demo games included

## Tech Stack

| Component | Purpose |
|-----------|---------|
| Go | HTTP API |
| Redis 7 | Real-time rankings (sorted sets) |
| PostgreSQL 16 | Users, games, score history |
| JWT | Authentication |
| Docker Compose | Local environment |

## Quick Start

Requires Docker (with Compose v2).

```bash
git clone https://github.com/bbrighty/leaderboard-system.git
cd leaderboard-system
docker compose up --build
```

The API is available at `http://127.0.0.1:8080`. On first start the database is created from `migrations/` and seeded with three demo games (`id` 1-3).

> **macOS note:** use `127.0.0.1` instead of `localhost` in the examples below. `localhost` may resolve to IPv6 first, and Docker Desktop can then answer with `Empty reply from server`.

## Usage Example

Register a player. The response already contains a token:
a1h5a
```bash
curl -s -X POST "http://127.0.0.1:8080/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","email":"p1@test.com","password":"pass123"}'
```

Log in later to get a fresh token (login uses **email**) and keep it in a variable:

```bash
TOKEN=$(curl -s -X POST "http://127.0.0.1:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"p1@test.com","password":"pass123"}' | jq -r .token)
```

List games, then submit a score:

```bash
curl -s "http://127.0.0.1:8080/api/games" -H "Authorization: Bearer $TOKEN"

curl -i -X POST "http://127.0.0.1:8080/api/scores" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"game_id":1,"score":1500}'
```

Read the leaderboards and your stats:

```bash
curl -s "http://127.0.0.1:8080/api/leaderboard/global?limit=10"
curl -s "http://127.0.0.1:8080/api/leaderboard/game/1?limit=10"
curl -s "http://127.0.0.1:8080/api/stats/me" -H "Authorization: Bearer $TOKEN"
```

Tip: pipe any response through `jq` for readable JSON.

## API

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | - | Register, returns user and token |
| POST | `/api/auth/login` | - | Log in with email and password, returns token |
| GET | `/api/games` | - | List games |
| POST | `/api/games` | admin | Create a game |
| DELETE | `/api/games/{id}` | admin | Delete a game |
| POST | `/api/scores` | user | Submit a score |
| GET | `/api/leaderboard/global?limit=10` | - | Global leaderboard |
| GET | `/api/leaderboard/game/{id}?limit=10` | - | Leaderboard of one game |
| GET | `/api/stats/me` | user | Stats of the current player |

Protected endpoints expect the header `Authorization: Bearer <token>`. Tokens are valid for 24 hours by default.

## Admin Access

New users get the `user` role. To create games through the API, promote a user directly in the database, then log in again (the role is stored inside the JWT):

```bash
docker compose exec postgres psql -U postgres -d leaderboard_system \
  -c "UPDATE users SET role='admin' WHERE email='p1@test.com';"
```

## Configuration

All settings are environment variables. Defaults for local development are set in `docker-compose.yml`; see `.env.example` for a template.

| Variable | Default in Compose | Description |
|----------|--------------------|-------------|
| `SERVER_HOST` / `SERVER_PORT` | `0.0.0.0` / `8080` | HTTP listen address |
| `DB_HOST` / `DB_PORT` | `postgres` / `5432` | PostgreSQL address |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | `postgres` / `password` / `leaderboard_system` | PostgreSQL credentials |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `REDIS_HOST` / `REDIS_PORT` | `redis` / `6379` | Redis address |
| `JWT_SECRET` | `demo-secret-change-me` | JWT signing key, **change in production** |
| `JWT_EXPIRY` | `24h` | Token lifetime |
| `ENV` | `staging` | Environment name |

The Compose credentials are for local demo only. Never reuse them in production.

## Project Structure

```
cmd/api/        application entry point
internal/       application code (config, handlers, storage, ...)
pkg/utils/      shared helpers
migrations/     SQL schema and seed data
Dockerfile      multi-stage build (Go builder, Alpine runtime)
docker-compose.yml
```

## Database

The schema lives in `migrations/001_init_schema.up.sql` and demo games in `migrations/002_seed_games.up.sql`. Compose mounts both into the Postgres init directory, so they run **only on the first start with an empty volume**. To reset the database:

```bash
docker compose down -v
docker compose up --build
```

## Troubleshooting

| Symptom | Cause and fix |
|---------|---------------|
| `curl: (52) Empty reply from server` | Use `127.0.0.1` instead of `localhost` |
| `zsh: no matches found: http://...?limit=10` | Put the URL in double quotes |
| `401 Unauthorized` | Missing, wrong or expired token: log in again |
| `403 Forbidden` on `/api/games` | The user is not an admin, see Admin Access |
| Empty leaderboard `[]` | No scores submitted yet |
| Cannot connect to the Docker daemon | Start Docker Desktop and retry |

## License

MIT
