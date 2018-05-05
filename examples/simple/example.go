package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yazver/pantry"
)

type subConfig struct {
	Weight int `default:"55"`
	Height int
}

type config struct {
	Age            int  `config:"flag:age|My age;env:AGE" default:"18" toml:"age"`
	DefaultAge     int  `default:"18"`
	DescriptionTag bool `config:"flag:d" description:"This is description"`
	ExpandedTags   bool `config.flag:"e|Expanded flag" config.env:"exp"`
	PointerToInt   *int `default:"100"`
	Cats           []string
	Pi             float64
	Perfection     []int
	DOB            time.Time
	Sub            []*subConfig
	SubMap         map[string]*subConfig
}

func main() {
	p := pantry.NewPantry("Pantry", pantry.LocationConfigDir, pantry.LocationApplicationDir)
	p.Options.Flags.Using = pantry.FlagsUseAll
	p.Options.Enviropment.Prefix = "TEST_"
	p.Options.Enviropment.Use = true

	cfg := &config{}

	box, err := p.Load("config.toml", cfg)
	if err != nil {
		log.Fatalln(err)
	}
	filePath := box.Path()
	fmt.Println("Config path: " + filePath)
	fmt.Printf("%#v", cfg)
	// scs := spew.ConfigState{Indent: "    ", DisableCapacities: true, DisablePointerAddresses: true}
	// scs.Dump(config)
	if _, err := p.Save(strings.Replace(filePath, ".toml", ".yaml", -1), cfg); err != nil {
		log.Fatalln(err)
	}
	if _, err := p.Save(strings.Replace(filePath, ".toml", ".json", -1), cfg); err != nil {
		log.Fatalln(err)
	}
}
