package users

import (
	"encoding/json"
	"net/http"
)

/*
RespondJSON writes the given payload as a JSON response with the specified HTTP status code.
It sets the "Content-Type" header to "application/json" and encodes the payload to the response writer.
If encoding fails, it sends an internal server error response.
*/
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

/*
RespondError writes an error message as a JSON response with the specified HTTP status code.
The response body contains a JSON object with an "error" field describing the error.
*/
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}
