package ratelimit

import (
	"errors"
	"testing"
	"time"

	"github.com/zekroTJA/shireikan"
)

func TestGetLimiter(t *testing.T) {
	m := newInternalManager()

	cmd := &testCmd{false, false, false}

	l1 := m.GetLimiter(cmd, "user1", "guild")
	l2 := m.GetLimiter(cmd, "user1", "guild")
	l3 := m.GetLimiter(cmd, "user2", "guild")

	if l1 != l2 {
		t.Errorf("l1 (%v) != l2 (%v)", l1, l2)
	}

	if l1 == l3 || l2 == l3 {
		t.Error("new limiter was a douplicate")
	}
}

// ----------------

type testCmd struct {
	wasExecuted bool
	fail        bool
	isGlobal    bool
}

func (c *testCmd) GetInvokes() []string {
	return []string{"ping", "p"}
}

func (c *testCmd) GetDescription() string {
	return "ping pong"
}

func (c *testCmd) GetHelp() string {
	return "`ping` - ping"
}

func (c *testCmd) GetGroup() string {
	return shireikan.GroupFun
}

func (c *testCmd) GetDomainName() string {
	return "test.fun.ping"
}

func (c *testCmd) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *testCmd) IsExecutableInDMChannels() bool {
	return true
}

func (c *testCmd) GetLimiterBurst() int {
	return 3
}

func (c *testCmd) GetLimiterRestoration() time.Duration {
	return time.Second
}

func (c *testCmd) IsLimiterGlobal() bool {
	return c.isGlobal
}

func (c *testCmd) Exec(ctx shireikan.Context) error {
	c.wasExecuted = true

	if c.fail {
		return errors.New("test error")
	}

	return nil
}
