package main

import (
	"fmt"
	"strings"
	"math/rand"
	"time"
	"github.com/bwmarrin/discordgo"
)


type Giveaway struct {
	UID          	string
	Creator      	*discordgo.User
	Message      	*discordgo.Message
	Channel      	*discordgo.Channel
	Content      	string
	WinnerCount  	int
	WinMessage   	string
	Timeout      	time.Duration
	Timer 		 	*time.Timer
	HandlerRemover  func()
	Participants    map[string]*discordgo.User
}


func NewGiveaway(session *discordgo.Session, creator *discordgo.User, channel *discordgo.Channel, winnerCount int, content, winMessage string, timeout time.Duration, emote string) (*Giveaway, error) {

	var giveaway *Giveaway
	
	expires := time.Now().Add(timeout).Format(time.RFC1123)

	embed := &discordgo.MessageEmbed{
		Title:  		"OPEN GIVEAWAY",
		Description:	content + "\n\n*Participate to this Giveaway by reacting to this message below.*",
		Color: 			COLOR_MAIN,
		Footer: &discordgo.MessageEmbedFooter{
			Text: 		"Expires on " + expires,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:		creator.Username,
			IconURL:	fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", creator.ID, creator.Avatar),
		},
	}

	message, err := session.ChannelMessageSendEmbed(channel.ID, embed)
	if err != nil {
		return nil, err
	}
	session.MessageReactionAdd(channel.ID, message.ID, emote)

	timer := time.NewTimer(timeout)
	go func() {
		<-timer.C

		giveaway.HandlerRemover()

		if len(giveaway.Participants) < winnerCount {
			privatechan, err := session.UserChannelCreate(creator.ID)
			if err != nil {
				return
			}
			SendEmbedError(session, privatechan.ID,
				fmt.Sprintf("Giveaway *(ID: `%s`)* ended with no result because to less people participated to it.", giveaway.UID))

			editembed := &discordgo.MessageEmbed{
				Title:  		"GIVEAWAY CLOSED",
				Description:	content + "\n\n**No winners. To less participants.**",
				Color: 			COLOR_CLOSED,
				Footer: &discordgo.MessageEmbedFooter{
					Text: 		"Expired",
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:		creator.Username,
					IconURL:	fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", creator.ID, creator.Avatar),
				},
			}
		
			session.ChannelMessageEditEmbed(giveaway.Channel.ID, giveaway.Message.ID, editembed)
			session.MessageReactionsRemoveAll(giveaway.Channel.ID, giveaway.Message.ID)

			return
		}

		rand.Seed(time.Now().UTC().UnixNano())

		keys := make([]string, len(giveaway.Participants))
		i := 0
		for k := range giveaway.Participants {
			keys[i] = k
			i++
		}

		fmt.Println(giveaway.Participants)

		winnerNames := make([]string, winnerCount)
		winners := make([]*discordgo.User, winnerCount)
		for i, _ := range winners {
			rnumb := randInt(0, len(giveaway.Participants) - 1)
			winners[i] = giveaway.Participants[keys[rnumb]]
			delete(giveaway.Participants, keys[rnumb])
		}

		for i, w := range winners {
			winnerNames[i] = w.Username
			privatechan, err := session.UserChannelCreate(w.ID)
			if err != nil {
				continue
			}
			SendEmbed(session, privatechan.ID, giveaway.WinMessage)
		}

		editembed := &discordgo.MessageEmbed{
			Title:  		"GIVEAWAY CLOSED",
			Description:	content + "\n\n**Winners: ``" + strings.Join(winnerNames, ", ") + "``**",
			Color: 			COLOR_CLOSED,
			Footer: &discordgo.MessageEmbedFooter{
				Text: 		"Expired",
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:		creator.Username,
				IconURL:	fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", creator.ID, creator.Avatar),
			},
		}

		session.ChannelMessageEditEmbed(giveaway.Channel.ID, giveaway.Message.ID, editembed)
		session.MessageReactionsRemoveAll(giveaway.Channel.ID, giveaway.Message.ID)
	}()

	remover := session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		if e.MessageID != giveaway.UID || e.UserID == session.State.User.ID {
			return
		}
		if e.Emoji.Name != emote {
			session.MessageReactionRemove(giveaway.Channel.ID, giveaway.UID, e.Emoji.Name, e.UserID)
			return
		}
		if _, ok := giveaway.Participants[e.UserID]; ok {
			pchan, err := session.UserChannelCreate(e.UserID)
			if err == nil {
				SendEmbedError(session, pchan.ID, "You can only participate once on a giveaway!")
				session.MessageReactionRemove(giveaway.Channel.ID, giveaway.UID, e.Emoji.Name, e.UserID)
			}
			return
		}
		giveaway.Participants[e.UserID], err = session.User(e.UserID)
		if err == nil {
			pchan, err := session.UserChannelCreate(e.UserID)
			if err == nil {
				SendEmbed(session, pchan.ID, "You have participated to the giveaway.\nIf you will win, you will get a notification via direct message.")
			}
		}
		session.MessageReactionRemove(giveaway.Channel.ID, giveaway.UID, e.Emoji.Name, e.UserID)
	})

	giveaway = &Giveaway{
		UID:		  	message.ID,
		Creator:      	creator,
		Message:      	message,
		Channel:      	channel,
		Content:      	content,
		WinnerCount:  	winnerCount,
		WinMessage:   	winMessage,
		Timeout:	  	timeout,
		Timer:		  	timer,
		HandlerRemover: remover,
		Participants: 	map[string]*discordgo.User{},
	}

	return giveaway, nil
} 


func randInt(min, max int) int {
	if min - max == 0 {
		return 0
	}
    return min + rand.Intn(max-min)
}