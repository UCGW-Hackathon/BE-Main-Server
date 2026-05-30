package repositories

import (
	"context"
	"gorm.io/gorm"
)

type ConnectionRepository interface {
	Connect(ctx context.Context) error
}

type connectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) ConnectionRepository {
	return &connectionRepository{db: db}
}

func (cr *connectionRepository) Connect(ctx context.Context) error {
	return nil
}

