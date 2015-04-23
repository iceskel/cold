package handlers

import (
	"encoding/json"
	irc "github.com/fluffle/goirc/client"
	"github.com/iceskel/cold/handlers"
	"github.com/lxn/win"
	"io/ioutil"
	"syscall"
	"unsafe"
)

const (
	vkA = 0x41 // win32 virtual key A code
	vkX = 0x42 // win32 virtual key B code
)

type WindowsBotHandler struct {
	Config handlers.Configuration
	Hwnd   win.HWND
	Total  int
}

func New(configFile *string) (*WindowsBotHandler, error) {
	var config handlers.Configuration
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	hwnd := win.FindWindow(syscall.StringToUTF16Ptr("{97E27FAA-C0B3-4b8e-A693-ED7881E99FC1}"),
		syscall.StringToUTF16Ptr("foobar2000 v1.3.3"))

	return &WindowsBotHandler{
		Config: config,
		Hwnd:   hwnd,
		Total:  0,
	}, nil
}

func (wbh *WindowsBotHandler) Foobar2kHandler(conn *irc.Conn, line *irc.Line) {
	if line.Args[1] == "!music next" || line.Args[1] == "!song next" {
		if unsafe.Pointer(wbh.Hwnd) == nil {
			return
		}
		wbh.Total++
		if !(wbh.Total >= 5) {
			return
		}
		win.PostMessage(wbh.Hwnd, win.WM_KEYDOWN, vkX, 1)
		win.PostMessage(wbh.Hwnd, win.WM_KEYUP, vkX, 1)
		wbh.Total = 0
	} else if line.Args[1] == "!music random" || line.Args[1] == "!song random" {
		if unsafe.Pointer(wbh.Hwnd) == nil {
			return
		}
		wbh.Total++
		if !(wbh.Total >= 5) {
			return
		}
		win.PostMessage(wbh.Hwnd, win.WM_KEYDOWN, vkA, 1)
		win.PostMessage(wbh.Hwnd, win.WM_KEYUP, vkA, 1)
		wbh.Total = 0
	}
}
