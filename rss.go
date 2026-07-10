package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
)

type rssDocument struct {
	XMLName xml.Name  `xml:"rss"`
	Items   []RssItem `xml:"channel>item"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

type RssFeed struct {
	Name     string
	Endpoint string
	Items    []RssItem
	Fetched  bool
}

func RssFeedInit(name, endpoint string) (*RssFeed, error) {
	feed := &RssFeed{
		Name:     name,
		Endpoint: endpoint,
		Items:    nil,
		Fetched:  false,
	}

	// TODO: check if the file was created today. if not, refetch.
	_, err := os.Stat(fmt.Sprintf("./storage/%s.xml", feed.Name))
	if errors.Is(err, os.ErrNotExist) {
		err = feed.Fetch()
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch feed %s: %w", feed.Name, err)
		}
		feed.Fetched = true
		if err = feed.Store(); err != nil {
			return nil, fmt.Errorf("Failed to store feed %s: %w", feed.Name, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Failed to check ./storage/%s.xml info: %w", feed.Name, err)
	}

	err = feed.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load feed %s: %w", feed.Name, err)
	}

	return feed, nil
}

// Fetch sends a get request to the RssFeed defined endpoint and
// parses the request into RssItems. It populates RssFeed.Items.
func (f *RssFeed) Fetch() error {
	resp, err := http.Get(f.Endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var doc rssDocument
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&doc); err != nil {
		return err
	}

	f.Items = doc.Items
	return nil
}

func (f *RssFeed) Store() error {
	// TODO: salvar em caminho absoluto
	document := rssDocument{
		XMLName: xml.Name{},
		Items:   f.Items,
	}
	data, err := xml.Marshal(document)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("./storage", 0755); err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("./storage/%s.xml", f.Name), data, 0644)
}

func (f *RssFeed) Load() error {
	// TODO: carregar o struct à partir do arquivo
	// gerado por f.Store()
	file, err := os.Open(fmt.Sprintf("./storage/%s.xml", f.Name))
	if err != nil {
		return err
	}
	var doc rssDocument
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&doc); err != nil {
		return err
	}
	f.Items = doc.Items
	return nil
}
