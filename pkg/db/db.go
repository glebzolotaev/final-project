package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // регистрирует драйвер "sqlite"
)

// db — глобальный указатель на открытое соединение с БД.
var db *sql.DB

// schema содержит SQL-команды для создания таблицы и индекса.
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT '',
    title   VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat  VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

// Init открывает базу данных и при необходимости создаёт в ней таблицу с индексом.
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := err != nil // если файла нет — нужна первичная установка схемы

	// Открываем (или создаём) базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файл не существовал — выполняем схему
	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}
