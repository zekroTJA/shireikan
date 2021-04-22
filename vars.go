package shireikan

import "errors"

const (
	VERSION = "0.6.0" // The package version.

	ObjectMapKeyHandler = "cmdhandler" // Handler instance object map key.
)

var (
	// Error returned when a command was not found.
	ErrCommandNotFound = errors.New("command not found")
	// Error returned when command was executed in a
	// DM channel which shall not be executed in a DM
	// channel.
	ErrCommandNotExecutableInDMs = errors.New("command is not executable in DM channels")

	EmbedColorDefault = 0x03A9F4 // Default Embed Color
	EmbedColorError   = 0xe53935 // Error Embed Color
)
