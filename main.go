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

func urlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("URLs GET")
	case "POST":
		fmt.Println("URLs POST")
	case "DELETE":
		fmt.Println("URLs DELETE")
	default:
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}
}

func patternHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("patterns GET")
	case "POST":
		fmt.Println("patterns POST")
	case "DELETE":
		fmt.Println("patterns DELETE")
	default:
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}
}
