package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handleListRootFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	rootPath := filepath.Join(LocalStoragePath, StorageAccount)
	files, err := os.ReadDir(rootPath)
	if err != nil {
		http.Error(w, "Failed to read storage root", http.StatusInternalServerError)
		log.Printf("Failed to list files in root: %v", err)
		return
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filenames)
	log.Printf("Root listing from: %s", rootPath)
log.Printf("Returned files: %+v", filenames)

}

// Export for use in main.go
var HandleListRootFiles = handleListRootFiles

