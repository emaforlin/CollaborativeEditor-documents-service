package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	document "github.com/emaforlin/ce-document-service/internal/document"
	"github.com/emaforlin/ce-document-service/pkg/config"
)

func main() {
	config.Load()
	configuration := config.GetConfig()

	repository := document.NewPostgresRepository(configuration.GetDatabaseConf())

	service, err := document.NewDocumentService(repository)
	if err != nil {
		log.Fatal("failed to start the service:", err)
	}

	server, err := document.NewAPIServer(service)
	if err != nil {
		log.Fatal("failed to initialize the server:", err)
	}

	if err := server.Start(configuration.GetServerConf()); err != nil {
		log.Fatal(err)
	}
	defer server.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server...")

}
