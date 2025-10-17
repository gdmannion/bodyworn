// ==============================
// File: main.go
// ==============================
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"BodyWornAPI/server_development_files"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Starting Axis Body Worn API Server on port", port)

	// Initialize logger
	server.SetLogger(&server.DefaultLogger{})

	// Initialize file structure and required objects
	server.CreateRequiredContainersAndObjects()

	// Serve the static index page
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	}))

	// Serve configuration (e.g. StorageAccount) to frontend
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"storageAccount": server.StorageAccount,
		})
	})

	// Authentication endpoint
	http.HandleFunc("/auth/v1.0", server.AuthHandler)

	// Storage + root file listing handler
	http.HandleFunc(fmt.Sprintf("/v1.0/%s/", server.StorageAccount), func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/v1.0/%s/", server.StorageAccount))

		if path == "" && r.Method == http.MethodGet {
			// Handle GET /v1.0/<account>/ — return list of root files
			server.HandleListRootFiles(w, r)
			return
		}

		// Handle standard Swift-style storage operations
		server.StorageHandler(w, r)
	})

	// Start server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

