package shireikan

import (
	"strings"
	"sync"
)

const minPrefixLen = 12

func hasPrefixMention(selfUserID string, content string) (ok bool, prefix string) {
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

	if content[cursor:ci] != selfUserID {
		return
	}

	ok = true
	prefix = content[0 : ci+1]

	return
}

func clearMap(m *sync.Map) {
	m.Range(func(key, _ interface{}) bool {
		m.Delete(key)
		return true
	})
}
