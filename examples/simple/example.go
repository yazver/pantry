package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/yazver/pantry"
)

type SubConfig struct {
	Weight int `default:"55"`
	Height int
}

type Config struct {
	Age            int  `config:"flag:age|My age;env:AGE" default:"18" toml:"age"`
	DefaultAge     int  `default:"18"`
	DescriptionTag bool `config:"flag:d" description:"This is description"`
	ExpandedTags   bool `config.flag:"e|Expanded flag" config.env:"exp"`
	Cats           []string
	Pi             float64
	Perfection     []int
	DOB            time.Time
	Sub            []*SubConfig
	SubMap         map[string]*SubConfig
}

func main() {
	p := pantry.NewPantry("Pantry", pantry.LocationConfigDir, pantry.LocationApplicationDir)
	p.Options.Flags.Using = pantry.FlagsUseAll
	p.Options.Enviropment.Prefix = "TEST_"
	p.Options.Enviropment.Use = true

	var config *Config = &Config{}

	filePath, err := p.Load("config.toml", config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Config path: " + filePath)
	//fmt.Printf("%#v", c)
	spew.Dump(config)
	p.Save(strings.Replace(filePath, ".toml", ".yaml", -1), config)
	p.Save(strings.Replace(filePath, ".toml", ".json", -1), config)
}
