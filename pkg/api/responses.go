package api

import (
	"encoding/json"
	"net/http"
)

// writeJson сериализует data в JSON и отправляет клиенту.
// Заголовок Content-Type обязателен — без него тесты и фронтенд
// не смогут корректно распарсить ответ.
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "ошибка формирования JSON-ответа", http.StatusInternalServerError)
	}
}
