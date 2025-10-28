package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/vtigo/news-catcher/fetcher"
	"github.com/vtigo/news-catcher/rss"
)

// TODO: pegar o endpoint do arquivo de config
// TODO: pegar múltiplas fontes

func main() {
	endpoints := []string{
		"https://revistacult.uol.com.br/feed/",
		"https://feeds.folha.uol.com.br/poder/rss091.xml",
		"https://catracalivre.com.br/feed/",
	}

	fetcher := fetcher.NewClient(
		fetcher.WithTimeout(10*time.Second),
		fetcher.WithMaxBytes(5*1024*1024),
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
			rss.CollectDocument(collection, *feed)
		}
	}

	close(results)

	// Convertemos os itens em um slice de bytes que representa uma string JSON formatada
	jsonData, err := json.MarshalIndent(collection.Items, "", " ")
	if err != nil {
		log.Fatalln("failed to marshal feed data into json: ", err)
	}

	// Gravamos os bytes em um arquivo json
	err = os.WriteFile("feed-items.json", jsonData, 0644)
	if err != nil {
		log.Fatalln("failed to write json data: ", err)
	}

}
