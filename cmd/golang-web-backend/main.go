package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("REST web server started")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/status was called")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		status := map[string]string{"status": "OK"}
		json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/", handleRoot)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("/ was called")
	w.WriteHeader(http.StatusNotImplemented)
}
