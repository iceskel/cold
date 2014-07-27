# IRC bot for twitch
Cold is a little irc bot made in Go that connects to twitch chat and performs command, but can do most other things
with changes. Uses the http://github.com/fluffle/goirc library and should be able to connect
to other irc networks besides twitch. Made just for fun

## How to use
It looks for a default `conf.json` config file in the current directory or one specified by the -c flag. The config file
format should be as follows

```
{
	"Channel": "#foo",
	"Botname": "robit",
	"Aouth": "oauth:32433ASDasd3easd3easdas2",
	"LastfmKey": "aq3asdasASDd213asdasdad",
	"LastfmSecret": "asdas3ASDzsdcz23rzdf",
	"LastfmUser": "bar",
	"RepeatMsg": "Don't forget to follow <twitter>, <facebook> etc",
	"TwitterConsumerKey": "QfOChVxhDyxBS8uh0QK2ylHA5",
	"TwitterConsumerSecret": "E9lTUrM09UwQUN9HN667LHAYTWHR04kCH8MEbDzHSqFSe0I8so",
	"TwitterAccessToken": "2435571859-fiPracoyGwTtypIrEKxPZ035lzwf3cLB3QSy8NI",
	"TwitterAccessSecret": "4LQSoE7HYvzckByuXV28qk2uxSvqda1F7rPRtHtBxYVgD"	
}
```
## Commands/Features
- `!song` or `!music` gets the current or last played track of the lastfm user specified in the config file
- Displays a message every 5 mins to the channel
- `!tweet` or `!twitter` displays your latest tweet
- `!song next` or `!song random` plays the next or random track from your playlist
- `!addtimeout <phrase>` adds that certain phrase for the bot to timeout people for 20s

more to come later

## Requirements
You will need to first register a twitch account and input that username into the ``Botname``
config. Next you will need to get an Aouth for that account, get it from here http://twitchapps.com/tmi/

For the !song/!music command, a last.fm api account is required. Get it here http://www.last.fm/api
Then input the Key and Secret into the config file and the User profile name you want to get
the information from.

For the !tweet/!twitter command, a twitter dev account is required (use your regular twitter login). Further instructions are here http://dev.twitter.com/
Then after creating the application, get the api keys. 

For the !song random command, foobar2k music player is required and two keyboard shortcuts need to be set. "A" key for Action "Playback / Next" and "B" key for Action "Playback / Random". This can be edited easily to whatever you wish, just need to have a look at http://msdn.microsoft.com/en-us/library/windows/desktop/dd375731%28v=vs.85%29.aspx and change the corresponding values in the constants. 


## Installing
Just `go get github.com/iceskel/cold` and cd to the directory and `go build cold.go`
next run it `./cold.exe -c configfile.json`. 
