// Package types contains type definitions for the application. VerificationTokenCheckRequest type will be removed in CRY-55.
package types

// VerificationTokenCheckRequest represents a request payload containing a verification token.
type VerificationTokenCheckRequest struct {
	UserToken   string `json:"userToken"`
	VerifyToken string `json:"verifyToken"`
}
