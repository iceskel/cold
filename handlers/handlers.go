package handlers

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/lastfm"
	"github.com/iceskel/twitch"
	"io/ioutil"
	"strings"
	"time"
)

// Configuration struct for the config file
type Configuration struct {
	Channel               string
	ChannelAouth          string
	Botname               string
	Aouth                 string
	LastfmKey             string
	LastfmSecret          string
	LastfmUser            string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

type CommandMessage struct {
	Commands []_commands
}

type _commands struct {
	Command string
	Message string
}

type BotHandler struct {
	Config         Configuration
	CommandsMap    map[string]string
	CommandsStruct CommandMessage
	Tweet          *anaconda.TwitterApi
	Twitch         *twitch.TwitchApi
	Lastfm         *lastfm.LastfmApi
	TimeoutList    map[string]bool
	OpList         map[string]bool
	Delay          time.Time
}

// New returns a new BotHandler instance
func New(configFile *string) (*BotHandler, error) {
	var config Configuration
	var cmd CommandMessage
	commands := make(map[string]string)

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	file, err = ioutil.ReadFile("commands.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &cmd); err != nil {
		return nil, err
	}

	for _, val := range cmd.Commands {
		commands[val.Command] = val.Message
	}

	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	op := make(map[string]bool)
	op[config.Channel[1:]] = true // op's for channel, gets op only commands

	return &BotHandler{
		Config:         config,
		CommandsMap:    commands,
		CommandsStruct: cmd,
		Tweet:          anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAccessSecret),
		Twitch:         twitch.New(config.Channel[1:], config.ChannelAouth),
		Lastfm:         lastfm.New(config.LastfmUser, config.LastfmKey),
		TimeoutList:    make(map[string]bool),
		OpList:         op,
		Delay:          time.Now(),
	}, nil
}

func saveCommandsToFile(bh *BotHandler) error {
	jsonfile, err := json.Marshal(bh.CommandsStruct)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("commands.json", jsonfile, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (bh *BotHandler) DeleteCommandHandler(conn *irc.Conn, line *irc.Line) {
	if !(bh.OpList[line.Nick]) {
		return
	}
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}

	delcmd := strings.Fields(line.Args[1])
	if len(delcmd) < 2 {
		return
	}
	if delcmd[0] != "!delcmd" {
		return
	}

	for i, v := range bh.CommandsStruct.Commands {
		if v.Command == delcmd[1] {
			bh.CommandsStruct.Commands = append(bh.CommandsStruct.Commands[:i], bh.CommandsStruct.Commands[i+1:]...)
		}
	}
	delete(bh.CommandsMap, delcmd[1])

	if err := saveCommandsToFile(bh); err != nil {
		conn.Privmsg(bh.Config.Channel, "► Delete command not available, please try later.")
		return
	}

	conn.Privmsg(bh.Config.Channel, "► Command "+delcmd[1]+" deleted!")
	bh.Delay = time.Now()

}

func (bh *BotHandler) AddCommandHandler(conn *irc.Conn, line *irc.Line) {
	if !(bh.OpList[line.Nick]) {
		return
	}
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}

	newcmd := strings.Fields(line.Args[1])
	if len(newcmd) < 3 {
		return
	}
	if newcmd[0] != "!addcmd" {
		return
	}
	var msg string
	for _, val := range newcmd[2:] {
		msg += " " + val
	}
	bh.CommandsMap[newcmd[1]] = msg

	bh.CommandsStruct.Commands = append(bh.CommandsStruct.Commands, _commands{Command: newcmd[1], Message: msg})

	if err := saveCommandsToFile(bh); err != nil {
		conn.Privmsg(bh.Config.Channel, "► Add command not available, please try later.")
		return
	}

	conn.Privmsg(bh.Config.Channel, "► New command "+newcmd[1]+" added!")

	bh.Delay = time.Now()
}

func (bh *BotHandler) CommandHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}

	if val, ok := bh.CommandsMap[line.Args[1]]; ok {
		conn.Privmsg(bh.Config.Channel, "► "+val)
		bh.Delay = time.Now()
	}
}

func (bh *BotHandler) ListCommandsHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!commands" && line.Args[1] != "!cmd" && line.Args[1] != "!command" {
		return
	}

	var cmds string
	for key := range bh.CommandsMap {
		cmds += " " + key
	}
	cmds += " !tweet !song !commands"
	conn.Privmsg(bh.Config.Channel, "► Command list: "+cmds)
	bh.Delay = time.Now()
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
		conn.Privmsg(bh.Config.Channel, "► Update status command not available, please try later.")
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
		conn.Privmsg(bh.Config.Channel, "► Update game command not available, please try later.")
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
		conn.Privmsg(bh.Config.Channel, "► Tweet command not available, please try later.")
		return
	}
	conn.Privmsg(bh.Config.Channel, "► "+thetweet[0].CreatedAt+": \""+thetweet[0].Text+"\"")
	bh.Delay = time.Now()
}

func (bh *BotHandler) UptimeHandler(conn *irc.Conn, line *irc.Line) {
	if !(time.Since(bh.Delay).Seconds() > 10) {
		return
	}
	if line.Args[1] != "!live" && line.Args[1] != "!uptime" {
		return
	}

	uptime, err := bh.Twitch.Uptime()
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Uptime command not available, please try later.")
		return
	}

	if uptime == "" {
		conn.Privmsg(bh.Config.Channel, "► Stream currently offline")
	} else {
		conn.Privmsg(bh.Config.Channel, "► Live for "+uptime)
	}

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
		conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later.")
		return
	}
	nwplay, err := bh.Lastfm.IsNowPlaying()
	if err != nil {
		conn.Privmsg(bh.Config.Channel, "► Song command not available, please try later.")
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
