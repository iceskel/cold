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
	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(genericPluginHandler.Config.Channel)
	})
	c.HandleFunc("PRIVMSG", genericPluginHandler.TweetHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.ListCommandsHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.AddCommandHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.DeleteCommandHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.CommandHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.UptimeHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.UpdateChannelGameHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.UpdateChannelStatusHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.SongHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.AddTimeoutListHandler)
	c.HandleFunc("PRIVMSG", genericPluginHandler.TimeoutHandler)
	c.HandleFunc("PRIVMSG", windowsPluginHandler.Foobar2kHandler)

	quit := make(chan bool)
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.ConnectTo("irc.twitch.tv", genericPluginHandler.Config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", genericPluginHandler.Config.Channel)

	<-quit
}
