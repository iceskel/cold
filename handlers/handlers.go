package handlers

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/lastfm"
	"github.com/iceskel/twitch"
	"io/ioutil"
	"time"
)

// struct for the config file
type Configuration struct {
	Channel               string
	ChannelAouth          string
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
	Twitch      *twitch.TwitchApi
	Lastfm      *lastfm.LastfmApi
	TimeoutList map[string]bool
	OpList      map[string]bool
	Delay       time.Time
}

// New returns a new BotHandler instance
func New(configFile *string) (*BotHandler, error) {
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
		Twitch:      twitch.New(config.Channel[1:], config.ChannelAouth),
		Lastfm:      lastfm.New(config.LastfmUser, config.LastfmKey),
		TimeoutList: make(map[string]bool),
		OpList:      op,
		Delay:       time.Now(),
	}, nil
}

func (bh *BotHandler) UpdateChannelStatusHandler(conn *irc.Conn, line *irc.Line) {
	if !(bh.OpList[line.Nick]) {
		return
	}
	if !(len(line.Args[1]) >= 9) {
		return
	}
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1][0:7] != "!status" {
		return
	}

	status := line.Args[1][8:]
	if err := bh.Twitch.UpdateStatus(status); err != nil {
		conn.Privmsg(bh.Config.Channel, "► Update status command not available, please try later")
		return
	}
	conn.Privmsg(bh.Config.Channel, "► Status changed to "+status)
	bh.Delay = time.Now()

}

func (bh *BotHandler) UpdateChannelGameHandler(conn *irc.Conn, line *irc.Line) {
	if !(bh.OpList[line.Nick]) {
		return
	}
	if !(len(line.Args[1]) >= 7) {
		return
	}
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1][0:5] != "!game" {
		return
	}

	game := line.Args[1][6:]
	if err := bh.Twitch.UpdateGame(game); err != nil {
		conn.Privmsg(bh.Config.Channel, "► Update game command not available, please try later")
		return
	}
	conn.Privmsg(bh.Config.Channel, "► Game changed to "+game)
	bh.Delay = time.Now()
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

	artist, trackName, err := bh.Lastfm.GetCurrentArtistAndTrackName()
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later")
		return
	}
	nwplay, err := bh.Lastfm.IsNowPlaying()
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later")
		return
	}

	if nwplay {
		conn.Privmsg(bh.Config.Channel, "► "+artist+" - "+trackName)
	} else {
		lastPlay, err := bh.Lastfm.GetLastPlayedDate()
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
