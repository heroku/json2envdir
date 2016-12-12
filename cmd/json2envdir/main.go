package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/heroku/json2envdir"
	"github.com/heroku/json2envdir/config"
	flag "github.com/ogier/pflag"
)

var (
	jsonFile    = flag.StringP("file", "f", "-", "file with JSON environment")
	configFile  = flag.StringP("config", "c", config.DefaultConfigFile, "config file")
	versionFlag = flag.BoolP("version", "v", false, "show version")
	version     = "dev"
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

	if *versionFlag {
		fmt.Println(version)
		return
	}

	cfg := config.LoadConfig(*configFile)

	var err error
	if *jsonFile == "-" {
		err = json2envdir.Process(cfg, readStdin())
	} else {
		err = json2envdir.Process(cfg, readFile(*jsonFile))
	}

	if err != nil {
		log.Fatal(err)
	}
}
