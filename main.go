package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"time"
)

const wikiUrl = "https://en.wikipedia.org"

func main() {
	startTopic := os.Args[1:][0]
	endTopic := os.Args[1:][1]

	chUrls := make(chan string)

	start := time.Now()
	go crawl("/wiki/"+startTopic, chUrls)

	for topic := range chUrls {
		if topic == "/wiki/"+endTopic {
			elapsed := time.Since(start)
			fmt.Printf("Elapsed time: %v\n", elapsed)
			close(chUrls)
		}
		go crawl(topic, chUrls)
	}
}

func crawl(topic string, ch chan string) {
	url := wikiUrl + topic
	res, err := http.Get(url)
	if err != nil {
		return
	}

	b := res.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data != "a" {
				continue
			}

			ok, url := getHref(t)
			if !ok {
				continue
			}

			isWiki := strings.Contains(url, "/wiki/")
			hasProto := strings.Contains(url, "http")
			if isWiki && !hasProto {
				ch <- url
			}
		}
	}
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}
