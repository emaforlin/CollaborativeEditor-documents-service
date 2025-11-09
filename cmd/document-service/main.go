package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/emaforlin/ce-document-service/api"
	"github.com/emaforlin/ce-document-service/config"
	"github.com/emaforlin/ce-document-service/core"
	"github.com/emaforlin/ce-document-service/repository"
)

func main() {
	config.Load()
	configuration := config.GetConfig()

	repository := repository.NewMockRepository()

	service, err := core.NewDocumentService(repository)
	if err != nil {
		log.Fatal("failed to start the service:", err)
	}

	server, err := api.NewAPIServer(service)
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
