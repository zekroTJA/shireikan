package shireikan

import "errors"

var (
	ErrCommandNotFound           = errors.New("command not found")
	ErrCommandNotExecutableInDMs = errors.New("command is not executable in DM channels")
)
