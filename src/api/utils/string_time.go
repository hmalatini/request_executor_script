package utils

import "time"

func GetStringTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02_15:04:05.000000")
}
