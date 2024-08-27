package dto

type UserInfo struct {
	Id    string `json:"user_id"`
	Name  string `json:"name"`
	Email string `email:"email"`
}
