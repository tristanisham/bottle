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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
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

	engine := html.New(filepath.Join(cwd, "themes", server.settings.Theme, "templates"), ".html")
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
		AppName:      fmt.Sprint("Bottle ", "v0.0.7"),
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
				"Locale": server.settings.Locale,
				"FbUserId": server.settings.FbUserId,
				"FbAppId": server.settings.FbAppId,

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
		posts := make([]Post, 0)
		for _, post := range server.router {
			posts = append(posts, post)
		}

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].PublishDate.After(posts[j].PublishDate)
		})

		return c.Render("index", fiber.Map{
			"Title":       server.settings.Title,
			"Keywords":    server.settings.Keywords,
			"Description": server.settings.Description,
			"Posts":       posts,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
