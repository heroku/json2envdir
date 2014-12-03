package main

import (
	"io/ioutil"
	"os"

	"github.com/heroku/json2envdir"
	flag "github.com/ogier/pflag"
)

var (
	jsonFile = flag.String("file", "-", "file with JSON environment")
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

	if *jsonFile == "-" {
		json2envdir.Parse(readStdin())
	} else {
		json2envdir.Parse(readFile(*jsonFile))
	}
}
