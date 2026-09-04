package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leaderboard-system/internal/config"
	"leaderboard-system/internal/database"
	"leaderboard-system/internal/handler"
	"leaderboard-system/internal/middleware"
	redisClient "leaderboard-system/internal/redis"
	"leaderboard-system/internal/repository"
	"leaderboard-system/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	redis, err := redisClient.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("Connected to Redis")

	userRepo := repository.NewUserRepository(db.DB)
	gameRepo := repository.NewGameRepository(db.DB)
	scoreRepo := repository.NewScoreRepository(db.DB)
	leaderboardRepo := repository.NewLeaderboardRepository(redis)

	authService := service.NewAuthService(userRepo, cfg)
	gameService := service.NewGameService(gameRepo, leaderboardRepo)
	leaderboardService := service.NewLeaderboardService(leaderboardRepo, scoreRepo, gameRepo, userRepo)

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	authHandler := handler.NewAuthHandler(authService)
	gameHandler := handler.NewGameHandler(gameService)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)

	router := mux.NewRouter()

	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	router.HandleFunc("/api/games", gameHandler.GetAllGames).Methods("GET")
	router.HandleFunc("/api/games/{id}", gameHandler.GetGame).Methods("GET")
	router.HandleFunc("/api/games", authMiddleware.RequireAdmin(gameHandler.CreateGame)).Methods("POST")
	router.HandleFunc("/api/games/{id}", authMiddleware.RequireAdmin(gameHandler.DeleteGame)).Methods("DELETE")

	router.HandleFunc("/api/scores", authMiddleware.RequireAuth(leaderboardHandler.SubmitScore)).Methods("POST")
	router.HandleFunc("/api/leaderboard/global", leaderboardHandler.GetGlobalLeaderboard).Methods("GET")
	router.HandleFunc("/api/leaderboard/game/{id}", leaderboardHandler.GetGameLeaderboard).Methods("GET")
	router.HandleFunc("/api/leaderboard/game/{id}/rank", authMiddleware.RequireAuth(leaderboardHandler.GetUserRanking)).Methods("GET")
	router.HandleFunc("/api/leaderboard/rank", authMiddleware.RequireAuth(leaderboardHandler.GetUserGlobalRanking)).Methods("GET")
	router.HandleFunc("/api/stats/me", authMiddleware.RequireAuth(leaderboardHandler.GetUserStats)).Methods("GET")
	router.HandleFunc("/api/reports/top-players", authMiddleware.RequireAdmin(leaderboardHandler.GetTopPlayersReport)).Methods("GET")

	router.Use(middleware.Logging)
	router.Use(corsMiddleware)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"leaderboard-api"}`))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
