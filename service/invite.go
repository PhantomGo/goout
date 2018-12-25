package service

import (
	"goout/model"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

var (
	characters = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
)

type InviteService struct {
	CodePool     []string
	currentIndex int
}

func NewInviteService() *InviteService {
	svc := &InviteService{
		CodePool:     make([]string, 100),
		currentIndex: 0,
	}
	svc.generateCode()
	return svc
}

func (svc *InviteService) generateCode() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := 0; index < 100; index++ {
		code := ""
		for j := 0; j < 4; j++ {
			num := r.Intn(10)
			code += characters[num]
		}
		//查重
		exist := false
		for _, v := range svc.CodePool {
			if v == code {
				exist = true
				break
			}
		}
		if !exist {
			svc.CodePool[index] = code
		}
	}
}

func (svc *InviteService) GetCode(invite model.Invite) (code string, err error) {
	if svc.currentIndex == 98 {
		svc.generateCode()
		svc.currentIndex = 0
	}
	code = svc.CodePool[svc.currentIndex]
	svc.currentIndex++
	key := "invite:" + code
	cache := redis.NewClient(options)
	cache.HSet(key, "openId", invite.CreatorOpenId)
	cache.HSet(key, "taskId", invite.TaskId)
	cache.HSet(key, "typeId", invite.TypeId)
	cache.HSet(key, "avatar", invite.Avatar)
	cache.HSet(key, "title", invite.Title)
	cache.HSet(key, "nickname", invite.NickName)
	cache.Expire(key, time.Minute*30)
	return
}

func (svc *InviteService) VerifyCode(code string) (task model.Task, err error) {
	key := "invite:" + code
	cache := redis.NewClient(options)
	exists := cache.Exists(key).Val()
	if exists > 0 {
		taskId := cache.HGet(key, "taskId").Val()
		creatorId := cache.HGet(key, "openId").Val()
		avatar := cache.HGet(key, "avatar").Val()
		typeId, _ := cache.HGet(key, "typeId").Int()
		title := cache.HGet(key, "title").Val()
		name := cache.HGet(key, "nickname").Val()
		task = model.Task{
			Id:        taskId,
			CreatorId: creatorId,
			Avatar:    avatar,
			TypeId:    typeId,
			NickName:  name,
			Title:     title,
		}
	}

	return
}
