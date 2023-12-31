package api2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// http://localhost:999/comment/save?userid=64&text=заманали%20комары&pubtime=12344134&ptype=A&pid=2345
func (api *API) commentSaveHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	//а должен быть POST, так как запрос будет формироваться на стороне фронта JS
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Commentator + "/comment/save?" + paramToPass
	fmt.Println("url:", url)
	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(id)
	log.Println("API_Gateway: API:commentSaveHandler: ", "ok ", r.URL.Query().Encode())
}

// http://localhost:999/comment/del?id=64
func (api *API) commentDelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Commentator + "/comment/del?" + paramToPass
	_, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("ok")
	log.Println("API_Gateway: API:commentDelHandler: ", "ok ", r.URL.Query().Encode())
}

// http://localhost:999/comment/comListP?pT=C&pId=47
func (api *API) commenListPHandler(w http.ResponseWriter, r *http.Request) {
	var c []Comment
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Commentator + "/comment/comListP?" + paramToPass
	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
	log.Println("API_Gateway: API:commenListPHandler: ", "ok ", r.URL.Query().Encode())
}
