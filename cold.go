package main

import (
	"encoding/json"
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/aybabtme/rgbterm"
	"github.com/iceskel/lastfm"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
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
	fm          *lastfm.Lastfm
	timeoutList = make(map[string]bool)
	opList      = make(map[string]bool)
)

func main() {
	configFile := flag.String("c", "conf.json", "config file")
	flag.Parse()

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
	fm, err = lastfm.NewLastfm(config.LastfmUser, config.LastfmKey)
	if err != nil {
		log.Fatal(err)
	}
	opList[config.Channel[1:]] = true // op's for channel, gets op only commands
	con := irc.IRC(config.Botname, config.Botname)
	con.Password = config.Aouth
	if err := con.Connect("irc.twitch.tv:6667"); err != nil {
		log.Fatal(err)
	}
	channel := config.Channel
	joinChannel(channel, con)
	con.Loop()
}

func joinChannel(channel string, con *irc.Connection) {
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
		var r, g, b uint8
		r, g, b = 216, 52, 52
		colored := rgbterm.String(channel, r, g, b)
		log.Print("Joined " + colored)
		songCommand(channel, con)
		tweetCommand(channel, con)
		timeoutCop(channel, 20, con)
		addTimeoutList(channel, con)
		rollCommand(channel, con)
		repeatMessenger(channel, con) // must be last
	})
}

func repeatMessenger(channel string, con *irc.Connection) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C
		con.Privmsgf(channel, "► %s", config.RepeatMsg)
	}
}

func tweetCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(time.Since(delay).Seconds() > 10) {
			return
		}
		if e.Message() != "!tweet" && e.Message() != "!twitter" {
			return
		}

		thetweet, err := tweet.GetUserTimeline(nil)
		if err != nil {
			log.Fatal(err)
		}
		con.Privmsgf(channel, "► %s: \"%s\"", thetweet[0].CreatedAt, thetweet[0].Text)
		delay = time.Now()
	})
}

func rollCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(len(e.Message()) >= 5 && time.Since(delay).Seconds() > 10) {
			return
		}
		if e.Message()[0:5] != "!roll" {
			return
		}

		num, err := strconv.Atoi(string(e.Message()[6:]))
		if err == nil && num >= 1 {
			randNum := rand.Intn(num)
			con.Privmsgf(channel, "► %s rolled %d!", e.Nick, randNum)
			delay = time.Now()
		}
	})
}

func songCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(time.Since(delay).Seconds() > 10) {
			return
		}
		if e.Message() != "!song" && e.Message() != "!music" {
			return
		}

		artist, trackName := fm.GetCurrentArtistAndTrackName()
		if fm.IsNowPlaying() {
			con.Privmsgf(channel, "► %s - %s", artist, trackName)
		} else {
			lastPlay, err := fm.GetLastPlayedDate()
			if err != nil {
				log.Fatal(err)
			}
			con.Privmsgf(channel, "► %s - %s. Last played %s", artist, trackName, lastPlay)
		}
		delay = time.Now()
	})
}

func addTimeoutList(channel string, con *irc.Connection) {
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(opList[e.Nick]) {
			return
		}
		if !(len(e.Message()) >= 13) {
			return
		}
		if e.Message()[0:11] != "!addtimeout" {
			return
		}

		timeoutList[e.Message()[12:]] = true
		con.Privmsg(channel, "Timeout word added!")
	})
}

func timeoutCop(channel string, length int, con *irc.Connection) {
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if timeoutList[e.Message()] {
			con.Privmsgf(channel, "/timeout %s %d", e.Nick, length)
		}
	})
}
