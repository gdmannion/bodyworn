package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"os"
	"path/filepath"
)


//In OpenStack Swift, the object storage system, the hierarchy is designed to organize and manage large volumes of unstructured data. 
// The basic structure in Swift follows a three-tier system: Account, Container, and Object.

//Please NOTE, this application is using the concept of OpenStack but it is not an OpenStack application
//It is using the concept of OpenStack but really is an OS File System application

const (
	LocalStoragePath   = "./"//root directory of appliaiton
	StorageAccount   = "WhateverStorageName"//this can be whatevery name you want.  In comparison this could be considered  your account in OpenStack Swift
	AuthPassword       = "WhateverPassWord"
	AuthUser           = "WhateverUserName"
	ConnectionFile     = "connection.json"
)


// CreateRequiredContainersAndObjects ensures required folders and objects are created in OS file system.
//In comparison the System, Users, Devices would be your Containers in OpenStack
func CreateRequiredContainersAndObjects() {
	log.Printf("Function CreateRequiredContainersAndObjects ensures required folders and objects are created in OS file system")
	createDirIfNotExists(filepath.Join(LocalStoragePath, StorageAccount))
	createDirIfNotExists(filepath.Join(LocalStoragePath, StorageAccount, "System"))
	createDirIfNotExists(filepath.Join(LocalStoragePath, StorageAccount, "Users"))
	createDirIfNotExists(filepath.Join(LocalStoragePath, StorageAccount, "Devices"))

	

	createLocalCapabilitiesFile()
	createLocalConnectionFile()
}

// createDirIfNotExists checks and creates a directory if it doesn't exist
func createDirIfNotExists(path string) {
	log.Printf("Function createDirIfNotExists being used to checks and creates a directory if it doesn't exist")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", path, err)
		}
		log.Printf("Directory %s created successfully", path)
	}
}


// writeFile writes content to a file
func writeFile(path string, content []byte) {
	// Check if the path is an existing directory
	log.Printf("Function writeFile being used to write content to a file")
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
		log.Fatalf("Failed to create file %s: Path is a directory", path)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", path, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", path, err)
	}
	log.Printf("File %s created successfully", path)
}

// GET handler with Swift-style headers
func getObject(w http.ResponseWriter, path string) {
	fullPath := filepath.Join(LocalStoragePath, StorageAccount, path)
	log.Printf("Function getObject being used to chandler with Swift-style headers")
	metaPath := fullPath + ".meta"

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "Object not found", http.StatusNotFound)
		log.Printf("GET: Object %s not found", path)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.Header().Set("Last-Modified", info.ModTime().UTC().Format(http.TimeFormat))
	w.Header().Set("ETag", generateETag(fullPath))
	w.Header().Set("Content-Type", "application/octet-stream")
	addMetadataHeaders(w, metaPath, "X-Object-Meta-")

	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		log.Printf("GET: Failed to open object %s", path)
		return
	}
	defer file.Close()

	io.Copy(w, file)
	log.Printf("GET: Object %s returned with headers", path)
}

// putObject stores a file or metadata in local storage
func putObject(w http.ResponseWriter, r *http.Request, path string) {
	filePath := ""
	metaPath := filepath.Join(LocalStoragePath, StorageAccount, path+".meta")
	log.Printf("Function putObject is being used stores a file or metadata in local storage")

	if path == "Users" || path == "Devices" || path == "System" {
		// Handle directories like Users and Devices correctly
		filePath = filepath.Join(LocalStoragePath, StorageAccount, path)
		createDirIfNotExists(filePath)
		w.WriteHeader(http.StatusCreated)
		log.Printf("Container %s created successfully", path)
		return
	} else if strings.HasSuffix(path, ".mkv") {
		// Store .mkv files in the root directory
		filePath = filepath.Join(LocalStoragePath, StorageAccount, filepath.Base(path))
	} else {
		// Handle all other files normally
		filePath = filepath.Join(LocalStoragePath, StorageAccount, path)
	}

	dirPath := filepath.Dir(filePath)

	// Check if parent path is a valid directory before adding objects
	parentInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create directory %s", dirPath), http.StatusInternalServerError)
			log.Printf("Failed to create parent directory %s: %v", dirPath, err)
			return
		}
		log.Printf("Parent directory %s created successfully", dirPath)
	} else if err == nil && !parentInfo.IsDir() {
		http.Error(w, fmt.Sprintf("Parent path %s is not a directory", dirPath), http.StatusInternalServerError)
		log.Printf("Parent path %s is not a directory", dirPath)
		return
	}

	// Create or overwrite the object file
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		log.Printf("Failed to create file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Failed to upload object", http.StatusInternalServerError)
		log.Printf("Failed to upload object %s: %v", filePath, err)
		return
	}

	// Log metadata headers
	logMetadata(r)

// Create metadata after storing objects in Users/, Devices/, or System/
if strings.HasPrefix(path, "Users/") || strings.HasPrefix(path, "Devices/") || strings.HasPrefix(path, "System/") {
	metadata := parseMetadata(r)
	if len(metadata) > 0 {
		metaContent, err := json.MarshalIndent(metadata, "", "  ")
		if err == nil {
			err = os.WriteFile(metaPath, metaContent, 0644)
			if err == nil {
				w.WriteHeader(http.StatusCreated)
				log.Printf("Metadata for %s created successfully", path)
				log.Printf("Response: %d Created", http.StatusCreated)
			} else {
				http.Error(w, "Failed to write metadata", http.StatusInternalServerError)
				log.Printf("Failed to create metadata for %s: %v", path, err)
				log.Printf("Response: %d Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Failed to marshal metadata", http.StatusInternalServerError)
			log.Printf("Failed to marshal metadata for %s: %v", path, err)
			log.Printf("Response: %d Internal Server Error", http.StatusInternalServerError)
		}
	}
}

    log.Printf("Function putObject stores a file or metadata in local storage")
	w.WriteHeader(http.StatusCreated)
	log.Printf("Object %s uploaded successfully", path)
}