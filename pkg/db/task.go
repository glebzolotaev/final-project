package db

import (
	"database/sql"
	"errors"
	"time"
)

const DateLayout = "20060102"

// Task — описание задачи из таблицы scheduler.
// Структура определена в пакете db, так как нужна и для записи в БД,
// и для сериализации/десериализации JSON в обработчиках api.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает
// идентификатор добавленной записи.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Tasks возвращает список ближайших задач, отсортированный по дате.
func Tasks(search string, limit int) ([]*Task, error) {
	base := `SELECT id, date, title, comment, repeat FROM scheduler`

	var (
		rows *sql.Rows
		err  error
	)

	switch {
	case search == "":
		rows, err = db.Query(base+` ORDER BY date ASC, id ASC LIMIT ?`, limit)
	case isDateSearch(search):
		// 08.02.2024 -> 20240208
		t, _ := time.Parse("02.01.2006", search)
		rows, err = db.Query(base+` WHERE date = ? ORDER BY date ASC, id ASC LIMIT ?`,
			t.Format(DateLayout), limit)
	default:
		pattern := "%" + search + "%"
		rows, err = db.Query(base+` WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC, id ASC LIMIT ?`,
			pattern, pattern, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// isDateSearch определяет, является ли строка датой формата 02.01.2006.
func isDateSearch(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

// GetTask возвращает задачу по её идентификатору.
// Если запись не найдена, возвращает ошибку.
func GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Задача не найдена")
		}
		return nil, err
	}
	return &t, nil
}

// UpdateTask обновляет задачу по её id.
// Если записи с таким id нет, возвращает ошибку.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Задача не найдена")
	}
	return nil
}
