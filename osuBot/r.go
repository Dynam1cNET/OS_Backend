package osuBot

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"os"
	"osuStatsBackend/helpers"
	"osuStatsBackend/osustatsDB"
	"strings"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Level:           log.Level(helpers.LogLevel),
	ReportTimestamp: true,
	ReportCaller:    true,
	Prefix:          "[BOT]",
})

func GetRandomMap(name string, msg string) ([]string, error) {
	logger.Debugf("Recived !r Command from %s! Message content=%s", name, msg)
	user, err := osustatsDB.GetUser(name)
	if err != nil {
		logger.Error("Error fetching user by name!", "Error", err)
	}

	if user.ID == 0 {
		//return "", errors.New("userID of 0 Provided, userid 0 does not exists")
		return []string{
			"It seems like you never Used this bot before. Please login with osu here: https://osubot.dynam1c.net/loginwithosu",
			"If you just signed up wait a minute! My bot is looking up your profile.",
			"I will send you a Message when i started calculating and when i am done!",
		}, nil
	} else {
		//msg = strings.TrimPrefix(msg, "!r")
		//TODO: Remove n prefix, debug command since i only have one account to test with. Thanks osu
		msg = strings.TrimPrefix(msg, "!nr")
		mods := "nm"
		searchParms, err := osustatsDB.GetUserStats(user.ID)
		if err != nil {
			logger.Error("Error getting search Parms", "Error", err)
		}
		if msg == "" {
			mods = "nm"
		} else if strings.Contains(strings.ToLower(msg), "nm") && len(msg) == 3 {
			mods = "nm"
		} else if strings.Contains(strings.ToLower(msg), "hd") && len(msg) == 3 {
			mods = "hd"
		} else if strings.Contains(strings.ToLower(msg), "hr") && len(msg) == 3 {
			mods = "hr"
		} else if strings.Contains(strings.ToLower(msg), "dt") && len(msg) == 3 {
			mods = "dt"
			searchParms.MinDiffSpeeddifficulty = searchParms.MinDiffSpeeddifficulty * 0.70
			//searchParms.MaxDiffSpeeddifficulty = searchParms.MaxDiffSpeeddifficulty * 0.90
			searchParms.MinPPSpeed = searchParms.MinPPSpeed * 0.70
			//searchParms.MaxPPSpeed = searchParms.MaxPPSpeed * 0.90
			searchParms.MinDiffSliderfactor = searchParms.MinDiffSliderfactor * 0.70
			//searchParms.MaxDiffSliderfactor = searchParms.MaxDiffSliderfactor * 0.90
			searchParms.MinOD = searchParms.MinOD * 0.70
			searchParms.MinAR = 6

		} else if strings.Contains(strings.ToLower(msg), "hr") && strings.Contains(strings.ToLower(msg), "dt") && len(msg) == 5 {
			mods = "hrdt"
		} else if strings.Contains(strings.ToLower(msg), "hr") && strings.Contains(strings.ToLower(msg), "hd") && len(msg) == 5 {
			mods = "hrhd"
		} else if strings.Contains(strings.ToLower(msg), "dt") && strings.Contains(strings.ToLower(msg), "hd") && len(msg) == 5 {
			mods = "dthd"
			searchParms.MinDiffSpeeddifficulty = searchParms.MinDiffSpeeddifficulty * 0.70
			//searchParms.MaxDiffSpeeddifficulty = searchParms.MaxDiffSpeeddifficulty * 0.90
			searchParms.MinPPSpeed = searchParms.MinPPSpeed * 0.70
			//searchParms.MaxPPSpeed = searchParms.MaxPPSpeed * 0.90
			searchParms.MinDiffSliderfactor = searchParms.MinDiffSliderfactor * 0.70
			//searchParms.MaxDiffSliderfactor = searchParms.MaxDiffSliderfactor * 0.90
			searchParms.MinOD = searchParms.MinOD * 0.70
			searchParms.MinAR = 6
		} else if strings.Contains(strings.ToLower(msg), "dt") && strings.Contains(strings.ToLower(msg), "hd") && strings.Contains(strings.ToLower(msg), "hr") && len(msg) == 7 {
			mods = "hrdthd"
		} else {
			return []string{"Typre !r mods example: \"!r hddt\" mods can be HD,HR,DT in any combination. !r is always nomod"}, nil

		}

		coll := osustatsDB.MongoClient.Database("users").Collection("fingerprints")
		curr, err := coll.Find(ctx, bson.D{{"UserID", id},{"mods", mods}})
		if err != nil {
			logger.Error("Error fetching map from Database!", "Error", err)
			return nil, err
		}
		var results []MapResult
		err = curr.All(context.Background(), &results)
		if err != nil {
			logger.Error("Error parsing results to Struct", "Error", err)
			return nil, err
		}
		if len(results) == 0 {
			return []string{"No maps found for your Settings. This can be because you just signed up for the bot and it is still Calculating your Stats. Try again in a Minute!"}, nil
		}

		result, remaining, noMaps := selectRandom(results, user.ID, mods, len(results))

		if noMaps != nil {
			return []string{noMaps.Error()}, nil
		}

		return []string{fmt.Sprintf("[%s %s [%s] +%s] 100%% %.2fpp | 99%% %.2fpp | 98%% %.2fpp | Results: %d | Remaining: %d", strings.Replace(result.Url, "beatmaps", "b", 1), result.Title, result.Version, strings.ToUpper(mods), result.PP100Percent, result.PP99Percent, result.PP98Percent, result.Results, remaining)}, nil
	}
}
func selectRandom(_maps []MapResult, _userid int, _mod string, maxCount int) (MapResult, int, error) {
	maps := _maps
	userid := _userid
	mod := _mod

	randomIndex := 0

	sessionId := findOrAddSession(userid, mod, "r")
	if len(maps) == 0 {
		sessions[sessionId].resetSession(sessionId)
		return MapResult{}, 0, errors.New("You have gone through all maps. Resetting Session! You can invoke !r again, but you will get old Requests again.")
	} else {
		randomIndex = rand.Intn(len(maps))
		if !sessions[sessionId].hasMap(maps[randomIndex].Id) {

			sessions[sessionId].addMap(maps[randomIndex].Id, sessionId)
			sessions[sessionId].resultCount = maxCount
			sessions[sessionId].remainingResultCount = maxCount - len(sessions[sessionId].sessionPlayed)
			return maps[randomIndex], sessions[sessionId].remainingResultCount, nil
		} else {
			remaining := removeMap(maps, randomIndex)
			return selectRandom(remaining, userid, mod, maxCount)
		}
	}
}

func removeMap(slice []MapResult, s int) []MapResult {
	return append(slice[:s], slice[s+1:]...)
}
