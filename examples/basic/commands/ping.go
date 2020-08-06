package commands

import (
	"log"

	"github.com/zekroTJA/shireikan"
)

type Ping struct {
}

func (c *Ping) GetInvokes() []string {
	return []string{"ping", "p"}
}

func (c *Ping) GetDescription() string {
	return "ping pong"
}

func (c *Ping) GetHelp() string {
	return "`ping` - ping"
}

func (c *Ping) GetGroup() string {
	return shireikan.GroupFun
}

func (c *Ping) GetDomainName() string {
	return "test.fun.ping"
}

func (c *Ping) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *Ping) IsExecutableInDMChannels() bool {
	return true
}

func (c *Ping) Exec(ctx shireikan.Context) error {
	log.Printf("%+v", ctx)
	_, err := ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID, "Pong! :ping_pong:")
	return err
}
