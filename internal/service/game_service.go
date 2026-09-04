package service

import (
	"fmt"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repository"
)

type GameService struct {
	gameRepo        *repository.GameRepository
	leaderboardRepo *repository.LeaderboardRepository
}

func NewGameService(gameRepo *repository.GameRepository, leaderboardRepo *repository.LeaderboardRepository) *GameService {
	return &GameService{gameRepo: gameRepo, leaderboardRepo: leaderboardRepo}
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
	if err := s.gameRepo.Delete(id); err != nil {
		return err
	}
	s.leaderboardRepo.DeleteGameKey(id)
	return nil
}
