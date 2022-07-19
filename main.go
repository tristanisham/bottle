package main

import (
	"bottle/server"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help()
	}

	switch args[0] {
	case "init":
		// Creates the default directories for the program to function.
		//	* posts
		//	* public

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		if err := initializeProject(cwd); err != nil {
			log.Fatalln(err)
		}

	// New is a catchall command for generating basic templates and folder structures that a typical user may not remember.
	// The command works $ bottle new <arg>. Ideally, supported args will be limited but have a high utility:
	//	* post
	//	* project
	case "new":
		if len(args) >= 2 {
			switch args[1] {
			default:
				if err := os.MkdirAll(args[1], 0775); err != nil {
					log.Fatalln(err)
				}
				if err := initializeProject(args[1]); err != nil {
					log.Fatalln(err)
				}
			}

		}
		// TODO: create configuration file and schema
	// case "build", "brew", "b":
	// Iterate through ever md file in /posts and render into HTML (memory).
	// If -o flag is supplied dump out to specified directory.
	case "serve":
		// Starts webserver using $BOTTLE_PORT, config.Port, or 8080 in that order.
		// Watches for file changes so users don't have to restart their blog everytime they add a new post.
		server.Start()
	case "help", "--h", "-h":
		help()
	}

}

func help() {
	fmt.Println("$ bottle v0.0.1")
	fmt.Println("init: creates a new website in your current directory")
	fmt.Println("serve: starts serving your website")
	fmt.Println("new: create a project in the specified path.")
}

func initializeProject(path string) error {
	settings, err := yaml.Marshal(server.Settings{
		Title:       "99 Bottles of ...",
		Author:      "Brewmaster",
		Description: "A blog by Brewmaster",
		Keywords:    []string{"blog", "bottle", "seo", "brewing"},
		Theme:       "default",
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

	if err := os.WriteFile(filepath.Join(path, "bottle.yaml"), settings, 0775); err != nil {
		return err
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
	return nil
}
