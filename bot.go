package main

import (
	"encoding/json"
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/iceskel/lastfm"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Configuration struct {
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
	config Configuration
	tweet  *anaconda.TwitterApi
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

	con := irc.IRC(config.Botname, config.Botname)
	con.Password = config.Aouth
	if err := con.Connect("irc.twitch.tv:6667"); err != nil {
		log.Fatal(err)
	}

	JoinChannel(config.Channel, con)
	con.Loop()
}

func JoinChannel(channel string, con *irc.Connection) {
	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
		log.Print("Joined " + channel + " \n")
		RollCommand(channel, con)
		SongCommand(channel, con)
		TweetCommand(channel, con)
		RepeatMessenger(channel, con)
	})
}

func RepeatMessenger(channel string, con *irc.Connection) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C
		con.Privmsg(channel, config.RepeatMsg)
	}
}

func TweetCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Arguments[0] == channel {
			if len(e.Message()) == 6 && time.Since(delay).Seconds() > 10 {
				if e.Message()[0:6] == "!tweet" {
					thetweet, err := tweet.GetUserTimeline(nil)
					if err != nil {
						log.Fatal(err)
					}
					con.Privmsg(channel, thetweet[0].CreatedAt+": \""+thetweet[0].Text+"\"")
				}
			}
		}
	})
}

func RollCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Arguments[0] == channel {
			if len(e.Message()) >= 5 && time.Since(delay).Seconds() > 10 {
				if e.Message()[0:5] == "!roll" {
					num, err := strconv.Atoi(string(e.Message()[6:]))
					if err == nil && num >= 1 {
						randNum := strconv.Itoa(rand.Intn(num))
						con.Privmsg(channel, e.Nick+" rolled "+randNum+"!")
						delay = time.Now()
					}
				}
			}
		}
	})
}

func SongCommand(channel string, con *irc.Connection) {
	delay := time.Now()
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Arguments[0] == channel {
			if len(e.Message()) == 5 && time.Since(delay).Seconds() > 10 {
				if e.Message()[0:5] == "!song" {
					fm := lastfm.NewLastfm(config.LastfmUser, config.LastfmKey)
					artist, trackName := fm.GetCurrentArtistAndTrackName()
					if fm.IsNowPlaying() {
						con.Privmsg(channel, artist+" - "+trackName)
					} else {
						lastPlay := fm.GetLastPlayedDate()
						con.Privmsg(channel, artist+" - "+trackName+". Last played "+lastPlay)
					}
					delay = time.Now()
				}
			}
		}
	})
}
