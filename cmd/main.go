package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/api/ping", pingHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err start srv")
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
