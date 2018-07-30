package main

import (
	"regexp"
	"fmt"
	"strings"
	"strconv"
	"time"
	"io/ioutil"
	"errors"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)


var OpenGiveaways map[string]*Giveaway

// CmdTest - FUnction for Test Command
func CmdTest(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	return nil
}

// CmdInfo - Function for Info command
func CmdInfo(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	embed := &discordgo.MessageEmbed{
		Title: "giveawayBot INFO",
		Color: COLOR_MAIN,
		Description: "© 2018 zekro Development\n[**zekro.de**](https://zekro.de)",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", s.State.User.ID, s.State.User.Avatar),
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Current Version",
				Value: "v." + VERSION,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name: "Licenced Under",
				Value: "MIT",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name: "GitHub Repository",
				Value: ":point_right:  [**zekroTJA/giveawayBot**](https://github.com/zekroTJA/giveawayBot)",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name: "3rd Party Dependencies",
				Value: ":white_small_square:  [discordgo](https://github.com/bwmarrin/discordgo)\n" + 
				       ":white_small_square:  [yaml](https://github.com/go-yaml/yaml)",
				Inline: false,
			},
		},
	}

	_, err := s.ChannelMessageSendEmbed(c.ID, embed)
	return err
}

// CmdHelp - Function for Help cpmmand
func CmdHelp(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {

	helpMsg := 	"**MISC**\n" +
				":white_small_square:  `help`  -  Display this help message\n" + 
				":white_small_square:  `info`  -  Display info about this bot\n" + 
				":white_small_square:  `authroles <role1> <role2>`  -  Add roles to authorized roles\n\n" + 			
				"**GIVEAWAYS**\n" +  
				":white_small_square:  `ga list`    -  List all open giveaways\n" +
				":white_small_square:  `ga close <UID>`    -  Close a giveaway (winners will be selected)\n" +
				":white_small_square:  `ga cancel <UID>`    -  Cancel a giveaway (no winners will be selected)\n"

	embed := &discordgo.MessageEmbed{
		Title: "HELP",
		Color: COLOR_MAIN,
		Description: helpMsg,
	}

	_, err := s.ChannelMessageSendEmbed(c.ID, embed)
	return err
}

// CmdSetAuthRoles - Function for Setauthroles Command
func CmdSetAuthRoles(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	if !CheckAdmin(config, a) {
		return errors.New("NO_PERMISSION")
	}
	
	authRoles := map[string][]string{}
	bData, err := ioutil.ReadFile("./.authroles.json")
	if err == nil {
		err = json.Unmarshal(bData, &authRoles)
	}

	var strRoles []string
	for _, a := range args {
		r, err := FetchRole(g, a)
		if err == nil {
			strRoles = append(strRoles, r.ID)
		}
	}

	authRoles[g.ID] = strRoles

	bData, err = json.Marshal(authRoles)
	if err == nil {
		err = ioutil.WriteFile("./.authroles.json", bData, 0644)
		fmt.Println(err)
	}

	if err == nil {
		SendEmbed(s, c.ID,
			fmt.Sprintf(Lang.Commands.Authrole.Added, func()string {
				marray := make([]string, len(strRoles))
				for i, rid := range strRoles {
					marray[i] = "<@&" + rid + ">"
				}
				return strings.Join(marray, ", ")
			}()))
	}

	return err
}

// CmdGiveaway - function for Giveaway Command
func CmdGiveaway(s *discordgo.Session, config *Config, args []string, m *discordgo.MessageCreate, c *discordgo.Channel, a *discordgo.User, g *discordgo.Guild) error {
	member, err := s.GuildMember(g.ID, a.ID)
	if err != nil {
		return err
	}

	if !CheckAutorized(config, g.ID, member) {
		return errors.New("NO_PERMISSION")
	}

	if len(args) > 0 {

		if args[0] == "list" || args[0] == "ls" {
			embed := &discordgo.MessageEmbed{
				Title: "Open Giveaways",
				Color: COLOR_MAIN,
			}
			for k, v := range OpenGiveaways {
				if v.Guild.ID == g.ID {
					embedText := fmt.Sprintf(
						"**Creator:** <@%s>\n" +
						"**Expires:** `%s`\n" +
						"**Winner Count:** `%d`\n" +
						"**Current Participants:** `%d`\n" +
						"**Content:**\n```\n%s\n```\n" +
						"**Winner Message:**\n```\n%s\n```\n",
						v.Creator.ID, 
						v.Expires.Format(time.RFC1123),
						v.WinnerCount,
						v.ParticipantsNumber,
						v.Content,
						v.WinMessage)
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name: k,
						Value: embedText,
					})
				}
			}
			_, err := s.ChannelMessageSendEmbed(c.ID, embed)
			return err
		}

		if args[0] == "close" || args[0] == "stop" {
			if len(args) < 2 {
				_, err := SendEmbedError(s, c.ID, Lang.Commands.Giveaway.CloseNoID)
				return err
			}
			uid := args[1]
			if ga, ok := OpenGiveaways[uid]; ok {
				if ga.Guild.ID != g.ID {
					_, err = SendEmbedError(s, c.ID, Lang.Commands.Giveaway.WrongGuild)
					return err
				}
				ga.Close(false)
				_, err = SendEmbed(s, c.ID, Lang.Commands.Giveaway.Closed)
			} else {
				_, err = SendEmbedError(s, c.ID, Lang.Commands.Giveaway.CloseInvalidID)
			}
			return err
		}

		if args[0] == "cancel" {
			if len(args) < 2 {
				_, err := SendEmbedError(s, c.ID, Lang.Commands.Giveaway.CloseNoID)
				return err
			}
			uid := args[1]
			if ga, ok := OpenGiveaways[uid]; ok {
				if ga.Guild.ID != g.ID {
					_, err = SendEmbedError(s, c.ID, Lang.Commands.Giveaway.WrongGuild)
					return err
				}
				ga.Close(true)
				_, err = SendEmbed(s, c.ID, Lang.Commands.Giveaway.Closed)
			} else {
				_, err = SendEmbedError(s, c.ID, Lang.Commands.Giveaway.Calceled)
			}
			return err
		}
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
			giveaway, err := NewGiveaway(s, a, channel, winnerCount, content, winMessage, timeout, config.Data.Emote)
			if OpenGiveaways == nil {
				OpenGiveaways = map[string]*Giveaway{ giveaway.UID: giveaway }
			} else {
				OpenGiveaways[giveaway.UID] = giveaway
			}
			remover()
			if err != nil {
				SendEmbedError(s, c.ID, fmt.Sprintf(Lang.Commands.Giveaway.CreatingFailed, err.Error()))
			} else {
				SendEmbed(s, c.ID, fmt.Sprintf(Lang.Commands.Giveaway.Created, giveaway.UID, channel.ID))
			}
		}
	})

	return nil
}