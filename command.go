package shireikan

// SubPermission wraps information about
// a sub permissions of commands.
type SubPermission struct {
	Term        string
	Explicit    bool
	Description string
}

// Command describes the functionalities of a
// command struct which can be registered
// in the CommandHandler.
type Command interface {
	// GetInvokes returns the unique strings udes to
	// call the command. The first invoke is the
	// primary command invoke and each following is
	// treated as command alias.
	GetInvokes() []string

	// GetDescription returns a brief description about
	// the functionality of the command.
	GetDescription() string

	// GetHelp returns detailed information on how to
	// use the command and their sub commands.
	GetHelp() string

	// GetGroup returns the group name of the command.
	GetGroup() string

	// GetDomainName returns the commands domain name.
	// The domain name is specified like following:
	//   sp.{group}(.{subGroup}...).{primaryInvoke}
	GetDomainName() string

	// GetSubPermissionRules returns optional sub
	// permissions of the command.
	GetSubPermissionRules() []SubPermission

	// IsExecutableInDMChannels returns true when
	// the command can be used in DM channels;
	// otherwise returns false.
	IsExecutableInDMChannels() bool

	// Exec is called when the command is executed and
	// is getting passed the command CommandArgs.
	// When the command was executed successfully, it
	// should return nil. Otherwise, the error
	// encountered should be returned.
	Exec(ctx Context) error
}

const (
	GroupGlobalAdmin = "GLOBAL ADMIN" // Global Admin Group
	GroupGuildAdmin  = "GUILD ADMIN"  // Guild Admin Group
	GroupModeration  = "MODERATION"   // Moderation Untilities Group
	GroupFun         = "FUN"          // Fun Group
	GroupGame        = "GAME"         // Game Group
	GroupChat        = "CHAT"         // Chat Group
	GroupEtc         = "ETC"          // Etc. Group
	GroupGeneral     = "GENERAL"      // General Group
	GroupGuildConfig = "GUILD CONFIG" // Guild Config Group
)
