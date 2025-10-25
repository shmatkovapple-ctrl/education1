package main

import (
	"fmt"
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет")
}

func main() {
	http.HandleFunc("/", home_page)
	http.ListenAndServe(":8000", nil)
}
