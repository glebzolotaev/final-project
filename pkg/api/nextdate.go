package api

import (
	"net/http"
	"time"

	"final-project/pkg/nextdate"
)

// nextDayHandler обрабатывает GET /api/nextdate.
//
// Формат запроса: /api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>
// Если now не задан, используется текущая дата.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(nextdate.DateLayout, nowStr)
		if err != nil {
			http.Error(w, "некорректное значение параметра now", http.StatusBadRequest)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	next, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
