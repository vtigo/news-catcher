package rss

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

func (fc *FeedCollection) CollectDocument(document RSSDocument) {
	start := len(fc.Items)
	fc.Items = append(fc.Items, document.Channel.Items...)

	fc.Feeds = append(fc.Feeds, FeedMetadata{
		Name:  document.Channel.Title,
		Start: start,
		Count: len(document.Channel.Items),
	})
}
