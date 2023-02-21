package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/datastx/FileTower/src/cli"
	"github.com/datastx/FileTower/src/config"
	"github.com/datastx/FileTower/src/filehash"
	"github.com/datastx/FileTower/src/schema"
	"github.com/datastx/FileTower/src/tower"
	"github.com/fsnotify/fsnotify"
)

// TODO: discuss front loading cached files?
// TODO: crash if branch from git changes?

func main() {
	var cmds cli.CLI
	ctx := kong.Parse(&cmds)
	config := config.GetConfig(cmds.Config)
	// Create a new watcher instance
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error creating watcher:", err)
		ctx.FatalIfErrorf(err)
	}
	defer watcher.Close()
	dirs := GetDirectories(cmds)
	ch := make(chan schema.Record)
	go tower.Run(dirs, watcher, ch)
	ShipFile(ch, config.Server.IntervalAmount)
}

func GetDirectories(cmd cli.CLI) []string {
	var directories []string
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
			directories = append(directories, path)
			log.Printf("Added directory to watch list  `%s`", path)
		}

		return nil
	})

	return directories
}

func ShipFile(ch <-chan schema.Record, secondsSleep int) {
	resetTime := time.Now()
	// TODO: Add locks on access to the maps?...records and lastProcessed
	// Stores the File Name and the last action we saw and gets reset every interval
	var records = make(map[string]schema.Record)
	// Stores the File Name and the last hash we sent
	var lastProcessed = make(map[string]string)
	interval := time.Duration(secondsSleep) * time.Second

	for {
		select {
		case rec, ok := <-ch:
			if !ok {
				log.Println("Channel closed")
				return
			}
			// We always want the last action on a file
			records[rec.FileName] = rec
		default:
			if time.Since(resetTime) >= interval {
				time.Sleep(interval)
				for _, record := range records {
					// TODO: parallelize this?
					// This protects against the case where we don't have a file to read
					// a current hash for
					if record.Operation == fsnotify.Remove || record.Operation == fsnotify.Rename {
						log.Printf("Removing file %s and operation: %s", record.FileName, record.Operation)
						continue
					}
					currentHash := filehash.GetCheckSum(record.FileName)
					if prevHash, ok := lastProcessed[record.FileName]; !ok {
						log.Printf("Sending File %s we haven't sent before on operation %s", record.FileName, record.Operation)
						lastProcessed[record.FileName] = currentHash
						continue
					} else if prevHash != currentHash {
						log.Printf("Sending File %s we have seen before but has been modified on operation %s", record.FileName, record.Operation)
						lastProcessed[record.FileName] = currentHash
						continue
					}
					log.Printf("File %s has not changed, skipping", record.FileName)

				}
				resetTime = time.Now()
				records = make(map[string]schema.Record)
			}
		}
	}
}
