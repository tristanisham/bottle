package main

import (
	"log"
	"os"

	"git.sr.ht/~atalocke/bottle/cli"
)

func main() {
	if err := cli.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
