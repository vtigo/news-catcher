# news-catcher

A terminal program for reading RSS feeds.

You add the address of a feed once, and the program downloads it, saves a copy
on your computer, and lets you browse the list of articles. Selecting an article
opens it in your web browser.

## Requirements

Go 1.26 or newer.

## Running

```
go run .
```

Or build it first:

```
go build
./news-catcher
```

## Using it

Feed list:

- `up`/`k` and `down`/`j` — move between feeds
- `enter` — open the selected feed
- `n` — add a new feed
- `r` — update the selected feed
- `R` — update all feeds
- `q` — quit

Inside a feed:

- `up`/`k` and `down`/`j` — move between articles
- `e` or `space` — open the article in the browser
- `esc`/`q` — go back to the feed list

## Stored files

Feeds and their downloaded articles are kept in the `storage` directory next to
the program. `storage/feeds.json` holds the list of feeds you added; each feed
also has its own `.xml` file with the last downloaded articles, so the list can
be read without an internet connection.

## License

MIT. See [LICENSE](LICENSE).
