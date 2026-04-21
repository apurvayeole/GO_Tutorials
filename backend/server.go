package main

import (
	"fmt"
	"net/http"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Users service working at port 3001 🚀")
}

func main() {
	http.HandleFunc("/users", usersHandler)

	fmt.Println("User service running on 3001")
	http.ListenAndServe(":3001", nil)
}