package service

import (
	"fmt"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repository"
)

type GameService struct {
	gameRepo *repository.GameRepository
}

func NewGameService(gameRepo *repository.GameRepository) *GameService {
	return &GameService{gameRepo: gameRepo}
}

type CreateGameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *GameService) CreateGame(req *CreateGameRequest) (*models.Game, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("game name is required")
	}

	game := &models.Game{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.gameRepo.Create(game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (s *GameService) GetGame(id int64) (*models.Game, error) {
	return s.gameRepo.GetByID(id)
}

func (s *GameService) GetAllGames() ([]models.Game, error) {
	return s.gameRepo.GetAll()
}

func (s *GameService) DeleteGame(id int64) error {
	return s.gameRepo.Delete(id)
}
