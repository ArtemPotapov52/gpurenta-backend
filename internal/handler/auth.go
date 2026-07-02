package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemPotapov52/gpurenta/internal/auth"
	"github.com/ArtemPotapov52/gpurenta/internal/db"
	"github.com/ArtemPotapov52/gpurenta/internal/middleware"
)

type AuthHandler struct {
	Store     *db.Store
	JWTSecret string
}

type authRequest struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type authResponse struct {
	Token string      `json:"token"`
	User  *authUser   `json:"user"`
}

type authUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *AuthHandler) DevLogin(w http.ResponseWriter, r *http.Request) {
	user, err := h.Store.FindOrCreateUser(r.Context(), "dev-user-1", "dev@test.com", "Dev User", "")
	if err != nil {
		middleware.JSONError(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		middleware.JSONError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{
		Token: token,
		User: &authUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *AuthHandler) Google(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var googleUser *auth.GoogleUser
	var err error

	if req.AccessToken != "" {
		googleUser, err = auth.VerifyToken(req.AccessToken)
	} else if req.IDToken != "" {
		googleUser, err = auth.VerifyIDToken(req.IDToken)
	} else {
		middleware.JSONError(w, "access_token or id_token is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		middleware.JSONError(w, "invalid google token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.Store.FindOrCreateUser(r.Context(),
		googleUser.Sub, googleUser.Email, googleUser.Name, googleUser.Picture,
	)
	if err != nil {
		middleware.JSONError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		middleware.JSONError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{
		Token: token,
		User: &authUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	})
}
