package service

import (
	"fmt"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repository"
	"time"
)

type LeaderboardService struct {
	leaderboardRepo *repository.LeaderboardRepository
	scoreRepo       *repository.ScoreRepository
	gameRepo        *repository.GameRepository
	userRepo        *repository.UserRepository
}

func NewLeaderboardService(
	leaderboardRepo *repository.LeaderboardRepository,
	scoreRepo *repository.ScoreRepository,
	gameRepo *repository.GameRepository,
	userRepo *repository.UserRepository,
) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
		scoreRepo:       scoreRepo,
		gameRepo:        gameRepo,
		userRepo:        userRepo,
	}
}

type SubmitScoreRequest struct {
	GameID int64 `json:"game_id"`
	Score  int64 `json:"score"`
}

func (s *LeaderboardService) SubmitScore(userID int64, req *SubmitScoreRequest) error {
	if req.Score < 0 {
		return fmt.Errorf("score must be non-negative")
	}

	_, err := s.gameRepo.GetByID(req.GameID)
	if err != nil {
		return fmt.Errorf("game not found")
	}

	score := &models.Score{
		UserID: userID,
		GameID: req.GameID,
		Score:  req.Score,
	}

	if err := s.scoreRepo.Create(score); err != nil {
		return fmt.Errorf("failed to save score: %w", err)
	}

	if err := s.leaderboardRepo.AddScore(userID, req.GameID, req.Score); err != nil {
		return fmt.Errorf("failed to update game leaderboard: %w", err)
	}

	if err := s.leaderboardRepo.AddToGlobal(userID, req.Score); err != nil {
		return fmt.Errorf("failed to update global leaderboard: %w", err)
	}

	return nil
}

func (s *LeaderboardService) GetGameLeaderboard(gameID int64, limit int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	entries, err := s.leaderboardRepo.GetTopPlayers(gameID, limit)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, len(entries))
	for i, entry := range entries {
		userIDs[i] = entry.UserID
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		if user, ok := users[entries[i].UserID]; ok {
			entries[i].Username = user.Username
		}
	}

	return entries, nil
}

func (s *LeaderboardService) GetGlobalLeaderboard(limit int) ([]models.GlobalLeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	entries, err := s.leaderboardRepo.GetGlobalTopPlayers(limit)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, len(entries))
	for i, entry := range entries {
		userIDs[i] = entry.UserID
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		if user, ok := users[entries[i].UserID]; ok {
			entries[i].Username = user.Username
		}
	}

	return entries, nil
}

func (s *LeaderboardService) GetUserRanking(userID, gameID int64) (*models.LeaderboardEntry, error) {
	rank, score, err := s.leaderboardRepo.GetUserRank(userID, gameID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.LeaderboardEntry{
		UserID:   userID,
		Username: user.Username,
		Score:    score,
		Rank:     rank,
	}, nil
}

func (s *LeaderboardService) GetUserGlobalRanking(userID int64) (*models.GlobalLeaderboardEntry, error) {
	rank, score, err := s.leaderboardRepo.GetUserGlobalRank(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.GlobalLeaderboardEntry{
		UserID:     userID,
		Username:   user.Username,
		TotalScore: score,
		Rank:       rank,
	}, nil
}

func (s *LeaderboardService) GetUserStats(userID int64) (*models.UserStats, error) {
	return s.scoreRepo.GetUserStats(userID)
}

func (s *LeaderboardService) GetTopPlayersReport(startDate, endDate time.Time, limit int) (*models.TopPlayerReport, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	players, err := s.scoreRepo.GetTopPlayersByPeriod(startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	return &models.TopPlayerReport{
		StartDate:  startDate,
		EndDate:    endDate,
		TopPlayers: players,
	}, nil
}
