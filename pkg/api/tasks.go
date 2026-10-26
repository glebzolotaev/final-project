package api

import (
	"net/http"

	"final-project/pkg/db"
)

// TasksResp — обёртка для ответа: список задач в поле tasks.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks.
// Необязательный параметр search: подстрока для поиска по title/comment
// или дата формата 02.01.2006 для выборки задач на конкретный день.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	search := r.FormValue("search")

	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка получения задач: "+err.Error())
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
