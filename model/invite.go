package model

type Invite struct {
	CreatorOpenId string `json:"token"`
	TaskId        string `json:"taskId"`
	TypeId        int    `json:"typeId"`
	Avatar        string `json:"avatar"`
	NickName      string `json:"nickName"`
	Title         string `json:"title"`
}
