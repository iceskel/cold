// +build linux darwin freebsd openbsd

package main

import (
	"flag"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/cold/handlers"
	"log"
)

func main() {
	configFile := flag.String("c", "conf.json", "config file")
	flag.Parse()

	pluginHandler, err := handlers.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	initIrcClient(pluginHandler)

}

func initIrcClient(pluginHandler *handlers.BotHandler) {
	c := irc.SimpleClient(pluginHandler.Config.Botname, pluginHandler.Config.Botname, "simple bot")
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(pluginHandler.Config.Channel)
	})
	c.AddHandler("PRIVMSG", pluginHandler.TweetHandler)
	c.AddHandler("PRIVMSG", pluginHandler.ListCommandsHandler)
	c.AddHandler("PRIVMSG", pluginHandler.AddCommandHandler)
	c.AddHandler("PRIVMSG", pluginHandler.DeleteCommandHandler)
	c.AddHandler("PRIVMSG", pluginHandler.CommandHandler)
	c.AddHandler("PRIVMSG", pluginHandler.UptimeHandler)
	c.AddHandler("PRIVMSG", pluginHandler.UpdateChannelGameHandler)
	c.AddHandler("PRIVMSG", pluginHandler.UpdateChannelStatusHandler)
	c.AddHandler("PRIVMSG", pluginHandler.SongHandler)
	c.AddHandler("PRIVMSG", pluginHandler.AddTimeoutListHandler)
	c.AddHandler("PRIVMSG", pluginHandler.TimeoutHandler)

	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.Connect("irc.twitch.tv", pluginHandler.Config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", pluginHandler.Config.Channel)

	<-quit
}
