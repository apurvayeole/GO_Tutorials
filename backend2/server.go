package main

import (
	"fmt"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Users service working 3003 🚀")
}

func main() {
	http.HandleFunc("/users", usersHandler)

	fmt.Println("User service running on 3003")
	http.ListenAndServe(":3003", nil)
}