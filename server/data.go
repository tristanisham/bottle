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
	Url         string   `json:"url" yaml:"url"`
	Locale      string   `json:"locale" yaml:"locale"`
	FbUserId    string   `json:"fb_user_id" yaml:"fb_user_id"`
	FbAppId     string   `json:"fb_app_id" yaml:"fb_app_id"`
	// TwPubHandle is your publication's Twitter handle
	TwPubHandle string `json:"tw_pub_handle" yaml:"tw_pub_handle"`
	// TwAuthorHandle is your publication's Twitter handle
	TwAuthorHandle string `json:"tw_author_handle" yaml:"tw_author_handle"`
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
	if len(post.TwAuthorHandle) == 0 {
		post.TwAuthorHandle = s.settings.TwAuthorHandle
	}
	if len(post.Author) == 0 {
		post.Author = s.settings.Author
	}

	post.ModifiedTime = time.Now()
	post.Slug = strings.Split(file.Name(), ".")[0]
	s.mu.Lock()
	s.router[file.Name()] = *post
	s.mu.Unlock()
}

type Post struct {
	Title          string    `json:"title" yaml:"title"`
	Subtitle       string    `json:"subtitle" yaml:"subtitle"`
	Author         string    `json:"author" yaml:"author"`
	Description    string    `json:"description" yaml:"description"`
	Keywords       []string  `json:"keywords" yaml:"keywords"`
	SiteImg        string    `json:"site_img" yaml:"site_img"`
	Body           string    `json:"body" yaml:"body"`
	Slug           string    `json:"slug" yaml:"slug"`
	PublishDate    time.Time `json:"publish_date" yaml:"publish_date"`
	OgVideo        string    `json:"og_video" yaml:"og_video"`
	Section        string    `json:"section" yaml:"section"`
	ModifiedTime   time.Time `json:"modifiedTime" yaml:"modifiedTime"`
	ExpirationTime time.Time `json:"expirationTime" yaml:"expirationTime"`
	// TwAuthorHandle is your publication's Twitter handle
	TwAuthorHandle string `json:"tw_author_handle" yaml:"tw_author_handle"`
	// TwPreviewImage defaults to SiteImg if not specified
	TwPreviewImage string `json:"tw_preview_image" yaml:"tw_preview_image"`
	TwVidAudPlayer string `json:"tw_vid_aud_player" yaml:"tw_vid_aud_player"`

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

	if len(result.TwPreviewImage) == 0 {
		result.TwPreviewImage = result.SiteImg
	}

	if result.PublishDate.IsZero() {
		result.PublishDate = time.Now()
	}
	
	return &result, nil
}
