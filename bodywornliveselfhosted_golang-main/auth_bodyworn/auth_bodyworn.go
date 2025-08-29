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

type TokenResponse struct {
	Token string `json:"token"`
}
type FetchConfig struct {
	IPAddress string `json:"ip_address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TargetID  string `json:"target_id"`
}


type TokenEnvelope struct {
	APIVersion string `json:"apiVersion"`
	Method     string `json:"method"`
	Data       struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
	} `json:"data"`
}


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

// FetchToken fetches the device token via Digest Auth over HTTP
func FetchToken(ip, username, password string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    nil,
			DisableCompression: true,
		},
	}

	url := fmt.Sprintf("http://%s/local/BodyWornLiveSelfHosted/auth.cgi", ip)
	fmt.Printf("🔎 Fetching token from device (POST): %s\n", url)

	
	payload := []byte(`{"apiVersion":"1.0","method":"getSignalingClientToken","params":{}}`)

	challenge, err := digest_auth.GetAuthChallenge(url, "POST", string(payload))
	if err != nil {
		return "", fmt.Errorf("failed to get digest challenge: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	err = digest_auth.SetDigestAuth(req, challenge, username, password)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("📜 Raw device response: %s\n", string(bodyBytes))

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if result.Data.Token == "" {
		return "", fmt.Errorf("Token fetch failed: received empty token from device")
	}

	return result.Data.Token, nil
}

