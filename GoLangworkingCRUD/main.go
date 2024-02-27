package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}

var users = []User{}
var idCounter int

func main() {
	r := mux.NewRouter()
	usersR := r.PathPrefix("/users").Subrouter()
	usersR.Path("").Methods(http.MethodGet).HandlerFunc(getAllUsers)
	usersR.Path("").Methods(http.MethodPost).HandlerFunc(createUser)
	usersR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(getUserByID)
	usersR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(updateUser)
	usersR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(deleteUser)

	fmt.Println("Start listening")
	fmt.Println(http.ListenAndServe(":8080", r))
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func indexByID(users []User, id string) int {
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return i
		}
	}
	return -1
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(users, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users[index]); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(users, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	users[index] = u
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(users, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	users = append(users[:index], users[index+1:]...)
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}
	users = append(users, u)
	response, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
