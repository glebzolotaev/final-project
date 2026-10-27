package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"final-project/pkg/api"
	"final-project/pkg/db"
	"final-project/tests"
)

func main() {
	if err := db.Init(getDBFile()); err != nil {
		log.Fatalf("Ошибка при инициализации БД: %v", err)
	}

	api.Init()

	port := getPort()
	webDir := "./web"

	log.Printf("Запуск веб-сервера на порту :%s, раздача файлов из: %s", port, webDir)

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func getPort() string {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		return envPort
	}
	return strconv.Itoa(tests.Port)
}

func getDBFile() string {
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		return envDBFile
	}
	return "scheduler.db"
}
