package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHelloRoute(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("Hello, Secure World!")); err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}).Methods("GET")
}
