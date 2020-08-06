package shireikan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type defaultHelpCommand struct {
}

func (c *defaultHelpCommand) GetInvokes() []string {
	return []string{"help", "h", "?", "man"}
}

func (c *defaultHelpCommand) GetDescription() string {
	return "display list of command or get help for a specific command"
}

func (c *defaultHelpCommand) GetHelp() string {
	return "`help` - display command list\n" +
		"`help <command>` - display help of specific command"
}

func (c *defaultHelpCommand) GetGroup() string {
	return GroupGeneral
}

func (c *defaultHelpCommand) GetDomainName() string {
	return "sp.etc.help"
}

func (c *defaultHelpCommand) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *defaultHelpCommand) IsExecutableInDMChannels() bool {
	return true
}

func (c *defaultHelpCommand) Exec(ctx Context) error {
	emb := &discordgo.MessageEmbed{
		Color:  EmbedColorDefault,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	handler, _ := ctx.Get("cmdhandler").(Handler)

	if len(ctx.GetArgs()) == 0 {
		cmds := make(map[string][]Command)
		for _, c := range handler.GetCommandInstances() {
			group := c.GetGroup()
			if _, ok := cmds[group]; !ok {
				cmds[group] = make([]Command, 0)
			}
			cmds[group] = append(cmds[group], c)
		}

		emb.Title = "Command List"

		for cat, catCmds := range cmds {
			commandHelpLines := ""
			for _, c := range catCmds {
				commandHelpLines += fmt.Sprintf("`%s` - *%s* `[%s]`\n", c.GetInvokes()[0], c.GetDescription(), c.GetDomainName())
			}
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
				Name:  cat,
				Value: commandHelpLines,
			})
		}
	} else {
		invoke := ctx.GetArgs().Get(0).AsString()
		cmd, ok := handler.GetCommand(invoke)
		if !ok {
			_, err := ctx.ReplyEmbedError(
				fmt.Sprintf("No command was found with the invoke `%s`.", invoke), "Error")
			return err
		}

		emb.Title = "Command Description"

		description := cmd.GetDescription()
		if description == "" {
			description = "`no description`"
		}

		help := cmd.GetHelp()
		if help == "" {
			help = "`no uage information`"
		}

		emb.Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Invokes",
				Value:  strings.Join(cmd.GetInvokes(), "\n"),
				Inline: true,
			},
			{
				Name:   "Group",
				Value:  cmd.GetGroup(),
				Inline: true,
			},
			{
				Name:   "Domain Name",
				Value:  cmd.GetDomainName(),
				Inline: true,
			},
			{
				Name:   "DM Capable",
				Value:  strconv.FormatBool(cmd.IsExecutableInDMChannels()),
				Inline: true,
			},
			{
				Name:  "Description",
				Value: description,
			},
			{
				Name:  "Usage",
				Value: help,
			},
		}

		if spr := cmd.GetSubPermissionRules(); spr != nil {
			txt := "*`[E]` in front of permissions means `Explicit`, which means that this " +
				"permission must be explicitly allowed and can not be wild-carded.\n" +
				"`[D]` implies that wildecards will apply to this sub permission.*\n\n"

			for _, rule := range spr {
				expl := "D"
				if rule.Explicit {
					expl = "E"
				}

				txt = fmt.Sprintf("%s`[%s]` %s.%s - *%s*\n",
					txt, expl, cmd.GetDomainName(), rule.Term, rule.Description)
			}

			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
				Name:  "Sub Permission Rules",
				Value: txt,
			})
		}
	}

	userChan, err := ctx.GetSession().UserChannelCreate(ctx.GetUser().ID)
	if err != nil {
		return err
	}
	_, err = ctx.GetSession().ChannelMessageSendEmbed(userChan.ID, emb)
	if err != nil {
		if strings.Contains(err.Error(), `{"code": 50007, "message": "Cannot send messages to this user"}`) {
			emb.Footer = &discordgo.MessageEmbedFooter{
				Text: "Actually, this message appears in your DM, but you have deactivated receiving DMs from" +
					"server members, so I can not send you this message via DM and you see this here right now.",
			}
			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
			return err
		}
	}

	return err
}
