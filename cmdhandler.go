package main

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type CmdFunction func(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild)error

type CmdHandler struct {
	Session        *discordgo.Session
	ConfigInstance *Config
	Prefix         string
	Commands       map[string]CmdFunction
}

func NewCmdHandler(session *discordgo.Session, config *Config, prefix string) *CmdHandler {
	return &CmdHandler{ session, config, prefix, map[string]CmdFunction{} }
}

func (c *CmdHandler) Register(invoke string, cmdf CmdFunction) {
	c.Commands[invoke] = cmdf
}

func (c *CmdHandler) Handle(m *discordgo.MessageCreate) {
	channel, err := c.Session.Channel(m.ChannelID)
	if err != nil {
		return
	}
	
	guild, err := c.Session.Guild(channel.GuildID)
	if err != nil {
		return
	}

	if m.Author.ID == c.Session.State.User.ID || channel.Type != discordgo.ChannelTypeGuildText || !strings.HasPrefix(m.Content, c.Prefix) {
		return
	}

	author := m.Author

	_contsplit := strings.Split(m.Content, " ")
	invoke := strings.ToLower(_contsplit[0])[len(c.Prefix):]
	args := _contsplit[1:]

	if cmdf, ok := c.Commands[invoke]; ok {
		fmt.Println(cmdf(c.Session, c.ConfigInstance, args, m, channel, author, guild))
	}
}