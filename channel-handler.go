package main

import "fmt"

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
					go crawler(url, "Find all the items on sale.", messageChan)
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
