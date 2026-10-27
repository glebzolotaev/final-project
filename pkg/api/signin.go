package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type signinRequest struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "request could not be read")
		return
	}

	var req signinRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON deserialization error: "+err.Error())
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != storedPassword {
		writeError(w, http.StatusUnauthorized, "incorrect password")
		return
	}

	token, err := createToken(storedPassword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "token creation error: "+err.Error())
		return
	}

	writeJson(w, map[string]string{"token": token})
}
