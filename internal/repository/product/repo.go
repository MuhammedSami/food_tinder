package session

import (
	"fmt"
	"foodtinder/internal/repository/models"
	"github.com/google/uuid"
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

func (r *Repo) GetVotesBySessionId(sessionID uuid.UUID) ([]models.ProductVote, error) {
	var votes []models.ProductVote

	err := r.db.
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Find(&votes).Error
	if err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *Repo) GetAverageScores() ([]models.ProductScore, error) {
	var results []models.ProductScore

	err := r.db.Model(&models.ProductVote{}).
		Select(`
			product_id,
			product_name,
			AVG(CASE WHEN liked THEN 1 ELSE 0 END) AS avg_score,
			COUNT(*) AS total_votes,
			SUM(CASE WHEN liked THEN 1 ELSE 0 END) AS likes
		`).
		Group("product_id, product_name").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get average scores: %w", err)
	}

	return results, nil
}
