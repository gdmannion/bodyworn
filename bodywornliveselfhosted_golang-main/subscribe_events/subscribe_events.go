package subscribe_events

import (
	"bodywornliveselfhosted/digest_auth"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/gorilla/websocket"
)

type eventMessage struct {
	APIVersion string `json:"apiVersion"`
	Method     string `json:"method"`
	Params     struct {
		EventFilterList []struct {
			TopicFilter string `json:"topicFilter"`
		} `json:"eventFilterList"`
	} `json:"params"`
}

type incomingEvent struct {
	Params struct {
		Notification struct {
			Message struct {
				Data struct {
					EventType string `json:"eventType"`
					Event     string `json:"event"`
				} `json:"data"`
			} `json:"message"`
		} `json:"notification"`
	} `json:"params"`
}

type Config struct {
	IPAddress string `json:"ip_address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type RecordingEntry struct {
	SessionID   string `json:"sessionId"`
	RecordingID string `json:"recordingId"`
	StartTime   string `json:"startTime"`
}



func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func StartEventSubscriptionWithCallback(callback func([]byte)) {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}

	retryDelay := 5 * time.Second

	for {
		sid, err := getSessionID(cfg.IPAddress, cfg.Username, cfg.Password)
		if err != nil {
			log.Printf("Failed to get session ID: %v", err)
			time.Sleep(retryDelay)
			continue
		}

		wsURL := fmt.Sprintf("wss://%s/vapix/ws-data-stream?wssession=%s&sources=events", cfg.IPAddress, sid)
		dialer := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

		log.Printf("Connecting to %s", wsURL)
		conn, resp, err := dialer.Dial(wsURL, nil)
		if err != nil {
			log.Printf("WebSocket dial failed: %v", err)
			if resp != nil {
				log.Printf("HTTP Status: %s", resp.Status)
			}
			time.Sleep(retryDelay)
			continue
		}
		log.Printf("Connected to Axis event stream")

		// Subscribe
		subscribePayload := eventMessage{
			APIVersion: "1.0",
			Method:     "events:configure",
		}
		subscribePayload.Params.EventFilterList = append(
			subscribePayload.Params.EventFilterList,
			struct {
				TopicFilter string `json:"topicFilter"`
			}{TopicFilter: "tns1:WebRTC/tnsaxis:Signaling/CloudEvent"},
		)

		msgBytes, _ := json.Marshal(subscribePayload)
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Failed to send subscription message: %v", err)
			conn.Close()
			time.Sleep(retryDelay)
			continue
		}

		// Listen
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				conn.Close()
				break
			}

			var incoming incomingEvent
			if err := json.Unmarshal(message, &incoming); err != nil {
				log.Printf("Failed to parse event: %v", err)
				continue
			}

			eventType := incoming.Params.Notification.Message.Data.EventType
			rawEvent := incoming.Params.Notification.Message.Data.Event
			log.Printf("Event type: %s", eventType)

			if strings.TrimSpace(rawEvent) == "" {
				log.Println("Empty event payload, skipping")
				continue
			}

			var detailed map[string]interface{}
			if err := json.Unmarshal([]byte(rawEvent), &detailed); err != nil {
				log.Printf("Failed to parse event body: %v", err)
				continue
			}

			if eventType == "com.axis.bodyworn.stream.started" {
				dataMap, ok := detailed["data"].(map[string]interface{})
				if !ok {
					log.Println("Invalid event data format")
					continue
				}

				entry := RecordingEntry{
					SessionID:   asString(dataMap["sessionId"]),
					RecordingID: asString(dataMap["recordingId"]),
					StartTime:   asString(detailed["time"]),
				}
				eventToSend := map[string]interface{}{
					"type":       eventType,
					"subject":    detailed["subject"],
					"sessionId":  entry.SessionID,
					"recordingId": entry.RecordingID,
					"bearerId":   dataMap["bearerId"],
					"bearerName": dataMap["bearerName"],
					"time":       entry.StartTime,
				}

				encoded, err := json.Marshal(eventToSend)
				if err != nil {
					log.Printf("Failed to marshal event: %v", err)
					continue
				}

				callback(encoded)
				log.Printf("Forwarded event to frontend: %s", encoded)
			}
		}
	}
}

func asString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}



func getSessionID(ip, username, password string) (string, error) {
	urlStr := fmt.Sprintf("http://%s/axis-cgi/wssession.cgi", ip)
	authState, err := digest_auth.GetAuthChallenge(urlStr, "POST", "")
	if err != nil {
		return "", fmt.Errorf("failed to get digest challenge: %w", err)
	}

	req, _ := http.NewRequest("POST", urlStr, nil)
	digest_auth.SetDigestAuth(req, authState, username, password)

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}

	return url.QueryEscape(strings.TrimSpace(buf.String())), nil
}
