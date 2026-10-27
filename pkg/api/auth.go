package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "scheduler_jwt_secret_key"

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func createToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"hash": hashPassword(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func validateToken(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hashClaim, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	return hashClaim == hashPassword(os.Getenv("TODO_PASSWORD"))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string
			if cookie, err := r.Cookie("token"); err == nil {
				jwtStr = cookie.Value
			}
			if jwtStr == "" || !validateToken(jwtStr) {
				writeError(w, http.StatusUnauthorized, "authentication is required")
				return
			}
		}
		next(w, r)
	})
}
