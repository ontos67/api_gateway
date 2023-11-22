// Сервер GoNews.
package main

import (
	"log"
	"net/http"
	"os"

	"cmd/pkg/api2"
)

// type config struct {
// 	URLS   []string `json:"657"`
// 	Period int      `json:"cicle"`
// }

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Запуск службы...")
	api := api2.New()

	log.Println("Запуск сервера. Порт: 80...")
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
