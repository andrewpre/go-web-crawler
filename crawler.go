package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Crawler struct {
	startURL  string
	urls      []string
	foundData interface{}
	isAlive   bool
}

func crawler(urlString string, prompt string, messages chan<- Crawler) {
	fmt.Printf("Starting the crawler for URL: %s\n", urlString)
	resp, err := http.Get(urlString)
	crawler := Crawler{startURL: urlString, urls: []string{}, foundData: []string{}, isAlive: true}
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		crawler.urls = []string{}
		crawler.foundData = []string{}
		crawler.isAlive = false
		messages <- crawler
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		crawler.urls = []string{}
		crawler.foundData = []string{}
		crawler.isAlive = false
		messages <- crawler
		return
	}
	// fmt.Println("Response Body:", string(body))
	llresponse, err := callLLM(process(body), prompt)
	if err != nil {
		fmt.Println("Error calling LLM:", err)
		crawler.urls = []string{}
		crawler.foundData = []string{}
		crawler.isAlive = false
		messages <- crawler
		return
	}
	crawler.urls = llresponse.Urls
	crawler.foundData = llresponse.FoundData
	crawler.isAlive = false
	time.Sleep(2 * time.Second) // simulate some processing time
	messages <- crawler
	fmt.Printf("Crawler finished for URL: %s\n", urlString)
}
