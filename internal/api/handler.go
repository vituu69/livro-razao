package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/vitu69/livro-razao/internal/db"
	"github.com/vitu69/livro-razao/internal/service"
	"github.com/vitu69/livro-razao/postgres/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	ledger *service.LedgerService
	store  *db.Store
}

func NewHandler(ledger *service.LedgerService, store *db.Store) *Handler {
	return &Handler{
		ledger: ledger,
		store:  store,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to decode registration request")
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		respondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user, err := h.store.CreateUser(r.Context(), sqlc.CreateUserParams{
		Email:          input.Email,
		HashedPassword: string(hashed),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		respondError(w, http.StatusConflict, "Email already exists")
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate JWT token")
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	log.Info().Str("user_id", user.ID.String()).Str("email", user.Email).Msg("User registered successfully")
	respondJSON(w, http.StatusCreated, RegisterResponse{
		UserID: user.ID.String(),
		Email:  user.Email,
		Token:  token,
	})

}
