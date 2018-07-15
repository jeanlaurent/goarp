package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Port int
}

func (server *Server) start() {
	router := mux.NewRouter()
	router.HandleFunc("/api/devices", errorHandler(server.allDevices)).Methods("GET")

	http.Handle("/", router)

	fmt.Println("listening on", server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}

func (server *Server) allDevices(response http.ResponseWriter, request *http.Request) error {
	json, err := json.Marshal(devices.all())
	if err != nil {
		return err
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
	return nil
}

func errorHandler(handler func(response http.ResponseWriter, request *http.Request) error) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		err := handler(response, request)
		if err != nil {
			fmt.Println(err)
			response.WriteHeader(500)
			params := map[string]string{"Error": err.Error()}
			templ, _ := template.ParseFiles("template/error.tpl") // FIXME: parse only once
			response.Header().Add("Content-Type", "text/html")
			templ.Execute(response, params)
		}
	}
}
