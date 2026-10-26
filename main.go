package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"final-project/pkg/db"
	"final-project/tests"
)

func main() {
	// Инициализация базы данных
	if err := db.Init(getDBFile()); err != nil {
		log.Fatalf("Ошибка при инициализации БД: %v", err)
	}

	port := getPort()
	webDir := "./web"

	log.Printf("Запуск веб-сервера на порту :%s, раздача файлов из: %s", port, webDir)

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// getPort определяет порт.
// Приоритет: переменная окружения TODO_PORT > tests.Port из tests/settings.go
func getPort() string {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		return envPort
	}
	return strconv.Itoa(tests.Port)
}

// getDBFile определяет путь к файлу базы данных.
// Приоритет: переменная окружения TODO_DBFILE > "scheduler.db" (по умолчанию)
func getDBFile() string {
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		return envDBFile
	}
	return "scheduler.db"
}
