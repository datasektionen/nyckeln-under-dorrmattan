package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/datasektionen/nyckeln-under-dorrmattan/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pls"
)

func main() {
	flag.Parse()

	loginIDs := make(chan string)

	go login.Listen(loginIDs)
	go pls.Listen()

	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		loginIDs <- strings.TrimSpace(line)
	}
}
