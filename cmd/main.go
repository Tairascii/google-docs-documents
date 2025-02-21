package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Tairascii/google-docs-documents/pkg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/internal/app/handler"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
	docRepo "github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"
	"github.com/Tairascii/google-docs-documents/internal/app/usecase"
	"github.com/dancannon/gorethink"
)

func main() {
	logger := pkg.NewLogger()
	cfg, err := app.LoadConfigs()
	if err != nil {
		log.Fatalf("[ERROR] load config: %s", err)
	}

	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address: fmt.Sprintf("%s:%s", cfg.Repo.Host, cfg.Repo.Port),
	})
	if err != nil {
		log.Fatalf("[ERROR] connect to rethink %s", err)
	}

	documentRepo := docRepo.NewRepo(docRepo.Params{
		Session:        session,
		DocumentsTable: cfg.Repo.DocumentsTable,
	})

	documentService := document.New(documentRepo)
	documents := usecase.NewDocumentsUseCase(documentService)

	useCases := app.UseCase{
		Documents: documents,
	}

	di := &app.DI{UseCase: useCases}
	handlers := handler.NewHandler(di)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[ERROR] listen %s", err)
		}
	}()

	logger.Info("listening on port", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	logger.Info("shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("[ERROR] shutdown %s", err)
	}
}
