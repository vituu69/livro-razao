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

type RegisterResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

type EntryResponse struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	Debit         string    `json:"debit"`
	Credit        string    `json:"credit"`
	TransactionID string    `json:"transaction_id"`
	OperationType string    `json:"operation_type"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
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
