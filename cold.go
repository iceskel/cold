package main

import (
	"encoding/json"
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"github.com/iceskel/lastfm"
	"github.com/lxn/win"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"syscall"
	"time"
	"unsafe"
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

const (
	vkA = 0x41 // win32 virtual key A code
	vkX = 0x42 // win32 virtual key B code
)

var (
	config      configuration
	tweet       *anaconda.TwitterApi
	timeoutList = make(map[string]bool)
	opList      = make(map[string]bool)
	hwnd        win.HWND
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

	foobar2kwindowclass := syscall.StringToUTF16Ptr("{97E27FAA-C0B3-4b8e-A693-ED7881E99FC1}")
	foobar2kwindowname := syscall.StringToUTF16Ptr("foobar2000 v1.2.9")
	hwnd = win.FindWindow(foobar2kwindowclass, foobar2kwindowname)
	if unsafe.Pointer(hwnd) == nil {
		log.Fatal("Foobar2k not open or not in default state (press the stop button)")
	}

	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	tweet = anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAccessSecret)
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
		log.Printf("Joined %s", channel)
		songCommand(channel, con)
		tweetCommand(channel, con)
		timeoutCop(channel, 20, con)
		addTimeoutList(channel, con)
		rollCommand(channel, con)
		foobar2kCommands(channel, 10, con)
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
			con.Privmsg(channel, "Tweet command not available, please try later")
			return
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
		fm, err := lastfm.NewLastfm(config.LastfmUser, config.LastfmKey)
		if err != nil {
			con.Privmsg(channel, "Song command not available, please try later")
			return
		}
		artist, trackName := fm.GetCurrentArtistAndTrackName()
		if fm.IsNowPlaying() {
			con.Privmsgf(channel, "► %s - %s", artist, trackName)
		} else {
			lastPlay, err := fm.GetLastPlayedDate()
			if err != nil {
				con.Privmsg(channel, "Song command not available, please try later")
				return
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

func foobar2kCommands(channel string, limit int, con *irc.Connection) {
	total := 0
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Message() == "!music next" || e.Message() == "!song next" {
			total++
			if !(total >= limit) {
				return
			}
			win.PostMessage(hwnd, win.WM_KEYDOWN, vkX, 1)
			win.PostMessage(hwnd, win.WM_KEYUP, vkX, 1)
			total = 0
		} else if e.Message() == "!music random" || e.Message() == "!song random" {
			total++
			if !(total >= limit) {
				return
			}
			win.PostMessage(hwnd, win.WM_KEYDOWN, vkA, 1)
			win.PostMessage(hwnd, win.WM_KEYUP, vkA, 1)
			total = 0
		}
	})
}
