package osuBot

func Ping() string {
	logger.Debug("Recived !ping Command!")
	return "pong!"
}
