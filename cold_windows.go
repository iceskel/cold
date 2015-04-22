// +build windows

package main

import (
	"flag"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/cold/handlers"
	windowsHandlers "github.com/iceskel/cold/handlers/windows"
	"log"
)

func main() {

	configFile := flag.String("c", "conf.json", "config file")
	flag.Parse()

	genericPluginHandler, err := handlers.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	windowsPluginHandler, err := windowsHandlers.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	initIrcClient(genericPluginHandler, windowsPluginHandler)

}

func initIrcClient(genericPluginHandler *handlers.BotHandler,
	windowsPluginHandler *windowsHandlers.WindowsBotHandler) {
	c := irc.SimpleClient(genericPluginHandler.Config.Botname, genericPluginHandler.Config.Botname, "simple bot")
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(genericPluginHandler.Config.Channel)
	})
	c.AddHandler("PRIVMSG", genericPluginHandler.TweetHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.ListCommandsHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.AddCommandHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.CommandHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.UptimeHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.UpdateChannelGameHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.UpdateChannelStatusHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.SongHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.AddTimeoutListHandler)
	c.AddHandler("PRIVMSG", genericPluginHandler.TimeoutHandler)
	c.AddHandler("PRIVMSG", windowsPluginHandler.Foobar2kHandler)

	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.Connect("irc.twitch.tv", genericPluginHandler.Config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", genericPluginHandler.Config.Channel)

	<-quit
}
