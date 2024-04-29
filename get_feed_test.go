package main

import (
	"testing"
)

func TestGetFeeds(t *testing.T) {
	fetchDataFromFeed("https://blog.boot.dev/index.xml")
}
