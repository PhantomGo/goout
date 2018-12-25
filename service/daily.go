package service

import (
	"goout/model"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/go-redis/redis"
	log "github.com/golang/glog"
)

type DailyService struct {
}

func NewDailyService() *DailyService {
	return &DailyService{}
}

func (svc *DailyService) CheckIn(openId string) (r model.DailyBonus, err error) {
	key := "dailybonus:" + openId
	ckey := "dailybonus:" + openId + ":checked"
	cache := redis.NewClient(options)

	isTodayChecked, _ := cache.Get(ckey).Result()
	if isTodayChecked == "1" {
		err = errors.New("checked")
		return
	}

	var (
		length int64
		bonus  string
	)
	length, err = cache.LLen(key).Result()
	if err != nil || length == 0 {
		cache.LPush(key, "5,50", "5,50", "5,50", "10,50", "20,100", "25,100", "50,200")
	}
	bonus, err = cache.RPop(key).Result()
	now := time.Now()
	today := now.Add(time.Hour * 24)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	cache.Set(ckey, "1", today.Sub(now)).Result()

	next := now.Add(time.Hour * 48)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	cache.ExpireAt(key, next)

	ukey := "user:" + openId
	bonusInt := convertBonus(bonus)
	coin := bonusInt[0]
	exp := bonusInt[1]
	_, err = cache.HIncrBy(ukey, "coin", int64(coin)).Result()
	_, err = cache.HIncrBy(ukey, "exp", int64(exp)).Result()
	log.Info("user:" + openId + ":dailybonus:coin:" + strconv.Itoa(coin) + ":exp:" + strconv.Itoa(exp))
	length = cache.LLen(key).Val()
	days := [7]int{0, 0, 0, 0, 0, 0, 0}
	for index := 0; index < int(7-length); index++ {
		days[index] = 1
	}

	r = model.DailyBonus{
		Daily: days,
		Bonus: model.Bonus{
			Coin: coin,
			EXP:  exp,
		},
	}
	return
}

func convertBonus(bonus string) (r []int) {
	sa := strings.Split(bonus, ",")
	r = make([]int, 2)
	r[0], _ = strconv.Atoi(sa[0])
	r[1], _ = strconv.Atoi(sa[1])
	return
}
