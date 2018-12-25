package http

import (
	"goout/errors"
	"goout/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createUser(c *gin.Context) {
	var user model.LoginUser
	if err := c.BindJSON(&user); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	t, _ := svc.TokenSvc.Generate(user.UserInfo.OpenId)
	svc.UserSvc.Create(user)
	r := struct {
		token string
	}{
		token: t,
	}
	result(c, &r, nil)
}

func myHome(c *gin.Context) {
	var u model.Token
	if err := c.BindJSON(&u); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}

	usr, _ := svc.UserSvc.Get(u.OpenId)
	tsks, _ := svc.TaskSvc.MyTasks(u.OpenId)
	home := &model.Home{
		UserInfo: usr,
		Tasks:    tsks,
	}
	result(c, home, nil)
}

func joinTask(c *gin.Context) {
	index, err := strconv.Atoi(c.Param("index"))
	var avatar model.Avatar
	if err := c.BindJSON(&avatar); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	svc.TaskSvc.Join(index, avatar.Avatar)
	r := struct {
		OK string `json:"ok"`
	}{
		OK: "ok",
	}
	result(c, &r, err)
}

func taskInfo(c *gin.Context) {
	index, err := strconv.Atoi(c.Param("index"))
	r, _ := svc.TaskSvc.Detail(index)
	result(c, &r, err)
}

func submit(c *gin.Context) {
	var r model.SubmitTask
	if err := c.BindJSON(&r); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	u, _ := svc.TaskSvc.Submit(r.OpenId, r.TaskId, int64(r.Coin), int64(r.EXP), int64(r.Medal), r.State)
	result(c, &u, nil)
}

func joinInvite(c *gin.Context) {
	code := c.Param("code")
	if code == "" || len(code) != 4 {
		result(c, nil, errors.ParamsErr)
		return
	}
	r, _ := svc.InviteSvc.VerifyCode(code)
	result(c, &r, nil)
}

func invite(c *gin.Context) {
	var i model.Invite
	if err := c.BindJSON(&i); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	code, _ := svc.InviteSvc.GetCode(i)
	r := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}
	result(c, &r, nil)
}

func dailybonus(c *gin.Context) {
	var u model.Token
	if err := c.BindJSON(&u); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	r, err := svc.DailySvc.CheckIn(u.OpenId)
	if err != nil {
		result(c, nil, errors.NotModified)
		return
	}
	result(c, &r, nil)
}

func fillTasks(c *gin.Context) {
	var tasks []interface{}
	if err := c.BindJSON(&tasks); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	svc.TaskSvc.FillTasks(tasks)
	r := struct {
		OK string `json:"ok"`
	}{
		OK: "ok",
	}
	result(c, &r, nil)
}

func consumeCoin(c *gin.Context) {
	var u model.Token
	if err := c.BindJSON(&u); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	var (
		coin int64
		err  error
	)
	if coin, err = svc.UserSvc.Consume(u.OpenId, u.Coin); err != nil {
		result(c, nil, errors.Conflict)
		return
	}
	r := struct {
		Coin int `json:"coin"`
	}{
		Coin: int(coin),
	}
	result(c, &r, nil)
}

func earnCoin(c *gin.Context) {
	var u model.Token
	if err := c.BindJSON(&u); err != nil {
		result(c, nil, errors.ParamsErr)
		return
	}
	var (
		coin int64
		err  error
	)
	if coin, err = svc.UserSvc.Earn(u.OpenId, u.Coin); err != nil {
		result(c, nil, errors.ServerErr)
		return
	}
	r := struct {
		Coin int `json:"coin"`
	}{
		Coin: int(coin),
	}
	result(c, &r, nil)
}

func ranking(c *gin.Context) {
	if err := svc.UserSvc.Ranking(); err != nil {
		result(c, nil, errors.ServerErr)
		return
	}
	r := struct {
		OK string `json:"ok"`
	}{
		OK: "ok",
	}
	result(c, &r, nil)
}

func rankList(c *gin.Context) {
	var r []model.User
	r, _ = svc.UserSvc.RankList()
	result(c, &r, nil)
}

func pong(c *gin.Context) {
	r := struct {
		Pong string `json:"pong"`
	}{
		Pong: "pong",
	}
	result(c, &r, nil)
}
