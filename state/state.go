package state

import "github.com/bwmarrin/discordgo"

type State interface {
	SelfUser(s *discordgo.Session) (*discordgo.User, error)
	Channel(s *discordgo.Session, id string) (*discordgo.Channel, error)
	Guild(s *discordgo.Session, id string) (*discordgo.Guild, error)
}
