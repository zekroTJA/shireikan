package shireikan

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const minPrefixLen = 12

func hasPrefixMention(s *discordgo.Session, content string) (ok bool, prefix string) {
	if len(content) < minPrefixLen || content[0] != '<' || content[1] != '@' {
		return
	}

	cursor := 2
	if content[2] == '!' {
		cursor++
	}

	ci := strings.IndexRune(content, '>')
	if ci < minPrefixLen {
		return
	}

	if content[cursor:ci] != s.State.User.ID {
		return
	}

	ok = true
	prefix = content[0 : ci+1]

	return
}
