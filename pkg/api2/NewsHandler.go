package api2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (api *API) lastHandler(w http.ResponseWriter, r *http.Request) {
	var a []Article
	var p Paging
	var err error
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	form := r.URL.Query()
	p.ItemPerPage, err = strconv.Atoi(form.Get("itemperpage"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.PageN, err = strconv.Atoi(form.Get("page"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := api.Conf.Agrigator + "/news/last?" + paramToPass
	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l := len(a)
	if p.ItemPerPage == 0 || p.PageN == 0 {
		p.Page = a
		json.NewEncoder(w).Encode(p)
		return
	}
	p.PageTotal = l / p.ItemPerPage
	if p.PageTotal*p.ItemPerPage < l {
		p.PageTotal = p.PageTotal + 1
	}
	p.Page = a[p.ItemPerPage*(p.PageN-1) : p.ItemPerPage*p.PageN-2]
	json.NewEncoder(w).Encode(p)
	log.Println("API_Gateway: API: lastHandler:", "ok ", r.URL.Query().Encode())
}

func (api *API) lastListHandler(w http.ResponseWriter, r *http.Request) {
	var a []Article
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Agrigator + "/news/lastlist?" + paramToPass
	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
	log.Println("API_Gateway: API:lastListHandler: ", "ok ", r.URL.Query().Encode())
}

func (api *API) filterHandler(w http.ResponseWriter, r *http.Request) {
	var a []Article
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Agrigator + "/news/filter?" + paramToPass
	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
	log.Println("API_Gateway: API:filterHandler:", "ok ", r.URL.Query().Encode())
}

// newsHandler асинхронно собирает статью и слайс комментариев первого уровня
func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	a := Article{}
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	form := r.URL.Query()
	id, err := strconv.Atoi(form.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aChan := make(chan interface{}, 2)
	go api.articleRequest(id, aChan)
	go api.commentsRequest(id, aChan)

	ex1, ex2 := true, true
	for ex1 || ex2 {
		switch v := <-aChan; v.(type) {
		case Article:
			a.Autor = v.(Article).Autor
			a.Content = v.(Article).Content
			a.Publisher = v.(Article).Publisher
			a.Title = v.(Article).Title
			a.ID = v.(Article).ID
			a.PubTime = v.(Article).PubTime
			a.Url = v.(Article).Url
			ex1 = false
		case []Comment:
			a.Comments = v.([]Comment)
			ex2 = false
		case error:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		default:
			http.Error(w, fmt.Errorf("непонятный тип").Error(), http.StatusInternalServerError)
			return
		}
	}
	close(aChan)
	json.NewEncoder(w).Encode(a)
	log.Println("API_Gateway: API:newsHandler: ", "ok ", r.URL.Query().Encode())
}

// Синхронная версия
func (api *API) newsHandlerSynh(w http.ResponseWriter, r *http.Request) {
	a := Article{}
	if r.Method != "GET" {
		errorHandler(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	paramToPass := r.URL.Query().Encode()
	url := api.Conf.Agrigator + "/news/news?" + paramToPass

	resp, err := callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url = api.Conf.Agrigator + "/comment/comListP?pT=A&pId=" + fmt.Sprintf("%d", a.ID)
	resp, err = callOtherAPI(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(resp, &a.Comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
	log.Println("API_Gateway: API: newsHandlerSynh", "ok ", r.URL.Query().Encode())
}

func (api *API) articleRequest(id int, ch chan<- interface{}) {
	agr := Article{}
	url := api.Conf.Agrigator + "/news/news?id=" + fmt.Sprintf("%d", id)
	resp, err := callOtherAPI(url)
	if err != nil {
		ch <- err
		return
	}
	err = json.Unmarshal(resp, &agr)
	if err != nil {
		ch <- err
		return
	}
	ch <- agr
}
func (api *API) commentsRequest(id int, ch chan<- interface{}) {

	var c []Comment
	url := api.Conf.Commentator + "comment/comListP?pT=A&pId=" + fmt.Sprintf("%d", id)
	resp, err := callOtherAPI(url)
	if err != nil {
		ch <- err
		return
	}
	err = json.Unmarshal(resp, &c)
	if err != nil {
		ch <- err
		return
	}
	ch <- c
}
