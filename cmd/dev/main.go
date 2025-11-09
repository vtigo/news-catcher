package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vtigo/news-catcher/config"
	"github.com/vtigo/news-catcher/fetcher"
	"github.com/vtigo/news-catcher/rss"
	"github.com/vtigo/news-catcher/storage"
)

func main() {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatalln("failed to load config: ", err)
	}

	endpoints := config.XMLEndpoints()

	fetcher := fetcher.NewClient(
		fetcher.WithMaxBytes(5 * 1024 * 1024),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make(chan *rss.RSSDocument, len(endpoints))

	for _, endpoint := range endpoints {
		go func(url string) {
			// Um feed RSS é um documento XML, temos o xml bruto (bytes) no corpo da resposta
			responseData, err := fetcher.Fetch(ctx, url)
			if err != nil {
				log.Printf("failed to fetch %s: %v\n", url, err)
				results <- nil
				return
			}

			// TODO: Checar se é um feed RSS antes de continuar

			feed, err := rss.UnmarshalXML(responseData)
			if err != nil {
				log.Printf("failed to unmarshall the RSS document for %s: %v\n", url, err)
				results <- nil
				return
			}

			results <- feed

		}(endpoint)
	}

	collection := rss.NewFeedCollection()
	for range endpoints {
		feed := <-results
		if feed != nil {
			collection.CollectDocument(feed)
		}
	}

	close(results)

	collectionJSON, err := collection.MapItemsJSON()
	if err != nil {
		log.Fatalln("failed to map items to json: ", err)
	}

	storage := storage.NewFileSystemStorage("data")

	now := time.Now()
	filename := fmt.Sprintf("feed-%v.json", now.Unix())

	filepath, err := storage.Store(filename, collectionJSON)
	if err != nil {
		log.Fatalln("failed to write json data: ", err)
	}

	log.Printf("saved at %s\n", filepath)
}
