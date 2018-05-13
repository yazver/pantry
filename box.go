package pantry

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"context"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

// Box
type Box interface {
	Path() string
	Format() string
	Get() ([]byte, error)
	Set(data []byte) error
	Watch(task func(err error)) error
	WatchContext(ctx context.Context, task func(err error)) error
}

type BaseBox struct {
	path string
}

// Format returns the file name extension used by path.
// The extension is the suffix beginning after the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func (bb *BaseBox) Format() string {
	return ext(bb.path)
}

func (bb *BaseBox) Path() string {
	return bb.path
}

// FileBox
type FileBox struct {
	BaseBox
	//lock      sync.Mutex
	//watchDone chan struct{}
}

func NewFileBox(path string) *FileBox {
	return &FileBox{BaseBox: BaseBox{filepath.Clean(path)}}
}

func (fb *FileBox) Get() ([]byte, error) {
	return ioutil.ReadFile(fb.path)
}

func (fb *FileBox) Set(data []byte) error {
	if path := filepath.Dir(fb.path); !pathExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(fb.path, data, os.ModePerm)
}

func (fb *FileBox) Watch(task func(err error)) error {
	return fb.WatchContext(context.Background(), task)
}
func (fb *FileBox) WatchContext(ctx context.Context, task func(err error)) error {
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

	go func() {
		defer watcher.Close()

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		go func() {
			defer close(done)
			ch := make(chan bool)
			defer close(ch)
			unjitterFunc(ch, time.Millisecond*500, func() { task(nil) })

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if filepath.Clean(event.Name) == filePath {
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
							ch <- true //task(nil)
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					task(err)
				case <-ctx.Done():
					return
				}
			}
		}()

		<-done
	}()
	return nil
}

// URLBox
type URLBox struct {
	url, format string
}

func NewURLBox(url, format string) *URLBox {
	return &URLBox{url, format}
}

func (ub *URLBox) Format() string {
	if ub.format == "" {
		return ext(ub.url)
	}
	return ub.format
}

func (ub *URLBox) Path() string {
	return ub.url
}

func (ub *URLBox) Get() ([]byte, error) {
	resp, err := http.Get(ub.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (ub *URLBox) Set(data []byte) error {
	return nil
}

func (ub *URLBox) Watch(task func(err error)) error {
	return nil
}

func (ub *URLBox) WatchContext(ctx context.Context, task func(err error)) error {
	return nil
}
