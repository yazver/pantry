package pantry

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// File
type Box interface {
	Path() string
	Get() ([]byte, error)
	Set(data []byte) error
	Watch(task func(err error)) error
	UnWatch()
}

type BaseBox struct {
	path string
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning after the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func (bb *BaseBox) Ext() string {
	return ext(bb.path)
}

func (bb *BaseBox) Path() string {
	return bb.path
}

// FileBox
type FileBox struct {
	BaseBox
	lock      sync.Mutex
	watchDone chan struct{}
}

func NewFileBox(path string) Box {
	return &FileBox{BaseBox: BaseBox{path}}
}

func (fb *FileBox) Get() ([]byte, error) {
	return ioutil.ReadFile(fb.path)
}

func (fb *FileBox) Set(data []byte) error {
	return ioutil.WriteFile(fb.path, data, os.ModePerm)
}

func (fb *FileBox) Watch(task func(err error)) error {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	fb.UnWatch()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	filePath := filepath.Clean(fb.path)
	fileDir, _ := filepath.Split(filePath)
	if err := watcher.Add(fileDir); err != nil {
		return err
	}

	done := make(chan struct{})
	fb.watchDone = done

	go func() {
		defer watcher.Close()

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if filepath.Clean(event.Name) == filePath {
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
							task(nil)
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					task(err)
				}
			}
		}()

		<-done
	}()
	return nil
}

func (fb *FileBox) UnWatch() {
	fb.lock.Lock()
	if fb.watchDone != nil {
		close(fb.watchDone)
		fb.watchDone = nil
	}
	defer fb.lock.Unlock()
}
