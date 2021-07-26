package shireikan

import (
	"testing"
)

func TestHasPrefixMention(t *testing.T) {
	const userID = "123456789012345"

	testPrefix := func(msg, prefix string, ok bool) {
		rOk, rPrefix := hasPrefixMention(userID, msg)
		if rOk != ok {
			t.Fatalf("ok was %t (expected: %t)", rOk, ok)
		}
		if rPrefix != prefix {
			t.Fatalf("prefix was %s (expected: %s)", rPrefix, prefix)
		}
	}

	testPrefix("", "", false)
	testPrefix("!a bc", "", false)
	testPrefix("<@123>", "", false)
	testPrefix("<@!123123123>", "", false)
	testPrefix("<@123456789012346>", "", false)
	testPrefix("<@!123456789012346>", "", false)
	testPrefix("<@"+userID, "", false)
	testPrefix("<@!"+userID, "", false)
	testPrefix("<@!"+userID+" asdasdh a asd hasj> >", "", false)

	testPrefix("<@"+userID+"> sad ada dasd ad", "<@"+userID+">", true)
	testPrefix("<@"+userID+">", "<@"+userID+">", true)
	testPrefix("<@!"+userID+"> asdjak jsdh a dasd ahd", "<@!"+userID+">", true)
	testPrefix("<@"+userID+">", "<@"+userID+">", true)
}
