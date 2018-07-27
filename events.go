package main

import (
	"github.com/bwmarrin/discordgo"
)


func ReadyEventHandler(s *discordgo.Session, e *discordgo.Ready) {
	LogInfo("Logged in as", e.User.Username + "#" + e.User.Discriminator, "(" + e.User.ID + ")")
	s.UpdateStatus(0, "Created by zekro | zekro.de")
}

func CommandEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Cmd.Handle(m)
}