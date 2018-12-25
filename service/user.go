package service

import (
	"encoding/json"
	"errors"
	"goout/model"
	"strconv"

	"github.com/go-redis/redis"
	log "github.com/golang/glog"
)

type UserService struct {
	levels []model.UserLevel
}

func NewUserService(levels []model.UserLevel) *UserService {
	u := &UserService{}
	u.levels = levels
	return u
}

func (svc *UserService) Create(user model.LoginUser) (err error) {
	ukey := "user:" + user.UserInfo.OpenId
	cache := redis.NewClient(options)
	cache.HSet(ukey, "avatar", user.UserInfo.AvatarUrl)
	cache.HSet(ukey, "nickname", user.UserInfo.NickName)
	cache.HSet(ukey, "coin", 0)
	cache.HSet(ukey, "exp", 0)
	cache.HSet(ukey, "medal", 0)

	ts, _ := json.Marshal(user.Tasks)

	cache.Set("usertask:"+user.UserInfo.OpenId, ts, 0)

	return
}

func (svc *UserService) Get(openId string) (user model.User, err error) {
	ukey := "user:" + openId
	cache := redis.NewClient(options)
	var usrMap map[string]string
	if usrMap, err = cache.HGetAll(ukey).Result(); err != nil {
		log.Error(err)
		return
	}
	user = model.MapUser(usrMap, svc.levels)
	return
}

func (svc *UserService) Consume(openId string, coin int) (left int64, err error) {
	ukey := "user:" + openId
	cache := redis.NewClient(options)
	left, err = cache.HGet(ukey, "coin").Int64()
	if left < int64(coin) {
		err = errors.New("not enough coin")
		return
	}
	left, err = cache.HIncrBy(ukey, "coin", -int64(coin)).Result()
	log.Info("user:" + openId + ":comsumecoin:" + strconv.Itoa(coin))
	return
}

func (svc *UserService) Earn(openId string, coin int) (left int64, err error) {
	ukey := "user:" + openId
	cache := redis.NewClient(options)
	left, err = cache.HIncrBy(ukey, "coin", int64(coin)).Result()
	log.Info("user:" + openId + ":earncoin:" + strconv.Itoa(coin))
	return
}

func (svc *UserService) Ranking() (err error) {
	ukeys := "user:????????????????????????????"
	cache := redis.NewClient(options)
	keys := cache.Keys(ukeys).Val()
	rankKey := "rank"
	err = cache.Del(rankKey).Err()
	for _, key := range keys {
		usrMap := cache.HGetAll(key).Val()
		user := model.MapUser(usrMap, svc.levels)
		data, _ := json.Marshal(user)
		err = cache.ZAdd(rankKey, redis.Z{
			Score:  float64(user.EXP),
			Member: data,
		}).Err()
		if err != nil {
			log.Error(err)
		}
	}
	return
}

func (svc *UserService) RankList() (r []model.User, err error) {
	key := "rank"
	cache := redis.NewClient(options)
	jstr := cache.ZRevRange(key, 0, -1).Val()
	r = make([]model.User, 0, len(jstr))
	for _, s := range jstr {
		var u model.User
		json.Unmarshal([]byte(s), &u)
		r = append(r, u)
	}
	return
}
