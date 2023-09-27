package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"os"
	"osuStatsBackend/helpers"
	"osuStatsBackend/osustatsDB"
	"osuStatsBackend/osustatsIRC"
	"osuStatsBackend/routines"
	"sync"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Level:           log.Level(helpers.LogLevel),
	ReportTimestamp: true,
	ReportCaller:    true,
	Prefix:          "[Main]",
})

/*
	.ENV File Contents:

	DBURI

	V1_API_KEY

	IRC_NAME
	IRC_PASSWORD

	V2_CLIENT_ID_OSUSTATS
	V2_CLIENT_SECRET_OSUSTATS
	V2_CALLBACK_OSUSTATS

	V2_CLIENT_ID_BOT
	V2_CLIENT_SECRET_BOT
	V2_CALLBACK_BOT

	PP_CALCULATOR_PATH
*/

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	//Load env file.
	envErr := godotenv.Load()
	if envErr != nil {
		logger.Fatal("Error loading .env file", "Error:", envErr)
	}
	logger.Debug(".env found and Loaded!")

	//Set up and connect to Database
	dbErr := osustatsDB.Connect(os.Getenv("DBURI"))
	if dbErr != nil {
		logger.Fatal("Can't connect to Database!", "Error:", dbErr)
	}

	//Set up and Connect to osu! IRC
	ircErr := osustatsIRC.Connect(os.Getenv("IRC_NAME"), os.Getenv("IRC_PASSWORD"))
	if ircErr != nil {
		logger.Fatal("Can't connect to osu irc!", "Error:", ircErr)
	}

	routines.StartInitialCalculator()
	routines.UpdateBotTops()

	//tell go to wait for workgroup, wich never finishes (block)
	wg.Wait()
}
