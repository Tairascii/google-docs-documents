package main

import (
	"context"
	"fmt"
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
	cfg, err := app.LoadConfigs()
	if err != nil {
		panic(err)
	}

	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address: fmt.Sprintf("%s:%s", cfg.Repo.Port, cfg.Repo.Port),
	})
	if err != nil {
		log.Fatal("Something went wrong connecting to rethink")
	}

	documentRepo := docRepo.NewRepo(docRepo.Params{
		Session:        session,
		DocumentsTable: cfg.Repo.DocumentsTable,
	})

	go documentRepo.WatchTableChange()
	documentService := document.New(documentRepo)
	documents := usecase.NewDocumentsUseCase(documentService)

	useCases := app.UseCase{
		Documents: documents,
	}

	di := &app.DI{UseCase: useCases}
	handlers := handler.NewHandler(di, session)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Something went wrong while runing server %s", err.Error())
		}
	}()

	log.Printf("Listening on port %s\n", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Something went wrong while shutting down server %s", err.Error())
	}
}
