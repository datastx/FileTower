package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/datastx/FileTower/src/cli"
	"github.com/datastx/FileTower/src/filehash"
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
		log.Println("Error creating watcher:", err)
		ctx.FatalIfErrorf(err)
	}
	defer watcher.Close()
	dirs := GetDirectories(cmds)
	ch := make(chan schema.Record)
	go tower.Run(dirs, watcher, ch)
	ShipFile(ch, 5)
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
	var records = make(map[string]schema.Record)
	var lastProcessed = make(map[string]bool)
	interval := time.Duration(secondsSleep) * time.Second

	for {
		select {
		case val, ok := <-ch:
			if !ok {
				log.Println("Channel closed")
				return
			}
			records = dedupeFiles(val, records)
		default:
			if time.Since(resetTime) >= interval {
				time.Sleep(interval)
				// TODO: Add Checksum logic for last write to file
				for _, record := range records {
					fileHash := filehash.GetCheckSum(record.FileName)
					if _, ok := lastProcessed[fileHash]; !ok {
						log.Printf("Sending file %s and operation: %s", record.FileName, record.Operation)
						lastProcessed[fileHash] = true
					} else {
						log.Printf("%q We have already processed that hash\n", record.FileName)
					}

				}
				resetTime = time.Now()
				records = make(map[string]schema.Record)
			}
		}
	}
}

// check if a string exists in the map
func dedupeFiles(record schema.Record, strMap map[string]schema.Record) map[string]schema.Record {
	if _, ok := strMap[record.FileName]; !ok {
		log.Printf("%q not found, adding to map\n", record.FileName)
		strMap[record.FileName] = record
	} else {
		log.Printf("%q already exists in map, skipping\n", record.FileName)
	}
	return strMap
}
