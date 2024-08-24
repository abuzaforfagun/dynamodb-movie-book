package request_model

type AddUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
