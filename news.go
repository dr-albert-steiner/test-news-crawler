package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type News struct {
	Title string
	Link string
}

type NewsRequest struct {
	Count int
	Find string
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Method %s is not supported", r.Method), http.StatusNotFound)
	}

	var requestBody NewsRequest
	err := decodeJSON(r.Body, &requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestBody.Find = fmt.Sprintf("%%%s%%", requestBody.Find)
	result, err := db.Query("select title, link from news where lower(title) like lower($2) limit $1",
		requestBody.Count,
		requestBody.Find)
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
	var newsList []News
	for result.Next() {
		var currentNews News
		err := result.Scan(&currentNews.Title, &currentNews.Link)
		if err != nil{
			log.Println(err)
			continue
		}
		newsList = append(newsList, currentNews)
	}
	jsonData, err := json.Marshal(newsList)
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
