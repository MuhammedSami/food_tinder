package tinderfood

import (
	"foodjiassignment/internal/repository/models"
	"time"
)

type repo interface {
	CreateSession(expiresAt *time.Time) (*models.Session, error)
}
