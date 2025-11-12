package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emaforlin/ce-document-service/pkg/config"
	"github.com/gin-gonic/gin"
)

type APIHTTPServer struct {
	router  *gin.Engine
	server  *http.Server
	handler *HTTPHandler
}

func (s *APIHTTPServer) Start(cfg config.ServerConfig) error {
	s.server = &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: s.router.Handler(),
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return nil
}

func (s *APIHTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown:", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}

func NewAPIServer(documentService *DocumentService) (*APIHTTPServer, error) {
	if documentService == nil {
		return nil, fmt.Errorf("documents service cannot be nil")
	}

	server := &APIHTTPServer{
		router:  gin.Default(),
		server:  &http.Server{},
		handler: NewHTTPHandler(documentService),
	}
	server.setupRoutes()
	return server, nil
}

func (s *APIHTTPServer) setupRoutes() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// ProtectedRoutes require the X-User-Id header
	protectedRoutes := s.router.Group("/")
	protectedRoutes.Use(UserHeaderMiddleware())
	{
		protectedRoutes.GET("/documents", s.handler.getDocuments)
		protectedRoutes.POST("/documents", s.handler.createDocument)
	}

	// CollaboratorRoutes require the X-User-Id user to be a document collaborator
	collaboratorRoutes := protectedRoutes.Group("/")
	collaboratorRoutes.Use(CollaboratorAccessMiddleware(s.handler.documentService))
	{
		collaboratorRoutes.GET("/documents/:id", s.handler.getOneDocument)
	}

	// OwnershipRoutes require the X-User-Id user to be a document owner
	ownerRoutes := protectedRoutes.Group("/documents/:id")
	ownerRoutes.Use(DocumentOwnershipMiddleware(s.handler.documentService))
	{
		ownerRoutes.POST("/collaborators", s.handler.addDocumentCollaborator)
		ownerRoutes.DELETE("/collaborators", s.handler.removeDocumentCollaborator)
		ownerRoutes.GET("/collaborators", s.handler.getDocumentCollaborators)
	}
}
