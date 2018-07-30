package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
)

// CmdFunction is the type defining command executable functions saved in the command handler
type CmdFunction func(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild)error

// CmdHandler saves the session and config instance, the prefix and all registered command functions
type CmdHandler struct {
	Session        *discordgo.Session
	ConfigInstance *Config
	Prefix         string
	Commands       map[string]CmdFunction
}

// NewCmdHandler creates a new instance of CmdHandler with discord session instance, config instance
// and command prefix as parameters.
func NewCmdHandler(session *discordgo.Session, config *Config, prefix string) *CmdHandler {
	return &CmdHandler{ session, config, prefix, map[string]CmdFunction{} }
}

// Register appends a command executable function with the invoke in the command
// map of the CmdHandler
func (c *CmdHandler) Register(invoke string, cmdf CmdFunction) {
	c.Commands[invoke] = cmdf
}

// Handle is the event handler figuring out if the executed message is a command and
// prepares all data to pass them to the command executable function.
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
		err = cmdf(c.Session, c.ConfigInstance, args, m, channel, author, guild)
		if err != nil {
			SendEmbedError(c.Session, channel.ID, "**Error:**\n```\n" + err.Error() + "\n```")
		}
	}
}