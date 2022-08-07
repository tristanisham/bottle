package cli

import (
	"fmt"

	"github.com/fatih/color"
)



func Help() {
	fmt.Printf("$ bottle %s\n", "v0.0.7")
	fmt.Printf("%s: creates a new website in your current directory.\n", color.GreenString("init"))
	fmt.Printf("%s: starts serving your website.\n", color.GreenString("serve"))
	fmt.Printf("%s %s: create a project in the specified path.\n", color.GreenString("new"), color.YellowString("$NAME"))
	fmt.Printf("\t- %s %s: creates a new post titled %s\n", color.GreenString("post"), color.YellowString("$NAME"), color.YellowString("$NAME"))
	fmt.Printf("\t- %s: creates a basic NGINX server configuration file\n", color.GreenString("nginx"))

}
