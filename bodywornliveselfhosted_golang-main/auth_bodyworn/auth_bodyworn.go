package auth_bodyworn

import (
	"bodywornliveselfhosted/digest_auth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// TokenResponse represents the JSON structure returned by the Axis device
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

// FetchConfig defines configuration for connecting to the Axis device
type FetchConfig struct {
	IPAddress string `json:"ip_address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TargetID  string `json:"target_id"`
}

// TokenEnvelope describes the structure of the Axis API token response
type TokenEnvelope struct {
	APIVersion string `json:"apiVersion"`
	Method     string `json:"method"`
	Data       struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
	} `json:"data"`
}

// LoadConfig loads credentials and device info from config.json
func LoadConfig() (*FetchConfig, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg FetchConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// FetchToken fetches a signaling token from the Axis Body Worn device using Digest Auth.
// It now returns both the token and its expiration time (if provided).
func FetchToken(ip, username, password string) (string, time.Time, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    nil,
			DisableCompression: true,
		},
	}

	url := fmt.Sprintf("http://%s/local/BodyWornLiveSelfHosted/auth.cgi", ip)
	fmt.Printf(" Fetching token from device (POST): %s\n", url)

	payload := []byte(`{"apiVersion":"1.0","method":"getSignalingClientToken","params":{}}`)

	// Step 1: Get Digest Auth challenge
	challenge, err := digest_auth.GetAuthChallenge(url, "POST", string(payload))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get digest challenge: %w", err)
	}

	// Step 2: Build request with Digest Auth
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	err = digest_auth.SetDigestAuth(req, challenge, username, password)
	if err != nil {
		return "", time.Time{}, err
	}

	// Step 3: Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf(" Raw device response: %s\n", string(bodyBytes))

	// Step 4: Parse token + expiry
	var result struct {
		Data struct {
			Token     string `json:"token"`
			ExpiresAt string `json:"expiresAt"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse token response: %w", err)
	}

	token := strings.TrimSpace(result.Data.Token)
	expiresAtStr := strings.TrimSpace(result.Data.ExpiresAt)

	if token == "" {
		return "", time.Time{}, fmt.Errorf("token fetch failed: received empty token from device")
	}

	// Step 5: Parse expiration (if present)
	var expiry time.Time
	if expiresAtStr != "" {
		expiry, err = time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			fmt.Printf(" Warning: failed to parse expiresAt: %v\n", err)
			expiry = time.Now().Add(60 * time.Second) // fallback 1 min lifetime
		}
	} else {
		fmt.Println(" No expiresAt field in response — using default 60 seconds.")
		expiry = time.Now().Add(60 * time.Second)
	}

	fmt.Printf(" Token fetched successfully (expires at %s)\n", expiry.Format(time.RFC3339))
	return token, expiry, nil
}


