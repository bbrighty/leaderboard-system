package handler

import (
	"encoding/json"
	"leaderboard-system/internal/middleware"
	"leaderboard-system/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
	authMiddleware     *middleware.AuthMiddleware
}

func NewLeaderboardHandler(leaderboardService *service.LeaderboardService, authMiddleware *middleware.AuthMiddleware) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
		authMiddleware:     authMiddleware,
	}
}

func (h *LeaderboardHandler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req service.SubmitScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.leaderboardService.SubmitScore(claims.UserID, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Score submitted successfully"})
}

func (h *LeaderboardHandler) GetGameLeaderboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	leaderboard, err := h.leaderboardService.GetGameLeaderboard(gameID, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, leaderboard)
}

func (h *LeaderboardHandler) GetGlobalLeaderboard(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	leaderboard, err := h.leaderboardService.GetGlobalLeaderboard(limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, leaderboard)
}

func (h *LeaderboardHandler) GetUserRanking(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	gameID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid game ID")
		return
	}

	ranking, err := h.leaderboardService.GetUserRanking(claims.UserID, gameID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ranking)
}

func (h *LeaderboardHandler) GetUserGlobalRanking(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ranking, err := h.leaderboardService.GetUserGlobalRanking(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ranking)
}

func (h *LeaderboardHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stats, err := h.leaderboardService.GetUserStats(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *LeaderboardHandler) GetTopPlayersReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")

	if startDateStr == "" || endDateStr == "" {
		respondWithError(w, http.StatusBadRequest, "start_date and end_date are required (format: YYYY-MM-DD)")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid end_date format")
		return
	}

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	report, err := h.leaderboardService.GetTopPlayersReport(startDate, endDate, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, report)
}
