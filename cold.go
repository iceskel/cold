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
	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(pluginHandler.Config.Channel)
	})
	c.HandleFunc("PRIVMSG", pluginHandler.TweetHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.ListCommandsHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.AddCommandHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.DeleteCommandHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.CommandHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.UptimeHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.UpdateChannelGameHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.UpdateChannelStatusHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.SongHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.AddTimeoutListHandler)
	c.HandleFunc("PRIVMSG", pluginHandler.TimeoutHandler)

	quit := make(chan bool)
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.ConnectTo("irc.twitch.tv", pluginHandler.Config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", pluginHandler.Config.Channel)

	<-quit
}
