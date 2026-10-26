package api

import "net/http"

// Init регистрирует все API-обработчики на стандартном маршрутизаторе.
// Вызывается один раз из main() до запуска сервера.
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
}

// taskHandler — промежуточный обработчик: раздаёт запросы по HTTP-методам
// (на следующих шагах добавятся GET и DELETE).
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

// writeError отправляет JSON вида {"error": "<текст>"} с указанным статусом.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	writeJson(w, map[string]string{"error": message})
}
