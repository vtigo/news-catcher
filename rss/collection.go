package rss

type FeedCollector interface {
	Collect(document RSSDocument) (*FeedCollection, error)
}

type FeedCollection struct {
	Feeds []Feed
}

type Feed struct {
	Name  string
	Items []Item
}

func (fc *FeedCollection) Collect(document RSSDocument) (*FeedCollection, error) {
	var name string
	var items []Item

	name = document.Channel.Title
	items = append(items, document.Channel.Items...)

	fc.Feeds = append(fc.Feeds, Feed{
		Name:  name,
		Items: items,
	})

	return fc, nil
}

// TODO: Criar a coleção a partir de um documento RSS, com a função Collect podemos adicionar mais documentos a coleção
