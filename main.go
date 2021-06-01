package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	connectDB()
	defer disconnectDB()

	fetchData()

	go runTicker(grabNews)

	http.HandleFunc("/api/news", newsHandler)
	http.HandleFunc("/api/url", urlHandler)
	http.HandleFunc("/api/pattern", patternHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchData(){
	fetchURLs()
	fetchPatterns()
}

func runTicker(target func ()){
	seconds, err := time.ParseDuration(os.Getenv("TIMER"))
	if err != nil {
		panic(err.Error())
	}
	ticker := time.NewTicker(time.Second * seconds)

	for range ticker.C {
		target()
	}
}

func grabNews(){
	fp := gofeed.NewParser()

	urlMutex.Lock()
	var cUrls map[int64]string
	cUrls = make(map[int64]string)
	for k, v := range urls {
		cUrls[k] = v
	}
	urlMutex.Unlock()

	patternMutex.Lock()
	var cPatterns map[int64]string
	cPatterns = make(map[int64]string)
	for k, v := range patterns {
		cPatterns[k] = v
	}
	patternMutex.Unlock()

	for urlID, rssURL := range cUrls {
		feed, err := fp.ParseURL(rssURL)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, item := range feed.Items {
			for patternID, pattern := range cPatterns {
				isMatch, _ := regexp.MatchString(pattern, item.Title)
				if isMatch {
					db.Exec("insert into news (url_id, pattern_id, title, link) values ($1, $2, $3, $4)",
						urlID, patternID, item.Title, item.Link)
				}
			}
		}
	}
}
