package main

import (
	"fmt"
	"strings"
	"math/rand"
	"time"
	"github.com/bwmarrin/discordgo"
)


type Giveaway struct {
	UID          	    string
	Session				*discordgo.Session
	Creator      	    *discordgo.User
	Message      	    *discordgo.Message
	Channel      	    *discordgo.Channel
	Content      	    string
	WinnerCount  	    int
	WinMessage   	    string
	Timeout      	    time.Duration
	Expires 			time.Time
	Timer 		 	    *time.Timer
	HandlerRemover      func()
	Participants        map[string]*discordgo.User
	ParticipantsNumber  int
}

// NewGiveaway creates new instance of Gieaway
func NewGiveaway(session *discordgo.Session, creator *discordgo.User, channel *discordgo.Channel, winnerCount int, content, winMessage string, timeout time.Duration, emote string) (*Giveaway, error) {

	var giveaway *Giveaway
	
	expires := time.Now().Add(timeout).Format(time.RFC1123)

	embed := &discordgo.MessageEmbed{
		Title:  		Lang.Classes.Giveaway.ActiveMessage.Title,
		Description:	content + "\n\n" + Lang.Classes.Giveaway.ActiveMessage.ParticipateInfo,
		Color: 			COLOR_MAIN,
		Footer: &discordgo.MessageEmbedFooter{
			Text: 		fmt.Sprintf(Lang.Classes.Giveaway.ActiveMessage.Expires, expires),
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
		giveaway.Close(false)
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
				SendEmbedError(session, pchan.ID, Lang.Classes.Giveaway.Notifications.MultiParticipation)
				session.MessageReactionRemove(giveaway.Channel.ID, giveaway.UID, e.Emoji.Name, e.UserID)
			}
			return
		}
		giveaway.Participants[e.UserID], err = session.User(e.UserID)
		giveaway.ParticipantsNumber++
		if err == nil {
			pchan, err := session.UserChannelCreate(e.UserID)
			if err == nil {
				SendEmbed(session, pchan.ID, Lang.Classes.Giveaway.Notifications.Participated)
			}
		}
		session.MessageReactionRemove(giveaway.Channel.ID, giveaway.UID, e.Emoji.Name, e.UserID)
	})

	giveaway = &Giveaway{
		UID:		  	    message.ID,
		Session:			session,
		Creator:      	    creator,
		Message:      	    message,
		Channel:      	    channel,
		Content:      	    content,
		WinnerCount:  	    winnerCount,
		WinMessage:   	    winMessage,
		Timeout:	  	    timeout,
		Expires:			time.Now().Add(timeout),
		Timer:		  	    timer,
		HandlerRemover:     remover,
		Participants: 	    map[string]*discordgo.User{},
		ParticipantsNumber: 0,
	}

	return giveaway, nil
} 

func (giveaway *Giveaway) Close(cancel bool) {
	giveaway.HandlerRemover()
	
	delete(OpenGiveaways, giveaway.UID)

	if cancel {
		giveaway.Timer.Stop()
		giveaway.Session.ChannelMessageDelete(giveaway.Channel.ID, giveaway.Message.ID)
		return
	}

	if len(giveaway.Participants) < giveaway.WinnerCount {
		privatechan, err := giveaway.Session.UserChannelCreate(giveaway.Creator.ID)
		if err != nil {
			return
		}
		SendEmbedError(giveaway.Session, privatechan.ID,
			fmt.Sprintf(Lang.Classes.Giveaway.CreatorDM.NoParticipations, giveaway.UID))
				
		editembed := &discordgo.MessageEmbed{
			Title:  		Lang.Classes.Giveaway.ClosedMessage.Title,
			Description:	giveaway.Content + "\n\n" + Lang.Classes.Giveaway.ClosedMessage.NoParticipants,
			Color: 			COLOR_CLOSED,
			Footer: &discordgo.MessageEmbedFooter{
				Text: 		Lang.Classes.Giveaway.ClosedMessage.Expired,
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:		giveaway.Creator.Username,
				IconURL:	fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", giveaway.Creator.ID, giveaway.Creator.Avatar),
			},
		}
		
		giveaway.Session.ChannelMessageEditEmbed(giveaway.Channel.ID, giveaway.Message.ID, editembed)
		giveaway.Session.MessageReactionsRemoveAll(giveaway.Channel.ID, giveaway.Message.ID)

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

	winnerNames := make([]string, giveaway.WinnerCount)
	winners := make([]*discordgo.User, giveaway.WinnerCount)
	for i, _ := range winners {
		rnumb := randInt(0, len(giveaway.Participants) - 1)
		winners[i] = giveaway.Participants[keys[rnumb]]
		delete(giveaway.Participants, keys[rnumb])
	}

	for i, w := range winners {
		winnerNames[i] = w.Username
		privatechan, err := giveaway.Session.UserChannelCreate(w.ID)
		if err != nil {
			continue
		}
		SendEmbed(giveaway.Session, privatechan.ID, giveaway.WinMessage)
	}

	editembed := &discordgo.MessageEmbed{
		Title:  		Lang.Classes.Giveaway.ClosedMessage.Title,
		Description:	giveaway.Content + "\n\n" + fmt.Sprintf(Lang.Classes.Giveaway.ClosedMessage.Winners, strings.Join(winnerNames, ", ")),
		Color: 			COLOR_CLOSED,
		Footer: &discordgo.MessageEmbedFooter{
			Text: 		"Expired",
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:		giveaway.Creator.Username,
			IconURL:	fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", giveaway.Creator.ID, giveaway.Creator.Avatar),
		},
	}

	giveaway.Session.ChannelMessageEditEmbed(giveaway.Channel.ID, giveaway.Message.ID, editembed)
	giveaway.Session.MessageReactionsRemoveAll(giveaway.Channel.ID, giveaway.Message.ID)

	privatechan, err := giveaway.Session.UserChannelCreate(giveaway.Creator.ID)
	if err != nil {
		return
	}
	SendEmbed(giveaway.Session, privatechan.ID, fmt.Sprintf(
		Lang.Classes.Giveaway.CreatorDM.Final,
		giveaway.UID, giveaway.ParticipantsNumber, strings.Join(winnerNames, ", ")))
}

func randInt(min, max int) int {
	if min - max == 0 {
		return 0
	}
    return min + rand.Intn(max-min)
}