package tower

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/datastx/FileTower/src/cli"
)

// TODO: add support for symbolic links
func Run(cmd cli.CLI, watcher *fsnotify.Watcher) {

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("File created:", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("File removed:", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("File renamed:", event.Name)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("File permissions modified:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()
	// TODO: Change this logic
	if cmd.Directory == "" {
		cmd.Directory = "/Users/brianmoore/githib.com/datastx/FileTower/src"
	}
	filepath.Walk(cmd.Directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			watcher.Add(path)
			log.Println("Added:", path)
		}
		return nil
	})
	<-done
}
