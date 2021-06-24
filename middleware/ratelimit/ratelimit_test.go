package ratelimit

import (
	"testing"
	"time"

	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/discordgo"
)

func TestGetLayer(t *testing.T) {
	m := New()

	if m.GetLayer() != shireikan.LayerBeforeCommand {
		t.Errorf("invalid middleware layer: %d", m.GetLayer())
	}
}

func TestHandle(t *testing.T) {
	m := New()

	cmd := &testCmd{false, false, false}
	ctx := &testContext{discordgo.ChannelTypeGuildText, "guild_id", "user_id"}

	shallPass := func() {
		ok, err := m.Handle(cmd, ctx, m.GetLayer())
		if err != nil {
			t.Errorf("middleware failed: %s", err.Error())
		}
		if !ok {
			t.Error("middleware stopped unexpectedly")
		}
	}

	shallFail := func() {
		ok, err := m.Handle(cmd, ctx, m.GetLayer())
		if err != nil {
			t.Errorf("middleware failed: %s", err.Error())
		}
		if ok {
			t.Error("middleware passed unexpectedly")
		}
	}

	for i := 0; i < cmd.GetLimiterBurst(); i++ {
		shallPass()
	}

	shallFail()

	time.Sleep(cmd.GetLimiterRestoration())
	shallPass()

	shallFail()
}

// ----------------------------

type testContext struct {
	chanType discordgo.ChannelType
	guildId  string
	userID   string
}

func (ctx *testContext) GetSession() *discordgo.Session {
	return nil
}

func (ctx *testContext) GetArgs() shireikan.ArgumentList {
	return nil
}

func (ctx *testContext) GetChannel() *discordgo.Channel {
	return &discordgo.Channel{
		Type: ctx.chanType,
	}
}

func (ctx *testContext) GetMessage() *discordgo.Message {
	return nil
}

func (ctx *testContext) GetGuild() *discordgo.Guild {
	return &discordgo.Guild{
		ID: ctx.guildId,
	}
}

func (ctx *testContext) GetUser() *discordgo.User {
	return &discordgo.User{
		ID: ctx.userID,
	}
}

func (ctx *testContext) GetMember() *discordgo.Member {
	return nil
}

func (ctx *testContext) IsDM() bool {
	return false
}

func (ctx *testContext) IsEdit() bool {
	return false
}

func (ctx *testContext) GetObject(key string) interface{} {
	return nil
}

func (ctx *testContext) SetObject(key string, val interface{}) {}

func (ctx *testContext) Reply(content string) (*discordgo.Message, error) {
	return nil, nil
}

func (ctx *testContext) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}

func (ctx *testContext) ReplyEmbedError(content, title string) (*discordgo.Message, error) {
	return nil, nil
}
