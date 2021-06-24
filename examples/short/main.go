package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/discordgo"
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
