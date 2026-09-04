# Leaderboard System

Real-time leaderboard service for games. Go + Redis + PostgreSQL.

## Quick Start

Requires Docker.

```bash
docker compose up --build
```

Server runs on `http://localhost:8080`.

## Example

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","email":"p1@test.com","password":"pass123"}'

# Login — returns a token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"p1@test.com","password":"pass123"}'

# Submit score (use the token from login)
curl -X POST http://localhost:8080/api/scores \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"game_id":1,"score":1500}'

# Global leaderboard
curl "http://localhost:8080/api/leaderboard/global?limit=10"
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register |
| POST | `/api/auth/login` | Login |
| GET | `/api/games` | List games |
| POST | `/api/games` | Create game (admin) |
| DELETE | `/api/games/{id}` | Delete game (admin) |
| POST | `/api/scores` | Submit score |
| GET | `/api/leaderboard/global?limit=10` | Global leaderboard |
| GET | `/api/leaderboard/game/{id}?limit=10` | Game leaderboard |
| GET | `/api/stats/me` | Your stats |

## Tech Stack

- Go
- Redis (sorted sets — real-time leaderboards)
- PostgreSQL (score history & reports)
- JWT auth

## License

MIT
