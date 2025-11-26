package tinder_food

import (
	"foodjiassignment/internal/repository/models"
	"gorm.io/gorm"
	"time"
)

type TinderFoodRepo struct {
	db *gorm.DB
}

func NewTinderFoodRepo(db *gorm.DB) *TinderFoodRepo {
	return &TinderFoodRepo{
		db: db,
	}
}

func (r *TinderFoodRepo) CreateSession(expiresAt *time.Time) (*models.Session, error) {
	session := &models.Session{
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}
