package tower

import (
	"log"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// TODO: add support for symbolic links
func Run(directories []string, watcher *fsnotify.Watcher) {
	new := color.New(color.FgBlue).SprintFunc()
	modified := color.New(color.FgYellow).SprintFunc()
	deleted := color.New(color.FgRed).SprintFunc()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println(new("File created:", event.Name))
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println(modified("File modified:", event.Name))
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println(deleted("File removed:", event.Name))
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println(deleted("File renamed:", event.Name))
				}
				// TODO: I think we don't want this
				// if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				// 	log.Println(modified("File permissions modified:", event.Name))
				// }
			case err := <-watcher.Errors:
				log.Fatalln("Error:", err)
			}
		}
	}()
	for _, directroy := range directories {
		watcher.Add(directroy)
		log.Printf("watching %s", directroy)
	}
	<-done
}
