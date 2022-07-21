package cli

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.sr.ht/~atalocke/bottle/server"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed client/post.html
	post_html []byte
	//go:embed client/index.html
	index_html []byte
	//go:embed client/nav.html
	nav_html []byte
	//go:embed client/style.css
	client_css []byte
	//go:embed client/error.html
	error_html []byte
)

func initializeProject(path string) error {
	settings, err := yaml.Marshal(server.Settings{
		Title:       "99 Bottles of ...",
		Author:      "Brewmaster",
		Description: "A blog by Brewmaster",
		Keywords:    []string{"blog", "bottle", "seo", "brewing"},
		Theme:       "default",
		Locale: "en_us",
	})

	if err != nil {
		return err
	}

	paths := []string{"posts", "public", "themes"}

	for i := range paths {
		os.Mkdir(filepath.Join(path, paths[i]), 0775) // creates each directory in paths
	}

	if err := createDefaultTheme(filepath.Join(path, "themes", "default")); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(path, "bottle.yaml")); os.IsNotExist(err) {
		if err := os.WriteFile(filepath.Join(path, "bottle.yaml"), settings, 0775); err != nil {
			return err
		}
	}

	return nil
}

// createDefaultTheme creates the default theme folder on the filesystem.
// Used in `$ bottle init` calls.
func createDefaultTheme(path string) error {
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}

	templates_dir := filepath.Join(path, "templates")

	for _, i := range []string{"templates", "public"} {
		if err := os.MkdirAll(filepath.Join(path, i), 0775); err != nil {
			return err
		}
	}

	if err := os.WriteFile(filepath.Join(path, "public", "style.css"), client_css, 0775); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(templates_dir, "post.html"), post_html, 0775); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(templates_dir, "index.html"), index_html, 0775); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(templates_dir, "nav.html"), nav_html, 0775); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(templates_dir, "error.html"), error_html, 0775); err != nil {
		return err
	}

	return nil
}

func createNewPost(name string) (*server.Post, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	posts := filepath.Join(cwd, "posts")
	if _, err := os.Stat(posts); os.IsNotExist(err) {
		return nil, fmt.Errorf("posts directory not detected. Try running $ bottle init")
	}

	settings, err := server.LoadSettings()
	if err != nil {
		return nil, err
	}

	post := server.Post{
		Title:       "My new post",
		Subtitle:    "A blank canvas--full of adventure",
		Author:      settings.Author,
		Slug:        name,
		PublishDate: time.Now(),
	}

	return &post, nil
}
