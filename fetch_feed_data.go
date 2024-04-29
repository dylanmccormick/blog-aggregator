package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func fetchDataFromFeed(url string) Rss {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting from URL: %v, ERR: %v\n", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR Unloading data from body: %v\n", err)
	}
	v := Rss{}
	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Fatalf("ERROR Unmarshaling %v\n", err)
	}
	fmt.Printf("BLOG NAME: %v", v.Channel.Title)
	for i, post := range v.Channel.Items {
		fmt.Printf("Number: %v, Title: %v\n", i, post.Title)
	}
	return v
}

func (cfg *apiConfig) feedWorker(n int) {

	for {
		feeds, err := cfg.DB.GetNextFeedsToFetch(context.TODO(), int32(n))
		if err != nil {
			log.Fatalf("ERROR Fetching feeds from DB: %v\n", err)
		}

		var wg sync.WaitGroup

		for _, feed := range feeds {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fetchDataFromFeed(feed.Url)
				log.Printf("Fetching feed from %v\n", feed.Url)
			}()
		}

		wg.Wait()

		time.Sleep(10 * time.Second)
	}
}
