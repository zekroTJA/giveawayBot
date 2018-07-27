package main

import (
	"regexp"
	"fmt"
	"strconv"
	"time"
	"github.com/bwmarrin/discordgo"
)


func CmdTest(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	return nil
}

func CmdGiveaway(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	member, err := s.GuildMember(g.ID, a.ID)
	if err != nil {
		return err
	}

	if !CheckAutorized(config, member) {
		return nil
	}

	currentStatus := 0
	var content, winMessage string
	var winnerCount int
	var timeout time.Duration

	SendEmbed(s, c.ID, Lang.Commands.Giveaway.EnterContent)

	var remover func()
	remover =  s.AddHandler(func(_ *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.ChannelID != c.ID || msg.Author.ID != a.ID {
			return
		}

		if msg.Content == "exit" {
			SendEmbed(s, c.ID, "Canceled.")
			remover()
			return
		}

		switch currentStatus {

		case 0:
			content = msg.Content
			currentStatus++
			SendEmbed(s, c.ID, Lang.Commands.Giveaway.EnterWinMessage)
		case 1:
			winMessage = msg.Content
			currentStatus++
			SendEmbed(s, c.ID, Lang.Commands.Giveaway.EnterParticipantsNumber)
		case 2:
			winnerCount, err = strconv.Atoi(msg.Content)
			if err != nil {
				SendEmbedError(s, c.ID, Lang.Commands.Giveaway.InvalidInput)
				return
			}
			currentStatus++
			SendEmbed(s, c.ID, Lang.Commands.Giveaway.EnterExpireTime)
		case 3:
			_nr, err := strconv.Atoi(regexp.MustCompile("\\d*").FindString(msg.Content))
			if err != nil {
				SendEmbedError(s, c.ID, Lang.Commands.Giveaway.InvalidInput)
				return
			}
			_type := regexp.MustCompile("[hm]").FindString(msg.Content)
			switch _type {
			case "h":
				timeout = time.Duration(_nr) * time.Hour
			case "m":
				timeout = time.Duration(_nr) * time.Minute
			default:
				SendEmbedError(s, c.ID, Lang.Commands.Giveaway.InvalidInput)
				return
			}
			currentStatus++
			SendEmbed(s, c.ID, Lang.Commands.Giveaway.EnterChannelResolvable)
		case 4:
			channel, err := FetchChannel(g, msg.Content)
			if err != nil {
				SendEmbedError(s, c.ID, Lang.Commands.Giveaway.InvalidInput)
				return
			}
			_, err = NewGiveaway(s, a, channel, winnerCount, content, winMessage, timeout, config.Data.Emote)
			remover()
			if err != nil {
				SendEmbedError(s, c.ID, fmt.Sprintf(Lang.Commands.Giveaway.CreatingFailed, err.Error()))
			} else {
				SendEmbed(s, c.ID, fmt.Sprintf(Lang.Commands.Giveaway.Created, channel.ID))
			}
		}
	})

	return nil
}