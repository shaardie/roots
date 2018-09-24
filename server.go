package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shaardie/roots/storage"
)

type server struct {
	db     *storage.Storage
	router *mux.Router
	static string
}

func (s *server) routes() {
	s.router = mux.NewRouter()
	s.router.
		HandleFunc("/api/map/{name}", s.withLogging(s.getWorldMap())).
		Methods("GET")
	s.router.
		HandleFunc("/api/map/{name}/{country}/", s.withLogging(s.AddCountry())).
		Methods("POST").
		Headers("Content-Type", "application/json")
	s.router.
		HandleFunc("/api/map/", s.withLogging(s.NewWorldMap())).
		Methods("POST").
		Headers("Content-Type", "application/json")
	if s.static != "" {
		s.router.
			PathPrefix("/").
			HandlerFunc(s.withLogging(http.FileServer(http.Dir(s.static)).ServeHTTP))
	}

}

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		log.Printf("%v on %v", r.Method, r.URL)
	}
}

func (s *server) getWorldMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			http.NotFound(w, r)
			return
		}

		// Get World Map
		worldMap, err := s.db.Get(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Create Response
		resp, err := json.Marshal(worldMap.Countries)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (s *server) NewWorldMap() http.HandlerFunc {
	type requestAndResponse struct {
		Name string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		rar := requestAndResponse{}
		err := json.NewDecoder(r.Body).Decode(&rar)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if rar.Name == "" {
			rar.Name = randStringBytes(16)
		}

		// Create
		_, err = s.db.New(rar.Name)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("Map with name %v already exists", rar.Name),
				http.StatusBadRequest)
			return
		}

		// Create Response
		resp, err := json.Marshal(rar)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (s *server) AddCountry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		country, ok := vars["country"]
		if !ok {
			http.NotFound(w, r)
			return
		}

		// Get World Map
		worldMap, err := s.db.Get(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Add Country
		err = worldMap.Add(country)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Create Response
		resp, err := json.Marshal(worldMap.Countries)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
