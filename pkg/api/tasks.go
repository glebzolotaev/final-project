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
		writeError(w, http.StatusMethodNotAllowed, "method is not supported")
		return
	}

	search := r.FormValue("search")

	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error in task retrieval: "+err.Error())
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
