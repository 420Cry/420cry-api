package types

// IVerificationResetPasswordForm is the type of Verify Reset Password
type IVerificationResetPasswordForm struct {
	NewPassword        string `json:"newPassword"`
	RepeatedPassword   string `json:"repeatedPassword"`
	ResetPasswordToken string `json:"resetPasswordToken"`
}
