package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RSSURL struct {
	URL string
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
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
	result, err := db.Query("select url from urls")
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
	var urlList []RSSURL
	for result.Next() {
		newUrl := RSSURL{}
		err := result.Scan(&newUrl.URL)
		if err != nil{
			log.Println(err)
			continue
		}
		urlList = append(urlList, newUrl)
	}
	jsonData, err := json.Marshal(urlList)
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
	var rssURL RSSURL
	err := decodeJSON(r.Body, &rssURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_, err = db.Exec("insert into urls (url) values ($1)", rssURL.URL)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func deleteURL(w http.ResponseWriter, r *http.Request){
	var rssURL RSSURL
	err := decodeJSON(r.Body, &rssURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_, err = db.Exec("delete from urls where url = $1", rssURL.URL)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
