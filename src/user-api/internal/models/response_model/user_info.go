package response_model

type UserInfo struct {
	Id    string `dynamodbav:"Id" json:"user_id"`
	Name  string `dynamodbav:"Name" json:"name"`
	Email string `dynamodabav:"Email" json:"email"`
}
