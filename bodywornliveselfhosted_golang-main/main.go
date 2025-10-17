package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"bodywornliveselfhosted/auth_bodyworn"
	"bodywornliveselfhosted/subscribe_events"
)

//go:embed static/*
var embeddedStatic embed.FS

var (
	staticFiles       fs.FS
	tmpl              *template.Template
	currentConfig     *auth_bodyworn.FetchConfig
	eventClients      = make(map[*websocket.Conn]bool)
	eventClientsMu    sync.Mutex
	targetID          string
	targetIDMu        sync.RWMutex
	cachedToken       string
	cachedTokenMu     sync.RWMutex
	lastTokenUpdate   time.Time
	tokenExpiresAt    time.Time
	tokenExpiresAtMu  sync.RWMutex
)

func main() {
	var err error

	staticFiles, err = fs.Sub(embeddedStatic, "static")
	if err != nil {
		log.Fatalf("Failed to load static files: %v", err)
	}

	tmpl, err = template.ParseFS(staticFiles, "index.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	config, err := auth_bodyworn.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}
	currentConfig = config

	// Background token refresh loop
	go func() {
		for {
			token, expiresAt, err := auth_bodyworn.FetchToken(currentConfig.IPAddress, currentConfig.Username, currentConfig.Password)
			if err != nil {
				log.Printf(" Failed to refresh token: %v", err)
				time.Sleep(30 * time.Second)
				continue
			}

			cachedTokenMu.Lock()
			cachedToken = token
			lastTokenUpdate = time.Now()
			cachedTokenMu.Unlock()

			tokenExpiresAtMu.Lock()
			tokenExpiresAt = expiresAt
			tokenExpiresAtMu.Unlock()

			log.Printf(" Token refreshed at %s (expires at %s)", lastTokenUpdate.Format(time.RFC3339), expiresAt.Format(time.RFC3339))

			// Calculate refresh delay (refresh 10 seconds before expiry)
			refreshDelay := time.Until(expiresAt.Add(-10 * time.Second))
			if refreshDelay < 10*time.Second {
				refreshDelay = 10 * time.Second // safety fallback
			}

			log.Printf(" Next token refresh in %v", refreshDelay)
			time.Sleep(refreshDelay)
		}
	}()

	// Start event subscription
	go subscribe_events.StartEventSubscriptionWithCallback(func(msg []byte) {
		var evt map[string]interface{}
		if err := json.Unmarshal(msg, &evt); err == nil {
			if subject, ok := evt["subject"].(string); ok {
				targetIDMu.Lock()
				targetID = subject
				targetIDMu.Unlock()
			}
		}
		broadcastEventToClients(msg)
	})

	// HTTP routes
	http.HandleFunc("/events", handleEvents)
	http.HandleFunc("/token", handleToken)
	http.HandleFunc("/api/auth", handleAuth)
	http.HandleFunc("/", serveIndex)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	http.HandleFunc("/ws-proxy", wsProxyHandler)

	log.Println(" Server running at http://0.0.0.0:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	cachedTokenMu.RLock()
	token := cachedToken
	lastUpdate := lastTokenUpdate
	cachedTokenMu.RUnlock()

	tokenExpiresAtMu.RLock()
	expiry := tokenExpiresAt
	tokenExpiresAtMu.RUnlock()

	if token == "" {
		http.Error(w, "Token not initialized yet. Please wait for background refresh.", http.StatusServiceUnavailable)
		return
	}

	targetIDMu.RLock()
	currentTargetID := targetID
	targetIDMu.RUnlock()

	log.Println(" Token endpoint called")
	log.Printf("Returning cached token (refreshed: %s, expires: %s)", lastUpdate.Format(time.RFC3339), expiry.Format(time.RFC3339))

	resp := map[string]string{
		"token":     token,
		"targetId":  currentTargetID,
		"refreshed": lastUpdate.Format(time.RFC3339),
		"expiresAt": expiry.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade /events WebSocket: %v", err)
		return
	}

	eventClientsMu.Lock()
	eventClients[conn] = true
	eventClientsMu.Unlock()
	log.Println("New frontend subscribed to /events")

	go func(c *websocket.Conn) {
		defer func() {
			eventClientsMu.Lock()
			delete(eventClients, c)
			eventClientsMu.Unlock()
			c.Close()
			log.Println(" Frontend disconnected from /events")
		}()
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}(conn)
}

func broadcastEventToClients(msg []byte) {
	eventClientsMu.Lock()
	defer eventClientsMu.Unlock()

	for conn := range eventClients {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Failed to write to client: %v", err)
			conn.Close()
			delete(eventClients, conn)
		}
	}
}

func wsProxyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade client connection: %v", err)
		return
	}
	defer clientConn.Close()

	signalURL := fmt.Sprintf("wss://%s:8082/client?authorization=%s", currentConfig.IPAddress, url.QueryEscape(token))
	log.Printf("Connecting to signaling server: %s", signalURL)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	serverConn, _, err := dialer.Dial(signalURL, nil)
	if err != nil {
		log.Printf("Failed to connect to signaling server: %v", err)
		clientConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"error":"%v"}`, err)))
		return
	}
	defer serverConn.Close()

	log.Println(" Proxy connection established")

	// Client → Server
	go func() {
		for {
			mt, msg, err := clientConn.ReadMessage()
			if err != nil {
				log.Println(" Client disconnected:", err)
				break
			}
			log.Printf("Client → Server: %s", msg)
			serverConn.WriteMessage(mt, msg)
		}
	}()

	// Server → Client
	for {
		mt, msg, err := serverConn.ReadMessage()
		if err != nil {
			log.Println("Signaling server disconnected:", err)
			break
		}
		log.Printf("Server → Client: %s", msg)
		clientConn.WriteMessage(mt, msg)
	}
}
