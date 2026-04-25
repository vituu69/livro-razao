package api

import "time"

type AccountResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   string    `json:"balance"`
	Currency  string    `json:"currency"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	IsSystem  bool      `json:"is_system"`
}

type EntryResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ReconcileResponse struct {
	Message string `json:"message"`
	Matched int    `json:"matched"`
}
