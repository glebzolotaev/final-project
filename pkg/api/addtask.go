package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"final-project/pkg/db"
	"final-project/pkg/nextdate"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "request could not be read")
		return
	}

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeError(w, http.StatusBadRequest, "JSON deserialization error: "+err.Error())
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "task title is not specified")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error when adding a task to the database: "+err.Error())
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(nextdate.DateLayout)
	}

	t, err := time.Parse(nextdate.DateLayout, task.Date)
	if err != nil {
		return errors.New("date is presented in an incorrect format, the expected format is YYYYMMDD")
	}

	var next string
	if task.Repeat != "" {
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(nextdate.DateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}

func afterNow(date, now time.Time) bool {
	d1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d1.After(d2)
}
