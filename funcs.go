package main

import (
	_ "strings"
	"github.com/bwmarrin/discordgo"
)


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