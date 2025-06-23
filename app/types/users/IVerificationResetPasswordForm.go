package types

// IVerificationResetPasswordForm is the type of Verify Reset Password
type IVerificationResetPasswordForm struct {
	NewPassword        string `json:"newPassword"`
	ResetPasswordToken string `json:"resetPasswordToken"`
}
