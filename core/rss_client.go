package core

import (
	"errors"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	"io"
	"log"
	"os"
	"time"
)

func PollFeed(uri string, timeout int, cr xmlx.CharsetFunc, out chan Record, limit int, source string) {
	feed := rss.New(timeout, true, chanHandler, makeItemHandler(out, limit, source))

	for {
		if err := feed.Fetch(uri, cr); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	log.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func makeItemHandler(out chan Record, limit int, source string) func(*rss.Feed, *rss.Channel, []*rss.Item) {
	return func(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
		log.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
		for i, v := range newitems {
			if i > limit {
				break
			}
			v := v
			body := v.Description
			if body == "" {
				body = v.Content.Text
			}
			link := feed.Url
			if v.Links != nil && len(v.Links) > 0 {
				link = v.Links[0].Href
			}
			defer func() { out <- Record{v.Title, body, source, link, v.PubDate} }()
		}
	}
}

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "iso-8859-1" {
		return r, nil
	}
	return nil, errors.New("Unsupported character set encoding: " + charset)
}
