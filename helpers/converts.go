package helpers

import "time"

const (
	layoutISO = "2006-01-02 15:04:05"
)

func StrToTime(timestr string) (time.Time, error) {
	date, err := time.Parse(layoutISO, timestr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
