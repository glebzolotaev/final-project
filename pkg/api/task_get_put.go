package api

import (
	"encoding/json"
	"io"
	"net/http"

	"final-project/pkg/db"
)

// getTaskHandler обрабатывает GET /api/task?id=<идентификатор> —
// возвращает JSON со всеми полями задачи.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJson(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task — сохраняет изменения задачи.
// Валидация полей та же, что и при добавлении: title обязателен,
// дата — формат 20060102, repeat — корректное правило.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось прочитать запрос")
		return
	}

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}
	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// успешное обновление — пустой JSON
	writeJson(w, map[string]any{})
}
