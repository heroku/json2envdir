package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/heroku/json2envdir"
	"github.com/heroku/json2envdir/config"
	flag "github.com/ogier/pflag"
)

var (
	jsonFile   = flag.String("file", "-", "file with JSON environment")
	configFile = flag.String("config", "", "config file")
)

func readStdin() string {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func main() {
	flag.Parse()

	cfg := config.LoadConfig(*configFile)

	var err error
	if *jsonFile == "-" {
		err = json2envdir.Parse(cfg, readStdin())
	} else {
		err = json2envdir.Parse(cfg, readFile(*jsonFile))
	}

	if err != nil {
		log.Fatal(err)
	}
}
