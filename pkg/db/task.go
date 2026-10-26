package db

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
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat
	          FROM scheduler
	          ORDER BY date ASC, id ASC
	          LIMIT ?`
	rows, err := db.Query(query, limit)
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
