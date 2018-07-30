package main

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

// ReadyEventHandler will be executed if Ready event fires 
func ReadyEventHandler(s *discordgo.Session, e *discordgo.Ready) {
	_guilds, _ := s.UserGuilds(-1, "", "")
	guildNumb := len(_guilds)
	guilds := make([]string, guildNumb)
	for i, g := range _guilds {
		guilds[i] = fmt.Sprintf("  - %s (%s)", g.Name, g.ID)
	}
	LogInfo(fmt.Sprintf(
		"Logged in as %s#%s (%s)\n" + 
		"Invite: https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=486464\n" +
		"Logged in on %d guilds:\n%s", 
		e.User.Username, e.User.Discriminator, e.User.ID, e.User.ID, guildNumb, strings.Join(guilds, "\n")))
	s.UpdateStatus(0, "Created by zekro | zekro.de")
}

// CommandEventHandler will be executed if Ready event fires
func CommandEventHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Cmd.Handle(m)
}