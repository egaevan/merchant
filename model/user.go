package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Gender   string `json:"Gender"`
	Password string `json:"Password"`
}
