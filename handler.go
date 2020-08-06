// Package shireikan provides a general command
// handler for discordgo.
package shireikan

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// ErrorType is the type of error occured in
// the command message handler.
type ErrorType int

const (
	ErrTypGuildPrefixGetter    ErrorType = iota // error from guild prefix getter function
	ErrTypGetChannel                            // error getting channel object
	ErrTypGetGuild                              // error getting guild object
	ErrTypCommandNotFound                       // command was not found by specified invoke
	ErrTypNotExecutableInDM                     // command which is specified as non-executable in DM got executed in a DM channel
	ErrTypMiddleware                            // middleware handler returned an error
	ErrTypCommandExec                           // command handler returned an error
	ErrTypDeleteCommandMessage                  // deleting command message failed
)

var (
	argsRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
)

// Config wraps configuration values for the CommandHandler.
type Config struct {
	GeneralPrefix         string `json:"general_prefix"`           // General and globally accessable prefix
	InvokeToLower         bool   `json:"invoke_to_lower"`          // Lowercase command invoke befor map matching
	AllowDM               bool   `json:"allow_dm"`                 // Allow commands to be executed in DM and GroupDM channels
	AllowBots             bool   `json:"allow_bots"`               // Allow bot accounts to execute commands
	ExecuteOnEdit         bool   `json:"execute_on_edit"`          // Execute command handler when a message was edited
	UseDefaultHelpCommand bool   `json:"use_default_help_command"` // Whether or not to use default help command
	DeleteMessageAfter    bool   `json:"delete_message_after"`     // Delete command message after command has processed

	// OnError is called when the command handler failed
	// or retrieved an error form a middleware or command
	// exec handler.
	//
	// The OnError handler is getting passed the context
	// (which may be incompletely initialized!), an
	// ErrorType and the error object.
	OnError func(ctx Context, errTyp ErrorType, err error)

	// GuildPrefixGetter is called to retrieve a guilds
	// specific prefix.
	//
	// The function is getting passed the guild's ID and
	// returns the guild prefix, when specified. The returned
	// string is empty when no guild prefix is specified.
	// An error is only returned when the retrieving of the
	// guild prefix failed unexpectedly.
	GuildPrefixGetter func(guildID string) (string, error)
}

// Handler specifies a command register and handler.
type Handler interface {

	// RegisterCommand registers the passed
	// Command instance.
	RegisterCommand(cmd Command)

	// RegisterMiddleware registers the
	// passed middleware instance.
	RegisterMiddleware(mw Middleware)

	// RegisterHandlers registers the message
	// handlers to the passed discordgo.Session
	// which are used to handle and parse commands.
	RegisterHandlers(session *discordgo.Session)

	// GetConfig returns the specified config
	// object which was specified on intialization.
	GetConfig() *Config

	// GetCommandMap returns the internal command
	// map.
	GetCommandMap() map[string]Command

	// GetCommandInstances returns an array of all
	// registered command instances.
	GetCommandInstances() []Command

	// GetCommand returns a command instance form
	// the command register by invoke. If the
	// command could not be found, false is returned.
	GetCommand(invoke string) (Command, bool)
}

// handler is the default implementation of Handler.
type handler struct {
	config       *Config
	cmdMap       map[string]Command
	cmdInstances []Command
	middlewares  []Middleware
}

// NewHandler returns a new instance of the default
// command Handler implementation.
func NewHandler(cfg *Config) Handler {
	if cfg.OnError == nil {
		cfg.OnError = func(Context, ErrorType, error) {}
	}

	if cfg.GuildPrefixGetter == nil {
		cfg.GuildPrefixGetter = func(string) (string, error) {
			return "", nil
		}
	}

	handler := &handler{
		config:       cfg,
		cmdMap:       make(map[string]Command),
		cmdInstances: make([]Command, 0),
	}

	if cfg.UseDefaultHelpCommand {
		handler.RegisterCommand(&defaultHelpCommand{})
	}

	return handler
}

func (h *handler) RegisterCommand(cmd Command) {
	h.cmdInstances = append(h.cmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		if h.config.InvokeToLower {
			invoke = strings.ToLower(invoke)
		}
		if _, ok := h.cmdMap[invoke]; ok {
			panic(fmt.Sprintf("invoke already '%s' already registered", invoke))
		}
		h.cmdMap[invoke] = cmd
	}
}

func (h *handler) RegisterMiddleware(mw Middleware) {
	h.middlewares = append(h.middlewares, mw)
}

func (h *handler) RegisterHandlers(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		h.messageHandler(s, e.Message, false)
	})

	if h.config.ExecuteOnEdit {
		session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageUpdate) {
			h.messageHandler(s, e.Message, false)
		})
	}
}

func (h *handler) GetConfig() *Config {
	return h.config
}

func (h *handler) GetCommandMap() map[string]Command {
	return h.cmdMap
}

func (h *handler) GetCommandInstances() []Command {
	return h.cmdInstances
}

func (h *handler) GetCommand(invoke string) (Command, bool) {
	if h.config.InvokeToLower {
		invoke = strings.ToLower(invoke)
	}

	cmd, ok := h.cmdMap[invoke]
	return cmd, ok
}

// messageHandler is called from the message create and
// message update events of discordgo.
func (h *handler) messageHandler(s *discordgo.Session, msg *discordgo.Message, isEdit bool) {
	if msg.Author == nil || msg.Author.ID == s.State.User.ID {
		return
	}

	if len(msg.Content) < 2 {
		return
	}

	if !h.config.AllowBots && msg.Author.Bot {
		return
	}

	ctx := &context{
		session: s,
		message: msg,
		member:  msg.Member,
		isEdit:  isEdit,
	}

	var err error

	usedPrefix := ""
	if strings.HasPrefix(msg.Content, h.config.GeneralPrefix) {
		usedPrefix = h.config.GeneralPrefix
	} else {
		guildPrefix, err := h.config.GuildPrefixGetter(msg.GuildID)
		if err != nil {
			h.config.OnError(ctx, ErrTypGuildPrefixGetter, err)
			return
		}
		if guildPrefix != "" && strings.HasPrefix(msg.Content, guildPrefix) {
			usedPrefix = guildPrefix
		}
	}

	if usedPrefix == "" {
		return
	}

	if ctx.channel, err = s.State.Channel(msg.ChannelID); err != nil {
		if ctx.channel, err = s.Channel(msg.ChannelID); err != nil {
			h.config.OnError(ctx, ErrTypGetChannel, err)
			return
		}
	}

	ctx.isDM = ctx.channel.Type == discordgo.ChannelTypeDM || ctx.channel.Type == discordgo.ChannelTypeGroupDM
	if !h.config.AllowDM && ctx.isDM {
		return
	}

	if !ctx.isDM {
		if ctx.guild, err = s.State.Guild(msg.GuildID); err != nil {
			if ctx.guild, err = s.Guild(msg.GuildID); err != nil {
				h.config.OnError(ctx, ErrTypGetGuild, err)
				return
			}
		}
	}

	args := argsRx.FindAllString(msg.Content, -1)
	for i, k := range args {
		if strings.Contains(k, "\"") {
			args[i] = strings.Replace(k, "\"", "", -1)
		}
	}

	invoke := args[0][len(usedPrefix):]
	args = args[1:]

	ctx.args = ArgumentList(args)

	cmd, ok := h.GetCommand(invoke)
	if !ok {
		h.config.OnError(ctx, ErrTypCommandNotFound, ErrCommandNotFound)
		return
	}

	if ctx.isDM && !cmd.IsExecutableInDMChannels() {
		h.config.OnError(ctx, ErrTypNotExecutableInDM, ErrCommandNotExecutableInDMs)
		return
	}

	ctx.objectMap = &sync.Map{}
	ctx.Set("cmdhandler", h)

	for _, mw := range h.middlewares {
		err, next := mw.Handle(cmd, ctx)
		if err != nil {
			h.config.OnError(ctx, ErrTypMiddleware, err)
			return
		}
		if !next {
			return
		}
	}

	if err = cmd.Exec(ctx); err != nil {
		h.config.OnError(ctx, ErrTypCommandExec, err)
		return
	}

	if h.config.DeleteMessageAfter {
		if err = s.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
			h.config.OnError(ctx, ErrTypDeleteCommandMessage, err)
			return
		}
	}
}
