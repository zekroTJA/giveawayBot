package main

import (
	"strings"
	"errors"
	"github.com/bwmarrin/discordgo"
)

//////////////// PUBLIC STUFF ////////////////

func SendEmbed(session *discordgo.Session, channelID string, cont string) (*discordgo.Message, error) {
	return session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Description: cont,
		Color: COLOR_MAIN,
	})
}

func SendEmbedError(session *discordgo.Session, channelID string, cont string) (*discordgo.Message, error) {
	return session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Description: cont,
		Color: COLOR_ERROR,
	})
}

func CheckAutorized(config *Config, member *discordgo.Member) bool {
	for _, r := range member.Roles {
		for _, a := range config.Data.Authorized {
			if r == a {
				return true
			}
		}
	}
	return false
}

func FetchChannel(g *discordgo.Guild, resolvable string) (*discordgo.Channel, error) {
	var channels []*discordgo.Channel
	for _, c := range g.Channels {
		if c.Type == discordgo.ChannelTypeGuildText {
			channels = append(channels, c)
		}
	} 
	for _, c := range channels {
		if c.ID == strings.Replace(resolvable, "<#", ">", -1) {
			return c, nil
		}
	}
	for _, c := range channels {
		if strings.ToLower(c.Name) == strings.ToLower(resolvable) {
			return c, nil
		}
	}
	for _, c := range channels {
		if strings.ToLower(c.Name) == strings.ToLower(resolvable) {
			return c, nil
		}
	}
	for _, c := range channels {
		if strings.HasPrefix(strings.ToLower(c.Name), strings.ToLower(resolvable)) {
			return c, nil
		}
	}
	return nil, errors.New("channel not found")
}