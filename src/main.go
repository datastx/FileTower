package main

import (
	"fmt"

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

	tower.Run(cmds, watcher)
}
