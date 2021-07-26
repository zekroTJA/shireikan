// Package shireikan provides a general command
// handler for discordgo.
package shireikan

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shireikan/state"
)

// ErrorType is the type of error occurred in
// the command message handler.
type ErrorType int

const (
	ErrTypGuildPrefixGetter    ErrorType = iota // Error from guild prefix getter function
	ErrTypGetChannel                            // Error getting channel object
	ErrTypGetGuild                              // Error getting guild object
	ErrTypCommandNotFound                       // Command was not found by specified invoke
	ErrTypNotExecutableInDM                     // Command which is specified as non-executable in DM got executed in a DM channel
	ErrTypMiddleware                            // Middleware handler returned an error
	ErrTypCommandExec                           // Command handler returned an error
	ErrTypDeleteCommandMessage                  // Deleting command message failed
	ErrTypeState                                // State action failed
)

var (
	argsRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
)

// Config wraps configuration values for the CommandHandler.
type Config struct {
	GeneralPrefix         string `json:"general_prefix"`           // General and globally accessible prefix
	InvokeToLower         bool   `json:"invoke_to_lower"`          // Lowercase command invoke befor map matching
	AllowDM               bool   `json:"allow_dm"`                 // Allow commands to be executed in DM and GroupDM channels
	AllowBots             bool   `json:"allow_bots"`               // Allow bot accounts to execute commands
	ExecuteOnEdit         bool   `json:"execute_on_edit"`          // Execute command handler when a message was edited
	UseDefaultHelpCommand bool   `json:"use_default_help_command"` // Whether or not to use default help command
	DeleteMessageAfter    bool   `json:"delete_message_after"`     // Delete command message after command has processed

	// Optionally, you can pass a di.Container to obtain
	// instances from in the command context.
	ObjectContainer di.Container `json:"-"`

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

	State state.State
}

// Handler specifies a command register and handler.
type Handler interface {
	ReadonlyObjectMap

	// Register is shorthand for RegisterMiddleware
	// or RegisterCommand and automatically choses
	// depending on the implementation the required
	// registration method.
	//
	// Panics if an instance is passed which neither
	// implements Command nor Middleware interface.
	Register(v interface{})

	// RegisterCommand registers the passed
	// Command instance.
	RegisterCommand(cmd Command)

	// RegisterMiddleware registers the
	// passed middleware instance.
	RegisterMiddleware(mw Middleware)

	// RegisterHandlers IS DEPRECATED!
	// Please use Setup(*Session) instead!
	//
	// RegisterHandlers registers the message
	// handlers to the passed discordgo.Session
	// which are used to handle and parse commands.
	RegisterHandlers(session *discordgo.Session)

	// Setup registers the message handlers to
	// the passed discordgo.Session which are
	// used to handle and parse commands.
	Setup(session *discordgo.Session)

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
	config          *Config
	cmdMap          map[string]Command
	cmdInstances    []Command
	middlewares     []Middleware
	objectContainer di.Container
	ctxPool         *sync.Pool
	state           state.State

	// DEPRECATED: will be removed on next update!
	objectMap *sync.Map
}

// New returns a new instance of the default
// command Handler implementation.
func New(cfg *Config) Handler {
	if cfg.OnError == nil {
		cfg.OnError = func(Context, ErrorType, error) {}
	}

	if cfg.GuildPrefixGetter == nil {
		cfg.GuildPrefixGetter = func(string) (string, error) {
			return "", nil
		}
	}

	if cfg.State == nil {
		cfg.State = state.NewInternal()
	}

	handler := &handler{
		config:          cfg,
		cmdMap:          make(map[string]Command),
		cmdInstances:    make([]Command, 0),
		objectContainer: cfg.ObjectContainer,
		ctxPool: &sync.Pool{
			New: func() interface{} {
				return &context{
					objectMap: &sync.Map{},
				}
			},
		},
		objectMap: &sync.Map{},
		state:     cfg.State,
	}

	if handler.objectContainer == nil {
		builder, _ := di.NewBuilder()
		handler.objectContainer = builder.Build()
	}

	if cfg.UseDefaultHelpCommand {
		handler.RegisterCommand(&defaultHelpCommand{})
	}

	return handler
}

func (h *handler) Register(v interface{}) {
	switch vi := v.(type) {
	case Command:
		h.RegisterCommand(vi)
	case Middleware:
		h.RegisterMiddleware(vi)
	default:
		panic("instance does neither implement Command nor Middleware interface")
	}
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

func (h *handler) Setup(session *discordgo.Session) {
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

func (h *handler) GetObject(key string) interface{} {
	val, err := h.objectContainer.SafeGet(key)
	if err != nil {
		val, _ = h.objectMap.Load(key)
	}

	return val
}

func (h *handler) SetObject(key string, val interface{}) {
	h.objectMap.Store(key, val)
}

// messageHandler is called from the message create and
// message update events of discordgo.
func (h *handler) messageHandler(s *discordgo.Session, msg *discordgo.Message, isEdit bool) {
	ctx := h.ctxPool.Get().(*context)
	ctx.handler = h
	ctx.session = s
	ctx.message = msg
	ctx.member = msg.Member
	ctx.isEdit = isEdit
	defer func() {
		clearMap(ctx.objectMap)
		h.ctxPool.Put(ctx)
	}()

	self, err := h.state.SelfUser(s)
	if err != nil {
		h.config.OnError(ctx, ErrTypeState, err)
		return
	}

	if msg.Author == nil || msg.Author.ID == self.ID {
		return
	}

	if len(msg.Content) < 1 {
		return
	}

	if !h.config.AllowBots && msg.Author.Bot {
		return
	}

	usedPrefix := ""
	if strings.HasPrefix(msg.Content, h.config.GeneralPrefix) {
		usedPrefix = h.config.GeneralPrefix
	} else if ok, prefix := hasPrefixMention(self.ID, msg.Content); ok {
		usedPrefix = prefix
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

	if ctx.channel, err = s.State.Channel(msg.ChannelID); err != nil {
		if ctx.channel, err = s.Channel(msg.ChannelID); err != nil {
			h.config.OnError(ctx, ErrTypGetChannel, err)
			return
		}
	}

	if usedPrefix == "" && ctx.channel.Type != discordgo.ChannelTypeDM {
		return
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

	content := msg.Content[len(usedPrefix):]

	args := argsRx.FindAllString(content, -1)

	if len(args) == 0 {
		return
	}

	for i, k := range args {
		if strings.Contains(k, "\"") {
			args[i] = strings.Replace(k, "\"", "", -1)
		}
	}

	invoke := args[0]
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

	ctx.SetObject(ObjectMapKeyHandler, h)

	if !h.executeMiddlewares(cmd, ctx, LayerBeforeCommand) {
		return
	}

	if err = cmd.Exec(ctx); err != nil {
		h.config.OnError(ctx, ErrTypCommandExec, err)
		return
	}

	if !h.executeMiddlewares(cmd, ctx, LayerAfterCommand) {
		return
	}

	if h.config.DeleteMessageAfter {
		if err = s.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
			h.config.OnError(ctx, ErrTypDeleteCommandMessage, err)
			return
		}
	}
}

func (h *handler) executeMiddlewares(cmd Command, ctx Context, layer MiddlewareLayer) bool {
	for _, mw := range h.middlewares {
		if mw.GetLayer()&layer == 0 {
			continue
		}

		next, err := mw.Handle(cmd, ctx, layer)
		if err != nil {
			h.config.OnError(ctx, ErrTypMiddleware, err)
			return false
		}
		if !next {
			return false
		}
	}

	return true
}
