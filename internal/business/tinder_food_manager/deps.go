package tinderfood

import (
	"foodtinder/internal/repository/models"
	"github.com/google/uuid"
	"time"
)

type sessionRepo interface {
	CreateSession(expiresAt *time.Time) (*models.Session, error)
	GetSession(sessionID uuid.UUID) (*models.Session, error)
}

type productVoteRepo interface {
	UpsertProductVote(vote *models.ProductVote) error
	GetVotesBySessionId(sessionID uuid.UUID) ([]models.ProductVote, error)
	GetAverageScores() ([]models.ProductScore, error)
}
