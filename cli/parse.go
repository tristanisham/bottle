package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tristanisham/bottle/server"
	"gopkg.in/yaml.v3"
)

func Parse(args []string) error {
	if len(args) == 0 {
		Help()
		return nil
	}

	switch args[0] {
	case "init", "reset":
		// Creates the default directories for the program to function.
		//	* posts
		//	* public

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := initializeProject(cwd); err != nil {
			return err
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
					return err
				}
				if err := initializeProject(args[1]); err != nil {
					return err
				}
			case "post":
				name := "my_new_post"
				if len(args) >= 3 {
					name = args[2]
				}

				post, err := createNewPost(name)
				if err != nil {
					return err
				}
				raw, err := yaml.Marshal(post)
				if err != nil {
					return err
				}

				buff := make([]byte, 0)
				buff = append(buff, []byte("---\n")...)
				buff = append(buff, raw...)
				buff = append(buff, []byte("---\n")...)

				if err := os.WriteFile(fmt.Sprintf("posts/%s.md", name), buff, 0775); err != nil {
					return err
				}

			case "nginx":
				if err := os.WriteFile("bottle.website", []byte(nginxService()), 0775); err != nil {
					log.Fatal(err)
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
		options := make([]string, 0)
		if len(args) >= 2 {
			options = append(options, args[1:]...)
		}
		multi_proc_server := false
		for _, opt := range options {
			if strings.Contains(opt, "-") {
				if strings.Contains(opt, "j") {
					multi_proc_server = true
				}
			}
		}

		server.Start(multi_proc_server)
	case "help", "--h", "-h", "version":
		Help()

	case "upgrade":
		Upgrade()
	}
	return nil
}
