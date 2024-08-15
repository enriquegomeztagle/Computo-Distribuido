// 1. Crear un modulo de go, será nuestro principal repositorio de trabajo
// 2. Crear un archivo llamado server.go
// 3. En el server.go
//     1. Agregar lógica para recibir peticiones http
//     2. Tener un endpoint para agregar un usuario
//     3. Otro endpoint para obtener la información
//     4. La información que se mande y que se reciba tiene que empaquetarse en un JSON

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

var tempUser User

var mu sync.Mutex

func main() {
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/getUser", getUser)
	fmt.Println("Server started on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}

// curl -X POST http://localhost:8080/addUser -d '{"name":"Iram Max", "username":"iramMax"}' -H "Content-Type: application/json"
func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		tempUser = user
		mu.Unlock()
		fmt.Fprintf(w, "User added successfully")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// curl http://localhost:8080/getUser
func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mu.Lock()
		user := tempUser
		mu.Unlock()
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
