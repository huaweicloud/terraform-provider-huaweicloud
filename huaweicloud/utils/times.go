package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// ConvertTimeStrToNanoTimestamp is a method that used to convert the time string into the corresponding timestamp (in
// nanosecond), e.g.
// The supported time formats are as follows:
//   - RFC3339 format:
//     2006-01-02T15:04:05Z (default time format, if you are missing customFormat input)
//     2006-01-02T15:04:05.000000Z
//     2006-01-02T15:04:05Z08:00
//   - Other time formats:
//     2006-01-02 15:04:05
//     2006-01-02 15:04:05+08:00
//     2006-01-02T15:04:05
//     ...
//
// Two common uses are shown below:
// - ConvertTimeStrToNanoTimestamp("2024-01-01T00:00:00Z")
// - ConvertTimeStrToNanoTimestamp("2024-01-01T00:00:00+08:00", "2006-01-02T15:04:05Z08:00")
func ConvertTimeStrToNanoTimestamp(timeStr string, customFormat ...string) int64 {
	// The default time format is RFC3339.
	timeFormat := time.RFC3339
	if len(customFormat) > 0 {
		timeFormat = customFormat[0]
	}
	t, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		log.Printf("error parsing the input time (%s), the time string does not match time format (%s): %s",
			timeStr, timeFormat, err)
		return 0
	}

	timestamp := t.UnixNano() / int64(time.Millisecond)
	// If the time is less than 1970-01-01T00:00:00Z, the timestamp is negative, such as: "0001-01-01T00:00:00Z"
	if timestamp < 0 {
		return 0
	}
	return timestamp
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

// CalculateNextWholeHourAfterFewTime method is used to calculate the minimum whole hour after X (interval, unit:
// hours/minutes/seconds) time.
// + inputTimeStr: UTC input time, not timestamp. For example: 2024-08-01T02:10:35Z
// + interval: Time interval
// + customOutputFormat: Custom output time format, the default output format is RFC3339 time
func CalculateNextWholeHourAfterFewTime(inputTimeStr string, interval time.Duration,
	customOutputFormat ...string) string {
	parsedTime, err := time.Parse(time.RFC3339, inputTimeStr)
	if err != nil {
		log.Printf("[ERROR] incorrect time format for the input time: %s", inputTimeStr)
		return ""
	}

	// Save the time after the target interval.
	nextTime := parsedTime.Add(interval)

	// Adjust the time to the next X (interval, unit: hours/minutes/seconds) time.
	// Pay attention to the case where the hour is 23, in which case it needs to be adjusted to 0 hours of the next day
	nextHour := nextTime.Hour()
	// Set minutes and seconds to 0 when formatting the target time.
	nextMinute := 0
	nextSecond := 0
	nextNano := 0

	// If it is already an hour, the current date and hour are used directly, but the minutes and seconds are set to 0.
	if nextTime.Minute() != 0 || nextTime.Second() != 0 || nextTime.Nanosecond() != 0 {
		if nextHour == 23 {
			nextHour = 0
			nextDate := nextTime.AddDate(0, 0, 1) // Plus one day to the next day.
			nextTime = time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(),
				nextHour, nextMinute, nextSecond, nextNano, nextTime.Location())
		} else {
			nextHour++
			nextTime = time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(),
				nextHour, nextMinute, nextSecond, nextNano, nextTime.Location())
		}
	}

	outputFormat := "2006-01-02T15:04:05.000Z" // The default time format.
	if len(customOutputFormat) > 0 {
		// Override the default time format using custom time format.
		outputFormat = customOutputFormat[0]
	}
	return nextTime.Format(outputFormat)
}

// GetCurrentTime is used to get the current time.
// + isUTC: Whether to generate UTC time
// + customFormat: Custom output time format, default RFC3339 format
func GetCurrentTime(isUTC bool, customFormat ...string) string {
	timeFormat := time.RFC3339
	if len(customFormat) > 0 {
		timeFormat = customFormat[0]
	}
	now := time.Now()
	if isUTC {
		now = now.UTC()
	}
	return now.Format(timeFormat)
}

// GetBeforeOrAfterDate is used to get a few days ago or a few days after.
// + day: The number of days forward or backward, `-` indicates backward.
// + customFormat: Custom output time format, default RFC3339 format
func GetBeforeOrAfterDate(inputTime time.Time, day int, customFormat ...string) string {
	timeFormat := time.RFC3339
	if len(customFormat) > 0 {
		timeFormat = customFormat[0]
	}
	outputTime := inputTime
	if day != 0 {
		outputTime = inputTime.AddDate(0, 0, day)
	}
	return outputTime.Format(timeFormat)
}
