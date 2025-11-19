package main

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Index Page"))
}

func main() {
	// This is a placeholder for the main function.
	http.HandleFunc("/", Index)
	fmt.Println("server run ing")
	http.ListenAndServe(":8081", nil)
}
