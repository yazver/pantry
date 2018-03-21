package pantry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// UnsupportedConfigError denotes encountering an unsupported
// configuration filetype.
type UnsupportedConfigError string

// Returns the formatted configuration error.
func (e UnsupportedConfigError) Error() string {
	return fmt.Sprintf("Unsupported Config Type %q", string(e))
}

// ConfigParseError denotes failing to parse configuration file.
type ConfigParseError struct {
	err error
}

// Returns the formatted configuration error.
func (e ConfigParseError) Error() string {
	return fmt.Sprintf("Parsing config failed: %s", e.err.Error())
}

// ConfigEncodeError denotes failing to encode configuration
type ConfigEncodeError struct {
	err error
}

// Returns the formatted configuration error.
func (e ConfigEncodeError) Error() string {
	return fmt.Sprintf("Encoding config failed: %s", e.err.Error())
}

// MarshalFunc is any marshaler.
type MarshalFunc func(v interface{}) ([]byte, error)

// UnmarshalFunc is any unmarshaler.
type UnmarshalFunc func(data []byte, v interface{}) error

type ConfigFormat struct {
	Marshal   MarshalFunc
	Unmarshal UnmarshalFunc
}

type ConfigFormats map[string]*ConfigFormat

// Register format marshaler and unmarphaler
func (c ConfigFormats) Register(format string, m MarshalFunc, um UnmarshalFunc) {
	c[strings.ToLower(format)] = &ConfigFormat{m, um}
}

func (c ConfigFormats) Search(format string) (*ConfigFormat, error) {
	format = strings.ToLower(format)
	f, ok := c[format]
	if !ok {
		if f, ok = c[filepath.Ext(format)]; !ok {
			return f, UnsupportedConfigError(format)
		}
	}
	return f, nil
}

// Formats contains marshalers and unmarshalers for different file formats
var Formats = ConfigFormats{}

func init() {

	Formats.Register("json", func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", "    ")
	}, json.Unmarshal)

	Formats.Register("yaml", yaml.Marshal, yaml.Unmarshal)
	Formats.Register("yml", yaml.Marshal, yaml.Unmarshal)

	Formats.Register("toml", func(v interface{}) ([]byte, error) {
		b := bytes.Buffer{}
		err := toml.NewEncoder(&b).Encode(v)
		return b.Bytes(), err
	}, toml.Unmarshal)

}
