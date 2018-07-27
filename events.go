package main

import (
	"github.com/bwmarrin/discordgo"
)

// ReadyEventHandler will be executed if Ready event fires 
func ReadyEventHandler(s *discordgo.Session, e *discordgo.Ready) {
	LogInfo("Logged in as", e.User.Username + "#" + e.User.Discriminator, "(" + e.User.ID + ")")
	s.UpdateStatus(0, "Created by zekro | zekro.de")
}

// CommandEventHandler will be executed if Ready event fires
func CommandEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Cmd.Handle(m)
}