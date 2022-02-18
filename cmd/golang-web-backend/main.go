package main

import "net/http"

type footballerHandler struct{}

func (h *footballerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/footballers/", &footballerHandler{})
	http.ListenAndServe(":8080", mux)
}
