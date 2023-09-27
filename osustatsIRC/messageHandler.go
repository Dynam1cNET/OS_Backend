package osustatsIRC

import (
	"gopkg.in/sorcix/irc.v2"
	"osuStatsBackend/osuBot"
	"strings"
)

func MessageHandler(_message *irc.Message) {
	message := _message
	switch {
	case strings.HasPrefix(message.Params[1], "!nr"):
		result, err := osuBot.GetRandomMap(message.Name, message.Params[1])
		if err != nil {
			logger.Error("Got error on GetRandomMap", "Error", err)
		}
		for _, messages := range result {
			SendIrcMessage(message.Name, messages)
		}
		break
	case strings.HasPrefix(message.Params[1], "!nping"):
		SendIrcMessage(message.Name, osuBot.Ping())
		break
	case strings.HasPrefix(message.Params[1], "!npong"):
		SendIrcMessage(message.Name, osuBot.Pong())
		break
	case strings.HasPrefix(message.Params[1], "!nhelp"):
		SendIrcMessage(message.Name, osuBot.Help())
		break
	case strings.HasPrefix(message.Params[1], "!n"):
		break
	}
}
