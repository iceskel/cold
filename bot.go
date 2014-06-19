package main

import (
	"encoding/json"
	"flag"
	"github.com/iceskel/lastfm"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Configuration struct {
	Channel			string
	Botname			string
	Aouth			string
	LastfmKey		string
	LastfmSecret	string
	LastfmUser		string
	RepeatMsg		string
}

var (
	config Configuration
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
		RepeatMessenger(channel, con)
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

func RepeatMessenger(channel string, con *irc.Connection) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C
		con.Privmsg(channel, config.RepeatMsg)
	}
}
