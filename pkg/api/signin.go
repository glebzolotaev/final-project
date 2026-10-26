package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// signinRequest — тело запроса на /api/signin.
type signinRequest struct {
	Password string `json:"password"`
}

// signinHandler обрабатывает POST /api/signin.
// Сверяет пароль с TODO_PASSWORD и при совпадении возвращает JWT-токен.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось прочитать запрос")
		return
	}

	var req signinRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON: "+err.Error())
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != storedPassword {
		writeError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	token, err := createToken(storedPassword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка создания токена: "+err.Error())
		return
	}

	writeJson(w, map[string]string{"token": token})
}
