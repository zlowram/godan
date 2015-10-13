package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	config, err := loadConfig(flag.Arg(0))
	if err != nil {
		log.Fatal("Couldn't read config file!")
	}
	s := newServer(config)
	log.Fatalln(s.start())
}

func usage() {
	fmt.Fprintln(os.Stderr, usageLine)
	flag.PrintDefaults()
}

const usageLine = `usage: godan <config>
  <config> path to config file`
