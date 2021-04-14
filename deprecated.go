package shireikan

import "github.com/bwmarrin/discordgo"

// NewHandler IS DEPRECATED!
// Please use New(*Config) instead!
func NewHandler(cfg *Config) Handler { return New(cfg) }

// RegisterHandlers IS DEPRECATED!
// Please use Handler#Setup(*Session) instead!
func (h *handler) RegisterHandlers(session *discordgo.Session) { h.Setup(session) }
