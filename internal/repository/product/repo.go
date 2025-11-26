package session

import (
	"foodjiassignment/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) UpsertProductVote(vote *models.ProductVote) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "session_id"},
			{Name: "product_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"liked":        vote.Liked,
			"product_name": vote.ProductName,
			"machine_id":   vote.MachineID,
			"updated_at":   time.Now(),
		}),
	}).Create(vote).Error
}
