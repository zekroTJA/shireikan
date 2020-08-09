package shireikan

import (
	"log"
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
		t.Error("invalid ammount of command instances were registered")
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

// -------------------------------
// --- HELPER ---

func makeConfig() *Config {
	return &Config{
		GeneralPrefix:         "!",
		AllowBots:             true,
		AllowDM:               true,
		DeleteMessageAfter:    true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: true,
		GuildPrefixGetter: func(string) (string, error) {
			return "", nil
		},
		OnError: func(Context, ErrorType, error) {},
	}
}

type testCmd struct {
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
	log.Printf("%+v", ctx)
	return nil
}

type testMiddleware struct {
}

func (m *testMiddleware) Handle(cmd Command, ctx Context) (error, bool) {
	return nil, true
}
