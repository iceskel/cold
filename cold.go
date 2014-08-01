// +build linux darwin freebsd openbsd

package main

import (
	"encoding/json"
	"flag"
	"github.com/ChimeraCoder/anaconda"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/lastfm"
	"io/ioutil"
	"log"
	"time"
)

type configuration struct {
	Channel               string
	Botname               string
	Aouth                 string
	LastfmKey             string
	LastfmSecret          string
	LastfmUser            string
	RepeatMsg             string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

var (
	config      configuration
	tweet       *anaconda.TwitterApi
	timeoutList = make(map[string]bool)
	opList      = make(map[string]bool)
	delay       = time.Now()
	total       = 0
)

func main() {
	configFile := flag.String("c", "conf.json", "config file")
	flag.Parse()

	initConfig(configFile)
	initIrcClient()

}

func initIrcClient() {
	c := irc.SimpleClient(config.Botname, config.Botname, "simple bot")
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(config.Channel)
	})
	c.AddHandler("PRIVMSG", repeatMessenger)
	c.AddHandler("PRIVMSG", tweetHandler)
	c.AddHandler("PRIVMSG", songHandler)
	c.AddHandler("PRIVMSG", addTimeoutListHandler)
	c.AddHandler("PRIVMSG", timeoutHandler)

	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	if err := c.Connect("irc.twitch.tv", config.Aouth); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s", config.Channel)

	<-quit
}

func initConfig(configFile *string) {
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}
	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	tweet = anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAccessSecret)

	opList[config.Channel[1:]] = true // op's for channel, gets op only commands
}

func repeatMessenger(conn *irc.Conn, line *irc.Line) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			<-ticker.C
			conn.Privmsg(config.Channel, "► "+config.RepeatMsg)
		}
	}()
	time.Sleep(5001 * time.Millisecond)
}

func tweetHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!tweet" && line.Args[1] != "!twitter" {
		return
	}

	thetweet, err := tweet.GetUserTimeline(nil)
	if err != nil {
		conn.Privmsg(config.Channel, "► Tweet command not available, please try later")
		return
	}
	conn.Privmsg(config.Channel, "► "+thetweet[0].CreatedAt+": \""+thetweet[0].Text+"\"")
	delay = time.Now()
}

func songHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!song" && line.Args[1] != "!music" {
		return
	}
	fm, err := lastfm.NewLastfm(config.LastfmUser, config.LastfmKey)
	if err != nil {
		conn.Privmsg(config.Channel, "► Song command not available, please try later")
		return
	}
	artist, trackName := fm.GetCurrentArtistAndTrackName()
	if fm.IsNowPlaying() {
		conn.Privmsg(config.Channel, "► "+artist+" - "+trackName)
	} else {
		lastPlay, err := fm.GetLastPlayedDate()
		if err != nil {
			conn.Privmsg(config.Channel, "► Song command not available, please try later")
			return
		}
		conn.Privmsg(config.Channel, "► "+artist+" - "+trackName+". Last played "+lastPlay)
	}
	delay = time.Now()
}

func addTimeoutListHandler(conn *irc.Conn, line *irc.Line) {
	if !(opList[line.Nick]) {
		return
	}
	if !(len(line.Args[1]) >= 13) {
		return
	}
	if line.Args[1][0:11] != "!addtimeout" {
		return
	}

	timeoutList[line.Args[1][12:]] = true
	conn.Privmsg(config.Channel, "Timeout word added!")
}

func timeoutHandler(conn *irc.Conn, line *irc.Line) {
	if timeoutList[line.Args[1]] {
		conn.Privmsg(config.Channel, "/timeout "+line.Nick+" 20")
	}
}

func foobar2kHandler(conn *irc.Conn, line *irc.Line) {
	if line.Args[1] == "!music next" || line.Args[1] == "!song next" {
		total++
		if !(total >= 5) {
			return
		}
		win.PostMessage(hwnd, win.WM_KEYDOWN, vkX, 1)
		win.PostMessage(hwnd, win.WM_KEYUP, vkX, 1)
		total = 0
	} else if line.Args[1] == "!music random" || line.Args[1] == "!song random" {
		total++
		if !(total >= 5) {
			return
		}
		win.PostMessage(hwnd, win.WM_KEYDOWN, vkA, 1)
		win.PostMessage(hwnd, win.WM_KEYUP, vkA, 1)
		total = 0
	}
}
