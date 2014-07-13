# IRC bot for twitch
Cold is a little irc bot made in Go that connects to twitch chat and performs command, but can do most other things
with changes. Uses the http://github.com/thoj/go-ircevent irc library and should be able to connect
to other irc networks besides twitch. Made just for fun

## How to use
It looks for a default `conf.json` config file in the current directory or one specified by the -c flag. The config file
format should be as follows

```
{
	"Channel": "#foo",
	"Botname": "robit",
	"Aouth": "oauth:32433asd3easd3easdas2",
	"LastfmKey": "aq3asdasd213asdasdad",
	"LastfmSecret": "asdas3zsdcz23rzdf",
	"LastfmUser": "bar",
	"RepeatMsg": "Don't forget to follow <twitter>, <facebook> etc",
	"TwitterConsumerKey": "QfOChVxhDyxBS8uh0QK2ylHA5",
	"TwitterConsumerSecret": "E9lTUrM09UwQUN9HN667LHAYTWHR04kCH8MEbDzHSqFSe0I8so",
	"TwitterAccessToken": "2435571859-fiPracoyGwTtypIrEKxPZ035lzwf3cLB3QSy8NI",
	"TwitterAccessSecret": "4LQSoE7HYvzckByuXV28qk2uxSvqda1F7rPRtHtBxYVgD"	
}
```
## Commands/Features
- !song gets the current or last played track of the lastfm user specified in the config file
- !roll rolls a random number the size of the int
- Displays a message every 5 mins to the channel
- !tweet displays your latest tweet

more to come later

## Requirements
You will need to first register a twitch account and input that username into the ``Botname``
config. Next you will need to get an Aouth for that account, get it from here http://twitchapps.com/tmi/

For the !song command, a last.fm api account is required. Get it here http://www.last.fm/api
Then input the Key and Secret into the config file and the User profile name you want to get
the information from.

## Installing
Just `go get github.com/iceskel/gobot` and cd to the directory and `go build bot.go`
next run it `./bot -c configfile.json`. 
