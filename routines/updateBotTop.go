package routines

import (
	"fmt"
	"golang.org/x/oauth2"
	"os"
	"osuStatsBackend/helpers"
	"osuStatsBackend/osuApi"
	"osuStatsBackend/osuPerformance"
	"osuStatsBackend/osustatsDB"
	"time"
)

const (
	layoutISO = "2006-01-02 15:04:05"
)

func UpdateBotTops() {
	for {
		tops, err := osustatsDB.GetBotTop(false, true)
		if err != nil {
			logger.Error("Error fetching toips from users!", "Error", err)
		}
		for _, top := range tops {
			expireTime, err := helpers.StrToTime(top.LastUpdated)
			if err != nil {
				logger.Errorf("Cannot convert time from String @user %s", top.Top[0].User.Username, "Error", err)
				continue
			}
			if checkIfExpired(expireTime) {
				authUser, err := osustatsDB.GetBotUserByID(top.ID)
				if err != nil {
					logger.Errorf("Error fetching Auth data for %s", top.Top[0].User.Username, "Error", err)
				}
				logger.Debug(authUser.Oauth2.Expiry.String())
				getTop, err := osuApi.GetTopBot(authUser, authUser.UserID, 100, oauth2.Config{
					ClientID:     os.Getenv("V2_CLIENT_ID_BOT"),
					ClientSecret: os.Getenv("V2_CLIENT_SECRET_BOT"),
					Endpoint:     osuApi.Osu,
				}, 0, "bot")
				if err != nil {
					logger.Errorf("Error fetching top 100 from %s", top.Top[0].User.Username, "Error", err)
					continue
				}

				logger.Debug(getTop[0])
				for j, topplay := range getTop {
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
					getTop[j].PpData = res.PerformanceAttributes
					getTop[j].DiffData = res.DifficultyAttributes
					//fmt.Printf("\rOn %d/100", j+1)
					logger.Debugf("\rCalculating scores for user %s on %d / %d", topplay.User.Username, j+1, len(getTop))
				}
				fmt.Printf("\n")
				top.Top = getTop
				updateTime := time.Now().Add(time.Hour * 168)
				top.LastUpdated = updateTime.Format(layoutISO)
				err = osustatsDB.SetBotTop(top.ID, top)
				if err != nil {
					logger.Errorf("Error setting new top plays for user %s", getTop[0].User.Username, "Error", err)
				}
				logger.Debugf("Successfully Calculated top 100 of %s", getTop[0].User.Username)
			}
			//

		}

		time.Sleep(1 * time.Hour)
	}
}

func checkIfExpired(t time.Time) bool {
	expiryDelta := 1 * time.Hour
	return t.Round(0).Add(-expiryDelta).Before(time.Now())
}
