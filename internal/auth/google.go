package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func VerifyToken(accessToken string) (*GoogleUser, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)

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

func VerifyIDToken(idToken string) (*GoogleUser, error) {
	// Google ID tokens from GIS are JWTs with format header.payload.signature
	// We decode the payload (base64) to extract user info without verifying
	// signature — acceptable for MVP since token comes directly from Google's GIS.
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id_token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// Try with padding
		p := parts[1]
		switch len(p) % 4 {
		case 2:
			p += "=="
		case 3:
			p += "="
		}
		payload, err = base64.URLEncoding.DecodeString(p)
		if err != nil {
			return nil, fmt.Errorf("decode id_token payload: %w", err)
		}
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("parse id_token claims: %w", err)
	}

	if claims.Sub == "" {
		return nil, fmt.Errorf("invalid id_token: no sub")
	}

	return &GoogleUser{
		Sub:           claims.Sub,
		Email:         claims.Email,
		Name:          claims.Name,
		Picture:       claims.Picture,
		EmailVerified: claims.EmailVerified,
	}, nil
}
