// Сервер GoNews.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cmd/pkg/api2"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Запуск службы...")
	api := api2.New()

	b, err := os.ReadFile("./cmd/config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &api.Conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(api.Conf)
	log.Println("Запуск сервера. Порт", api.Conf.Port)
	err = http.ListenAndServe(api.Conf.Port, api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
