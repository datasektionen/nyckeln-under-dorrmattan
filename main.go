package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/hive"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/pls"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/sso"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/kthldap"

	"golang.org/x/term"
)

func main() {

	cfg := config.GetConfig()

	dao := dao.New(cfg)

	loginIDs := make(chan string)

	go login.Listen(cfg, loginIDs)
	go pls.Listen(cfg, dao)
	go sso.Listen(cfg, dao)
	go hive.Listen(cfg, dao)
	go kthldap.Listen(cfg, dao)

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
