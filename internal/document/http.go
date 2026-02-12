package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emaforlin/ce-document-service/pkg/config"
	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-User-Id"}
	s.router.Use(cors.New(config))

	// ProtectedRoutes require the X-User-Id header
	protectedRoutes := s.router.Group("/")
	protectedRoutes.Use(UserHeaderMiddleware())
	{
		protectedRoutes.GET("/documents", s.handler.getDocuments)
		protectedRoutes.POST("/documents", s.handler.createDocument)
	}

	// Document routes with specific permission requirements
	documentRoutes := protectedRoutes.Group("/documents/:id")
	{
		// Routes that require viewer access (read-only)
		documentRoutes.GET("", RequireViewerAccess(s.handler.documentService), s.handler.getOneDocument)

		// Routes that require editor access (can modify content)
		documentRoutes.PATCH("", RequireEditorAccess(s.handler.documentService), s.handler.updateDocument)

		// Routes that require owner access (can manage permissions)
		documentRoutes.DELETE("", RequireOwnerAccess(s.handler.documentService), s.handler.deleteDocument)
		documentRoutes.POST("/collaborators", RequireOwnerAccess(s.handler.documentService), s.handler.addDocumentCollaborator)
		documentRoutes.DELETE("/collaborators", RequireOwnerAccess(s.handler.documentService), s.handler.removeDocumentCollaborator)
		documentRoutes.GET("/collaborators", RequireOwnerAccess(s.handler.documentService), s.handler.getDocumentCollaborators)
	}
}
