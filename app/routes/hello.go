package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHelloRoute(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Secure World!"))
	}).Methods("GET")
}
