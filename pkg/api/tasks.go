package api

import (
	"net/http"

	"final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка получения задач: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writeJson(w, TasksResp{Tasks: tasks})
}
