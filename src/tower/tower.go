package tower

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/fsnotify/fsnotify"

	"github.com/datastx/FileTower/src/cli"
)

func Run(cmd cli.CLI, ctx *kong.Context, watcher *fsnotify.Watcher) error {
	// Add the directory to the watcher's watch list
	err := watcher.Add(cmd.Directory)
	if err != nil {
		fmt.Println("Error adding directory to watcher:", err)
		ctx.FatalIfErrorf(err)
	}

	// Start an infinite loop to wait for events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Check the event type
			switch event.Op {
			case fsnotify.Create:
				fmt.Println("File created:", event.Name)
			case fsnotify.Write:
				fmt.Println("File modified:", event.Name)
			case fsnotify.Remove:
				fmt.Println("File deleted:", event.Name)
			case fsnotify.Rename:
				fmt.Println("File renamed:", event.Name)
			case fsnotify.Chmod:
				fmt.Println("File permissions changed:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
			fmt.Println("Error watching directory:", err)
			ctx.FatalIfErrorf(err)
		}
	}
}
