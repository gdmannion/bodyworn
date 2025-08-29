package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// logMetadata prints all X-Container-Meta and X-Object-Meta headers
func logMetadata(r *http.Request) {
	log.Println("---- Metadata Details ----")
	log.Printf("Function logMetadata prints all X-Container-Meta and X-Object-Meta headers ")
	for k, v := range r.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-container-meta-") || strings.HasPrefix(strings.ToLower(k), "x-object-meta-") {
			log.Printf("%s: %s", k, v)
		}
	}
	log.Println("--------------------------")
}



// handlePostMetadata updates metadata for a user or device
func handlePostMetadata(w http.ResponseWriter, r *http.Request, path string) {
	metaPath := filepath.Join(LocalStoragePath, StorageAccount, path+".meta")

	// Check if the parent object exists before updating metadata
	objPath := filepath.Join(LocalStoragePath, StorageAccount, path)
	log.Printf("Function handlePostMetadata updates metadata for a user or device ")
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		http.Error(w, "Object not found", http.StatusNotFound)
		log.Printf("Attempted to update metadata for non-existent object %s", path)
		return
	}

	// Parse metadata from request headers
	metadata := parseMetadata(r)

	// Log metadata headers
	logMetadata(r)

	// Convert metadata to JSON and save to .meta file
	metaContent, err := json.MarshalIndent(metadata, "", "  ")
	log.Printf("Function handlePostMetadata updates metadata for a user or device ")
	if err != nil {
		http.Error(w, "Failed to marshal metadata", http.StatusInternalServerError)
		log.Printf("Failed to marshal metadata for %s: %v", path, err)
		return
	}

	err = os.WriteFile(metaPath, metaContent, 0644)
	if err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		log.Printf("Failed to update metadata for %s: %v", path, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	log.Printf("Function handlePostMetadata being used to updates metadata for a user or device")
	log.Printf("Metadata for %s updated successfully", path)
}

// parseMetadata extracts metadata from HTTP headers
func parseMetadata(r *http.Request) map[string]string {
	metadata := make(map[string]string)

	for k, v := range r.Header {
		lowerKey := strings.ToLower(k)

		if strings.Contains(lowerKey, "-meta-") {
			// Split key like "x-object-meta-name" or "x-container-meta-model"
			parts := strings.SplitN(lowerKey, "-meta-", 2)
			if len(parts) == 2 {
				metaKey := parts[1]
				metadata[metaKey] = v[0] // Take first value (you can join all if needed)
			}
		}
	}
	log.Printf("Function parseMetadata capture all -meta- ")
	log.Printf("Captured metadata: %+v", metadata)
	return metadata
}





// HEAD handler with metadata
func handleHeadRequest(w http.ResponseWriter, r *http.Request, path string) {
	fullPath := filepath.Join(LocalStoragePath, StorageAccount, path)
	metaPath := fullPath + ".meta"

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Printf("HEAD: %s not found", fullPath)
		return
	}

	// Return object or container headers
	if info.IsDir() {
		w.Header().Set("X-Container-Object-Count", countObjects(fullPath))
		w.Header().Set("X-Container-Bytes-Used", calculateSize(fullPath))
		addMetadataHeaders(w, metaPath, "X-Container-Meta-")
		w.WriteHeader(http.StatusNoContent)
		log.Printf("Function handleHeadRequest with metadata")
		log.Printf("HEAD: Container metadata returned for %s", fullPath)
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
		w.Header().Set("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))
		w.Header().Set("ETag", generateETag(fullPath))
		w.Header().Set("Content-Type", "application/octet-stream")
		addMetadataHeaders(w, metaPath, "X-Object-Meta-")
		w.WriteHeader(http.StatusOK)
		log.Printf("Function handleHeadRequest with metadata")
		log.Printf("HEAD: Object metadata returned for %s", fullPath)
	}
}



// Add metadata headers to response
func addMetadataHeaders(w http.ResponseWriter, metaPath, prefix string) {
	file, err := os.Open(metaPath)
	if err != nil {
		return
	}
	defer file.Close()

	var meta map[string]string
	if err := json.NewDecoder(file).Decode(&meta); err != nil {
		return
	}

	// Create a case converter for title case
	title := cases.Title(language.English)

	for k, v := range meta {
		// Use the Title case converter from x/text/cases
		headerKey := prefix + title.String(k)
		w.Header().Set(headerKey, v)
		log.Printf("Function addMetadataHeaders add metadata headers to response")
	}
}

// Helpers
func generateETag(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	hash := md5.New()
	io.Copy(hash, f)
	log.Printf("Function generateETag being used")
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func countObjects(path string) string {
	files, _ := os.ReadDir(path)
	log.Printf("Function countObjects being used")
	return fmt.Sprintf("%d", len(files))
}

func calculateSize(path string) string {
	var total int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	log.Printf("Function calculateSize being used")
	return fmt.Sprintf("%d", total)
}



// StorageHandler handles storage requests for GET, PUT, POST, and HEAD operations
func StorageHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("/v1.0/%s/", StorageAccount)
	path := strings.TrimPrefix(r.URL.Path, prefix)

	switch r.Method {
	case http.MethodPut:
		putObject(w, r, path)
	case http.MethodGet:
		getObject(w, path)
	case http.MethodPost:
		handlePostMetadata(w, r, path)
	case http.MethodHead:
		handleHeadRequest(w, r, path)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		log.Printf("Unsupported method %s for path %s", r.Method, path)
	}
}
