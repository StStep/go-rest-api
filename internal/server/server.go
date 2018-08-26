package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/StStep/go-test-server/internal/auth"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Status struct {
	PeopleNumber int    `json:"peoplenumber,omitempty"`
	Status       string `json:"status,omitempty"`
}

var people []Person
var status Status

func Start() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})
	status = Status{len(people), "Ready"}

	router.HandleFunc("/", getHome).Methods("GET")
	router.HandleFunc("/login", postLogin).Methods("POST")
	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/people/{id}", getPerson).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router)))
}

func getHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(status)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	token, err := auth.MakeToken()
	if err != nil {
		fmt.Println("Failed to generate token with error " + err.Error())
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(token))
	}
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	if !auth.AuthRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	if !auth.AuthRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
