package types

// IVerificationResetPasswordForm is the type of Verify Reset Password
type IVerificationResetPasswordForm struct {
	Password string `json:"password"`
	NewPassword string `json:"newPassword"`
}