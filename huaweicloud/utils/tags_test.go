package utils_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func TestTagsFunc_FlattenSameKeyTagsToMap(t *testing.T) {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	d := resource.TestResourceData()
	var (
		rawArray = []map[string]interface{}{
			{"a": "b"},
			{"a": "b"},
			{"a": "b", "c": "d"},
			{"a": "b", "c": "d"},
			{"a": "b"},
			{},
			{},
		}

		remoteTagArray = []interface{}{
			[]interface{}{
				map[string]interface{}{"key": "a", "value": "d"},
			},
			[]interface{}{
				map[string]interface{}{"key": "a", "value": "d"},
				map[string]interface{}{"key": "m", "value": "n"},
			},
			[]interface{}{
				map[string]interface{}{"key": "a", "value": "d"},
				map[string]interface{}{"key": "m", "value": "n"},
			},
			[]interface{}{
				map[string]interface{}{"key": "a", "value": "d"},
				map[string]interface{}{"key": "c", "value": "a"},
				map[string]interface{}{"key": "m", "value": "n"},
			},
			make([]interface{}, 0),
			[]interface{}{
				map[string]interface{}{"key": "m", "value": "n"},
			},
			make([]interface{}, 0),
		}

		expectedArray = []map[string]interface{}{
			{"a": "d"},
			{"a": "d"},
			{"a": "d"},
			{"a": "d", "c": "a"},
			{},
			{},
			{},
		}
	)

	for i := 0; i < 7; i++ {
		if err := d.Set("tags", rawArray[i]); err != nil {
			t.Fatalf("error setting tags attribute: %s", utils.Yellow(err))
		}

		if !reflect.DeepEqual(d.Get("tags"), rawArray[i]) {
			t.Fatalf("error setting tags attribute, want '%v', but got '%v'", utils.Green(rawArray[i]),
				utils.Yellow(d.Get("tags")))
		}

		remoteTags := remoteTagArray[i]
		expectedMap := expectedArray[i]
		result := utils.FlattenSameKeyTagsToMap(d, remoteTags)

		if !reflect.DeepEqual(result, expectedMap) {
			t.Fatalf("The processing result of the function 'FlattenSameKeyTagsToMap' is not as expected, want '%v', "+
				"but got '%v'", utils.Green(expectedMap), utils.Yellow(result))
		}
		t.Logf("The processing result of `FlattenSameKeyTagsToMap` method meets expectation: %s", utils.Green(expectedMap))
	}
}
