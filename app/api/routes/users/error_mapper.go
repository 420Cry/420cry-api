package users

import "net/http"

// mapErrorToStatusCode maps the error message to an HTTP status code
func mapUserCreationErrorToStatusCode(errMessage string) int {
	switch errMessage {
	case "username is already taken":
		return http.StatusConflict
	case "email is already taken":
		return http.StatusConflict
	case "failed to generate signup token":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
