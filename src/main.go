package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/datastx/FileTower/src/cli"
	"github.com/datastx/FileTower/src/tower"
	"github.com/fsnotify/fsnotify"
)

func main() {
	var cmds cli.CLI
	ctx := kong.Parse(&cmds)

	// Create a new watcher instance
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		ctx.FatalIfErrorf(err)
	}
	defer watcher.Close()

	// Handle errors and run your app
	if err := tower.Run(cmds, ctx, watcher); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
