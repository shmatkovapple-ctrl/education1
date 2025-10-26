package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Println("Привет:", name)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/name", nameHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Произошла ошибка", err.Error())
	}
}
