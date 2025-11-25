package app

import (
	tinderFoodRepo "foodjiassignment/internal/repository/tinder_food"
)

func (a *APP) GetTinderFoodRepo() *tinderFoodRepo.TinderFoodRepo {
	return tinderFoodRepo.NewTinderFoodRepo(a.DB)
}
