package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NjiruClinton/reproxy-mobile-bff/initializations"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Auth struct {
	cfg *config.Config
}

func (a *Auth) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body.", initializations.CODE_BAD_REQUEST)
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Email and password are required.", initializations.CODE_BAD_REQUEST)
	}

	loginResp, err := gql.Do[loginData](r.Context(), a.gql)
}
