package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func main() {
	http.Handle("/stat/", http.StripPrefix("/stat/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
