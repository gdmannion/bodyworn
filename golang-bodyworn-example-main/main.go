// ==============================
// File: main.go
// ==============================
package main

import (
	"fmt"
	"log"
	"net/http"

	"BodyWornAPI/server_development_files"
)

const (
	port = ":8080"
)

func main() {

	log.Println("Starting Axis Body Worn API Server on port", port)

	// Explicitly initialize the logger before any other operations
	server.SetLogger(&server.DefaultLogger{})

	// Create required containers and objects
	server.CreateRequiredContainersAndObjects()

	// Route handlers for authentication and storage
	http.HandleFunc("/auth/v1.0", server.AuthHandler)
	http.HandleFunc(fmt.Sprintf("/v1.0/%s/", server.StorageAccount), server.StorageHandler)



	// Start the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
