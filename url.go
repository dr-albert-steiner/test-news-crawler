package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var urls map[int64]string
var urlMutex sync.Mutex

type URLRequest struct {
	URL string
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	urlMutex.Lock()
	defer urlMutex.Unlock()
	switch r.Method {
	case "GET":
		getURL(w)
	case "POST":
		postURL(w, r)
	case "DELETE":
		deleteURL(w, r)
	default:
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}
}

func getURL(w http.ResponseWriter){
	jsonData, err := json.Marshal(urls)
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

func postURL(w http.ResponseWriter, r *http.Request) {
	var rssURL URLRequest
	err := decodeJSON(r.Body, &rssURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, item := range urls {
		if item == rssURL.URL {
			http.Error(w, "URL already exists", http.StatusNotFound)
			return
		}
	}

	result, err := db.Exec("insert into urls (url) values ($1)", rssURL.URL)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	rowID, _ := result.LastInsertId()
	urls[rowID] = rssURL.URL
}

func deleteURL(w http.ResponseWriter, r *http.Request){
	var rssURL URLRequest
	err := decodeJSON(r.Body, &rssURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := db.Exec("delete from urls where url = $1", rssURL.URL)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	rowID, _ := result.LastInsertId()
	delete(urls, rowID)
}

func fetchURLs(){
	if db == nil {
		panic("Database is not connected")
	}

	rows, err := db.Query("select id, url from urls")
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
	urls = make(map[int64]string)
	for rows.Next() {
		var newUrl string
		var id int64
		err := rows.Scan(&id, &newUrl)
		if err != nil{
			log.Println(err)
			continue
		}
		urls[id] = newUrl
	}
}
