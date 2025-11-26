package tinderfood

import (
	"context"
	"fmt"
	apiModels "foodjiassignment/internal/api/models"
	"foodjiassignment/internal/repository/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	sessionRepo     sessionRepo
	productVoteRepo productVoteRepo
	redis           *redis.Client // I thought to use redis for caching session but I rolled back implementation
}

func NewManager(
	repo sessionRepo,
	productVoteRepo productVoteRepo,
) *Manager {
	return &Manager{
		sessionRepo:     repo,
		productVoteRepo: productVoteRepo,
	}
}

func (m *Manager) CreateSession(ctx context.Context) (string, error) {
	session, err := m.sessionRepo.CreateSession(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create session err:%w", err)
	}

	return session.ID.String(), nil
}

func (m *Manager) GetSession(sessionId string) (string, error) {
	session, err := m.sessionRepo.GetSession(sessionId)
	if err != nil {
		return "", fmt.Errorf("failed to create session err:%w", err)
	}

	return session.ID.String(), nil
}

func (m *Manager) UpsertVote(ctx context.Context, productVote apiModels.UpsertProductVoteRequest) error {
	sessId, err := uuid.Parse(productVote.SessionId)
	if err != nil {
		return fmt.Errorf("failed to parse session id")
	}

	productId, err := uuid.Parse(productVote.ProductId)
	if err != nil {
		return fmt.Errorf("failed to parse product id")
	}

	productVoteModel := models.ProductVote{
		SessionID:   sessId,
		ProductID:   productId,
		Liked:       productVote.Liked,
		ProductName: productVote.ProductName,
	}

	if productVote.MachineId != nil {
		machineId, err := uuid.Parse(*productVote.MachineId)
		if err != nil {
			return fmt.Errorf("failed to parse machine id")
		}

		productVoteModel.MachineID = &machineId
	}

	return m.productVoteRepo.UpsertProductVote(&productVoteModel)
}
