package pantry

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	LocationConfigDir      = "${configdir}"
	LocationApplicationDir = "${appdir}"
	LocationCurrentDir     = "${curdir}"
	LocationHomeDir        = "${home}"
)

// UnfoundLocation denotes that location is not found.
type UnfoundLocation string

// Returns the formatted configuration error.
func (str UnfoundLocation) Error() string {
	return fmt.Sprintf("Location is not found: %q", string(str))
}

func IsUnfoundLocation(err error) bool {
	_, ok := err.(UnfoundLocation)
	return ok
}

type Locations struct {
	baseLocations   map[string]string
	ApplicationName string
	List            []string
}

func (l *Locations) Init(applicationName string, locations ...string) *Locations {
	if applicationName == "" {
		panic("Plantry: application name not defined.")
	}

	l.baseLocations = map[string]string{}
	if path, err := getAppConfigDir(applicationName); err == nil {
		l.baseLocations[LocationConfigDir] = path
	}
	if path, err := getAppDir(); err == nil {
		l.baseLocations[LocationApplicationDir] = path
	}
	if path, err := os.Getwd(); err == nil {
		l.baseLocations[LocationCurrentDir] = path
	}
	if path, err := getHomeDir(); err == nil {
		l.baseLocations[LocationHomeDir] = path
	}

	l.ApplicationName = applicationName
	l.Add(locations...)
	return l
}

func NewLocations(applicationName string, locations ...string) *Locations {
	return new(Locations).Init(applicationName, locations...)
}

func (l *Locations) ExpandLocation(path string) string {
	// path = filepath.FromSlash(path)
	// for key, location := range l.baseLocations {
	// 	path = strings.Replace(path, key, location, -1)
	// }
	// if strings.HasPrefix(path, "$") {
	// 	if end := strings.Index(path, string(os.PathSeparator)); end != -1 {
	// 		path = os.Getenv(path[1:end]) + path[end:]
	// 	} else {
	// 		path = os.Getenv(path[1:])
	// 	}

	// }
	if len(path) >= 2 && path[:2] == "./" {
		path = filepath.Join(LocationCurrentDir, path[2:])
	}
	if len(path) >= 2 && path[:2] == "~/" {
		path = filepath.Join(LocationHomeDir, path[2:])
	}
	path = os.Expand(filepath.FromSlash(path), func(key string) string {
		return l.baseLocations["${"+key+"}"]
	})
	return os.ExpandEnv(path)
}

func (l *Locations) add(location string) {
	locations := filepath.SplitList(l.ExpandLocation(location))
	for _, location := range locations {
		if !stringInSlice(location, l.List) {
			l.List = append(l.List, location)
		}
	}
}

func (l *Locations) AddJoin(locationParts ...string) {
	l.add(path.Join(locationParts...))
}

func (l *Locations) Add(locations ...string) {
	for _, location := range locations {
		l.add(location)
	}
}

func (l *Locations) LocatePath(filename string) (string, error) {
	box, err := l.Locate(filename)
	if box != nil {
		return box.Path(), err
	}
	return "", err
}

func (l *Locations) Locate(filename string) (box Box, err error) {
	err = nil
	filename = l.ExpandLocation(filename)
	if strings.HasPrefix(filename, "https://") || strings.HasPrefix(filename, "http://") {
		return NewURLBox(filename, ext(filename)), nil
	}
	if filepath.IsAbs(filename) {
		box = NewFileBox(filename)
		if pathExists(filename) {
			return
		}
		return box, UnfoundLocation(filename)
	}
	for _, location := range l.List {
		if filePath := filepath.Join(location, filename); pathExists(filePath) {
			return NewFileBox(filePath), nil
		}
	}
	if len(l.List) > 0 {
		return NewFileBox(filepath.Join(l.List[0], filename)), UnfoundLocation(filename)
	}
	return nil, UnfoundLocation(filename)
}
