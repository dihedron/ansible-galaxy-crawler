package main

import (
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Collections *Collections `long:"collections" short:"c" description:"The list of collections to sync." required:"yes" json:"collections" yaml:"collections"`
	Directory   string       `long:"directory" short:"d" description:"The directory into which collections are copies." optional:"yes" default:"_dowloads" json:"directory" yaml:"directory"`
}

func main() {

	options := Options{}
	if _, err := flags.Parse(&options); err != nil {
		os.Exit(1)
	}

	//fmt.Println(helpers.ToPrettyJSON(options))

	client := resty.New()
	for _, collection := range options.Collections.Entries {
		collection.Download(client, options.Directory)
	}
}
