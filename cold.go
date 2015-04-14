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

	genericHandler, err := handlers.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	initIrcClient(genericHandler)

}

func initIrcClient(genericHandler *handlers.BotHandler) {
	c := irc.SimpleClient(genericHandler.Config.Botname, genericHandler.Config.Botname, "simple bot")
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(genericHandler.Config.Channel)
	})
	c.AddHandler("PRIVMSG", genericHandler.TweetHandler)
	c.AddHandler("PRIVMSG", genericHandler.UpdateChannelGameHandler)
	c.AddHandler("PRIVMSG", genericHandler.UpdateChannelStatusHandler)
	c.AddHandler("PRIVMSG", genericHandler.SongHandler)
	c.AddHandler("PRIVMSG", genericHandler.AddTimeoutListHandler)
	c.AddHandler("PRIVMSG", genericHandler.TimeoutHandler)

	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.Connect("irc.twitch.tv", genericHandler.Config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", genericHandler.Config.Channel)

	<-quit
}
