package tower

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fsnotify/fsnotify"

	"github.com/datastx/FileTower/src/cli"
)

func TestRun(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a new watcher instance
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Add the temporary directory to the watcher's watch list
	err = watcher.Add(dir)
	if err != nil {
		t.Fatalf("failed to add directory to watcher: %v", err)
	}

	ctx := &kong.Context{}

	// Start the run function in a goroutine
	errs := make(chan error, 1)
	go func(test *testing.T) {
		if err := run(cli.CLI{Directory: dir}, ctx, watcher); err != nil {
			test.Fatalf("run returned error: %v", err)
			errs <- err
		}
	}(t)

	// Create a new file in the temporary directory
	f, err := os.Create(filepath.Join(dir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	f.Close()

	// Wait for the event to be processed
	time.Sleep(time.Millisecond * 100)
}
