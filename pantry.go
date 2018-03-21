package pantry

import (
	"io"
	"io/ioutil"

	"github.com/yazver/golibs/reflect"
)

//var (
//	ErrInvalid    = errors.New("invalid argument") // methods on File will return this error when the receiver is nil
//	ErrPermission = errors.New("permission denied")
//	ErrExist      = errors.New("file already exists")
//	ErrNotExist   = errors.New("file does not exist")
//)

// Tag options
type Tag struct {
	Use  bool   // Use tag
	Name string // Name of tag
}

// Tags options
type Tags struct {
	Default     Tag
	Config      Tag
	Description Tag
}

// Config processing options
type Options struct {
	Flags         Flags
	Enviropment   Enviropment
	Tags          Tags
	DefaultFormat string
}

// Global options used for initiasation
var DefaultOptions = &Options{
	Flags:       Flags{Using: FlagsDontUse},
	Enviropment: Enviropment{Use: true, Prefix: ""},
	Tags: Tags{
		Config:      Tag{true, "config"},
		Default:     Tag{true, "default"},
		Description: Tag{false, "description"},
	},
	DefaultFormat: "",
}

type Pantry struct {
	Locations *Locations
	Options   *Options
}

func (p *Pantry) Init(applicationName string, locations ...string) *Pantry {
	p.Locations = NewLocations(applicationName, locations...)
	p.Options = &Options{}
	*p.Options = *DefaultOptions
	return p
}

func NewPantry(applicationName string, locations ...string) *Pantry {
	return new(Pantry).Init(applicationName, locations...)
}

func (p *Pantry) AddLocation(locationParts ...string) {
	p.Locations.AddJoin(locationParts...)
}

func (p *Pantry) AddLocations(locations ...string) {
	p.Locations.Add(locations...)
}

func (p *Pantry) LocatePath(filename string) (string, error) {
	return p.Locations.LocatePath(filename)
}

func (p *Pantry) Locate(filename string) (Box, error) {
	return p.Locations.Locate(filename)
}

func (p *Pantry) searchFormat(s string) (*ConfigFormat, error) {
	f, err := Formats.Search(s)
	if err != nil && p.Options.DefaultFormat != "" {
		return Formats.Search(p.Options.DefaultFormat)
	}
	return f, err
}

func (p *Pantry) UnmarshalWith(b []byte, v interface{}, unmarshaler UnmarshalFunc) error {
	reflect.Clear(v)
	if p.Options.Tags.Default.Use {
		if err := unmarshaler(b, v); err != nil {
			return ConfigParseError{err}
		}
		if err := processDefaultValues(v, p.Options); err != nil {
			return err
		}
	}
	if err := unmarshaler(b, v); err != nil {
		return ConfigParseError{err}
	}
	return processTags(v, p.Options)
}

func (p *Pantry) MarshalWith(v interface{}, marshaler MarshalFunc) (b []byte, err error) {
	b, err = marshaler(v)
	if err != nil {
		return nil, ConfigEncodeError{err}
	}
	return
}

func (p *Pantry) Unmarshal(b []byte, v interface{}, format string) error {
	f, err := Formats.Search(format)
	if err != nil {
		return err
	}
	return p.UnmarshalWith(b, v, f.Unmarshal)
}

func (p *Pantry) Marshal(v interface{}, format string) (b []byte, err error) {
	f, err := Formats.Search(format)
	if err != nil {
		return nil, err
	}
	return p.MarshalWith(v, f.Marshal)
}

func (p *Pantry) Decode(r io.Reader, v interface{}, format string) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.Unmarshal(b, v, format)
}

func (p *Pantry) Encode(w io.Writer, v interface{}, format string) error {
	b, err := p.Marshal(v, format)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (p *Pantry) UnBoxWith(box Box, v interface{}, unmarshaler UnmarshalFunc) error {
	b, err := box.Get()
	if err != nil {
		return err
	}
	return p.UnmarshalWith(b, v, unmarshaler)
}

func (p *Pantry) BoxWith(box Box, v interface{}, marshaler MarshalFunc) error {
	b, err := p.MarshalWith(v, marshaler)
	if err != nil {
		return err
	}
	return box.Set(b)
}

func (p *Pantry) LoadWith(path string, v interface{}, unmarshaler UnmarshalFunc) (string, error) {
	box, err := p.Locate(path)
	if box == nil {
		return "", err
	}
	return box.Path(), p.UnBoxWith(box, v, unmarshaler)
}

func (p *Pantry) SaveWith(path string, v interface{}, marshaler MarshalFunc) (string, error) {
	box, err := p.Locate(path)
	if box == nil {
		return "", err
	}
	return box.Path(), p.BoxWith(box, v, marshaler)
}

func (p *Pantry) UnBox(box Box, v interface{}) error {
	b, err := box.Get()
	if err != nil {
		return err
	}
	return p.Unmarshal(b, v, box.Format())
}

func (p *Pantry) Box(box Box, v interface{}) error {
	b, err := p.Marshal(v, box.Format())
	if err != nil {
		return err
	}
	return box.Set(b)
}

func (p *Pantry) Load(path string, v interface{}) (string, error) {
	box, err := p.Locate(path)
	if box == nil {
		return "", err
	}
	return box.Path(), p.UnBox(box, v)
}

func (p *Pantry) LoadBox(path string, v interface{}) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}
	return box, p.UnBox(box, v)
}

func (p *Pantry) Save(path string, v interface{}) (string, error) {
	box, err := p.Locate(path)
	if box == nil {
		return "", err
	}
	return box.Path(), p.Box(box, v)
}
