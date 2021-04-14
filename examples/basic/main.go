package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/shireikan/examples/basic/commands"
	"github.com/zekroTJA/shireikan/examples/basic/middleware"
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

	handler.Register(&middleware.Test{})

	handler.Register(&commands.Ping{})
	handler.Register(&commands.Object{})

	handler.Setup(session)
}
