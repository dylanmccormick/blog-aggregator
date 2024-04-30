package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dylanmccormick/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	log.Printf("Fetching posts for url: %v", url)
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
				v := fetchDataFromFeed(feed.Url)

				for _, post := range v.Channel.Items {
					input := database.CreatePostParams{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Title:     post.Title,
						Url:       post.Link,
						Description: sql.NullString{
							String: post.Description,
							Valid:  (post.Description != ""),
						},
						PublishedAt: sql.NullTime{},
						FeedID:      feed.ID,
					}
					cfg.DB.CreatePost(context.TODO(), input)

				}
				log.Printf("Fetching feed from %v\n", feed.Url)
			}()
		}

		wg.Wait()

		time.Sleep(300 * time.Second)
	}
}
