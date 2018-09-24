package main

import (
	"log"
	"net/http"

	"github.com/shaardie/roots/storage"
)

func main() {
	if err := mainWithErrors(); err != nil {
		log.Fatal(err)
	}
}

func mainWithErrors() error {

	// Create Server
	s := server{
		db:     storage.New(),
		static: "./static/",
	}
	s.routes()

	log.Println("Listening...")
	return http.ListenAndServe("0.0.0.0:8080", s.router)
}
