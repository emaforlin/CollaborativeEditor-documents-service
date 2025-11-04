package repository

import "gorm.io/gorm"

type PostgresDocumentRepositoryImpl struct {
	db *gorm.DB
}

func (r *PostgresDocumentRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}
