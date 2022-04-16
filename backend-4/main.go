package main

import (
	"encoding/json"
	"golang-web-backend/models"
	"net/http"
)

func getFootballers(w http.ResponseWriter, r *http.Request) {
	footballers, err := models.GetAllFootballers()
	writeFootballers(err, w, footballers)
}

func getFootballer(w http.ResponseWriter, r *http.Request) {
	urlID := getUrlPart(r, 2)
	footballer, err := models.GetFootballerByID(urlID)
	writeFootballer(err, w, footballer)
}

func postFootballer(w http.ResponseWriter, r *http.Request) {
	bodyBytes := readBody(w, r)
	checkContentType(w, r)
	footballer := unmarshalFootballer(bodyBytes, w)

	createdFootballer, err := models.CreateFootballer(&footballer)

	writeFootballer(err, w, createdFootballer)
}

func updateFootballer(w http.ResponseWriter, r *http.Request) {
	urlID := getUrlPart(r, 2)
	bodyBytes := readBody(w, r)
	checkContentType(w, r)
	footballer := unmarshalFootballer(bodyBytes, w)

	updatedFootballer, err := models.UpdateFootballer(urlID, &footballer)
	writeFootballer(err, w, updatedFootballer)
}

func deleteFootballer(w http.ResponseWriter, r *http.Request) {
	urlID := getUrlPart(r, 2)
	status, err := models.DeleteFootballer(urlID)

	if err != nil {
		JSONRes(w, 500, err)
	}
	JSONRes(w, 200, status)
}

func unmarshalFootballer(bodyBytes []byte, w http.ResponseWriter) models.Footballer {
	var account models.Footballer
	err := json.Unmarshal(bodyBytes, &account)
	if err != nil {
		JSONRes(w, 400, err)
	}
	return account
}

func writeFootballer(err error, w http.ResponseWriter, footballer *models.Footballer) {
	if err != nil {
		JSONRes(w, 500, err)
	}
	JSONRes(w, 200, footballer)
}

func writeFootballers(err error, w http.ResponseWriter, footballers []*models.Footballer) {
	if err != nil {
		JSONRes(w, 500, err)
	}
	JSONRes(w, 200, footballers)
}

func main() {
	setRouter()
	models.SetMongoConnection()
	Listen(8080)
}
