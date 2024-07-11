package common_test

import (
	"reflect"
	"testing"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
)

func TestErrorFunc_ConvertExpected400ErrInto404Err(t *testing.T) {
	input400Err := golangsdk.ErrDefault400{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0002\", \"error_msg\": \"Resource not found\"}"),
		},
	}
	input403Err := golangsdk.ErrDefault403{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0003\", \"error_msg\": \"Authentication Failed\"}"),
		},
	}

	// Step1: Check whether the function can normally identify the expected error code under follow input and convert
	// the 400 error into a 404 error.
	parseResult1 := common.ConvertExpected400ErrInto404Err(input400Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult1.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 400 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step1: Check whether the function can normally recognize unexpected 403 error under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertExpected400ErrInto404Err(input403Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult2.(golangsdk.ErrDefault404); ok {
		t.Fatalf("The expected 403 error was not recognized and was incorrectly converted")
	}
	// Step2: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertExpected400ErrInto404Err(input400Err, "err_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult3, input400Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step1: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 400 error into a 404 error.
	parseResult4 := common.ConvertExpected400ErrInto404Err(input400Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult4.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 400 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step1: Check whether the function can normally recognize unexpected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult5 := common.ConvertExpected400ErrInto404Err(input400Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult5, input400Err) {
		t.Fatalf("error converting 400 error to 404 error via a non-exist error code")
	}
}
