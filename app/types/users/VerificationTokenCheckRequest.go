package types

type VerificationTokenCheckRequest struct {
	UserToken   string `json:"userToken"`
	VerifyToken string `json:"verifyToken"`
}
