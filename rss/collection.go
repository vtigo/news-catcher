package rss

import "encoding/json"

type FeedCollection struct {
	Items []Item
	Feeds []FeedMetadata
}

type FeedMetadata struct {
	Name  string
	Start int
	Count int
}

func NewFeedCollection() *FeedCollection {
	return &FeedCollection{
		Items: make([]Item, 0, 100),
		Feeds: make([]FeedMetadata, 0, 10),
	}
}

func (fc *FeedCollection) CollectDocument(document *RSSDocument) {
	start := len(fc.Items)
	fc.Items = append(fc.Items, document.Channel.Items...)

	fc.Feeds = append(fc.Feeds, FeedMetadata{
		Name:  document.Channel.Title,
		Start: start,
		Count: len(document.Channel.Items),
	})
}

func (fc *FeedCollection) MapItems() (map[string][]Item, error) {
	items := make(map[string][]Item, len(fc.Feeds))

	for _, feed := range fc.Feeds {
		items[feed.Name] = fc.Items[feed.Start : feed.Start+feed.Count]
	}

	return items, nil
}

func (fc *FeedCollection) MapItemsJSON() ([]byte, error) {
	items, err := fc.MapItems()
	if err != nil {
		return nil, err
	}

	json, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return nil, err
	}

	return json, nil
}
