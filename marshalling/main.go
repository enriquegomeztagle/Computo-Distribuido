package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	user := User{
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Convert JSON data to a string and print it
	fmt.Println(string(jsonData))
}
