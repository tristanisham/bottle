package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	"github.com/gorilla/feeds"
	"github.com/tristanisham/bottle/utils"
	"gopkg.in/yaml.v3"
)

// Build creates the inital map of existing posts before watcher takes over to handle individual changes.
func (s *server) Build(path string) error {
	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, _ error) error {
		if ok, err := regexp.MatchString("^[a-z0-9_-].*?(.md)$", d.Name()); ok {
			post, err := postFromFile(path)
			if err != nil {
				log.Fatal(err)
			}
			post.Slug = strings.Split(d.Name(), ".")[0]
			s.mu.Lock()
			s.router[d.Name()] = *post
			s.mu.Unlock()

		} else {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func Start(multiProc bool) {

	server := server{
		settings: Settings{},
		router:   map[string]Post{},
		mu:       sync.Mutex{},
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	configFile, err := os.ReadFile(filepath.Join(cwd, "bottle.yaml"))
	if err != nil {
		log.Fatalln(err)
	}

	if err := yaml.Unmarshal(configFile, &server.settings); err != nil {
		log.Panicln(err) // Need to add recursive watch to remove necessary server restart.
	}

	server.Build(filepath.Join(cwd, "posts")) // need to change to pull from config file.

	go fswatch(&server)

	//   _    _      _       _____
	// | |  | |    | |     /  ___|
	// | |  | | ___| |__   \ `--.  ___ _ ____   _____ _ __
	// | |/\| |/ _ \ '_ \   `--. \/ _ \ '__\ \ / / _ \ '__|
	// \  /\  /  __/ |_) | /\__/ /  __/ |   \ V /  __/ |
	//  \/  \/ \___|_.__/  \____/ \___|_|    \_/ \___|_|

	engine := html.New(filepath.Join(cwd, "themes", server.settings.Theme, "templates"), ".html").Reload(true)
	engine.AddFunc("html", func(copy string) template.HTML {
		return template.HTML(copy)
	})
	engine.AddFunc("ymd", func(t time.Time) string {
		year, month, date := t.Date()
		return fmt.Sprintf("%d/%d/%d", year, month, date)
	})
	engine.AddFunc("join", func(in []string, sep string) string {
		return strings.Join(in, sep)
	})

	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "Bottle",
		AppName:      fmt.Sprint("Bottle ", utils.VERSION),
		Prefork:      multiProc,
	})

	// Middleware
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))
	app.Use(compress.New())

	// End middleware

	app.Static("/", filepath.Join(cwd, "public"))
	app.Static("/", filepath.Join(cwd, "themes", server.settings.Theme, "public"))

	app.Get("/feed", func(c *fiber.Ctx) error {
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

		headers := c.GetReqHeaders()
		switch headers["Accept"] {
		case "application/json":
			return c.JSON(json)
		case "application/rss+xml":
			return c.SendString(rss)
		case "application/atom+xml":
			return c.SendString(atom)
		}

		return c.SendString(rss)
	})

	app.Get("/feed.:datatype", FeedType(&server))

	app.Get("/:slug", func(c *fiber.Ctx) error {
		slug := url.PathEscape(c.Params("slug"))
		server.mu.Lock()
		defer server.mu.Unlock()

		if val, ok := server.router[fmt.Sprintf("%s.md", slug)]; ok {
			return c.Render("post", fiber.Map{
				"Title":       server.settings.Title,
				"Keywords":    strings.Join(server.settings.Keywords, ", "),
				"Description": server.settings.Description,
				"Author":      server.settings.Author,
				"Post":        val,
				"TwPubHandle": server.settings.TwPubHandle,
				"URL":         fmt.Sprintf("%s/%s", server.settings.Url, slug),
				"Locale":      server.settings.Locale,
				"FbUserId":    server.settings.FbUserId,
				"FbAppId":     server.settings.FbAppId,
			})
		}

		c.Status(404)
		return c.Render("error", fiber.Map{
			"Title":       server.settings.Title,
			"Keywords":    strings.Join(server.settings.Keywords, ", "),
			"Description": server.settings.Description,
			"Author":      server.settings.Author,
			"Error":       "404",
			"Message":     "Post not found",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		
		posts := server.router.GetPosts()

		return c.Render("index", fiber.Map{
			"Title":       server.settings.Title,
			"Keywords":    server.settings.Keywords,
			"Description": server.settings.Description,
			"Posts":       posts,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
