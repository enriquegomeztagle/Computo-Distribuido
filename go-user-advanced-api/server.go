package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
}

var users []User

var mu sync.Mutex

func main() {
	http.HandleFunc("/user/search", searchUser)
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addUser(w, r)
		case http.MethodGet:
			getUser(w, r)
		default:
			http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

// curl -X POST http://localhost:8080/user -d '{"name":"Iram Max", "username":"iramMax"}' -H "Content-Type: application/json"
func addUser(w http.ResponseWriter, r *http.Request) {
	var tempUser User

	err := json.NewDecoder(r.Body).Decode(&tempUser)
	if err != nil {
		http.Error(w, "You sent an invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	users = append(users, tempUser)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User added successfully: %s", tempUser.Name)
}

// curl http://localhost:8080/user
func getUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	userList := users
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

// Scapes zsh
// curl http://localhost:8080/user/search\?username=iramMax

// curl 'http://localhost:8080/user/search?username=iramMax'
func searchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Requires `username` field ", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, user := range users {
		if user.UserName == username {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User was not found in list", http.StatusNotFound)
}

// curl -X POST http://localhost:8080/user -d '{"name":"Iram Max", "username":"iramMax"}' -H "Content-Type: application/json"
// curl -X POST http://localhost:8080/user -d '{"name":"Maria Garcia", "username":"mgarcia"}' -H "Content-Type: application/json"
