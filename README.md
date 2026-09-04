# Leaderboard System

Real-time leaderboard service built with Go, Redis and PostgreSQL.

## Tech Stack

- Go 1.23
- Redis (sorted sets)
- PostgreSQL
- JWT auth

## Features

- User registration & login
- Submit scores for different games
- Real-time leaderboard updates via Redis sorted sets
- Game-specific and global leaderboards
- User rankings
- Score history tracking
- Top players reports by time period
- User statistics

## Quick Start (Docker — recommended)

### Prerequisites

- Docker & Docker Compose
- Go 1.23+

### Setup

```bash
# 1. Start PostgreSQL and Redis
docker-compose up -d

# 2. Create .env file (edit if needed)
cp .env.example .env

# 3. Build and run
make dev
```

Server runs on `http://localhost:8080`

### Manual setup (without Docker)

```bash
# 1. Start PostgreSQL and Redis on your machine

# 2. Create database and run migrations
createdb leaderboard_system
psql -d leaderboard_system -f migrations/001_init_schema.up.sql

# 3. Create .env file
cp .env.example .env
# edit .env with your credentials

# 4. Run
make dev
```

## API Endpoints

### Auth
```
POST /api/auth/register
POST /api/auth/login
```

### Games
```
GET    /api/games           - list all games
GET    /api/games/{id}      - get game
POST   /api/games           - create game (admin)
DELETE /api/games/{id}      - delete game (admin)
```

### Leaderboard
```
POST /api/scores                        - submit score (auth)
GET  /api/leaderboard/global?limit=10   - global leaderboard
GET  /api/leaderboard/game/{id}?limit=10 - game leaderboard
GET  /api/leaderboard/game/{id}/rank    - your rank in game (auth)
GET  /api/leaderboard/rank               - your global rank (auth)
GET  /api/stats/me                       - your stats (auth)
```

### Reports
```
GET /api/reports/top-players?start_date=2026-01-01&end_date=2026-12-31&limit=10 (admin)
```

## Step-by-Step Usage Example

### 1. Register an admin user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

Response will contain a `token` — save it. Then manually update the user role in the database:

```bash
psql -d leaderboard_system -c "UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';"
```

### 2. Login and get a fresh token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

Save the `token` from the response as `ADMIN_TOKEN`.

### 3. Create a game

```bash
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Chess",
    "description": "Classic chess game"
  }'
```

### 4. Register a regular player

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "password123"
  }'
```

Save the `token` as `PLAYER_TOKEN`.

### 5. Submit scores

```bash
curl -X POST http://localhost:8080/api/scores \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -d '{
    "game_id": 1,
    "score": 1500
  }'
```

### 6. Check the leaderboard

```bash
# Global leaderboard
curl http://localhost:8080/api/leaderboard/global?limit=10

# Game-specific leaderboard
curl http://localhost:8080/api/leaderboard/game/1?limit=10

# Your rank in a game
curl http://localhost:8080/api/leaderboard/game/1/rank \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Your global rank
curl http://localhost:8080/api/leaderboard/rank \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Your stats
curl http://localhost:8080/api/stats/me \
  -H "Authorization: Bearer $PLAYER_TOKEN"
```

### 7. Top players report (admin)

```bash
curl "http://localhost:8080/api/reports/top-players?start_date=2026-01-01&end_date=2026-12-31&limit=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Config

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=leaderboard_system

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-secret-here
```

## Makefile Commands

```bash
make build       # Build binary
make run         # Build and run
make dev         # Run without building (go run)
make test        # Run tests
make docker-up   # Start PostgreSQL + Redis
make docker-down # Stop PostgreSQL + Redis
make tidy        # go mod tidy
```

## Architecture

```
cmd/api/              - entry point
internal/
  ├── handler/        - HTTP handlers
  ├── service/        - business logic
  ├── repository/     - data access (DB + Redis)
  ├── models/         - domain models
  ├── middleware/     - auth + logging middleware
  ├── config/         - configuration
  ├── database/       - PostgreSQL connection
  └── redis/          - Redis connection
```

## Redis Keys

- `leaderboard:game:{gameID}` - per-game leaderboard (sorted set)
- `leaderboard:global` - global leaderboard (sorted set)

Redis sorted sets provide O(log N) complexity for:
- Adding scores
- Getting rankings
- Fetching top N players

## How It Works

1. User submits score via `/api/scores`
2. Score saved to PostgreSQL (history)
3. Score added to Redis sorted sets:
   - Game-specific: `ZINCRBY leaderboard:game:1 {score} {userID}`
   - Global: `ZINCRBY leaderboard:global {score} {userID}`
4. Leaderboard queries hit Redis directly (fast)
5. User data enriched from PostgreSQL

This hybrid approach gives:
- Fast real-time leaderboard queries (Redis)
- Persistent score history (PostgreSQL)
- Time-based reports (PostgreSQL)

## Notes

- Leaderboards limited to 100 entries per request
- Scores must be non-negative (enforced at DB and app level)
- User stats calculated from PostgreSQL
- Graceful shutdown on SIGINT/SIGTERM
- Request logging on all endpoints

## License

MIT
