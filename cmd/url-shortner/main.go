package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortner/internal/config"
	"url-shortner/internal/http/handlers/shorten"
	"url-shortner/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()

	router := http.NewServeMux()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router.HandleFunc("POST /api/shorten", shorten.New(storage))
	router.HandleFunc("GET /api/shorten", shorten.GetList(storage))
	router.HandleFunc("GET /api/shorten/id/{id}", shorten.GetById(storage))
	router.HandleFunc("GET /api/shorten/redirect", shorten.GetByShortenUrl(storage))

	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
