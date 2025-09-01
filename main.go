package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")
	messages := make(chan string)
	go start_crawler("https://example.com", messages)
	go handleChannel(messages)
	for {
	}
}

func start_crawler(url string, messages chan<- string) {
	fmt.Printf("Starting the crawler for URL: %s\n", url)
	for {
		messages <- "Crawling started"
		time.Sleep(2 * time.Second)
	}
}

func handleChannel(messages <-chan string) {
	for msg := range messages {
		fmt.Println("Received message:", msg)
	}
}
