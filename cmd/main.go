package main

import (
	"flag"
	"fmt"
)

var configString string

func init() {
	flag.StringVar(&configString, "config", "./configs/test.toml", "Path to configuration file")
}

func main() {
	fmt.Println("Hello, World")
	fmt.Println(configString)
}
