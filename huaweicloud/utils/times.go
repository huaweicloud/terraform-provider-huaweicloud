package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

// GetTimezoneCode calculates the time zone code and returns a signed number.
// For example, the time zone code for 'Asia/Shanghai' is 8, and the time zone code for 'America/Alaska' is -4.
func GetTimezoneCode() int {
	timeStr := strings.Split(time.Now().String(), " ")[2]
	timezoneNum, _ := strconv.Atoi(timeStr)
	return timezoneNum / 100
}

// FormatTimeStampRFC3339 is used to unify the time format to RFC-3339 and return a time string.
// We can use "isUTC" parameter to reset the timezone. If omitted, the method will return local time.
// Parameter "customFormat" allows you to use a custom RFC3339 format, such as: "2006-01-02T15:04:05.000Z", this
// parameter can be omitted.
func FormatTimeStampRFC3339(timestamp int64, isUTC bool, customFormat ...string) string {
	if timestamp == 0 {
		return ""
	}

	createTime := time.Unix(timestamp, 0)
	if isUTC {
		createTime = createTime.UTC()
	}
	if len(customFormat) > 0 {
		return createTime.Format(customFormat[0])
	}
	return createTime.Format(time.RFC3339)
}

// FormatTimeStampUTC is used to unify the unix second time to UTC time string, format: YYYY-MM-DD HH:MM:SS.
func FormatTimeStampUTC(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02 15:04:05")
}

// FormatTimeStampUTC is used to unify the unix second time to UTC time string, format: YYYY-MM-DD HH:MM:SS.
func FormatUTCTimeStamp(utcTime string) (int64, error) {
	timestamp, err := time.Parse("2006-01-02 15:04:05", utcTime)
	if err != nil {
		return 0, fmt.Errorf("unable to prase the time: %s", utcTime)
	}
	return timestamp.Unix(), nil
}
