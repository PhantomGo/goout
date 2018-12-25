package service

import (
	"encoding/json"
	"goout/model"
	"strconv"

	"github.com/go-redis/redis"
	log "github.com/golang/glog"
)

type TaskService struct {
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (svc *TaskService) MyTasks(openId string) (task []model.Task, err error) {
	key := "usertask:" + openId
	cache := redis.NewClient(options)
	taskData := cache.Get(key)
	json.Unmarshal([]byte(taskData.Val()), &task)
	for i, t := range task {
		tkey := "user:" + openId + ":task:" + t.Id + ":state"
		v := cache.Exists(tkey).Val()
		if v == 1 {
			t.State, _ = cache.Get(tkey).Int()
		} else {
			t.State = int(v)
		}
		task[i] = t
	}
	return
}

func (svc *TaskService) Submit(openId, taskId string, coin, exp, medal int64, state int) (user model.User, err error) {
	key := "user:" + openId + ":task:" + taskId + ":state"
	cache := redis.NewClient(options)
	v := cache.Get(key).Val()

	cache.Set(key, state, 0)
	ukey := "user:" + openId
	user = model.User{
		OpenId: openId,
	}
	if state == 1 && v != "1" {
		user.Coin = int(cache.HIncrBy(ukey, "coin", coin).Val())
		user.EXP = int(cache.HIncrBy(ukey, "exp", exp).Val())
		user.Medal = int(cache.HIncrBy(ukey, "medal", medal).Val())
		log.Info("user:" + openId + ":submitTask:" + taskId + ":increase:coin:" + strconv.FormatInt(coin, 10) + ":exp:" + strconv.FormatInt(exp, 10))
	}

	return
}

func (svc *TaskService) Join(index int, avatar string) (err error) {
	key := "tasks:" + strconv.Itoa(index)
	cache := redis.NewClient(options)
	if _, err = cache.SAdd(key, avatar).Result(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (svc *TaskService) Detail(index int) (task model.TaskDetail, err error) {
	key := "tasks"
	tkey := "tasks:" + strconv.Itoa(index)
	cache := redis.NewClient(options)
	var taskStr string
	if taskStr, err = cache.LIndex(key, int64(index-1)).Result(); err != nil {
		return
	}
	var detail interface{}
	json.Unmarshal([]byte(taskStr), &detail)
	count := cache.SCard(tkey).Val()
	var avatars []string
	if count > 0 {
		avatars = cache.SRandMemberN(tkey, 4).Val()
	}
	task = model.TaskDetail{
		Task: detail,
	}
	task.Members = model.Member{
		Count:  int(count),
		Avatar: avatars,
	}
	return
}

func (svc *TaskService) FillTasks(tasks []interface{}) (err error) {
	key := "tasks"
	cache := redis.NewClient(options)
	for _, item := range tasks {
		data, _ := json.Marshal(item)
		cache.RPush(key, data)
	}
	err = nil
	return
}
