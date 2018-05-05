package main

import (
	"fmt"
	"log"
	"net"

	"github.com/yazver/pantry"
)

type SubConfig struct {
	Weight int `default:"55"`
	Height int
}

type Config struct {
	IP   net.IP `config:"flag:ip|IP address;env:IP" default:"0.0.0.0" toml:"ip"`
	Port int16  `config:"flag:port|Port;env:PORT" default:"1080" toml:"port"`
}

func main() {
	p := pantry.NewPantry("TestApp", pantry.LocationConfigDir, pantry.LocationApplicationDir)
	p.Options.Flags.Using = pantry.FlagsUseAll

	config := Config{}

	filePath, err := p.Load("config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Config path: " + filePath)
	fmt.Printf("%#v", c)
	//spew.Dump(config)

	// 1
	box, err := p.Load("config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	box.Watch(func(err error) { p.UnBox(box, &config) })

}
