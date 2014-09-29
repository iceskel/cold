package handlers

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/lastfm"
	"io/ioutil"
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

type BotHandler struct {
	Config      Configuration
	Tweet       *anaconda.TwitterApi
	TimeoutList map[string]bool
	OpList      map[string]bool
	Delay       time.Time
}

func NewBotHandler(configFile *string) (*BotHandler, error) {
	var config Configuration
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	op := make(map[string]bool)
	op[config.Channel[1:]] = true // op's for channel, gets op only commands

	return &BotHandler{
		Config:      config,
		Tweet:       anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAccessSecret),
		TimeoutList: make(map[string]bool),
		OpList:      op,
		Delay:       time.Now(),
	}, nil
}

func (bh *BotHandler) TweetHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!tweet" && line.Args[1] != "!twitter" {
		return
	}

	thetweet, err := bh.Tweet.GetUserTimeline(nil)
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Tweet command not available, please try later")
		return
	}
	conn.Privmsg(bh.Config.Channel, "► "+thetweet[0].CreatedAt+": \""+thetweet[0].Text+"\"")
	bh.Delay = time.Now()
}

func (bh *BotHandler) SongHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!song" && line.Args[1] != "!music" {
		return
	}
	fm, err := lastfm.NewLastfm(bh.Config.LastfmUser, bh.Config.LastfmKey)
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later")
		return
	}
	artist, trackName := fm.GetCurrentArtistAndTrackName()
	if fm.IsNowPlaying() {
		conn.Privmsg(bh.Config.Channel, "► "+artist+" - "+trackName)
	} else {
		lastPlay, err := fm.GetLastPlayedDate()
		if err != nil {
			conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later")
			return
		}
		conn.Privmsg(bh.Config.Channel, "► "+artist+" - "+trackName+". Last played "+lastPlay)
	}
	bh.Delay = time.Now()
}

func (bh *BotHandler) AddTimeoutListHandler(conn *irc.Conn, line *irc.Line) {
	if !(bh.OpList[line.Nick]) {
		return
	}
	if !(len(line.Args[1]) >= 13) {
		return
	}
	if line.Args[1][0:11] != "!addtimeout" {
		return
	}

	bh.TimeoutList[line.Args[1][12:]] = true
	conn.Privmsg(bh.Config.Channel, "Timeout word added!")
}

func (bh *BotHandler) TimeoutHandler(conn *irc.Conn, line *irc.Line) {
	if bh.TimeoutList[line.Args[1]] {
		conn.Privmsg(bh.Config.Channel, "/timeout "+line.Nick+" 20")
	}
}

func (bh *BotHandler) RepeatMessenger(conn *irc.Conn, line *irc.Line) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			<-ticker.C
			conn.Privmsg(bh.Config.Channel, "► "+bh.Config.RepeatMsg)
		}
	}()
	time.Sleep(5001 * time.Millisecond)
}
