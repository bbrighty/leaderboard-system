package repository

import (
	"database/sql"
	"leaderboard-system/internal/models"
	"time"
)

type ScoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Create(score *models.Score) error {
	query := `
		INSERT INTO scores (user_id, game_id, score)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, score.UserID, score.GameID, score.Score).
		Scan(&score.ID, &score.CreatedAt)
}

func (r *ScoreRepository) GetUserScores(userID int64) ([]models.Score, error) {
	query := `SELECT id, user_id, game_id, score, created_at FROM scores WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make([]models.Score, 0)
	for rows.Next() {
		var score models.Score
		if err := rows.Scan(&score.ID, &score.UserID, &score.GameID, &score.Score, &score.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func (r *ScoreRepository) GetGameScores(gameID int64) ([]models.Score, error) {
	query := `SELECT id, user_id, game_id, score, created_at FROM scores WHERE game_id = $1 ORDER BY score DESC, created_at ASC`
	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make([]models.Score, 0)
	for rows.Next() {
		var score models.Score
		if err := rows.Scan(&score.ID, &score.UserID, &score.GameID, &score.Score, &score.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func (r *ScoreRepository) GetUserStats(userID int64) (*models.UserStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT game_id) as total_games,
			COALESCE(SUM(score), 0) as total_score,
			COALESCE(AVG(score), 0) as average_score,
			COALESCE(MAX(score), 0) as best_score
		FROM scores
		WHERE user_id = $1
	`
	stats := &models.UserStats{UserID: userID}
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalGames,
		&stats.TotalScore,
		&stats.AverageScore,
		&stats.BestScore,
	)
	return stats, err
}

func (r *ScoreRepository) GetTopPlayersByPeriod(startDate, endDate time.Time, limit int) ([]models.GlobalLeaderboardEntry, error) {
	query := `
		SELECT u.id, u.username, COALESCE(SUM(s.score), 0) as total_score
		FROM users u
		INNER JOIN scores s ON u.id = s.user_id AND s.created_at BETWEEN $1 AND $2
		GROUP BY u.id, u.username
		HAVING SUM(s.score) > 0
		ORDER BY total_score DESC
		LIMIT $3
	`
	rows, err := r.db.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]models.GlobalLeaderboardEntry, 0)
	rank := int64(1)
	for rows.Next() {
		var entry models.GlobalLeaderboardEntry
		if err := rows.Scan(&entry.UserID, &entry.Username, &entry.TotalScore); err != nil {
			return nil, err
		}
		entry.Rank = rank
		entries = append(entries, entry)
		rank++
	}
	return entries, rows.Err()
}
