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
	config configuration
	tweet  *anaconda.TwitterApi
	fm     *lastfm.Lastfm
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

	con := irc.IRC(config.Botname, config.Botname)
	con.Password = config.Aouth
	if err := con.Connect("irc.twitch.tv:6667"); err != nil {
		log.Fatal(err)
	}

	joinChannel(config.Channel, con)
	con.Loop()
}

func joinChannel(channel string, con *irc.Connection) {
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
		var r, g, b uint8
		r, g, b = 216, 52, 52
		colored := rgbterm.String(channel, r, g, b)
		log.Print("Joined " + colored)
		rollCommand(channel, con)
		songCommand(channel, con)
		tweetCommand(channel, con)
		repeatMessenger(channel, con)
	})
}

func repeatMessenger(channel string, con *irc.Connection) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C
		con.Privmsg(channel, config.RepeatMsg)
	}
}

func tweetCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(len(e.Message()) == 6 && time.Since(delay).Seconds() > 10) {
			return
		}
		if e.Message()[0:6] != "!tweet" {
			return
		}

		thetweet, err := tweet.GetUserTimeline(nil)
		if err != nil {
			log.Fatal(err)
		}
		con.Privmsg(channel, thetweet[0].CreatedAt+": \""+thetweet[0].Text+"\"")
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
			randNum := strconv.Itoa(rand.Intn(num))
			con.Privmsg(channel, e.Nick+" rolled "+randNum+"!")
			delay = time.Now()
		}
	})
}

func songCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if !(len(e.Message()) == 5 && time.Since(delay).Seconds() > 10) {
			return
		}
		if e.Message()[0:5] != "!song" {
			return
		}

		artist, trackName := fm.GetCurrentArtistAndTrackName()
		if fm.IsNowPlaying() {
			con.Privmsg(channel, artist+" - "+trackName)
		} else {
			lastPlay, err := fm.GetLastPlayedDate()
			if err != nil {
				log.Fatal(err)
			}
			con.Privmsg(channel, artist+" - "+trackName+". Last played "+lastPlay)
		}
		delay = time.Now()
	})
}
