package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
)

var _ State = (*Dgrs)(nil)

type Dgrs struct {
	st *dgrs.State
}

func NewDgrs(st *dgrs.State) *Dgrs {
	return &Dgrs{st}
}

func (s *Dgrs) SelfUser(_ *discordgo.Session) (u *discordgo.User, err error) {
	u, err = s.st.SelfUser()
	return
}

func (s *Dgrs) Channel(_ *discordgo.Session, id string) (c *discordgo.Channel, err error) {
	c, err = s.st.Channel(id)
	return
}

func (s *Dgrs) Guild(_ *discordgo.Session, id string) (g *discordgo.Guild, err error) {
	g, err = s.st.Guild(id)
	return
}
