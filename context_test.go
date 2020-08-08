package shireikan

import (
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestContextGetSession(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetSession() != ctx.session {
		t.Error("recovered sesion is invalid")
	}
}

func TestContextGetArgs(t *testing.T) {
	ctx := makeContext(false)
	for i, v := range ctx.GetArgs() {
		if v != ctx.args[i] {
			t.Error("recovered argument is invalid")
		}
	}
}

func TestContextGetChannel(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetChannel() != ctx.channel {
		t.Error("recovered channel is invalid")
	}
}

func TestContextGetMessage(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetMessage() != ctx.message {
		t.Error("recovered message is invalid")
	}
}

func TestContextGetGuild(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetGuild() != ctx.guild {
		t.Error("recovered guild is invalid")
	}
}

func TestContextGetUser(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetUser() != ctx.message.Author {
		t.Error("recovered user is invalid")
	}
}

func TestContextGetMember(t *testing.T) {
	ctx := makeContext(false)
	if ctx.GetMember() != ctx.member {
		t.Error("recovered member is invalid")
	}
}

func TestContextIsDM(t *testing.T) {
	ctx := makeContext(false)
	if !ctx.IsDM() {
		t.Error("recovered isDM is invalid")
	}
}

func TestContextIsEdit(t *testing.T) {
	ctx := makeContext(false)
	if !ctx.IsEdit() {
		t.Error("recovered isEdit is invalid")
	}
}

func TestContextSetGetObject(t *testing.T) {
	ctx := makeContext(true)

	ctx.SetObject("key", 123)
	v, ok := ctx.GetObject("key").(int)
	if !ok {
		t.Error("recovered value is not of type int")
	}
	if v != 123 {
		t.Error("recovered value is invalid")
	}

	v, ok = ctx.GetObject("invalid_key").(int)
	if ok || v != 0 {
		t.Error("non-existent key shall not be recovered")
	}
}

func TestContextGetObjectGlobal(t *testing.T) {
	ctx := makeContext(true)

	v, ok := ctx.GetObject("test_handler").(string)
	if !ok {
		t.Error("recovered global value shall be type of string")
	}
	if v != "test_value" {
		t.Error("recovered global value was invalid")
	}
}

func TestContextInitObjectMap(t *testing.T) {
	ctx := makeContext(false)

	ctx.SetObject("test", "test")
	if ctx.objectMap == nil {
		t.Error("object map shall be initialized after setting key-value")
	}
	v, ok := ctx.GetObject("test").(string)
	if !ok {
		t.Error("recovered value shall be type of string")
	}
	if v != "test" {
		t.Error("recovered value was invalid")
	}

	ctx = makeContext(false)
	vi := ctx.GetObject("test")
	if vi != nil {
		t.Error("recovered value when object map is not initialized shall be nil")
	}
}

// -------------------------------
// --- HELPER ---

func makeContext(initObjectMap bool) *context {
	ctx := &context{
		session: &discordgo.Session{
			Token: "test_token",
		},
		args: ArgumentList([]string{"a", "b", "c"}),
		channel: &discordgo.Channel{
			ID: "test_channel",
		},
		guild: &discordgo.Guild{
			ID: "test_guild",
		},
		isDM:   true,
		isEdit: true,
		member: &discordgo.Member{
			Nick: "test_nick",
		},
		message: &discordgo.Message{
			ID: "test_message",
			Author: &discordgo.User{
				ID: "test_user",
			},
		},
	}

	if initObjectMap {
		handler := &handler{
			objectMap: &sync.Map{},
		}
		handler.objectMap.Store("test_handler", "test_value")

		ctx.objectMap = &sync.Map{}
		ctx.objectMap.Store(ObjectMapKeyHandler, handler)
	}

	return ctx
}
