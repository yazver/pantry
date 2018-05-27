package main

import (
	"fmt"
	"log"
	"net"

	"github.com/yazver/pantry"
)

type config struct {
	IP   net.IP `config:"flag:ip|IP address;env:IP" default:"0.0.0.0" toml:"ip"`
	Port int16  `config:"flag:port|Port;env:PORT" default:"1080" toml:"port"`
}

func main() {
	p := pantry.NewPantry("testapp", pantry.LocationConfigDir, pantry.LocationApplicationDir)
	p.Options.Flags.Using = pantry.FlagsUseAll

	cfg := config{}

	// box, err := p.Load("config.toml", &cfg)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("Config path: " + box.Path())
	// fmt.Printf("%+v\n", cfg)

	// done := make(chan struct{})
	// box.Watch(func(err error) {
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	fmt.Println("Config file changed.")
	// 	close(done)
	// })

	done := make(chan struct{})
	box, err := p.Load("config.toml", &cfg, pantry.Watch(func(err error) {
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Config file changed.")
		close(done)
	}))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Config path: " + box.Path())
	fmt.Printf("%+v\n", cfg)

	if err := p.Box(box, cfg); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Write config.")

	<-done
}
