package server

import (
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

// fswatch watches the post directory for file changes and rebuilds the in-memory database.
func fswatch(router *server) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Remove, watcher.Rename, watcher.Create, watcher.Move, watcher.Write)

	filter := regexp.MustCompile("^[a-z0-9_-].*?(.md)$")
	w.AddFilterHook(watcher.RegexFilterHook(filter, false))

	go func(router *server) {
		for {
			select {
			case event := <-w.Event:
				// fmt.Println(event)
				router.render(event)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}(router)

	if err := w.AddRecursive("./posts"); err != nil {
		log.Fatalln(err)
	}

	
	// Print a list of all of the files and folders currently
	// being watched and their paths.
	// for path, f := range w.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }


	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
