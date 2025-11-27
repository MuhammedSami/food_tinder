package api

import (
	"context"
	apiModels "foodjiassignment/internal/api/models"
	"foodjiassignment/internal/repository/models"
	"github.com/google/uuid"
)

type manager interface {
	CreateSession(ctx context.Context) (string, error)
	UpsertVote(ctx context.Context, productVote apiModels.UpsertProductVoteRequest) error

	GetSession(sessionId uuid.UUID) (*models.Session, error)
	GetVotesBySession(ctx context.Context, sessionId uuid.UUID) ([]models.ProductVote, error)
}
