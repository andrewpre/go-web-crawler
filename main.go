package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go"
)

type Crawler struct {
	startURL  string
	urls      []string
	foundData []string
	isAlive   bool
}

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("Type in the start URL:")
	var url string
	fmt.Scanln(&url)
	sampleURL := "https://www.zappos.com/kratos/c/zappos-homepage"
	if sampleURL == "" {
		fmt.Println("No URL provided, exiting.")
		return
	}
	// messages channel to handle communication
	messages := make(chan Crawler)

	var m map[string]bool = make(map[string]bool)
	m[sampleURL] = true
	// initial crawler
	go crawler(sampleURL, messages)

	// channel handler
	go handleChannel(messages, m)

	// later you will keep going until theres no more crawlers active
	// you need heartbeat to check if they are alive
	for {
		time.Sleep(10 * time.Second)
	}
}

func crawler(urlString string, messages chan<- Crawler) {
	fmt.Printf("Starting the crawler for URL: %s\n", urlString)
	resp, err := http.Get(urlString)
	crawler := Crawler{startURL: urlString, urls: []string{}, foundData: []string{}, isAlive: true}
	if err != nil {
		// messages <- fmt.Sprintf("Invalid URL: %s", err)
		fmt.Println("Error fetching URL:", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
	callLLM(body)
	i := 0
	for {

		messages <- crawler
		time.Sleep(2 * time.Second)
		i++
		if i >= 3 {
			crawler.isAlive = false
			crawler.urls = []string{"https://example.com"}
			messages <- crawler
			break
		}
	}
}

func handleChannel(messageChan chan Crawler, mapper map[string]bool) {
	// keep track of all the crawlers active -> turn off when no crawlers are active anymore
	// when no crawlers are active anymore -> send data to database
	for len(mapper) > 0 {
		// blocks until theres something in the channel
		msg := <-messageChan
		_, ok := mapper[msg.startURL]
		if !msg.isAlive {
			if ok { // still expecting the crawler to be alive
				fmt.Println("Crawler is not alive anymore for URL:", msg.startURL)
				for _, url := range msg.urls {
					mapper[url] = true
					go crawler(url, messageChan)
				}

				delete(mapper, msg.startURL)
			}
		} else {
			fmt.Println("Received message from crawler:", msg.startURL)
		}
	}
	fmt.Println("No more active crawlers. Exiting channel handler.")
	close(messageChan)
}

func callLLM(responseBody []byte) []string {
	client := openai.NewClient()
	context := context.Background()
	fmt.Println("Client Key:", client)
	chatCompletion, err := client.Chat.Completions.New(context, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Say this is a test")},
		Model:    openai.ChatModelGPT4o,
	})
	if err != nil {
		fmt.Println("Error creating chat completion:", err)
		return []string{}
	}
	fmt.Println("Chat Completion:", chatCompletion)
	return []string{}
}
