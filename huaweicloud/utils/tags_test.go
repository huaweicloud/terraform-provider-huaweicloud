package utils_test

import (
	"reflect"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestTagsFunc_ExpandResourceTagsMap(t *testing.T) {
	var (
		testInput = map[string]interface{}{
			"foo": "bar",
		}
		emptyInput = make(map[string]interface{})

		expected = []map[string]interface{}{
			{
				"key":   "foo",
				"value": "bar",
			},
		}
		nilExpectedWithType = []map[string]interface{}(nil)
		EmptyExpected       = make([]map[string]interface{}, 0)
	)

	testOutput := utils.ExpandResourceTagsMap(testInput)
	if !reflect.DeepEqual(testOutput, expected) {
		t.Fatalf("[1]The processing result of ExpandResourceTagsMap method is not as expected, want %s, but %s",
			utils.Green(expected), utils.Yellow(testOutput))
	}

	testOutput = utils.ExpandResourceTagsMap(emptyInput)
	if !reflect.DeepEqual(testOutput, nilExpectedWithType) {
		t.Fatalf("[2]The processing result of ExpandResourceTagsMap method is not as expected, want %s, but %s",
			utils.Green(nilExpectedWithType), utils.Yellow(testOutput))
	}

	testOutput = utils.ExpandResourceTagsMap(emptyInput, false)
	if !reflect.DeepEqual(testOutput, nilExpectedWithType) {
		t.Fatalf("[3]The processing result of ExpandResourceTagsMap method is not as expected, want %s, but %s",
			utils.Green(nilExpectedWithType), utils.Yellow(testOutput))
	}

	testOutput = utils.ExpandResourceTagsMap(emptyInput, true)
	if !reflect.DeepEqual(testOutput, EmptyExpected) {
		t.Fatalf("The processing result of ExpandResourceTagsMap method is not as expected, want %s, but %s",
			utils.Green(EmptyExpected), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the ExpandResourceTagsMap method meets expectation")
}
