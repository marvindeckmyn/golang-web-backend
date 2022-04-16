package main

import (
	"encoding/json"
	"fmt"
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

func (h *footballerHandlers) getFootballers(w http.ResponseWriter, r *http.Request) {
	footballers := make([]Footballer, len(h.store))

	h.Lock()
	i := 0
	for _, footballer := range h.store {
		footballers[i] = footballer
		i++
	}
	h.Unlock()

	JSONRes(w, 200, footballers)
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

	JSONRes(w, 200, footballer)
}

func (h *footballerHandlers) postFootballer(w http.ResponseWriter, r *http.Request) {
	bodyBytes := readBody(w, r)
	checkContentType(w, r)

	var footballer Footballer
	err := json.Unmarshal(bodyBytes, &footballer)
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
	setRouter()
	Listen(8080)
}
