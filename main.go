package main

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/o1egl/paseto"
	"os"
	"time"
	"net/http"
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
	PeopleNumber int `json:"peoplenumber,omitempty"`
	Status string `json:"status,omitempty"`
}

var symmetricKey = []byte("YELLOW SUBMARINE, BLACK WIZARDRY")
var people []Person
var status Status

// our main function
func main() {
	startServer()
}

func startServer() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})
	status = Status{len(people), "Ready"}

	router.HandleFunc("/", getHome).Methods("GET")
	router.HandleFunc("/login", postLogin).Methods("POST")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router))
}

func makeToken() (string, error) {
	now := time.Now()
	exp := now.Add(8 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test_service",
		Jti:        "123",
		Subject:    "test_subject",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	// Add custom claim to the token
	jsonToken.Set("data", "this is a signed message")
	footer := "some footer"

	v2 := paseto.NewV2()

	// Encrypt data
	return v2.Encrypt(symmetricKey, jsonToken, paseto.WithFooter(footer))
}

func getHome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(status)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	token, err := makeToken()
	if  err != nil {
		fmt.Println("Failed to generate token with error " + err.Error())
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(token))
	}
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
