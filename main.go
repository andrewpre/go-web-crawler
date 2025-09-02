package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Type in the start URL:")
	var url string
	fmt.Scanln(&url)
	sampleURL := "https://www.zappos.com/kratos/c/zappos-homepage"
	if sampleURL == "" {
		fmt.Println("No URL provided, exiting.")
		return
	}
	prompt := "Find all items that are on sale."
	// messages channel to handle communication
	messages := make(chan Crawler)

	var mapper map[string]bool = make(map[string]bool)
	var accessedURLS map[string]bool = make(map[string]bool)
	var answers []string = []string{}
	mapper[sampleURL] = true

	// initial crawler
	accessedURLS[sampleURL] = true
	go crawler(sampleURL, prompt, messages)

	// later you will keep going until theres no more crawlers active
	// you need heartbeat to check if they are alive
	for len(mapper) > 0 {
		// blocks until theres something in the channel
		msg := <-messages
		_, ok := mapper[msg.startURL]
		if !msg.isAlive {
			if ok { // still expecting the crawler to be alive
				fmt.Println("Terminating crawler with url:", msg.startURL)
				fmt.Println("Found URLs:", msg.urls, "Total URLS:", len(msg.urls))
				answers = append(answers, msg.foundData.([]string)...)
				for _, url := range msg.urls {
					mapper[url] = true
					if _, seen := accessedURLS[url]; !seen && len(accessedURLS) < 20 {
						accessedURLS[url] = true
						go crawler(url, prompt, messages)
						time.Sleep(1 * time.Second) // slight delay to avoid overwhelming
					}
				}

				delete(mapper, msg.startURL)
			}
		} else {
			fmt.Println("Received message from crawler:", msg.startURL)
		}
	}
	fmt.Println("No more active crawlers. Exiting channel handler.")
	close(messages)
	getUserAnswer(answers, prompt)
}
