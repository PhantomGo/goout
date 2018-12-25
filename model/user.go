package model

import "strconv"

type User struct {
	OpenId    string `json:"openId"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
	NickName  string `json:"nickName",omitempty`
	Coin      int    `json:"coin",omitempty`
	EXP       int    `json:"exp",omitempty`
	Level     int    `json:"level",omitempty`
	Medal     int    `json:"medal",omitempty`
}

type LoginUser struct {
	UserInfo User   `json:"user"`
	Tasks    []Task `json:"tasks"`
}

type Avatar struct {
	Avatar string `json:avatar`
}

type UserLevel struct {
	Level int `json:"level"`
	Exp   int `json:"exp"`
}

func MapUser(usrMap map[string]string, levels []UserLevel) (user User) {
	coin, _ := strconv.Atoi(usrMap["coin"])
	exp, _ := strconv.Atoi(usrMap["exp"])
	medal, _ := strconv.Atoi(usrMap["medal"])
	level := 1
	for index := len(levels) - 1; index >= 0; index-- {
		if exp >= levels[index].Exp {
			level = levels[index].Level
			break
		}
	}
	user = User{
		OpenId:    "",
		AvatarUrl: usrMap["avatar"],
		NickName:  usrMap["nickname"],
		Coin:      coin,
		Level:     level,
		Medal:     medal,
		EXP:       exp,
	}
	return
}
