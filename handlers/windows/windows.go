package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ChimeraCoder/anaconda"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/cold/handlers"
	"github.com/lxn/win"
	"io/ioutil"
	"syscall"
	"time"
	"unsafe"
)

const (
	vkA = 0x41 // win32 virtual key A code
	vkX = 0x42 // win32 virtual key B code
)

type WindowsBotHandler struct {
	Config      handlers.Configuration
	Tweet       *anaconda.TwitterApi
	TimeoutList map[string]bool
	OpList      map[string]bool
	Delay       time.Time
	Hwnd        win.HWND
	Total       int
}

func NewWindowsBotHandler(configFile *string) (*WindowsBotHandler, error) {
	var config handlers.Configuration
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

	hwnd := win.FindWindow(syscall.StringToUTF16Ptr("{97E27FAA-C0B3-4b8e-A693-ED7881E99FC1}"),
		syscall.StringToUTF16Ptr("foobar2000 v1.3.3"))
	if unsafe.Pointer(hwnd) == nil {
		return nil, errors.New("Foobar2k not open or not in default state (press the stop button)")
	}

	return &WindowsBotHandler{
		Config:      config,
		Tweet:       anaconda.NewTwitterApi(config.TwitterAccessToken, config.TwitterAccessSecret),
		TimeoutList: make(map[string]bool),
		OpList:      op,
		Delay:       time.Now(),
		Hwnd:        hwnd,
		Total:       0,
	}, nil
}

func (wbh *WindowsBotHandler) Foobar2kHandler(conn *irc.Conn, line *irc.Line) {
	if line.Args[1] == "!music next" || line.Args[1] == "!song next" {
		wbh.Total++
		if !(wbh.Total >= 5) {
			return
		}
		win.PostMessage(wbh.Hwnd, win.WM_KEYDOWN, vkX, 1)
		win.PostMessage(wbh.Hwnd, win.WM_KEYUP, vkX, 1)
		wbh.Total = 0
	} else if line.Args[1] == "!music random" || line.Args[1] == "!song random" {
		wbh.Total++
		if !(wbh.Total >= 5) {
			return
		}
		win.PostMessage(wbh.Hwnd, win.WM_KEYDOWN, vkA, 1)
		win.PostMessage(wbh.Hwnd, win.WM_KEYUP, vkA, 1)
		wbh.Total = 0
	}
}
