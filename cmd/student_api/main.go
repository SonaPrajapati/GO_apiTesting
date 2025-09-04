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

	"github.com/SonaPrajapati/GO_apiTesing/internal/config"
	"github.com/SonaPrajapati/GO_apiTesing/internal/config/http/handlers/student"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// setup database
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server is started", slog.String("address", cfg.HTTPServer.Address))
	// fmt.Printf("Server is started %s", cfg.HTTPServer.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// gracefull shutdown
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
