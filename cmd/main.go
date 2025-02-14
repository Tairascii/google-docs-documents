package main

import (
	"context"
	"github.com/Tairascii/google-docs-documents/internal/app"
	"github.com/Tairascii/google-docs-documents/internal/app/handler"
	"github.com/Tairascii/google-docs-documents/internal/app/service/document"
	docRepo "github.com/Tairascii/google-docs-documents/internal/app/service/document/repo"
	"github.com/Tairascii/google-docs-documents/internal/app/usecase"
	"github.com/dancannon/gorethink"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address: "localhost:28015", // TODO change to .env
	})
	if err != nil {
		log.Fatal("Something went wrong connecting to rethink")
	}

	documentRepo := docRepo.NewRepo(docRepo.Params{
		Session:        session,
		DocumentsTable: "documents",
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
		Addr:         ":8000", // TODO add .env
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
