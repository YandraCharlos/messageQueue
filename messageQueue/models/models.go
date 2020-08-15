package models

type MessageIn struct {
	Desc string `json:"desc"`
}

type Configure struct {
	IsAuth bool `json:"isAuth"`
}
