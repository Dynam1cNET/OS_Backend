package osuBot

import "slices"

var sessions []session

type session struct {
	userID               int
	modifier             string
	command              string
	sessionPlayed        []int
	remainingResultCount int
	resultCount          int
}

func (_session session) resetSession(index int) {
	sessions[index].sessionPlayed = nil
}

func (_session session) addMap(beatmap int, index int) {
	mysession := _session
	sessions[index].sessionPlayed = append(mysession.sessionPlayed, beatmap)
}

func (_session session) hasMap(beatmap int) bool {
	mysession := _session
	return slices.Contains(mysession.sessionPlayed, beatmap)
}

func addSession(_userid int, _mod string, _command string) {
	userid := _userid
	mod := _mod
	command := _command
	sessions = append(sessions, session{
		userID:               userid,
		modifier:             mod,
		command:              command,
		sessionPlayed:        nil,
		remainingResultCount: 0,
		resultCount:          0,
	})
}

func findOrAddSession(_userid int, _mod string, _command string) int {
	userid := _userid
	mod := _mod
	command := _command
	sessionIndex := slices.IndexFunc(sessions, func(c session) bool { return c.userID == userid && c.modifier == mod && c.command == command })
	if sessionIndex > -1 {
		return sessionIndex
	} else {
		logger.Debug("No Session Found! Create Session...")
		addSession(userid, mod, command)
		return findOrAddSession(userid, mod, command)
	}
}
