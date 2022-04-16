package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Footballer struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FootballClub string `json:"football_club"`
}

type footballerHandlers struct {
	sync.Mutex
	store map[string]Footballer
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

func (h *footballerHandlers) get(w http.ResponseWriter, r *http.Request) {
	footballers := make([]Footballer, len(h.store))

	h.Lock()
	i := 0
	for _, footballer := range h.store {
		footballers[i] = footballer
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(footballers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *footballerHandlers) getFootballer(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusFound)
		return
	}

	h.Lock()
	footballer, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
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

	var footballer Footballer
	err = json.Unmarshal(bodyBytes, &footballer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	footballer.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[footballer.ID] = footballer
	defer h.Unlock()
}

func newFootballerHandlers() *footballerHandlers {
	return &footballerHandlers{
		store: map[string]Footballer{},
	}
}

func main() {
	footballerHandlers := newFootballerHandlers()
	http.HandleFunc("/footballer", footballerHandlers.footballers)
	http.HandleFunc("/footballer/", footballerHandlers.getFootballer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
