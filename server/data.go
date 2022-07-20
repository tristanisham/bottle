package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

func LoadSettings() (*Settings, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(filepath.Join(cwd, "bottle.yaml"))
	if err != nil {
		return nil, err
	}

	settings := Settings{}

	if err := yaml.Unmarshal(raw, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

type server struct {
	settings Settings
	router   map[string]Post
	mu       sync.Mutex
}

func (s *server) render(file watcher.Event) {
	post, err := postFromFile(file.Path)
	if err != nil {
		if err.Error() == "fde" {
			s.mu.Lock()
			delete(s.router, file.Name())
			s.mu.Unlock()
			return
		} else {
			log.Panicln(err)
		}
	}
	post.Slug = strings.Split(file.Name(), ".")[0]
	s.mu.Lock()
	s.router[file.Name()] = *post
	s.mu.Unlock()
}

type Post struct {
	Title       string    `json:"title" yaml:"title"`
	Subtitle    string    `json:"subtitle" yaml:"subtitle"`
	Author      string    `json:"author" yaml:"author"`
	Description string    `json:"description" yaml:"description"`
	Keywords    []string  `json:"keywords" yaml:"keywords"`
	Body        string    `json:"body" yaml:"body"`
	Slug        string    `json:"slug" yaml:"slug"`
	PublishDate time.Time `json:"publish_date" yaml:"publish_date"`
}

func postFromFile(path string) (*Post, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("fde")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := Post{}

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
		return nil, err
	}

	result.Body = string(blackfriday.Run([]byte(strings.Join(body, "\n")))) // Not safe by design to let user's hack

	return &result, nil
}
