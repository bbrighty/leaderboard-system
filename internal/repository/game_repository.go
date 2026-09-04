package repository

import (
	"database/sql"
	"fmt"
	"leaderboard-system/internal/models"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(game *models.Game) error {
	query := `
		INSERT INTO games (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, game.Name, game.Description).
		Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)
}

func (r *GameRepository) GetByID(id int64) (*models.Game, error) {
	game := &models.Game{}
	query := `SELECT id, name, description, created_at, updated_at FROM games WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&game.ID, &game.Name, &game.Description, &game.CreatedAt, &game.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("game not found")
	}
	return game, err
}

func (r *GameRepository) GetAll() ([]models.Game, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM games ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]models.Game, 0)
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.Name, &game.Description, &game.CreatedAt, &game.UpdatedAt); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, rows.Err()
}

func (r *GameRepository) Delete(id int64) error {
	query := `DELETE FROM games WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("game not found")
	}
	return nil
}
