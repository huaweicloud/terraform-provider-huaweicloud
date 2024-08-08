package utils_test

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestTimeFunc_ConvertTimeStrToNanoTimestamp(t *testing.T) {
	var (
		// Test the function and without the custom time format.
		defaultTimeInput       = "2024-01-01T00:00:00Z"
		defaultExpected  int64 = 1704067200000

		timeInputs = [][]interface{}{
			// timeInput, timeFormat, expected
			{"2024-01-01T00:00:00Z", "2006-01-02T15:04:05Z", int64(1704067200000)},
			{"2024-01-01T00:00:00.000000Z", "2006-01-02T15:04:05.000000Z", int64(1704067200000)},
			{"2024-01-01T00:00:00+08:00", "2006-01-02T15:04:05Z07:00", int64(1704038400000)},
			{"2024-01-01 00:00:00", "2006-01-02 15:04:05", int64(1704067200000)},
			{"2024-01-01 00:00:00+08:00", "2006-01-02 15:04:05Z07:00", int64(1704038400000)},
			{"2024-01-01T00:00:00", "2006-01-02T15:04:05", int64(1704067200000)},
			{"", "2006-01-02T15:04:05", int64(0)},
		}
	)

	// Compatibility test
	testOutput := utils.ConvertTimeStrToNanoTimestamp(defaultTimeInput)
	if !reflect.DeepEqual(testOutput, defaultExpected) {
		t.Fatalf("The processing result of the ConvertTimeStrToNanoTimestamp method is not as expected, want %s, but got %s",
			utils.Green(defaultExpected), utils.Yellow(testOutput))
	}

	for _, v := range timeInputs {
		testOutput := utils.ConvertTimeStrToNanoTimestamp(v[0].(string), v[1].(string))
		if !reflect.DeepEqual(testOutput, v[2].(int64)) {
			t.Fatalf("The processing result of the ConvertTimeStrToNanoTimestamp method is not as expected, want %s, but got %s",
				utils.Green(v[2]), utils.Yellow(testOutput))
		}
	}
	t.Logf("All processing results of the ConvertTimeStrToNanoTimestamp method meets expectation")
}

func TestTimeFunc_CalculateNextWholeHourAfterFewTime(t *testing.T) {
	var (
		timeInput = "2024-01-01T01:01:01Z"
		interval  = 24 * time.Hour
	)

	// Compatibility test
	testOutput1 := utils.CalculateNextWholeHourAfterFewTime(timeInput, interval)
	if testOutput1 != "2024-01-02T02:00:00.000Z" {
		t.Fatalf("The processing result of the CalculateNextWholeHourAfterFewTime method is not as expected, want %s, but got %s",
			utils.Green("2024-01-02T02:00:00.000Z"), utils.Yellow(testOutput1))
	}

	testOutput2 := utils.CalculateNextWholeHourAfterFewTime(timeInput, interval, "2006-01-02T15:04:05Z")
	if testOutput2 != "2024-01-02T02:00:00Z" {
		t.Fatalf("The processing result of the CalculateNextWholeHourAfterFewTime method is not as expected, want %s, but got %s",
			utils.Green("2024-01-02T02:00:00Z"), utils.Yellow(testOutput2))
	}
	t.Logf("All processing results of the CalculateNextWholeHourAfterFewTime method meets expectation")
}

func TestTimeFunc_GetCurrentTimeRFC3339(t *testing.T) {
	// Compatibility test
	testOutput1 := utils.GetCurrentTime(true)
	re1 := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`
	if !regexp.MustCompile(re1).MatchString(testOutput1) {
		t.Fatalf(`The processing result of the GetCurrentTime method is not as expected, the time format want %s, but the value got %s`,
			utils.Green(re1), utils.Yellow(testOutput1))
	}

	// Compatibility test
	testOutput2 := utils.GetCurrentTime(false)
	re2 := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$`
	if !regexp.MustCompile(re2).MatchString(testOutput2) {
		t.Fatalf(`The processing result of the GetCurrentTime method is not as expected, the time format want %s, but the value got %s`,
			utils.Green(re2), utils.Yellow(testOutput2))
	}

	// Compatibility test
	testOutput3 := utils.GetCurrentTime(false, "2006-01-02T15:04:05.000Z")
	re3 := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`
	if !regexp.MustCompile(re3).MatchString(testOutput3) {
		t.Fatalf(`The processing result of the GetCurrentTime method is not as expected, the time format want %s, but the value got %s`,
			utils.Green(re3), utils.Yellow(testOutput3))
	}

	t.Logf("All processing results of the GetCurrentTime method meets expectation")
}
