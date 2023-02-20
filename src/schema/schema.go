package schema

import (
	"github.com/fsnotify/fsnotify"
)

type Record struct {
	Operation fsnotify.Op
	FileName  string
}
