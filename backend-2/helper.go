package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	POST   string = "POST"
	GET    string = "GET"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

type Handler struct {
	Method   string
	Callback func(http.ResponseWriter, *http.Request)
}

var handlers = map[string]map[string]func(http.ResponseWriter, *http.Request){}

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func JSONRes(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	if res == nil {
		res = map[string]interface{}{}
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(jsonBytes)
}

func readBody(w http.ResponseWriter, r *http.Request) []byte {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		JSONRes(w, 500, err)
	}
	return bodyBytes
}

func checkContentType(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		JSONRes(w, 415, nil)
	}
}

func returnError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed"))
}

func callback(
	w http.ResponseWriter,
	r *http.Request,
	methodMap map[string]func(http.ResponseWriter, *http.Request),
	method string,
) {
	callback, ok := methodMap[method]
	if !ok {
		returnError(w)
		return
	}

	callback(w, r)
}

func execute(url string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		methodMap, ok := handlers[url]
		if !ok {
			returnError(w)
			return
		}

		switch r.Method {
		case "GET":
			callback(w, r, methodMap, GET)
			return
		case "POST":
			callback(w, r, methodMap, POST)
			return
		case "PUT":
			callback(w, r, methodMap, PUT)
			return
		case "DELETE":
			callback(w, r, methodMap, DELETE)
			return
		default:
			returnError(w)
			return
		}
	}
}

func assignHandler(url string, method string, handler func(http.ResponseWriter, *http.Request)) {
	if urlMap, ok := handlers[url]; ok {
		_, ok := urlMap[method]
		if !ok {
			handlers[url][method] = handler
		}

		return
	}

	methodMap := map[string]func(http.ResponseWriter, *http.Request){}
	methodMap[method] = handler
	handlers[url] = methodMap
	http.HandleFunc(url, execute(url))
}

func Post(url string, callback func(http.ResponseWriter, *http.Request)) {
	assignHandler(url, POST, callback)
}

func Get(url string, callback func(http.ResponseWriter, *http.Request)) {
	assignHandler(url, GET, callback)
}

func Put(url string, callback func(http.ResponseWriter, *http.Request)) {
	assignHandler(url, PUT, callback)
}

func Delete(url string, callback func(http.ResponseWriter, *http.Request)) {
	assignHandler(url, DELETE, callback)
}
