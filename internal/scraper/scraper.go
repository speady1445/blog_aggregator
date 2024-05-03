package scraper

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/speady1445/blog_aggregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Desctiption string `xml:"description"`
		Items       []Post `xml:"item"`
	} `xml:"channel"`
}

type time_ time.Time

type Post struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     time_  `xml:"pubDate"`
	Description string `xml:"description"`
}

func (t *time_) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	parsed, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", s)
	if err != nil {
		return err
	}

	*t = time_(parsed)
	return nil
}

func fetchFeed(url_ string) (*RSSFeed, error) {
	res, err := http.Get(url_)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rss := RSSFeed{}

	if err := xml.NewDecoder(res.Body).Decode(&rss); err != nil {
		return nil, err
	}

	return &rss, nil
}

func Start(db *database.Queries, fetchLimit int, fetchInterval time.Duration) {
	log.Printf("Feed scraper started, scraping %d feeds every %s", fetchLimit, fetchInterval)
	ticker := time.NewTicker(fetchInterval)
	wg := &sync.WaitGroup{}

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(fetchLimit))
		if err != nil {
			log.Printf("Failed to fetch feeds: %s", err)
			continue
		}
		log.Printf("Fetching %d feeds", len(feeds))

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Failed to mark feed %s as fetched: %s", feed.Name, err)
		return
	}

	rss, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Failed to fetch feed %s: %s", feed.Url, err)
		return
	}

	for _, post := range rss.Channel.Items {
		savePost(db, feed, post)
	}
	log.Printf("Saved %d posts from %s", len(rss.Channel.Items), rss.Channel.Title)
}

func savePost(db *database.Queries, feed database.Feed, post Post) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error generating UUID:", err)
		return
	}

	now := time.Now().UTC()

	_, err = db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid,
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       post.Title,
		Url:         post.Link,
		Description: post.Description,
		PublishedAt: time.Time(post.PubDate),
		FeedID:      feed.ID,
	})
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		log.Printf("Could not create post: %s from feed %s with error: %s", post.Title, feed.Name, err)
	}
}
