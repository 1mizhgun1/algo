package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"algo/algorithms/a_star"
	"algo/config"
	"algo/handlers"
	"algo/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg := config.MustLoadConfig(os.Getenv("CONFIG_FILE"), logger)
	logger.Info("Config file loaded")

	reqIDMiddleware := middleware.CreateRequestIDMiddleware(logger)

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	r.Use(reqIDMiddleware, middleware.CorsMiddleware, middleware.RecoverMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.Handle("/calc_path", http.HandlerFunc(handlers.SolveMazeHandler)).Methods(http.MethodPost, http.MethodOptions)

	a_star.TestAStar()

	http.Handle("/", r)
	server := http.Server{
		Handler:           middleware.PathMiddleware(r),
		Addr:              fmt.Sprintf(":%s", cfg.Main.Port),
		ReadTimeout:       cfg.Main.ReadTimeout,
		WriteTimeout:      cfg.Main.WriteTimeout,
		ReadHeaderTimeout: cfg.Main.ReadHeaderTimeout,
		IdleTimeout:       cfg.Main.IdleTimeout,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Info("Server stopped")
		}
	}()
	logger.Info("Server started")

	sig := <-signalCh
	logger.Info("Received signal: " + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Main.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(errors.Wrap(err, "failed to gracefully shutdown").Error())
	}
}
