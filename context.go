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

	// GetMember returns the user object
	// of the author of the command message.
	GetUser() *discordgo.User

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

	// GetObject tries to retrieve an object from the
	// context's object map. The retrieved value
	// is returned. If no value could be retrieved,
	// nil is returned.
	//
	// This function must return a value when it
	// is not available in the contexts object
	// map but available in the command handlers
	// global object map.
	GetObject(key string) interface{}

	// SetObject sets an object or value to the context's
	// object map with the passed key.
	SetObject(key string, val interface{})

	// Reply sends a message with the passed content
	// to the channel where the command was sent into.
	Reply(content string) (*discordgo.Message, error)

	// Reply sends a message with the passed embed
	// to the channel where the command was sent into.
	ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error)

	// ReplyEmbedError sends a pre-constructed embed
	// as error message to the channel where the command
	// was sent into.
	ReplyEmbedError(content, title string) (*discordgo.Message, error)
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

func (ctx *context) GetUser() *discordgo.User {
	return ctx.message.Author
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

func (ctx *context) GetObject(key string) interface{} {
	if ctx.objectMap == nil {
		return nil
	}

	val, ok := ctx.objectMap.Load(key)

	if !ok {
		vHandler, _ := ctx.objectMap.Load(ObjectMapKeyHandler)
		handler, ok := vHandler.(Handler)
		if ok {
			val = handler.GetObject(key)
		}
	}

	return val
}

func (ctx *context) SetObject(key string, val interface{}) {
	if ctx.objectMap == nil {
		ctx.objectMap = &sync.Map{}
	}

	ctx.objectMap.Store(key, val)
}

func (ctx *context) Reply(content string) (*discordgo.Message, error) {
	return ctx.session.ChannelMessageSend(ctx.channel.ID, content)
}

func (ctx *context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.session.ChannelMessageSendEmbed(ctx.channel.ID, embed)
}

func (ctx *context) ReplyEmbedError(content, title string) (*discordgo.Message, error) {
	return ctx.ReplyEmbed(&discordgo.MessageEmbed{
		Title:       title,
		Description: content,
		Color:       EmbedColorError,
	})
}
