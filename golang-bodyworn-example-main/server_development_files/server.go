// ==============================
// File: server.go
// ==============================
package server

import (
	"fmt"
	"log"
	"net/http"

)

// AuthHandler validates credentials and returns a token if successful
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	// Extract auth headers
	username := r.Header.Get("X-Auth-User")
	password := r.Header.Get("X-Auth-Key")

	// Validate username and password
	if username != AuthUser || password != AuthPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Send token if authentication is successful
	w.Header().Set("X-Auth-Token", "dummy-token")
	w.Header().Set("X-Storage-Url", fmt.Sprintf("http://%s/v1.0/%s", r.Host, StorageAccount))

	w.WriteHeader(http.StatusOK)
	log.Printf("Function AuthHandler being used to validates and returns token if authenticated successfully")
	log.Printf("Authenticated user=%s from %s — token and storage URL returned", username, r.RemoteAddr)
}
