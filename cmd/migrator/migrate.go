package main

import (
	"log"

	document "github.com/emaforlin/ce-document-service/internal/document"
	"github.com/emaforlin/ce-document-service/pkg/config"
)

type DatabaseMigrator struct {
	repo *document.PostgresDocumentRepositoryImpl
}

func NewMigrator(cfg config.DatabaseConfig) *DatabaseMigrator {
	repo := document.NewPostgresRepository(cfg)
	return &DatabaseMigrator{
		repo: repo,
	}
}

// RunMigrations executes all database migrations
func (m *DatabaseMigrator) RunMigrations() error {
	log.Println("Starting database migrations...")

	// Add all models that need to be migrated here
	models := []interface{}{
		&document.Document{},
	}

	if err := m.repo.GetDB().AutoMigrate(models...); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully!")
	return nil
}

func main() {
	cfg := config.GetConfig()
	migrator := NewMigrator(cfg.GetDatabaseConf())

	if err := migrator.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
