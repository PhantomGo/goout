package model

type Task struct {
	Id        string `json:"id"`
	State     int    `json:"state",omitempty`
	TypeId    int    `json:"typeId",omitempty`
	Avatar    string `json:"avatar",omitempty`
	CreatorId string `json:"creatorId",omitempty`
	NickName  string `json:"nickName",omitempty`
	Title     string `json:"title",omitempty`
}

type TaskDetail struct {
	Task    interface{} `json:"task"`
	Members Member      `json:"member"`
}

type SubmitTask struct {
	OpenId string `json:"token"`
	TaskId string `json:"taskId"`
	Coin   int    `json:"coin"`
	EXP    int    `json:"exp"`
	Medal  int    `json:"medal"`
	State  int    `json:"state"`
}

type Member struct {
	Count  int      `json:"count"`
	Avatar []string `json:"avatar"`
}
