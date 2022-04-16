package main

import (
	"encoding/json"
	"fmt"
	"golang-web-backend/models"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type footballerHandlers struct {
	sync.Mutex
	store map[string]models.Footballer
}

func (h *footballerHandlers) footballers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *footballerHandlers) footballer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getOne(w, r)
		return
	case "PUT":
		h.update(w, r)
		return
	case "DELETE":
		h.delete(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *footballerHandlers) get(w http.ResponseWriter, r *http.Request) {
	footballers, err := models.GetAllFootballers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	jsonBytes, err := json.Marshal(footballers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *footballerHandlers) getOne(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	urlID := parts[2]

	footballer, err := models.GetFootballerByID(urlID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	jsonBytes, err := json.Marshal(footballer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *footballerHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var footballer models.Footballer
	err = json.Unmarshal(bodyBytes, &footballer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	createdFootballer, err := models.CreateFootballer(&footballer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(createdFootballer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *footballerHandlers) update(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	urlID := parts[2]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var footballer models.Footballer
	err = json.Unmarshal(bodyBytes, &footballer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	updatedFootballer, err := models.UpdateFootballer(urlID, &footballer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	jsonBytes, err := json.Marshal(updatedFootballer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *footballerHandlers) delete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	urlID := parts[2]

	status, err := models.DeleteFootballer(urlID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func newFootballerHandlers() *footballerHandlers {
	return &footballerHandlers{
		store: map[string]models.Footballer{},
	}
}

func main() {
	footballerHandlers := newFootballerHandlers()
	models.SetMongoConnection()
	http.HandleFunc("/footballer", footballerHandlers.footballers)
	http.HandleFunc("/footballer/", footballerHandlers.footballer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
