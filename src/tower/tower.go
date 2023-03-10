package tower

import (
	"log"

	"github.com/datastx/FileTower/src/schemas"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// TODO: add support for symbolic links
func Run(directories []string, watcher *fsnotify.Watcher, ch chan<- schemas.Record) {
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
					ch <- schemas.Record{Operation: fsnotify.Create, FileName: event.Name}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println(modified("File modified:", event.Name))
					ch <- schemas.Record{Operation: fsnotify.Write, FileName: event.Name}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println(deleted("File removed:", event.Name))
					ch <- schemas.Record{Operation: fsnotify.Remove, FileName: event.Name}
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println(deleted("File renamed:", event.Name))
					ch <- schemas.Record{Operation: fsnotify.Rename, FileName: event.Name}
				}
				// TODO: I don't think we care about this. This proves out a race condition maybe talk to Nik
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println(modified("File permissions modified:", event.Name))
					ch <- schemas.Record{Operation: fsnotify.Chmod, FileName: event.Name}
				} //
			case err := <-watcher.Errors:
				log.Fatalln("Error:", err)
				close(ch)
			}
		}
	}()
	// You have to add the directories you want to watch
	// after starting the watcher.
	for _, directroy := range directories {
		watcher.Add(directroy)
		log.Printf("watching %s", directroy)
	}
	<-done
}
