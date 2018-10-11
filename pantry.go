package pantry

import (
	"context"
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

// Locker interface is used for exlusive access to a variable if it realized.
type Locker interface {
	Lock()
	Unlock()
}

// Tag options.
type Tag struct {
	Use  bool   // Use tag
	Name string // Name of tag
}

// Tags options.
type Tags struct {
	Config      Tag
	Flag        Tag
	Env         Tag
	Default     Tag
	Description Tag
}

// Config processing options.
type Options struct {
	Flags         Flags
	Enviropment   Enviropment
	Tags          Tags
	DefaultFormat string
}

// Pantry is used to load config data from different sources.
type Pantry struct {
	Locations *Locations
	Options   *Options
}

// Init Pantry.
func (p *Pantry) Init(applicationName string, locations ...string) *Pantry {
	p.Locations = NewLocations(applicationName, locations...)
	p.Options = &Options{
		Flags:       Flags{Using: FlagsDontUse},
		Enviropment: Enviropment{Use: true, Prefix: ""},
		Tags: Tags{
			Config:      Tag{true, "pantry"},
			Flag:        Tag{true, "flag"},
			Env:         Tag{true, "env"},
			Default:     Tag{true, "default"},
			Description: Tag{false, "desc"},
		},
		DefaultFormat: "",
	}
	p.Options.Flags.Init(nil, nil)
	return p
}

// NewPantry creates new Pantry.
func NewPantry(applicationName string, locations ...string) *Pantry {
	return new(Pantry).Init(applicationName, locations...)
}

// AddLocation adds searching location.
func (p *Pantry) AddLocation(locationParts ...string) {
	p.Locations.AddJoin(locationParts...)
}

// AddLocations adds searching locations.
func (p *Pantry) AddLocations(locations ...string) {
	p.Locations.Add(locations...)
}

// LocatePath looks for the file in previously added locations.
func (p *Pantry) LocatePath(filename string) (string, error) {
	if s := p.Options.Flags.GetConfigPath(); s != "" {
		filename = s
	}
	return p.Locations.LocatePath(filename)
}

// Locate looks for the file in previously added locations.
func (p *Pantry) Locate(filename string) (Box, error) {
	if s := p.Options.Flags.GetConfigPath(); s != "" {
		filename = s
	}
	return p.Locations.Locate(filename)
}

func (p *Pantry) searchFormat(s string) (*ConfigFormat, error) {
	f, err := Formats.Search(s)
	if err != nil && p.Options.DefaultFormat != "" {
		return Formats.Search(p.Options.DefaultFormat)
	}
	return f, err
}

// UnmarshalWith unmarshals data by "unmarshaler" and applays enviropment variables and command line flags.
func (p *Pantry) UnmarshalWith(b []byte, v interface{}, unmarshaler UnmarshalFunc) error {
	if l, ok := v.(Locker); ok {
		l.Lock()
		defer l.Unlock()
	}

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
	if l, ok := v.(Locker); ok {
		l.Lock()
		defer l.Unlock()
	}

	b, err = marshaler(v)
	if err != nil {
		return nil, ConfigEncodeError{err}
	}
	return
}

// Unmarshal unmarshals data as diffened format and applays enviropment variables and command line flags.
func (p *Pantry) Unmarshal(b []byte, v interface{}, format string) error {
	f, err := p.searchFormat(format)
	if err != nil {
		return err
	}
	return p.UnmarshalWith(b, v, f.Unmarshal)
}

func (p *Pantry) Marshal(v interface{}, format string) (b []byte, err error) {
	f, err := p.searchFormat(format)
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
	return p.Unmarshal(b, v, box.Path())
}

func (p *Pantry) Box(box Box, v interface{}) error {
	b, err := p.Marshal(v, box.Path())
	if err != nil {
		return err
	}
	return box.Set(b)
}

func (p *Pantry) UnBoxAs(box Box, v interface{}, format string) error {
	b, err := box.Get()
	if err != nil {
		return err
	}
	if format == "" {
		format = box.Path()
	}
	return p.Unmarshal(b, v, format)
}

func (p *Pantry) BoxAs(box Box, v interface{}, format string) error {
	if format == "" {
		format = box.Path()
	}
	b, err := p.Marshal(v, format)
	if err != nil {
		return err
	}
	return box.Set(b)
}

func (p *Pantry) LoadAs(path string, v interface{}, format string, opts ...func(*LoadOptions)) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}

	lo := LoadOptions{Context: context.Background()}
	for _, opt := range opts {
		opt(&lo)
	}

	err = p.UnBoxAs(box, v, format)
	if err != nil {
		return nil, err
	}

	if lo.Watcher != nil || lo.Reload {
		err := box.WatchContext(lo.Context, func(err error) {
			if err == nil {
				if lo.Reload {
					err = p.UnBoxAs(box, v, format)
				}
			}
			for lo.Watcher != nil {
				lo.Watcher(err)
			}
		})
		if err != nil {
			return nil, err
		}
	}

	return box, nil
}

func (p *Pantry) Load(path string, v interface{}, opt ...func(*LoadOptions)) (Box, error) {
	return p.LoadAs(path, v, "", opt...)
}

func (p *Pantry) Save(path string, v interface{}) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}
	return box, p.Box(box, v)
}

func (p *Pantry) SaveAs(path string, v interface{}, format string) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}
	return box, p.BoxAs(box, v, format)
}
