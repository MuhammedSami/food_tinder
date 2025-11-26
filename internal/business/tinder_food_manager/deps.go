package tinderfood

import (
	"foodjiassignment/internal/repository/models"
	"time"
)

type sessionRepo interface {
	CreateSession(expiresAt *time.Time) (*models.Session, error)
	GetSession(sessionID string) (*models.Session, error)
}

type productVoteRepo interface {
	UpsertProductVote(vote *models.ProductVote) error
}
