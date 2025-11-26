package session

import (
	"errors"
	"fmt"
	modelsErr "foodjiassignment/internal/repository/errors"
	"foodjiassignment/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Repo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateSession(expiresAt *time.Time) (*models.Session, error) {
	session := &models.Session{
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repo) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	if err := r.db.First(&session, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, modelsErr.NewNotFoundError(fmt.Sprintf("session id : %s, not found", id))
		}
		return nil, err
	}

	return &session, nil
}
