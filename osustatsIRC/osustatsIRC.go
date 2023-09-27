package osustatsIRC

import (
	"github.com/charmbracelet/log"
	"gopkg.in/sorcix/irc.v2"
	"os"
	"osuStatsBackend/helpers"
	"time"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Level:           log.Level(helpers.LogLevel),
	ReportTimestamp: true,
	ReportCaller:    true,
	Prefix:          "[IRC]",
})
var pingRecived = true
var quit = make(chan bool)
var IrcClient *irc.Conn

func Connect(username string, password string) error {
	var err error
	IrcClient, err = irc.Dial("irc.ppy.sh:6667")
	if err != nil {
		return err
	}
	go loop()
	login(username, password)
	go ping(username, password)
	logger.Info("IRC connection up!")
	return nil
}
func loop() {
	var lasterr string
	for {
		select {
		case <-quit:
			return
		default:
			message, err := IrcClient.Decode()
			if err != nil {
				if err.Error() != lasterr { //Stop spamming log.
					logger.Error("Error decoding IRC Meaasge!", "Error:", err)
					lasterr = err.Error()
				}
				continue
			}

			switch message.Command {
			case "QUIT":
				continue
			case "PONG":
				pingRecived = true
				logger.Debug("Recived Pong!")
				continue
			case "PRIVMSG":
				go MessageHandler(message)
				logger.Debug(message)
			default:
				logger.Debug(message)
			}
		}
	}
}
func ping(username string, password string) {
	for {
		if pingRecived {
			pingRecived = false
			err := IrcClient.Encode(&irc.Message{
				Prefix:  nil,
				Command: irc.PING,
				Params:  nil,
			})
			if err != nil {
				logger.Error("Error encoding IRC Meaasge!", "Error:", err)
			}
			logger.Debug("Sending Ping!")
		} else {
			logger.Debug("Last ping was not recived, reconnecting!")
			quit <- true
			var err error
			err = IrcClient.Close()
			if err != nil {
				logger.Error("Error disconnecting from osu IRC!", "Error:", err)
			}
			IrcClient, err = irc.Dial("irc.ppy.sh:6667")
			time.Sleep(time.Second * 5)
			if err != nil {
				for {
					logger.Error("Error connecting to osu IRC, try again in 5 Seconds...", "Error:", err)
					time.Sleep(time.Second * 5)
					IrcClient, err = irc.Dial("irc.ppy.sh:6667")
					if err == nil {
						break
					}
				}
			}
			go loop()
			login(username, password)
			err = IrcClient.Encode(&irc.Message{
				Prefix:  nil,
				Command: irc.PING,
				Params:  nil,
			})
			if err != nil {
				logger.Error("Error encoding IRC Meaasge!", "Error:", err)
			}
		}
		time.Sleep(time.Second * 30)
	}
}
func login(ircNAME string, ircPW string) {
	passwd := &irc.Message{
		Prefix:  nil,
		Command: irc.PASS,
		Params:  []string{ircPW},
	}
	username := &irc.Message{
		Prefix:  nil,
		Command: irc.NICK,
		Params:  []string{ircNAME},
	}

	err := IrcClient.Encode(passwd)
	if err != nil {
		logger.Error("Error encoding IRC Meaasge!", "Error:", err)
	}
	err = IrcClient.Encode(username)
	if err != nil {
		logger.Error("Error encoding IRC Meaasge!", "Error:", err)
	}
}
func SendIrcMessage(target string, message string) {
	err := IrcClient.Encode(&irc.Message{
		Prefix:  nil,
		Command: irc.PRIVMSG,
		Params:  []string{target, message},
	})
	if err != nil {
		logger.Error("Error encoding IRC Meaasge!", "Error:", err)
	}
}
