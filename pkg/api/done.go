package api

import (
	"net/http"
	"time"

	"final-project/pkg/db"
	"final-project/pkg/nextdate"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "task identifier is not specified")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJson(w, map[string]any{})
		return
	}

	now := time.Now()
	next, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error in calculating the next date: "+err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "task identifier is not specified")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJson(w, map[string]any{})
}
