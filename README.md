# IRC bot for twitch
A little irc bot made in Go that connects to twitch chat and performs command, but can do most other things
with changes. Uses the http://github.com/thoj/go-ircevent irc library and should be able to connect
to other irc networks besides twitch. Made just for fun

## How to use
It looks for a default `conf.json` config file or one specified by the -c flag. The config file
format should be as follows

```
{
	"Channels": ["#foo"],
	"Botname": "robit",
	"Aouth": "oauth:32433asd3easd3easdas2",
	"LastfmKey": "aq3asdasd213asdasdad",
	"LastfmSecret": "asdas3zsdcz23rzdf",
	"LastfmUser": "bar"
}
```
## Commands
- !song gets the current or last played track of the lastfm user specified in the config file
- !roll rolls a random number the size of the int

## Requirements
You will need to first register a twitch account and input that username into the ``Botname``
config. Next you will need to get an Aouth for that account, get it from here http://twitchapps.com/tmi/

For the !song command, a last.fm api account is required. Get it here http://www.last.fm/api
Then input the Key and Secret into the config file and the User profile name you want to get
the information from.

## Installing
Just `go get github.com/iceskel/gobot` and cd to the directory and `go build bot.go`
next run it `./bot -c configfile.json`. 
