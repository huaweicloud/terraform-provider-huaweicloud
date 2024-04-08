package utils

import (
	"log"
	"time"
)

// ConvertTimeStrToNanoTimestamp is a method that used to convert the time in RFC3339 format into the corresponding
// timestamp (in nanosecond), e.g.
// Before: 2007-09-02T00:00:00.00000Z
// After:  1188691200000
func ConvertTimeStrToNanoTimestamp(timeStr string) int64 {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("error parsing time string, it's not RFC3339 format: %s", err)
	}

	return t.UnixNano() / int64(time.Millisecond)
}
