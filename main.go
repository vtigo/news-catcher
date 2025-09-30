package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vtigo/news-catcher/internal/config"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title        string `xml:"title"`
	Link         string `xml:"link"`
	Description  string `xml:"description"`
	Language     string `xml:"language"`
	Copyright    string `xml:"copyright"`
	Docs         string `xml:"docs"`
	WebMaster    string `xml:"webMaster"`
	CategorySlug string `xml:"categorySlug"`
	Image        Image  `xml:"image"`
	Items        []Item `xml:"item"`
}

type Image struct {
	Title       string `xml:"title"`
	URL         string `xml:"url"`
	Link        string `xml:"link"`
	Width       int    `xml:"width"`
	Height      int    `xml:"height"`
	Description string `xml:"description"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalln("failed to load config file: ", err)
	}
	fmt.Println(config.Sources)

	endpoint := "https://catracalivre.com.br/feed/"

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalln("failed to get: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("failed to read response body: ", err)

	}

	// Create a decoder with charset reader
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charsetReader

	var rss RSS
	err = decoder.Decode(&rss)
	if err != nil {
		log.Fatalln("failed to unmarshall the XML response:", err)
	}

	jsonData, err := json.MarshalIndent(rss, "", " ")
	if err != nil {
		log.Fatalln("failed to marshall the RSS object into json:", err)
	}

	os.WriteFile("feed.json", jsonData, 0644)
}

// charsetReader handles different character encodings
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "iso-8859-1", "latin1":
		// Convert ISO-8859-1 to UTF-8
		buf := new(bytes.Buffer)
		for {
			b := make([]byte, 1)
			_, err := input.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			// ISO-8859-1 maps directly to Unicode code points 0-255
			if b[0] < 128 {
				buf.WriteByte(b[0])
			} else {
				buf.WriteString(string(rune(b[0])))
			}
		}
		return buf, nil
	default:
		return input, nil
	}
}
