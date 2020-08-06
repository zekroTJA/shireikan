package shireikan

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Context wraps information about a message
// and the environment where the message was
// created which is passed to middleware and
// command handlers.
type Context interface {

	// GetSession returns the current discordgo.Session.
	GetSession() *discordgo.Session

	// GetArgs returns an ArgumentList of
	// the parsed command arguments.
	GetArgs() ArgumentList

	// GetChannel returns the channel where
	// the command message was sent into.
	GetChannel() *discordgo.Channel

	// GetMessage returns the original
	// message object of the command.
	GetMessage() *discordgo.Message

	// GetGuild returns the guild object
	// where the command was sent.
	GetGuild() *discordgo.Guild

	// GetMember returns the member object
	// of the author of the command message.
	GetMember() *discordgo.Member

	// IsDM returns true when the command
	// message was sent into a DM or
	// GroupDM channel.
	IsDM() bool

	// IsEdit returns true if the event which
	// invoked the command was a
	// discordgo.MessageUpdate event.
	IsEdit() bool

	// Get tries to retrieve an object from the
	// context's object map. The retrieved value
	// is returned. If no value could be retrieved,
	// nil is returned.
	Get(key string) interface{}

	// Set sets an object or value to the context's
	// object map with the passed key.
	Set(key string, val interface{})
}

// context is the default implementation of Context.
type context struct {
	session   *discordgo.Session
	args      ArgumentList
	message   *discordgo.Message
	guild     *discordgo.Guild
	channel   *discordgo.Channel
	member    *discordgo.Member
	isDM      bool
	isEdit    bool
	objectMap *sync.Map
}

func (ctx *context) GetSession() *discordgo.Session {
	return ctx.session
}

func (ctx *context) GetArgs() ArgumentList {
	return ctx.args
}

func (ctx *context) GetChannel() *discordgo.Channel {
	return ctx.channel
}

func (ctx *context) GetMessage() *discordgo.Message {
	return ctx.message
}

func (ctx *context) GetGuild() *discordgo.Guild {
	return ctx.guild
}

func (ctx *context) GetMember() *discordgo.Member {
	return ctx.member
}

func (ctx *context) IsDM() bool {
	return ctx.isDM
}

func (ctx *context) IsEdit() bool {
	return ctx.isEdit
}

func (ctx *context) Get(key string) interface{} {
	if ctx.objectMap == nil {
		return nil
	}

	val, _ := ctx.objectMap.Load(key)
	return val
}

func (ctx *context) Set(key string, val interface{}) {
	if ctx.objectMap == nil {
		ctx.objectMap = &sync.Map{}
	}

	ctx.objectMap.Store(key, val)
}
