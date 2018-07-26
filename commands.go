package main

import (
	"regexp"
	"fmt"
	"strconv"
	"time"
	"github.com/bwmarrin/discordgo"
)


func CmdTest(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	ga, _ := NewGiveaway(s, a, c, 1, "test content", "", time.Duration(10) * time.Second, config.Data.Emote)
	// user, _ := s.User("221905671296253953")
	// ga.Participants["1"] = user
	// ga.Participants["2"] = user
	// ga.Participants["3"] = user
	fmt.Println(ga)
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

	SendEmbed(s, c.ID, "**Enter the content message of the giveaway:**")

	var remover func()
	remover =  s.AddHandler(func(_ *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.ChannelID != c.ID || msg.Author.ID != a.ID {
			return
		}

		switch currentStatus {

		case 0:
			content = msg.Content
			currentStatus++
			SendEmbed(s, c.ID, "**Enter the message, which will appear in the direct message of the winner after expire:**")
		case 1:
			winMessage = msg.Content
			currentStatus++
			SendEmbed(s, c.ID, "**Now, enter the number of participant who can win in the giveaway:**")
		case 2:
			winnerCount, err = strconv.Atoi(msg.Content)
			if err != nil {
				SendEmbedError(s, c.ID, "Invalid input.\n**Please enter again:**")
				return
			}
			currentStatus++
			SendEmbed(s, c.ID, "**Now, enter the expire time of the giveaway:**\n*(i.e. '30m', '4h' or '48 h')*")
		case 3:
			_nr, err := strconv.Atoi(regexp.MustCompile("\\d*").FindString(msg.Content))
			if err != nil {
				SendEmbedError(s, c.ID, "Invalid input.\n**Please enter again:**")
				return
			}
			_type := regexp.MustCompile("[hm]").FindString(msg.Content)
			switch _type {
			case "h":
				timeout = time.Duration(_nr) * time.Hour
			case "m":
				timeout = time.Duration(_nr) * time.Minute
			default:
				SendEmbedError(s, c.ID, "Invalid input.\n**Please enter again:**")
				return
			}
			ga, err := NewGiveaway(s, a, c, winnerCount, content, winMessage, timeout, config.Data.Emote)
			fmt.Println(ga, err)
			remover()
		}
	})

	return nil
}