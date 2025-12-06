package app

import (
	"fmt"
	"foodtinder/config"
	tinderfoodMgr "foodtinder/internal/business/tinder_food_manager"
	productRepo "foodtinder/internal/repository/product"
	sessionRepo "foodtinder/internal/repository/session"
	"github.com/redis/go-redis/v9"
)

func (a *App) GetSessionRepo() *sessionRepo.Repo {
	return sessionRepo.NewSessionRepo(a.DB)
}

func (a *App) GetProductVoteRepo() *productRepo.Repo {
	return productRepo.NewProductRepo(a.DB)
}

func (a *App) GetTinderFoodManager() *tinderfoodMgr.Manager {
	return tinderfoodMgr.NewManager(
		a.GetSessionRepo(),
		a.GetProductVoteRepo(),
	)
}

func (a *App) GetRedisClient(cfg config.RedisConn) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
}
