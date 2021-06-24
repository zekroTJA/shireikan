package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/shireikan/examples/di/commands"
	"github.com/zekroTJA/shireikan/examples/di/database"
	"github.com/zekrotja/discordgo"
)

func main() {
	token := os.Getenv("TOKEN")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	err = session.Open()
	if err != nil {
		panic(err)
	}

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()

	diBuilder, _ := di.NewBuilder()
	diBuilder.Add(di.Def{
		Name: "db",
		Build: func(ctn di.Container) (interface{}, error) {
			return &database.TestDB{}, nil
		},
	})

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
		ObjectContainer: diBuilder.Build(),
	})

	handler.Register(&commands.Object{})

	handler.Setup(session)
}
