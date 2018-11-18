package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/radio-t/utils/cmd"
)

var opts struct {
	Mp3Tags cmd.Mp3TagsCommand `command:"mp3tags"`
}

func main() {
	log.Printf("rt-utils, %+v", opts)

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
