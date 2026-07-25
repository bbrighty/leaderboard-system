package repository

import (
	"fmt"
	"leaderboard-system/internal/models"
	redisClient "leaderboard-system/internal/redis"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type LeaderboardRepository struct {
	redis *redisClient.RedisClient
}

func NewLeaderboardRepository(redis *redisClient.RedisClient) *LeaderboardRepository {
	return &LeaderboardRepository{redis: redis}
}

func (r *LeaderboardRepository) getGameKey(gameID int64) string {
	return fmt.Sprintf("leaderboard:game:%d", gameID)
}

func (r *LeaderboardRepository) getGlobalKey() string {
	return "leaderboard:global"
}

func (r *LeaderboardRepository) AddScore(userID, gameID int64, score int64) error {
	key := r.getGameKey(gameID)
	member := fmt.Sprintf("%d", userID)
	
	return r.redis.Client.ZIncrBy(r.redis.Ctx, key, float64(score), member).Err()
}

func (r *LeaderboardRepository) AddToGlobal(userID int64, score int64) error {
	key := r.getGlobalKey()
	member := fmt.Sprintf("%d", userID)
	
	return r.redis.Client.ZIncrBy(r.redis.Ctx, key, float64(score), member).Err()
}

func (r *LeaderboardRepository) GetTopPlayers(gameID int64, limit int) ([]models.LeaderboardEntry, error) {
	key := r.getGameKey(gameID)
	
	results, err := r.redis.Client.ZRevRangeWithScores(r.redis.Ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]models.LeaderboardEntry, 0, len(results))
	for i, result := range results {
		userID, _ := strconv.ParseInt(result.Member.(string), 10, 64)
		entries = append(entries, models.LeaderboardEntry{
			UserID: userID,
			Score:  int64(result.Score),
			Rank:   int64(i + 1),
		})
	}
	
	return entries, nil
}

func (r *LeaderboardRepository) GetGlobalTopPlayers(limit int) ([]models.GlobalLeaderboardEntry, error) {
	key := r.getGlobalKey()
	
	results, err := r.redis.Client.ZRevRangeWithScores(r.redis.Ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]models.GlobalLeaderboardEntry, 0, len(results))
	for i, result := range results {
		userID, _ := strconv.ParseInt(result.Member.(string), 10, 64)
		entries = append(entries, models.GlobalLeaderboardEntry{
			UserID:     userID,
			TotalScore: int64(result.Score),
			Rank:       int64(i + 1),
		})
	}
	
	return entries, nil
}

func (r *LeaderboardRepository) GetUserRank(userID, gameID int64) (int64, int64, error) {
	key := r.getGameKey(gameID)
	member := fmt.Sprintf("%d", userID)
	
	rank, err := r.redis.Client.ZRevRank(r.redis.Ctx, key, member).Result()
	if err == redis.Nil {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	
	score, err := r.redis.Client.ZScore(r.redis.Ctx, key, member).Result()
	if err != nil {
		return 0, 0, err
	}
	
	return rank + 1, int64(score), nil
}

func (r *LeaderboardRepository) GetUserGlobalRank(userID int64) (int64, int64, error) {
	key := r.getGlobalKey()
	member := fmt.Sprintf("%d", userID)
	
	rank, err := r.redis.Client.ZRevRank(r.redis.Ctx, key, member).Result()
	if err == redis.Nil {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	
	score, err := r.redis.Client.ZScore(r.redis.Ctx, key, member).Result()
	if err != nil {
		return 0, 0, err
	}
	
	return rank + 1, int64(score), nil
}

func (r *LeaderboardRepository) GetUserScoreInGame(userID, gameID int64) (int64, error) {
	key := r.getGameKey(gameID)
	member := fmt.Sprintf("%d", userID)
	
	score, err := r.redis.Client.ZScore(r.redis.Ctx, key, member).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	
	return int64(score), nil
}
