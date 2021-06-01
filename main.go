package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	connectDB()
	defer disconnectDB()

	http.HandleFunc("/api/news", newsHandler)
	http.HandleFunc("/api/url", urlHandler)
	http.HandleFunc("/api/pattern", patternHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}
	fmt.Println("news GET")
}

