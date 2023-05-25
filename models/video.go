package models

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

type Video struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Size     int    `json:"size"`
}
