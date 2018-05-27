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

type Locker interface {
	Lock()
	Unlock()
}

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

type Pantry struct {
	Locations *Locations
	Options   *Options
}

func (p *Pantry) Init(applicationName string, locations ...string) *Pantry {
	p.Locations = NewLocations(applicationName, locations...)
	p.Options = &Options{
		Flags:       Flags{Using: FlagsDontUse},
		Enviropment: Enviropment{Use: true, Prefix: ""},
		Tags: Tags{
			Config:      Tag{true, "config"},
			Default:     Tag{true, "default"},
			Description: Tag{false, "description"},
		},
		DefaultFormat: "",
	}
	p.Options.Flags.Init(nil, nil)
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
	if s := p.Options.Flags.GetConfigPath(); s != "" {
		filename = s
	}
	return p.Locations.LocatePath(filename)
}

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

func (p *Pantry) UnBoxWithFormat(box Box, v interface{}, format string) error {
	b, err := box.Get()
	if err != nil {
		return err
	}
	if format == "" {
		format = box.Path()
	}
	return p.Unmarshal(b, v, format)
}

func (p *Pantry) BoxWithFormat(box Box, v interface{}, format string) error {
	if format == "" {
		format = box.Path()
	}
	b, err := p.Marshal(v, format)
	if err != nil {
		return err
	}
	return box.Set(b)
}

// func (p *Pantry) Load(path string, v interface{}) (Box, error) {
// 	box, err := p.Locate(path)
// 	if box == nil {
// 		return nil, err
// 	}
// 	return box, p.UnBox(box, v)
// }

func (p *Pantry) Load(path string, v interface{}, opt ...LoadOptions) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}

	format := box.Path()
	watchers := make([]Watcher, 0, len(opt))
	ctx := context.Background()
	reload := false
	if len(opt) > 0 {
		for _, o := range opt {
			if o.Context != nil {
				ctx = o.Context
			}
			if o.Format != "" {
				format = o.Format
			}
			reload = reload || o.Reload
			if o.Watcher != nil {
				watchers = append(watchers, o.Watcher)
			}
		}

	}

	err = p.UnBoxWithFormat(box, v, format)
	if err != nil {
		return nil, err
	}

	if len(watchers) > 0 || reload {
		err := box.WatchContext(ctx, func(err error) {
			if err == nil {
				if reload {
					err = p.UnBoxWithFormat(box, v, format)
				}
			}
			for _, w := range watchers {
				w(err)
			}
		})
		if err != nil {
			return nil, err
		}
	}

	return box, nil
}

func (p *Pantry) Save(path string, v interface{}) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}
	return box, p.Box(box, v)
}

func (p *Pantry) SaveWithFormat(path string, v interface{}, format string) (Box, error) {
	box, err := p.Locate(path)
	if box == nil {
		return nil, err
	}
	return box, p.BoxWithFormat(box, v, format)
}
