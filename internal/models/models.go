package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Role      UserRole  `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Game struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Score struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	GameID    int64     `json:"game_id" db:"game_id"`
	Score     int64     `json:"score" db:"score"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Score    int64  `json:"score"`
	Rank     int64  `json:"rank"`
}

type GlobalLeaderboardEntry struct {
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	TotalScore int64  `json:"total_score"`
	Rank       int64  `json:"rank"`
}

type UserStats struct {
	UserID      int64 `json:"user_id"`
	TotalGames  int   `json:"total_games"`
	TotalScore  int64 `json:"total_score"`
	AverageScore float64 `json:"average_score"`
	BestScore   int64 `json:"best_score"`
}

type TopPlayerReport struct {
	StartDate time.Time                `json:"start_date"`
	EndDate   time.Time                `json:"end_date"`
	TopPlayers []GlobalLeaderboardEntry `json:"top_players"`
}
