package model

type Home struct {
	UserInfo User   `json:"user"`
	Tasks    []Task `json:"tasks"`
}

type Token struct {
	OpenId string `json:"token"`
	TaskId string `json:"taskId",omitempty`
	TypeId int    `json:"typeId",omitempty`
	Avatar string `json:"avatar",omitempty`
	Coin   int    `json:"coin",omitempty`
}
