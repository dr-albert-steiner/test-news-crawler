package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var patterns map[int64]string
var patternMutex sync.Mutex

type PatternRequest struct {
	Pattern string
}

func patternHandler(w http.ResponseWriter, r *http.Request) {
	patternMutex.Lock()
	defer patternMutex.Unlock()
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



func getPattern(w http.ResponseWriter){
	jsonData, err := json.Marshal(patterns)
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
	var newPattern PatternRequest
	err := decodeJSON(r.Body, &newPattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, item := range patterns {
		if item == newPattern.Pattern {
			http.Error(w, "Pattern is already exists", http.StatusNotFound)
			return
		}
	}

	result, err := db.Exec("insert into patterns (pattern) values ($1)", newPattern.Pattern)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	rowID, _ := result.LastInsertId()
	patterns[rowID] = newPattern.Pattern
}

func deletePattern(w http.ResponseWriter, r *http.Request) {
	var patternToDelete PatternRequest
	err := decodeJSON(r.Body, &patternToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := db.Exec("delete from patterns where pattern = $1", patternToDelete.Pattern)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	rowID, _ := result.LastInsertId()
	delete(patterns, rowID)
}

func fetchPatterns(){
	if db == nil {
		panic("Database is not connected")
	}

	rows, err := db.Query("select id, pattern from patterns")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	patterns = make(map[int64]string)
	for rows.Next() {
		var pattern string
		var id int64
		err := rows.Scan(&id, &pattern)
		if err != nil{
			log.Println(err)
			continue
		}
		patterns[id] = pattern
	}
}
