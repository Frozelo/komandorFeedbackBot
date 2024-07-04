package entity

type User struct {
	Id       int    `json:"id"`
	TgId     int    `json:"tg_id"`
	Username string `json:"username"`
}
