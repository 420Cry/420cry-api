package types

type UserResponse struct {
	JWT      string `json:"jwt"`
	UUID     string `json:"uuid"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
