package pantry

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func stringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ext returns the file name extension used by path.
// The extension is the suffix beginning after the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/' && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			if (i + 1) < len(path) {
				return path[i+1:]
			}
		}
	}
	return ""
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func getHomeDir() (string, error) {
	var homeDir string

	switch runtime.GOOS {
	case "windows":
		homeDir = os.Getenv("USERPROFILE")
	default:
		homeDir = os.Getenv("HOME")
	}

	if homeDir == "" {
		return "", errors.New("No home directory found - set $HOME (or the platform equivalent).")
	}

	return homeDir, nil
}

func getHomeSubDir(subDirs ...string) (dir string, err error) {
	if dir, err = getHomeDir(); err != nil {
		return "", err
	}
	dir = filepath.Join(dir, filepath.Join(subDirs...))
	return
}

func getConfigDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("LOCALAPPDATA"), nil
	case "darwin":
		return getHomeSubDir("Library", "Application Support")
	default:
		return getHomeSubDir(".config")
	}
}

func getAppConfigDir(appName string) (dir string, err error) {
	if dir, err = getConfigDir(); err != nil {
		return "", err
	}
	dir = filepath.Join(dir, appName)
	return
}

func getAppDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// accumulate signals
func unjitter(in <-chan bool, d time.Duration) <-chan bool {
	out := make(chan bool)

	go func() {
		defer close(out)

		timer := time.NewTimer(time.Hour * 100000)
		defer timer.Stop()
		for {
			changed := false
			resetTimer := false
			select {
			case <-timer.C:
				changed = false
				out <- true
			case _, ok := <-in:
				if ok {
					if !changed {
						changed = true
						resetTimer = true
					}
				} else {
					return
				}
			}
			if resetTimer {
				timer.Stop()
				changed = true
				timer.Reset(time.Second)
				resetTimer = false
			}
		}
	}()
	return out
}

func unjitterFunc(in <-chan bool, d time.Duration, f func()) {
	out := unjitter(in, d)
	go func() {
		for _ = range out {
			f()
		}
	}()
}
