# shireikan &nbsp; [![](https://img.shields.io/badge/docs-pkg.do.dev-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/zekroTJA/shireikan?tab=doc) ![CI](https://github.com/zekroTJA/shireikan/workflows/CI/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/zekroTJA/shireikan/badge.svg?branch=master)](https://coveralls.io/github/zekroTJA/shireikan?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/zekroTJA/shireikan)](https://goreportcard.com/report/github.com/zekroTJA/shireikan)

司令官 - A prototype, object-oriented command handler for discordgo.

This is a work-in-progress command handler which shall later replace the current command handler used in [shinpuru](https://github.com/zekroTJA/shinpuru).

If you are interested how to use this package, see the [basic example](examples/basic).

This command handler is strongly inspired by [Lukaesebrot's](https://github.com/Lukaesebrot) package [dgc](https://github.com/Lukaesebrot/dgc).

---

## Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekrotja/discordgo"
	"github.com/zekroTJA/shireikan"
)

type UsernameMiddleware struct{}

func (m *UsernameMiddleware) Handle(cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (bool, error) {
	ctx.SetObject("username", ctx.GetUser().Username)
	return true, nil
}

func (m *UsernameMiddleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}

type HelloCommand struct {
}

func (c *HelloCommand) GetInvokes() []string {
	return []string{"hello", "hey"}
}

func (c *HelloCommand) GetDescription() string {
	return "Greets you or another user!"
}

func (c *HelloCommand) GetHelp() string {
	return "`hello` - greets you\n" +
		"`hello <username>` - greets someone else"
}

func (c *HelloCommand) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *HelloCommand) GetDomainName() string {
	return "test.etc.hellp"
}

func (c *HelloCommand) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *HelloCommand) IsExecutableInDMChannels() bool {
	return true
}

func (c *HelloCommand) Exec(ctx shireikan.Context) error {
	var username string
	if username = ctx.GetArgs().Get(0).AsString(); username == "" {
		// Get username from middleware
		username, _ = ctx.GetObject("username").(string)
	}
	_, err := ctx.Reply(fmt.Sprintf("Hello, %s!", username))
	return err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	token := os.Getenv("TOKEN")

	session, err := discordgo.New("Bot " + token)
	must(err)

	must(session.Open())

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()

	handler := shireikan.New(&shireikan.Config{
		GeneralPrefix:         "!",
		AllowBots:             false,
		AllowDM:               true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: true,
		OnError: func(ctx shireikan.Context, typ shireikan.ErrorType, err error) {
			log.Printf("[ERR] [%d] %s", typ, err.Error())
		},
	})

	handler.Register(&UsernameMiddleware{})
	handler.Register(&HelloCommand{})

	handler.Setup(session)
}
```

You can find more complex examples in the [`examples/`](examples) directory.

shireikan is used in my own Discord bot [**shinpuru**](https://github.com/zekroTJA/shinpuru).
Take a look if you want to see a "real world" implementation example. 😉

---

© 2020-2021 Ringo Hoffmann (zekro Development).  
Covered by the MIT Licence.
