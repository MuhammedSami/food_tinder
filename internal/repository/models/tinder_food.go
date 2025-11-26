package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ProductVote struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null"`
	SessionID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductName string    `gorm:"type:text;not null"`
	Liked       bool      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`

	Session Session `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE"`
}
