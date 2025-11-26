package api

import (
	"context"
	apiModels "foodjiassignment/internal/api/models"
)

type manager interface {
	CreateSession(ctx context.Context) (string, error)
	GetSession(sessionId string) (string, error)
	UpsertVote(ctx context.Context, productVote apiModels.UpsertProductVoteRequest) error
}
