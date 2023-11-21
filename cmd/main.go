// Сервер GoNews.
package main

import (
	"log"
	"net/http"

	"cmd/pkg/api2"
)

// type config struct {
// 	URLS   []string `json:"657"`
// 	Period int      `json:"cicle"`
// }

func main() {
	api := api2.New()
	// b, err := ioutil.ReadFile("./config.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var config config
	// err = json.Unmarshal(b, &config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
