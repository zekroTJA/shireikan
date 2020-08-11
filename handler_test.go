package shireikan

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHandlerNewHandler(t *testing.T) {
	h := NewHandler(makeConfig())
	if h == nil {
		t.Error("recovered handler instance is nil")
	}
}

func TestHandlerRegisterCommand(t *testing.T) {
	h := NewHandler(makeConfig())

	cmdInstance := &testCmd{}
	h.RegisterCommand(cmdInstance)

	if h.(*handler).cmdMap["ping"] != cmdInstance {
		t.Error("command instance was not registered correctly")
	}

	if h.(*handler).cmdMap["p"] != cmdInstance {
		t.Error("command instance was not registered correctly")
	}
}

func TestHandlerRegisterMiddleware(t *testing.T) {
	h := NewHandler(makeConfig())

	mwInstance := &testMiddleware{}
	h.RegisterMiddleware(mwInstance)

	if len(h.(*handler).middlewares) == 0 {
		t.Error("middleware did not got registered")
	}

	if h.(*handler).middlewares[0] != mwInstance {
		t.Error("recovered middleware returned false")
	}
}

func TestHandlerRegisterHandlers(t *testing.T) {
	h := NewHandler(makeConfig())

	session, _ := discordgo.New("")

	h.RegisterHandlers(session)
}

func TestHandlerGetConfig(t *testing.T) {
	cfg := makeConfig()
	h := NewHandler(cfg)

	if h.GetConfig() != cfg {
		t.Error("recovered config is invalid")
	}
}

func TestHandlerGetCommandMap(t *testing.T) {
	h := NewHandler(makeConfig())

	cmdInstance := &testCmd{}
	h.RegisterCommand(cmdInstance)

	cmdMap := h.GetCommandMap()

	if cmdMap["ping"] != cmdInstance {
		t.Error("recovered map entry for test command is invalid")
	}

	if cmdMap["p"] != cmdInstance {
		t.Error("recovered map entry for test command is invalid")
	}
}

func TestHandlerGetCommandInstances(t *testing.T) {
	h := NewHandler(makeConfig())

	cmdInstance := &testCmd{}
	h.RegisterCommand(cmdInstance)

	cmdInstances := h.GetCommandInstances()

	if len(cmdInstances) != 1 {
		t.Error("invalid amount of command instances were registered")
	}

	if cmdInstances[0] != cmdInstance {
		t.Error("recovered invalid command instance")
	}
}

func TestHandlerGetCommand(t *testing.T) {
	h := NewHandler(makeConfig())

	cmdInstance := &testCmd{}
	h.RegisterCommand(cmdInstance)

	if cmd, ok := h.GetCommand("ping"); !ok || cmd != cmdInstance {
		t.Error("recovered invalid command instance")
	}

	if cmd, ok := h.GetCommand("p"); !ok || cmd != cmdInstance {
		t.Error("recovered invalid command instance")
	}
}

func TestHandlerSetObject(t *testing.T) {
	h := NewHandler(makeConfig())

	h.SetObject("test", 123)

	rec, ok := h.(*handler).objectMap.Load("test")
	if !ok {
		t.Error("set object map value does not exist in object map")
	}
	if v, _ := rec.(int); v != 123 {
		t.Error("recovered object map value is invalid")
	}
}

func TestHandlerGetObject(t *testing.T) {
	h := NewHandler(makeConfig())

	h.(*handler).objectMap.Store("test", 456)

	val, _ := h.GetObject("test").(int)
	if val != 456 {
		t.Error("recovered object map value is invalid")
	}
}

func TestHandlerMessageHandler(t *testing.T) {
	testMessageHandler(t, true, func(msg *discordgo.Message) {
		msg.Content = "!ping"
	})

	testMessageHandler(t, true, func(msg *discordgo.Message) {
		msg.Content = "test!ping"
	})

	testMessageHandler(t, false, func(msg *discordgo.Message) {
		msg.Content = "!abc"
	})

	testMessageHandler(t, false, func(msg *discordgo.Message) {
		msg.Author.Bot = true
		msg.Content = "!ping"
	})
}

func TestHandlerMiddleware(t *testing.T) {
	s, _ := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))

	cExit := make(chan bool, 1)

	cmd := &testCmd{}
	cfg := makeConfig()
	cfg.DeleteMessageAfter = false
	h := NewHandler(cfg)
	h.RegisterCommand(cmd)

	mwBefore := &testMiddleware{}
	mwBefore.layer = LayerBeforeCommand
	h.RegisterMiddleware(mwBefore)

	mwAfter := &testMiddleware{}
	mwAfter.layer = LayerAfterCommand
	h.RegisterMiddleware(mwAfter)

	msg := &discordgo.Message{
		ChannelID: getEnvOrDefault("CHANNEL_ID", "549871005321920513"),
		GuildID:   getEnvOrDefault("GUILD_ID", "526196711962705925"),
		Author: &discordgo.User{
			ID:  getEnvOrDefault("AUTHOR_ID", "221905671296253953"),
			Bot: false,
		},
		Member: &discordgo.Member{
			GuildID: getEnvOrDefault("GUILD_ID", "526196711962705925"),
			User: &discordgo.User{
				ID:  getEnvOrDefault("AUTHOR_ID", "221905671296253953"),
				Bot: false,
			},
		},
		Content: "!ping",
	}

	s.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		cmd.fail = true
		h.(*handler).messageHandler(s, msg, false)

		if !mwBefore.executed {
			t.Error("brefore middleware was not executed")
		}

		if mwAfter.executed {
			t.Error("after middleware was executed even if it should not be")
		}

		mwBefore.executed = false
		mwAfter.executed = false

		// ----------------------------------------------------

		cmd.fail = false
		h.(*handler).messageHandler(s, msg, false)

		if !mwBefore.executed {
			t.Error("brefore middleware was not executed")
		}

		if !mwAfter.executed {
			t.Error("after middleware was not executed")
		}

		mwBefore.executed = false
		mwAfter.executed = false

		// ----------------------------------------------------

		msg.Content = ""
		h.(*handler).messageHandler(s, msg, false)

		if mwBefore.executed {
			t.Error("before middleware was executed even if it should not be")
		}

		if mwAfter.executed {
			t.Error("after middleware was executed even if it should not be")
		}

		cExit <- true
	})

	err := s.Open()
	if err != nil {
		t.Error(err)
	}

	<-cExit
}

// -------------------------------
// --- HELPER ---

func testMessageHandler(t *testing.T,
	cmdShallbeexecuted bool,
	configurator func(msg *discordgo.Message)) {

	t.Helper()

	s, _ := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))

	cExit := make(chan bool, 1)

	cmd := &testCmd{}
	h := NewHandler(makeConfig())
	h.RegisterCommand(cmd)

	msg := &discordgo.Message{
		ChannelID: getEnvOrDefault("CHANNEL_ID", "549871005321920513"),
		GuildID:   getEnvOrDefault("GUILD_ID", "526196711962705925"),
		Author: &discordgo.User{
			ID:  getEnvOrDefault("AUTHOR_ID", "221905671296253953"),
			Bot: false,
		},
		Member: &discordgo.Member{
			GuildID: getEnvOrDefault("GUILD_ID", "526196711962705925"),
			User: &discordgo.User{
				ID:  getEnvOrDefault("AUTHOR_ID", "221905671296253953"),
				Bot: false,
			},
		},
	}

	configurator(msg)

	s.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		h.(*handler).messageHandler(s, msg, false)

		if !cmd.wasExecuted && cmdShallbeexecuted {
			t.Error("command was not executed")
		} else if cmd.wasExecuted && !cmdShallbeexecuted {
			t.Error("command was executed")
		}

		cExit <- true
	})

	err := s.Open()
	if err != nil {
		t.Error(err)
	}

	<-cExit
}

func makeConfig() *Config {
	return &Config{
		GeneralPrefix:         "!",
		AllowBots:             false,
		AllowDM:               false,
		DeleteMessageAfter:    true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: false,
		GuildPrefixGetter: func(string) (string, error) {
			return "test!", nil
		},
		OnError: func(_ Context, t ErrorType, err error) {
			fmt.Printf("[%d] %s\n", t, err.Error())
		},
	}
}

func getEnvOrDefault(envKey, def string) string {
	v := os.Getenv(envKey)
	if v == "" {
		v = def
	}
	return v
}

type testCmd struct {
	wasExecuted bool
	fail        bool
}

func (c *testCmd) GetInvokes() []string {
	return []string{"ping", "p"}
}

func (c *testCmd) GetDescription() string {
	return "ping pong"
}

func (c *testCmd) GetHelp() string {
	return "`ping` - ping"
}

func (c *testCmd) GetGroup() string {
	return GroupFun
}

func (c *testCmd) GetDomainName() string {
	return "test.fun.ping"
}

func (c *testCmd) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *testCmd) IsExecutableInDMChannels() bool {
	return true
}

func (c *testCmd) Exec(ctx Context) error {
	c.wasExecuted = true

	if c.fail {
		return errors.New("test error")
	}

	return nil
}

type testMiddleware struct {
	executed bool
	layer    MiddlewareLayer
}

func (m *testMiddleware) Handle(cmd Command, ctx Context) (bool, error) {
	m.executed = true
	return true, nil
}

func (m *testMiddleware) GetLayer() MiddlewareLayer {
	return m.layer
}
