package main

import (
	//"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var Cmd *CmdHandler
var Lang *Language

func main() {

	LogInfo("Starting up...")

	config, err := NewConfig("config.yaml")
	CheckError(err, false)

	Lang, err = NewLanguage(config)
	CheckError(err, false)

	session, err := discordgo.New("Bot " + config.Data.Token)
	CheckError(err, false)

	//////////// COMMAND REGISTRATION ///////////
	Cmd = NewCmdHandler(session, config, config.Data.Prefix)
	Cmd.Register("test", CmdTest)
	Cmd.Register("ga", CmdGiveaway)
	////////////////////////////////////////////

	///////////// EVENT REGISTRATION ////////////
	event := NewEvents(session)
	event.Register(ReadyEventHandler)
	event.Register(CommandEventHandler)
	////////////////////////////////////////////

	err = session.Open()
	CheckError(err, false)

	LogInfo("Logged in. Waiting for response...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}