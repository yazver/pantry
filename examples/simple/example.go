package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/yazver/pantry"
)

type subConfig struct {
	Weight int `default:"55"`
	Height int
}

type config struct {
	Age            int   `pantry:"flag:age|My age;env:AGE" default:"18" toml:"age"`
	DefaultAge     int   `default:"18"`
	DescriptionTag bool  `pantry:"flag:d" description:"This is description"`
	ExpandedTags   bool  `pantry.flag:"e|Expanded flag" pantry.env:"exp"`
	PointerToInt   *int  `default:"100"`
	PointerToUInt  *uint `pantry:"env:PTU"`
	Cats           []string
	Pi             float64
	Perfection     []int
	DOB            time.Time
	Sub            []*subConfig
	SubMap         map[string]*subConfig
}

func main() {
	os.Setenv("TEST_AGE", "25")
	os.Setenv("TEST_PTU", "111")

	p := pantry.NewPantry("Pantry", pantry.LocationConfigDir, pantry.LocationApplicationDir)
	p.Options.Flags.Using = pantry.FlagsUseAll
	p.Options.Flags.Args = []string{"-e"} // Only for test
	p.Options.Enviropment.Prefix = "TEST_"
	p.Options.Enviropment.Use = true

	cfg := &config{}

	box, err := p.Load("config.toml", cfg)
	if err != nil {
		log.Fatalln(err)
	}
	filePath := box.Path()
	//fmt.Println("Config path: " + filePath)
	//fmt.Printf("%#v", cfg)
	spew.Dump(cfg)

	if _, err := p.Save(strings.Replace(filePath, ".toml", ".yaml", -1), cfg); err != nil {
		log.Fatalln(err)
	}
	if _, err := p.Save(strings.Replace(filePath, ".toml", ".json", -1), cfg); err != nil {
		log.Fatalln(err)
	}
}
