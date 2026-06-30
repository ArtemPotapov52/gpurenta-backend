package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// VerifyToken verifies a Google-issued access_token via the userinfo API.
func VerifyToken(accessToken string) (*GoogleUser, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)
	return callGoogleAPI(url)
}

// VerifyIDToken verifies a Google id_token (credential from GIS) via the tokeninfo endpoint.
func VerifyIDToken(idToken string) (*GoogleUser, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
	return callGoogleAPI(url)
}

func callGoogleAPI(url string) (*GoogleUser, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("google api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status %d", resp.StatusCode)
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode google response: %w", err)
	}

	if user.Sub == "" {
		return nil, fmt.Errorf("invalid google token: no sub")
	}

	return &user, nil
}
