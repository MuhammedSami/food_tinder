package tinder_food

import (
	"gorm.io/gorm"
)

type TinderFoodRepo struct {
	db *gorm.DB
}

func NewTinderFoodRepo(db *gorm.DB) *TinderFoodRepo {
	return &TinderFoodRepo{
		db: db,
	}
}
