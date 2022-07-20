package cli

import (
	"fmt"
)

func Help() {
	fmt.Println("$ bottle v0.0.1")
	fmt.Println("init: creates a new website in your current directory")
	fmt.Println("serve: starts serving your website")
	fmt.Println("new: create a project in the specified path.")
}
