package server

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"
	
)

const (
	LocalConnectionFilePath = "./connection.json"
)

//this function is used to automatically assign local IP address that will be used in connection file
func getServerIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error fetching interfaces: %v", err)
	}

	for _, iface := range interfaces {
		// Only considering interfaces that are up and not a loopback
		if iface.Flags&net.FlagUp != 0 && iface.Name != "lo" {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Fatalf("Error getting addresses for interface %s: %v", iface.Name, err)
			}
			for _, addr := range addrs {
				// Try to find an IPv4 address
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					log.Printf("Function getServerIP providing IP address for connection file and authentication requests")
					return ipnet.IP.String()
				}
			}
		}
	}
	log.Fatal("No valid non-loopback interface found")
	
	return ""
	
}





// getCapabilitiesJSON returns Capabilities.json content as a JSON byte slice
func getCapabilitiesJSON() []byte {
	log.Printf("Function getCapabilitiesJSON returing Capabilities.json content as a JSON byte slice")
	
	capabilities := map[string]interface{}{
		"StoreAndRead": map[string]bool{
			"StoreReadSystemID":       true,
			"StoreUserIDKey":      	   true,
			"StoreBookmarks":          true,
			"StoreSignedVideo":        true,
			"StoreGNSSTrackRecording": true,
			"StoreRejectedContent": true,
		},
	}
	file, _ := json.MarshalIndent(capabilities, "", "  ")
	return file
}

// getConnectionJSON returns connection.json content as a JSON byte slice
func getConnectionJSON() []byte {
	logger.Infof("Generating content for connection.json...")

	// Get the server IP dynamically
	serverIP := getServerIP()

	// Update the AuthenticationTokenURI to use the dynamically fetched IP address
	connection := map[string]interface{}{
		"ConnectionFileVersion":        "1.0",
		"SiteName":                     "Axis Body Worn",
		"ApplicationName":              "BodyWornAPI",
		"ApplicationVersion":           "1.0",
		"AuthenticationTokenURI":       []string{"http://" + serverIP + ":8080/auth/v1.0"},
		"BlobAPIKey":                   AuthPassword,
		"BlobAPIUserName":              AuthUser,
		"ContainerType":                "mkv",
		"WantEncryption":      false,
		"PublicKey":      "",
		"PublicKeyId":      "",
		"FullStoreAndReadSupport":      true,  
	}

	//Note FullStoreAndReadSupport cannot be true unless using HEAD and GET configure in applcaiton. If you try to load a connection without you will recieve an error
	// FullStoreAndReadSupport if true will allow you to receive the meta for your system folder providing connectionID and name of the W800 that loaded the connection 
	// file to connect to content destination but if you don't have HEAD support in your application you will receive an error




	// Marshal the connection map to JSON
	file, _ := json.MarshalIndent(connection, "", "  ")
	return file
}


// createLocalCapabilitiesFile creates Capabilities.json in the local filesystem
func createLocalCapabilitiesFile() {
	log.Printf("Function createLocalCapabilitiesFile creatingCapabilities.json in the local filesystem ")
	content := getCapabilitiesJSON()
	capabilitiesPath := filepath.Join(LocalStoragePath, StorageAccount, "System", "Capabilities.json")
	writeFile(capabilitiesPath, content)
}

// createLocalConnectionFile creates connection.json in the root directory
func createLocalConnectionFile() {
	content := getConnectionJSON()

	// Dynamically get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get current directory: %v", err)
		return
	}

	connectionFilePath := filepath.Join(rootDir, "connection.json")

	logger.Infof("Creating local %s...", ConnectionFile)
	err = os.WriteFile(connectionFilePath, content, 0644)
	if err != nil {
		logger.Errorf("Failed to create local %s: %v", ConnectionFile, err)
		return
	}
	logger.Infof("Local %s created successfully at %s", ConnectionFile, connectionFilePath)
}