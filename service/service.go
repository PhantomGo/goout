package service

import (
	"context"
	"goout/conf"
	"goout/lib/http"
	"goout/model"

	"github.com/go-redis/redis"
	log "github.com/golang/glog"
)

var (
	cache      *redis.Client
	httpClient *http.Client
	options    = &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
)

type Service struct {
	TokenSvc  *TokenService
	UserSvc   *UserService
	TaskSvc   *TaskService
	InviteSvc *InviteService
	DailySvc  *DailyService
}

func New(c *conf.Config) *Service {
	httpClient = http.NewClient(c.HTTPClient)
	levels := make([]model.UserLevel, 0, 10)
	httpClient.Get(context.TODO(), "http://localhost:3000/levels", "127.0.0.1", nil, &levels)
	log.Error(levels)
	return &Service{
		TokenSvc:  NewTokenService(),
		UserSvc:   NewUserService(levels),
		TaskSvc:   NewTaskService(),
		InviteSvc: NewInviteService(),
		DailySvc:  NewDailyService(),
	}
}
