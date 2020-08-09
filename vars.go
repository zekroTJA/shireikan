package shireikan

import "errors"

const (
	VERSION = "0.2.0"

	ObjectMapKeyHandler = "cmdhandler"
)

var (
	ErrCommandNotFound           = errors.New("command not found")
	ErrCommandNotExecutableInDMs = errors.New("command is not executable in DM channels")

	EmbedColorDefault = 0x03A9F4
	EmbedColorError   = 0xe53935
)
