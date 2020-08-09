package commands

import (
	"log"

	"github.com/zekroTJA/shireikan"
)

// Ping is a command responding with a ping
// message in the commands channel.
type Ping struct {
}

// GetInvoke returns the command invokes.
func (c *Ping) GetInvokes() []string {
	return []string{"ping", "p"}
}

// GetDescription returns the commands description.
func (c *Ping) GetDescription() string {
	return "ping pong"
}

// GetHelp returns the commands help text.
func (c *Ping) GetHelp() string {
	return "`ping` - ping"
}

// GetGroup returns the commands group.
func (c *Ping) GetGroup() string {
	return shireikan.GroupFun
}

// GetDomainName returns the commands domain name.
func (c *Ping) GetDomainName() string {
	return "test.fun.ping"
}

// GetSubPermissionRules returns the commands sub
// permissions array.
func (c *Ping) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

// IsExecutableInDMChannels returns whether
// the command is executable in DM channels.
func (c *Ping) IsExecutableInDMChannels() bool {
	return true
}

// Exec is the commands execution handler.
func (c *Ping) Exec(ctx shireikan.Context) error {
	log.Printf("%+v", ctx)
	_, err := ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID, "Pong! :ping_pong:")
	return err
}
