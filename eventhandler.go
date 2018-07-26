package main

import (
	"github.com/bwmarrin/discordgo"
)


type Events struct {
	Session *discordgo.Session
}

func (e *Events) Register(handler interface{}) {
	e.Session.AddHandler(handler)
}

func NewEvents(session *discordgo.Session) *Events {
	return &Events{ session }
}
