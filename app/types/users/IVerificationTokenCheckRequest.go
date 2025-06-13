// Package types contains type definitions for the application. VerificationTokenCheckRequest type will be removed in CRY-55.
package types

// IVerificationTokenCheckRequest represents a request payload containing a verification token.
type IVerificationTokenCheckRequest struct {
	UserToken   string `json:"userToken"`
	VerifyToken string `json:"verifyToken"`
}
