package main

import (
	"flag"

	"github.com/datasektionen/nyckeln-under-dorrmattan/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pls"
)

func main() {
	flag.Parse()

	go login.Listen()
	pls.Listen()
}
