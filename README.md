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

## Quick Start

### Prerequisites

- PostgreSQL
- Redis
- Go 1.23+

### Setup

```bash
createdb leaderboard_system
psql -d leaderboard_system -f migrations/001_init_schema.up.sql

cp .env.example .env
# edit .env with your credentials

go run cmd/api/main.go
```

Server runs on `localhost:8080`

### Config

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=leaderboard_system

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-secret-here
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

## Examples

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "password123"
  }'
```

### Submit Score
```bash
curl -X POST http://localhost:8080/api/scores \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "game_id": 1,
    "score": 1500
  }'
```

### Get Global Leaderboard
```bash
curl http://localhost:8080/api/leaderboard/global?limit=10
```

## Architecture

```
cmd/api/              - entry point
internal/
  ├── handler/        - HTTP handlers
  ├── service/        - business logic
  ├── repository/     - data access (DB + Redis)
  ├── models/         - domain models
  ├── middleware/     - auth middleware
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
- Scores must be non-negative
- User stats calculated from PostgreSQL (not real-time cached yet)

## License

MIT
