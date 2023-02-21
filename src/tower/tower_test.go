package tower

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/datastx/FileTower/src/schemas"
	"github.com/fsnotify/fsnotify"
)

func TestRun(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "tower-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	ch := make(chan schemas.Record)

	done := make(chan bool)
	go func() {
		Run([]string{tempDir}, watcher, ch)
		done <- true
	}()

	// create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := ioutil.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// wait for events to be processed
	expected := []schemas.Record{
		{Operation: fsnotify.Create, FileName: testFile},
	}
	for i := 0; i < len(expected); i++ {
		select {
		case actual := <-ch:
			if actual != expected[i] {
				t.Errorf("Unexpected record received: got %v, want %v", actual, expected[i])
			}
		}
	}

	// write to the test file
	if err := ioutil.WriteFile(testFile, []byte("test2"), 0644); err != nil {
		t.Fatal(err)
	}

	// wait for events to be processed
	expected = append(expected, schemas.Record{Operation: fsnotify.Write, FileName: testFile})
	for i := 1; i < len(expected); i++ {
		select {
		case actual := <-ch:
			if actual != expected[i] {
				t.Errorf("Unexpected record received: got %v, want %v", actual, expected[i])
			}
		}
	}

	// remove the test file
	if err := os.Remove(testFile); err != nil {
		t.Fatal(err)
	}

	// wait for events to be processed
	expected = append(expected, schemas.Record{Operation: fsnotify.Remove, FileName: testFile})
	for i := 2; i < len(expected); i++ {
		select {
		case actual := <-ch:
			if actual != expected[i] {
				t.Errorf("Unexpected record received: got %v, want %v", actual, expected[i])
			}
		}
	}

	// stop the watcher
	watcher.Close()
	<-done
}
