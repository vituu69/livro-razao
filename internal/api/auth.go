package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]interface{}{
		"sub": userID.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encoder := base64.RawURLEncoding
	unsigned := encoder.EncodeToString(headerJSON) + "." + encoder.EncodeToString(claimsJSON)

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		secret = "development-secret"
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))

	return unsigned + "." + encoder.EncodeToString(mac.Sum(nil)), nil
}
