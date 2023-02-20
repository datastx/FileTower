package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/datastx/FileTower/src/cli"
	"github.com/datastx/FileTower/src/schema"
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
	dirs := GetDirectories(cmds)
	ch := make(chan schema.Record)
	go tower.Run(dirs, watcher, ch)
	ShipFile(ch)
}

func GetDirectories(cmd cli.CLI) []string {
	var directorys []string
	if cmd.Directory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("no directory specified, using current working directory %s", cwd)
		cmd.Directory = cwd
	}
	filepath.Walk(cmd.Directory, func(path string, info os.FileInfo, err error) error {
		// TODO: add support for symbolic links
		// TODO: add support for hidden directories
		if info.IsDir() && filepath.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}
		if info.IsDir() {
			directorys = append(directorys, path)
			log.Printf("Added directory to watch list  `%s`", path)
		}

		return nil
	})

	return directorys
}

func ShipFile(ch <-chan schema.Record) {
	resetTime := time.Now()

	for {
		select {
		case val, ok := <-ch:
			if !ok {
				log.Println("Channel closed")
				return
			}
			log.Printf("Got File %s and operation: %s", val.FileName, val.Operation)
			// TODO: add support for sending files to a remote server
		default:
			if time.Since(resetTime) >= 5*time.Second {
				log.Printf("No files to send, sleeping for 5 seconds")
				time.Sleep(5 * time.Second)
				resetTime = time.Now()
			}
		}
	}
}
