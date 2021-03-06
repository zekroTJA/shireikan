package ratelimit

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shireikan"
)

// Middleware command implements the ratelimit middleware.
type Middleware struct {
	manager Manager
}

// New returns a new instance of Middleware.
//
// Optionally, you can pass a custom Manager instance
// if you want to handle limiters differently than the
// standard Manager implementation.
func New(manager ...Manager) *Middleware {
	var m Manager

	if len(manager) > 0 && manager[0] != nil {
		m = manager[0]
	} else {
		m = newInternalManager()
	}

	return &Middleware{m}
}

func (m *Middleware) Handle(cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (bool, error) {
	c, ok := cmd.(LimitedCommand)
	if !ok {
		return true, nil
	}

	var guildID string
	if c.IsLimiterGlobal() {
		guildID = "__global__"
	} else if ctx.GetChannel().Type == discordgo.ChannelTypeDM || ctx.GetChannel().Type == discordgo.ChannelTypeGroupDM {
		guildID = "__dm__"
	} else {
		guildID = ctx.GetGuild().ID
	}

	limiter := m.manager.GetLimiter(cmd, ctx.GetUser().ID, guildID)
	if ok, next := limiter.Take(); !ok {
		_, err := ctx.ReplyEmbedError(fmt.Sprintf(
			"You are being ratelimited.\nWait %s until you can use this command again.",
			next.String()), "Rate Limited")
		return false, err
	}

	return true, nil
}

func (m *Middleware) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}
