package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func patternHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getPattern(w)
	case "POST":
		postPattern(w, r)
	case "DELETE":
		deletePattern(w, r)
	default:
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}
}

type RegexPattern struct {
	Pattern string
}

func getPattern(w http.ResponseWriter){
	result, err := db.Query("select pattern from patterns")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	defer func() {
		err = result.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	var patternsList []RegexPattern
	for result.Next() {
		pattern := RegexPattern{}
		err := result.Scan(&pattern.Pattern)
		if err != nil{
			log.Println(err)
			continue
		}
		patternsList = append(patternsList, pattern)
	}
	jsonData, err := json.Marshal(patternsList)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println(err.Error())
	}
}

func postPattern(w http.ResponseWriter, r *http.Request){
	var newPattern RegexPattern
	err := decodeJSON(r.Body, &newPattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_, err = db.Exec("insert into patterns (pattern) values ($1)", newPattern.Pattern)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func deletePattern(w http.ResponseWriter, r *http.Request) {
	var patternToDelete RegexPattern
	err := decodeJSON(r.Body, &patternToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_, err = db.Exec("delete from patterns where pattern = $1", patternToDelete.Pattern)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
