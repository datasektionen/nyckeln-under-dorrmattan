package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/datasektionen/nyckeln-under-dorrmattan/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pls"

	"golang.org/x/term"
)

func main() {
	flag.Parse()

	loginIDs := make(chan string)

	go login.Listen(loginIDs)
	go pls.Listen()

	if term.IsTerminal(int(os.Stdin.Fd())) {
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
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Wait()
	}
}
