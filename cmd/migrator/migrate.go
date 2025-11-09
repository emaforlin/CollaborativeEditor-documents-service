package main

import (
	"log"

	"github.com/emaforlin/ce-document-service/config"
	"github.com/emaforlin/ce-document-service/models"
	"github.com/emaforlin/ce-document-service/repository"
)

type DatabaseMigrator struct {
	repo *repository.PostgresDocumentRepositoryImpl
}

func NewMigrator(cfg config.DatabaseConfig) *DatabaseMigrator {
	repo := repository.NewPostgresRepository(cfg)
	return &DatabaseMigrator{
		repo: repo,
	}
}

// RunMigrations executes all database migrations
func (m *DatabaseMigrator) RunMigrations() error {
	log.Println("Starting database migrations...")

	// Add all models that need to be migrated here
	models := []interface{}{
		&models.Document{},
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
