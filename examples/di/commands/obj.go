package commands

import (
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/shireikan/examples/di/database"
)

// Object is a command returning an object
// recovered from the context's object map.
type Object struct {
}

// GetInvoke returns the command invokes.
func (c *Object) GetInvokes() []string {
	return []string{"object", "obj", "o"}
}

// GetDescription returns the commands description.
func (c *Object) GetDescription() string {
	return "retrieve the object from the object map set by the middleware"
}

// GetHelp returns the commands help text.
func (c *Object) GetHelp() string {
	return "`object`"
}

// GetGroup returns the commands group.
func (c *Object) GetGroup() string {
	return shireikan.GroupEtc
}

// GetDomainName returns the commands domain name.
func (c *Object) GetDomainName() string {
	return "test.etc.object"
}

// GetSubPermissionRules returns the commands sub
// permissions array.
func (c *Object) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

// IsExecutableInDMChannels returns whether
// the command is executable in DM channels.
func (c *Object) IsExecutableInDMChannels() bool {
	return true
}

// Exec is the commands execution handler.
func (c *Object) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject("db").(database.Database)

	_, err := ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID,
		db.GetData())
	return err
}
