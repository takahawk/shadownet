package main

import (
	"fmt"
	"net/http"
)

func gateway(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Work in Progress")
}


func main() {
	port := 1337
	http.HandleFunc("/", gateway)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}