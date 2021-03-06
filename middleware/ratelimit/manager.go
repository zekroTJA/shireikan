package ratelimit

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/timedmap"
)

// Manager provides limiter instances by command instance,
// userID and guildID.
type Manager interface {
	// GetLimiter returns a limiter instance from the given
	// cmd instance, userID and guildID.
	//
	// cmd is guaranteed to also implement the LimitedCommand
	// interface when calling this function.
	GetLimiter(cmd shireikan.Command, userID, guildID string) *Limiter
}

type internalManager struct {
	limiters *timedmap.TimedMap
}

func newInternalManager() *internalManager {
	return &internalManager{
		limiters: timedmap.New(10 * time.Minute),
	}
}

func (c *internalManager) GetLimiter(cmd shireikan.Command, userID, guildID string) *Limiter {
	key := fmt.Sprintf("%s:%s:%s", cmd.GetDomainName(), guildID, userID)

	// It is ensured, that the passed command instance also
	// implements LimitedCommand.
	lcmd := cmd.(LimitedCommand)
	expireDuration := time.Duration(lcmd.GetLimiterBurst()) * lcmd.GetLimiterRestoration()

	limiter, ok := c.limiters.GetValue(key).(*Limiter)
	if ok {
		c.limiters.SetExpire(key, expireDuration)
		return limiter
	}

	limiter = NewLimiter(lcmd.GetLimiterBurst(), lcmd.GetLimiterRestoration())
	c.limiters.Set(key, limiter, expireDuration)

	return limiter
}
