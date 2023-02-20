package tower

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// TODO: add support for symbolic links
func Run(directories []string, watcher *fsnotify.Watcher) {

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("File created:", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File modified:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("File removed:", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("File renamed:", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("File permissions modified:", event.Name)
				}
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
