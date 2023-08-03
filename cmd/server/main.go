package main

import (
	"fmt"
	"log"
	"net/http"
)

const serverPort = 8080

func main() {

	port := fmt.Sprintf(":%d", serverPort)
	// Запуск сервера на порту 8080
	log.Printf("Серевер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
