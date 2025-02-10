package utils_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestTypeConvertFunc_StringToJson(t *testing.T) {
	var (
		emptyInput     = "{}"
		correctInput   = "{\"foo\":\"bar\"}"
		incorrectInput = `func() {
			fmt.Println("Hello, this is a function!")
		}`
		emptyInputExpected   = make(map[string]interface{})
		correctInputExpected = map[string]interface{}{
			"foo": "bar",
		}
	)

	testOutput := utils.StringToJson(emptyInput)
	if !reflect.DeepEqual(testOutput, emptyInputExpected) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want %s, but got %s",
			utils.Green(emptyInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJson(correctInput)
	if !reflect.DeepEqual(testOutput, correctInputExpected) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want %s, but got %s",
			utils.Green(correctInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJson(incorrectInput)
	if !reflect.DeepEqual(testOutput, make(map[string]interface{})) {
		t.Fatalf("The processing result of the StringToJson method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the JsonToString method meets expectation")
}

func TestTypeConvertFunc_StringToJsonArray(t *testing.T) {
	var (
		emptyInput           = "[]"
		correctInput         = "[{\"key1\":\"value1\"},{\"key2\":\"value2\"}]"
		incorrectInput       = "{\"foo\":\"bar\"}"
		emptyInputExpected   = make([]map[string]interface{}, 0)
		correctInputExpected = []map[string]interface{}{
			{
				"key1": "value1",
			},
			{
				"key2": "value2",
			},
		}
	)

	testOutput := utils.StringToJsonArray(emptyInput)
	if !reflect.DeepEqual(testOutput, emptyInputExpected) {
		t.Fatalf("The processing result of the StringToJsonArray method is not as expected, want %s, but got %s",
			utils.Green(emptyInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJsonArray(correctInput)
	if !reflect.DeepEqual(testOutput, correctInputExpected) {
		t.Fatalf("The processing result of the StringToJsonArray method is not as expected, want %s, but got %s",
			utils.Green(correctInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.StringToJsonArray(incorrectInput)
	if !reflect.DeepEqual(testOutput, make([]map[string]interface{}, 0)) {
		t.Fatalf("The processing result of the StringToJsonArray method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the StringToJsonArray method meets expectation")
}

func TestTypeConvertFunc_JsonToString(t *testing.T) {
	type Test struct {
		Foo string `json:"foo,omitempty"`
	}

	var (
		emptyInput   = Test{}
		correctInput = Test{
			Foo: "bar",
		}
		emptyInputExpected   = "{}"
		correctInputExpected = "{\"foo\":\"bar\"}"
		// Function is an unsupported type for JsonToString() function input and an error will be returned.
		functionInput = func() {
			fmt.Println("Hello, this is a function!")
		}
	)

	testOutput := utils.JsonToString(emptyInput)
	if !reflect.DeepEqual(testOutput, emptyInputExpected) {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want %s, but got %s",
			utils.Green(emptyInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(correctInput)
	if !reflect.DeepEqual(testOutput, correctInputExpected) {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want %s, but got %s",
			utils.Green(correctInputExpected), utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(functionInput)
	if !reflect.DeepEqual(testOutput, "") {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	testOutput = utils.JsonToString(nil)
	if !reflect.DeepEqual(testOutput, "") {
		t.Fatalf("The processing result of the JsonToString method is not as expected, want \"\", but got %s",
			utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the JsonToString method meets expectation")
}
