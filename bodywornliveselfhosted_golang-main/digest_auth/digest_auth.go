package digest_auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// AuthState structure holds authentication state values
type AuthState struct {
	nonce     string
	opaque    string
	qop       string
	realm     string
	algorithm string
	nc        int
}

// Function to get the challenge and parse the WWW-Authenticate header
func GetAuthChallenge(url string, method string, payload string) (AuthState, error) {
	state := AuthState{}
	client := &http.Client{}

	// Function to make the challenge request
	makeChallengeRequest := func(method string) (*http.Response, error) {
		var req *http.Request
		var err error

		// Create request based on method
		if method == "POST" {
			req, err = http.NewRequest("POST", url, strings.NewReader(payload))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			// Default to GET request
			req, err = http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, err
			}
		}

		// Send request and get response
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	// Attempt with the given method first
	resp, err := makeChallengeRequest(method)
	if err != nil {
		return state, err
	}
	defer resp.Body.Close()

	// Debugging: Output the status code
	

	// Check if the challenge was issued
	if resp.StatusCode == http.StatusUnauthorized {
		// Challenge received, continue with Digest Auth flow
		
	} else {
		// Retry with the other method if no challenge received
		

		// Retry with the opposite method
		alternateMethod := "GET"
		if method == "GET" {
			alternateMethod = "POST"
		}

		resp, err = makeChallengeRequest(alternateMethod)
		if err != nil {
			return state, err
		}
		defer resp.Body.Close()

		// Debugging: Output the status code for retry
		fmt.Printf("Retry Challenge Attempt with %s: %d\n", alternateMethod, resp.StatusCode)

		// If still no challenge, return an error
		if resp.StatusCode != http.StatusUnauthorized {
			return state, fmt.Errorf("expected 401 Unauthorized, got %d with %s and %s", resp.StatusCode, method, alternateMethod)
		}
	}

	// Extract WWW-Authenticate header
	authHeader := resp.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		return state, fmt.Errorf("WWW-Authenticate header not found")
	}

	// Use regular expressions to extract needed values
	r := regexp.MustCompile(`([a-zA-Z]+)="([^"]+)"`)
	matches := r.FindAllStringSubmatch(authHeader, -1)

	for _, match := range matches {
		switch strings.ToLower(match[1]) {
		case "nonce":
			state.nonce = match[2]
		case "realm":
			state.realm = match[2]
		case "qop":
			state.qop = match[2]
		case "opaque":
			state.opaque = match[2]
		case "algorithm":
			state.algorithm = match[2]
		}
	}

	// Set default algorithm to MD5 if not specified
	if state.algorithm == "" {
		state.algorithm = "MD5"
	}

	if state.nonce == "" || state.realm == "" {
		return state, fmt.Errorf("failed to parse nonce or realm from WWW-Authenticate header")
	}

	state.nc = 1 // Initialize nonce count

	return state, nil
}


// SetDigestAuth sets the Digest Authentication header for the request
func SetDigestAuth(req *http.Request, state AuthState, username, password string) error {
	uri := req.URL.Path // Only the path, no query params

	parameters, err := digest(state, req.Method, uri, username, password)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Digest "+parameters)
	return nil
}

// digest creates digest parameters
func digest(
	state AuthState,
	method string,
	uri string,
	username string,
	password string,
) (digest string, err error) {
	nc := fmt.Sprintf("%08x", state.nc)
	cnonce := cnonceStr()

	return digestParams(
		digestParamsArgs{
			algorithm: state.algorithm,
			cnonce:    cnonce,
			method:    method,
			nc:        nc,
			nonce:     state.nonce,
			opaque:    state.opaque,
			password:  password,
			qop:       state.qop,
			realm:     state.realm,
			uri:       uri,
			username:  username,
		},
	)
}

type digestParamsArgs struct {
	algorithm string
	cnonce    string
	method    string
	nc        string
	nonce     string
	opaque    string
	password  string
	qop       string
	realm     string
	uri       string
	username  string
}

func digestParams(args digestParamsArgs) (digest string, err error) {
	var h hash.Hash
	var response string

	// Use MD5 hashing by default
	switch args.algorithm {
	case "", "MD5":
		h = md5.New()
	default:
		return "", fmt.Errorf("[digest]: unsupported algorithm: %q", args.algorithm)
	}

	ha1 := hashStr(h, args.username+":"+args.realm+":"+args.password)
	ha2 := hashStr(h, args.method+":"+args.uri)

	// Handle qop cases (Quality of Protection)
	switch {
	case args.qop == "auth":
		response = hashStr(h, strings.Join([]string{ha1, args.nonce, args.nc, args.cnonce, args.qop, ha2}, ":"))
	default:
		return "", fmt.Errorf("[Digest]: unsupported qop: %q", args.qop)
	}

	// Construct digest parameters
	params := make([]string, 0, 10)

	params = append(params,
		`username="`+args.username+`"`,
		`realm="`+args.realm+`"`,
		`uri="`+args.uri+`"`,
		`algorithm="`+args.algorithm+`"`,
		`nonce="`+args.nonce+`"`,
		`nc=`+args.nc,
		`cnonce="`+args.cnonce+`"`,
		`qop=`+args.qop,
		`response="`+response+`"`,
	)
	if args.opaque != "" {
		params = append(params, `opaque="`+args.opaque+`"`)
	}

	digest = strings.Join(params, ", ")

	return digest, nil
}

func hashStr(h hash.Hash, s string) string {
	h.Reset()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func cnonceStr() string {
	b := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

