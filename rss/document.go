package rss

import "encoding/xml"

type RSSDocument struct {
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
