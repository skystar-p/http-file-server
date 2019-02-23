package main

import (
	// "fmt"
	"net/http"
	// "github.com/gorilla/mux"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if false {
			// perform authentication on here...
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
