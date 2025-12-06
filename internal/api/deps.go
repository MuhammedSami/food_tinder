package api

import (
	"context"
	apiModels "foodtinder/internal/api/models"
	"foodtinder/internal/repository/models"
	"github.com/google/uuid"
)

type manager interface {
	CreateSession(ctx context.Context) (string, error)
	UpsertVote(ctx context.Context, productVote apiModels.UpsertProductVoteRequest) error

	GetSession(sessionId uuid.UUID) (*models.Session, error)
	GetAverageScores() ([]models.ProductScore, error)
	GetVotesBySession(ctx context.Context, sessionId uuid.UUID) ([]models.ProductVote, error)
}
