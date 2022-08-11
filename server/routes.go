package server

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/feeds"
)

type FiberHandler func(c *fiber.Ctx) error

func FeedType(server *server) FiberHandler {
	return func(c *fiber.Ctx) error {
		dataType := url.PathEscape(c.Params("datatype", "rss"))
		now := time.Now()
		server.mu.Lock()
		defer server.mu.Unlock()
		feed := &feeds.Feed{
			Title:       server.settings.Title,
			Link:        &feeds.Link{Href: server.settings.Url},
			Description: server.settings.Description,
			Author:      &feeds.Author{Name: server.settings.Author, Email: server.settings.Email},
			Created:     now,
		}

		posts := server.router.GetPosts()
		for _, post := range posts {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: fmt.Sprintf("%s/%s", server.settings.Url, post.Slug)},
				Author:      &feeds.Author{Name: post.Author},
				Description: post.Description,
				Created:     post.PublishDate,
			})
		}

		atom, err := feed.ToAtom()
		if err != nil {
			return err
		}

		rss, err := feed.ToRss()
		if err != nil {
			return err
		}

		json, err := feed.ToJSON()
		if err != nil {
			return err
		}

		switch dataType {
		case "json":
			return c.SendString(json)
		case "atom":
			return c.SendString(atom)
		case "rss":
			return c.SendString(rss)
		default:
			return fmt.Errorf("invalid type. Acceptable types are: 'json', 'rss', and 'atom'")
		}
	}
}
