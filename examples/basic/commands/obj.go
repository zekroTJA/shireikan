package commands

import (
	"fmt"

	"github.com/zekroTJA/shireikan"
)

type Object struct {
}

func (c *Object) GetInvokes() []string {
	return []string{"object", "obj", "o"}
}

func (c *Object) GetDescription() string {
	return "retrieve the object from the object map set by the middleware"
}

func (c *Object) GetHelp() string {
	return "`object`"
}

func (c *Object) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *Object) GetDomainName() string {
	return "test.etc.object"
}

func (c *Object) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *Object) IsExecutableInDMChannels() bool {
	return true
}

func (c *Object) Exec(ctx shireikan.Context) error {
	obj, _ := ctx.Get("test").(string)
	_, err := ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID,
		fmt.Sprintf("Retrieved Object:\n```%s```", obj))
	return err
}
