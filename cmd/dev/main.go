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
	endpoint := "https://revistacult.uol.com.br/feed/"

	fetcher := fetcher.NewClient(
		fetcher.WithTimeout(15*time.Second),
		fetcher.WithMaxBytes(5*1024*1024),
	)

	// Um feed RSS é um documento XML, temos o xml bruto (bytes) no corpo da resposta
	ctx := context.Background()

	responseData, err := fetcher.Fetch(ctx, endpoint)
	if err != nil {
		log.Fatalln("failed to fetch feed: ", err)
	}

	feed, err := rss.UnmarshalXML(responseData)
	if err != nil {
		log.Fatalln("failed to unmarshall the RSS document: ", err)
	}

	// Agora que temos o feed em memória, podemos decidir o que fazer com os dados

	// Convertemos o feed em um slice de bytes que representa uma string JSON formatada
	jsonData, err := json.MarshalIndent(feed, "", " ")
	if err != nil {
		log.Fatalln("failed to marshal feed data into json: ", err)
	}

	// Gravamos os bytes em um arquivo json
	err = os.WriteFile("feed.json", jsonData, 0644)
	if err != nil {
		log.Fatalln("failed to write json data: ", err)
	}

}

// func fetchConcurrent(endpoints []string) {
//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//     defer cancel()

//     fetcher := fetcher.NewClient(fetcher.WithTimeout(15*time.Second))

//     // Launch goroutines for each feed
//     for _, endpoint := range endpoints {
//         go func(url string) {
//             responseData, err := fetcher.Fetch(ctx, url)
//             if err != nil {
//                 log.Printf("failed to fetch %s: %v", url, err)
//                 return
//             }
//             // Process feed...
//         }(endpoint)
//     }

//     // Wait for completion or timeout...
// }
