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
	"github.com/gofiber/template/html"
	"github.com/radovskyb/watcher"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Title       string   `json:"title" yaml:"title"`
	Author      string   `json:"author" yaml:"author"`
	Description string   `json:"description" yaml:"description"`
	Keywords    []string `json:"keywords" yaml:"keywords"`
	Theme       string   `json:"theme" yaml:"theme"`
}

type server struct {
	settings Settings
	router   map[string]post
	mu       sync.Mutex
}

func (s *server) render(file watcher.Event) {
	post, err := postFromFile(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	post.Slug = strings.Split(file.Name(), ".")[0]
	s.mu.Lock()
	s.router[file.Name()] = *post
	s.mu.Unlock()
}

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

type post struct {
	Title       string    `json:"title" yaml:"title"`
	Subtitle    string    `json:"subtitle" yaml:"subtitle"`
	Author      string    `json:"author" yaml:"author"`
	Description string    `json:"description" yaml:"description"`
	Keywords    []string  `json:"keywords" yaml:"keywords"`
	Body        string    `json:"body" yaml:"body"`
	Slug        string    `json:"slug" yaml:"slug"`
	PublishDate time.Time `json:"publish_date" yaml:"publish_date"`
}

func postFromFile(path string) (*post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := post{}

	frontMatter := make([]string, 0)
	var body []string

	header_closed := false
	lines := strings.Split(string(raw), "\n")
	if lines[0] == "---" {
		var j int
		for i := 1; !header_closed; i++ {
			if lines[i] == "---" {
				header_closed = true
				j = i
				break
			}
			frontMatter = append(frontMatter, lines[i])
		}

		body = lines[j+1:] // Plus 1 to cut out the secondary closing HR

	}

	metadata := strings.Join(frontMatter, "\n")
	if err := yaml.Unmarshal([]byte(metadata), &result); err != nil {
		log.Fatalln(err)
	}

	result.Body = string(blackfriday.Run([]byte(strings.Join(body, "\n")))) // Not safe by design to let user's hack

	return &result, nil
}

func Start() {

	server := server{
		settings: Settings{},
		router:   map[string]post{},
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

	app := fiber.New(fiber.Config{
		Views: engine,
	})

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
			})
		}
		
		return c.SendStatus(404)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		posts := make([]post, 0)
		for _, post := range server.router {
			posts = append(posts, post)
		}

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].PublishDate.Before(posts[j].PublishDate)
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
