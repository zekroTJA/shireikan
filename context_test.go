package shireikan

import (
	"os"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
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
	ctx := makeContext(true)

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

	ctx = makeContext(true)
	vi := ctx.GetObject("test")
	if vi != nil {
		t.Error("recovered value when object map is not initialized shall be nil")
	}
}

func TestContextReply(t *testing.T) {
	ctx := makeContext(false)
	ctx.session = makeSession()

	msg, err := ctx.Reply("test reply")
	if err != nil {
		t.Error("reply failed: ", err)
	}

	if msg.Content != "test reply" {
		t.Error("reply had invalid content")
	}
}

func TestContextReplyEmbed(t *testing.T) {
	ctx := makeContext(false)
	ctx.session = makeSession()

	emb := &discordgo.MessageEmbed{
		Title:       "test embed",
		Description: "test description",
	}

	msg, err := ctx.ReplyEmbed(emb)
	if err != nil {
		t.Error("reply failed: ", err)
	}

	if len(msg.Embeds) != 1 {
		t.Error("reply had no embed")
	}

	if msg.Embeds[0].Title != emb.Title || msg.Embeds[0].Description != emb.Description {
		t.Error("reply had invalid embed")
	}
}

func TestContextReplyEmbedError(t *testing.T) {
	ctx := makeContext(false)
	ctx.session = makeSession()

	msg, err := ctx.ReplyEmbedError("test content", "test title")
	if err != nil {
		t.Error("reply failed: ", err)
	}

	if len(msg.Embeds) != 1 {
		t.Error("reply had no embed")
	}

	if msg.Embeds[0].Title != "test title" ||
		msg.Embeds[0].Description != "test content" ||
		msg.Embeds[0].Color != EmbedColorError {

		t.Error("reply had invalid embed")
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
			ID: getEnvOrDefault("CHANNEL_ID", "549871005321920513"),
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
		b, _ := di.NewBuilder()
		b.Set("test_handler", "test_value")

		handler := &handler{
			objectContainer: b.Build(),
			objectMaps: &sync.Pool{
				New: func() interface{} { return &sync.Map{} },
			},
			objectMap: &sync.Map{},
		}
		ctx.handler = handler
	}

	return ctx
}

func makeSession() *discordgo.Session {
	s, _ := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	return s
}
