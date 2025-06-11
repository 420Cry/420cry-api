package types

// UserResponse represents the response payload containing user authentication and profile information.
type UserResponse struct {
	JWT      string `json:"jwt"`
	UUID     string `json:"uuid"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
