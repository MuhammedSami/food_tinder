package app

import (
	tinderfoodMgr "foodjiassignment/internal/business/tinder_food_manager"
	tinderFoodRepo "foodjiassignment/internal/repository/tinder_food"
)

func (a *App) GetTinderFoodRepo() *tinderFoodRepo.TinderFoodRepo {
	return tinderFoodRepo.NewTinderFoodRepo(a.DB)
}

func (a *App) GetTinderFoodManager() *tinderfoodMgr.Manager {
	return tinderfoodMgr.NewManager(a.GetTinderFoodRepo())
}
