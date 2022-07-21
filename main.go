package main

import (
	"log"
	"os"

	"github.com/tristanisham/bottle/cli"
)

func main() {
	if err := cli.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
