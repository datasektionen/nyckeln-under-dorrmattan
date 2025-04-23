package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/doi"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/pls"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/sso"

	"golang.org/x/term"
)

func main() {

	cfg := config.GetConfig()

	doi := doi.New(cfg)

	loginIDs := make(chan string)

	go login.Listen(cfg, loginIDs)
	go pls.Listen(cfg, doi)
	go sso.Listen(cfg, doi)

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
