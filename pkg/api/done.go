package api

import (
	"net/http"
	"time"

	"final-project/pkg/db"
	"final-project/pkg/nextdate"
)

// doneHandler обрабатывает POST /api/task/done?id=<идентификатор>.
// Помечает задачу выполненной:
//   - одноразовая (repeat пустой) — удаляется;
//   - периодическая — её дата заменяется на дату следующего выполнения.
func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if task.Repeat == "" {
		// одноразовая задача — удаляем
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJson(w, map[string]any{})
		return
	}

	// периодическая задача — вычисляем следующую дату
	now := time.Now()
	next, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка вычисления следующей даты: "+err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, map[string]any{})
}

// deleteTaskHandler обрабатывает DELETE /api/task?id=<идентификатор>.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор задачи")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJson(w, map[string]any{})
}
