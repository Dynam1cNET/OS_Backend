package routines

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"osuStatsBackend/helpers"
	"osuStatsBackend/osuPerformance"
	"osuStatsBackend/osustatsDB"
	"osuStatsBackend/osustatsIRC"
	"strings"
	"time"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Level:           log.Level(helpers.LogLevel),
	ReportTimestamp: true,
	ReportCaller:    true,
	Prefix:          "[Routines]",
})

func StartInitialCalculator() {
	go func() {
		for {
			signupCalculator()
			time.Sleep(time.Second * 5)
		}
	}()
}
func signupCalculator() {
	results, err := osustatsDB.GetBotTop(false, false)
	if err != nil {
		logger.Error("Error Fetching users from Database", "Error", err)
	}
	for _, result := range results {
		osustatsIRC.SendIrcMessage(strings.Replace(result.Top[0].User.Username, " ", "_", -1), "I've started the calculation process. Please wait a Minute")
		for j, topplay := range result.Top {
			args := []string{os.Getenv("PP_CALCULATOR_PATH"), "simulate"}
			args = append(args, topplay.Mode)
			args = append(args, fmt.Sprintf("%d", topplay.Beatmap.ID))
			args = append(args, "-j")
			for _, mod := range topplay.Mods {
				args = append(args, fmt.Sprintf("-m %s", mod))
			}
			res, err := osuPerformance.GetPP(args, fmt.Sprintf("%d", topplay.Beatmap.ID))
			if err != nil {
				logger.Errorf("Error Calculating PP for user %s", topplay.User.Username, "Error", err)
			}
			result.Top[j].PpData = res.PerformanceAttributes
			result.Top[j].DiffData = res.DifficultyAttributes
			//fmt.Printf("\rOn %d/100", j+1)
			fmt.Printf("\rCalculating scores for user %s on %d / %d", topplay.User.Username, j+1, len(result.Top))
			if len(result.Top) > 5 {
				if (j+1)%(len(result.Top)/5) == 0 {
					osustatsIRC.SendIrcMessage(strings.Replace(topplay.User.Username, " ", "_", -1), fmt.Sprintf("Calculating top plays %d / %d", j+1, len(result.Top)))
				}
			}

		}
		fmt.Printf("\n")
		result.IsPPCalculated = true
		err := osustatsDB.SetBotTop(result.ID, result)
		if err != nil {
			logger.Errorf("Error setting new top plays for user %s", result.Top[0].User.Username, "Error", err)
		}
		osustatsIRC.SendIrcMessage(strings.Replace(result.Top[0].User.Username, " ", "_", -1), "I am done! You can start using this bot now.")
		logger.Debugf("Successfully Calculated top 100 of %s", result.Top[0].User.Username)
	}
}
