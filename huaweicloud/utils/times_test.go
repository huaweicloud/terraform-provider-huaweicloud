package utils_test

import (
	"reflect"
	"testing"

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
