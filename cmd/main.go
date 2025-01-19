package main

import (
	"context"
	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/internal/app/handler"
	"github.com/Tairascii/google-docs-documents/internal/app/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	documents := usecase.NewDocumentsUseCase()
	useCases := app.UseCase{
		Documents: documents,
	}

	di := &app.DI{UseCase: useCases}
	handlers := handler.NewHandler(di)

	srv := &http.Server{
		Addr:         ":8080", // TODO add .env
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Something went wrong while runing server %s", err.Error())
		}
	}()

	log.Println("Listening on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Something went wrong while shutting down server %s", err.Error())
	}
}
