package model

type DailyBonus struct {
	Daily   [7]int `json:"daily"`
	Bonus   Bonus  `json:"bonus"`
	Checked bool   `json:"checked"`
}

type Bonus struct {
	Coin int `json:"coin"`
	EXP  int `json:"exp"`
}
