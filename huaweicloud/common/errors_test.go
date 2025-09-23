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
	// Step2: Check whether the function can normally recognize unexpected 403 error under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertExpected400ErrInto404Err(input403Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult2.(golangsdk.ErrDefault404); ok {
		t.Fatalf("The expected 403 error was not recognized and was incorrectly converted")
	}
	// Step3: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertExpected400ErrInto404Err(input400Err, "err_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult3, input400Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step4: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 400 error into a 404 error.
	parseResult4 := common.ConvertExpected400ErrInto404Err(input400Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult4.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 400 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step5: Check whether the function can normally recognize unexpected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult5 := common.ConvertExpected400ErrInto404Err(input400Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult5, input400Err) {
		t.Fatalf("error converting 400 error to 404 error via a non-exist error code")
	}
}

func TestErrorFunc_ConvertExpected401ErrInto404Err(t *testing.T) {
	input401Err := golangsdk.ErrDefault401{
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
	// the 401 error into a 404 error.
	parseResult1 := common.ConvertExpected401ErrInto404Err(input401Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult1.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 401 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step2: Check whether the function can normally recognize unexpected 403 error under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertExpected401ErrInto404Err(input403Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult2.(golangsdk.ErrDefault404); ok {
		t.Fatalf("The expected 403 error was not recognized and was incorrectly converted")
	}
	// Step3: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertExpected401ErrInto404Err(input401Err, "err_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult3, input401Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step4: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 401 error into a 404 error.
	parseResult4 := common.ConvertExpected401ErrInto404Err(input401Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult4.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 401 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step5: Check whether the function can normally recognize unexpected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult5 := common.ConvertExpected401ErrInto404Err(input401Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult5, input401Err) {
		t.Fatalf("error converting 401 error to 404 error via a non-exist error code")
	}
}

func TestErrorFunc_ConvertExpected403ErrInto404Err(t *testing.T) {
	input403Err := golangsdk.ErrDefault403{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0002\", \"error_msg\": \"Resource not found\"}"),
		},
	}
	input400Err := golangsdk.ErrDefault400{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0003\", \"error_msg\": \"Authentication Failed\"}"),
		},
	}

	// Step1: Check whether the function can normally identify the expected error code under follow input and convert
	// the 403 error into a 404 error.
	parseResult1 := common.ConvertExpected403ErrInto404Err(input403Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult1.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 403 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step2: Check whether the function can normally recognize unexpected 400 error under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertExpected403ErrInto404Err(input400Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult2.(golangsdk.ErrDefault404); ok {
		t.Fatalf("The expected 403 error was not recognized and was incorrectly converted")
	}
	// Step3: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertExpected403ErrInto404Err(input403Err, "err_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult3, input403Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step4: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 403 error into a 404 error.
	parseResult4 := common.ConvertExpected403ErrInto404Err(input403Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult4.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 403 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step5: Check whether the function can normally recognize unexpected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult5 := common.ConvertExpected403ErrInto404Err(input403Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult5, input403Err) {
		t.Fatalf("error converting 403 error to 404 error via a non-exist error code")
	}
}

func TestErrorFunc_ConvertExpected500ErrInto404Err(t *testing.T) {
	input500Err := golangsdk.ErrDefault500{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0002\", \"error_msg\": \"Resource not found\"}"),
		},
	}
	input400Err := golangsdk.ErrDefault400{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0003\", \"error_msg\": \"Authentication Failed\"}"),
		},
	}

	// Step1: Check whether the function can normally identify the expected error code under follow input and convert
	// the 500 error into a 404 error.
	parseResult1 := common.ConvertExpected500ErrInto404Err(input500Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult1.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 500 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step2: Check whether the function can normally recognize unexpected 400 error under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertExpected500ErrInto404Err(input400Err, "error_code", "TESTERR.0002")
	if _, ok := parseResult2.(golangsdk.ErrDefault404); ok {
		t.Fatalf("The expected 500 error was not recognized and was incorrectly converted")
	}
	// Step3: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertExpected500ErrInto404Err(input500Err, "err_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult3, input500Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step4: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 500 error into a 404 error.
	parseResult4 := common.ConvertExpected500ErrInto404Err(input500Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult4.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 500 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step5: Check whether the function can normally recognize unexpected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult5 := common.ConvertExpected500ErrInto404Err(input500Err, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult5, input500Err) {
		t.Fatalf("error converting 500 error to 404 error via a non-exist error code")
	}
}

func TestErrorFunc_ConvertUndefinedErrInto404Err(t *testing.T) {
	input403Err := golangsdk.ErrDefault403{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte("{\"error_code\": \"TESTERR.0002\", \"error_msg\": \"Authentication Failed\"}"),
		},
	}
	input409Err := golangsdk.ErrUnexpectedResponseCode{
		Actual: 409,
		Body:   []byte("{\"error_code\": \"TESTERR.0001\", \"error_msg\": \"Conflict, please try again later\"}"),
	}
	input415Err := golangsdk.ErrUnexpectedResponseCode{
		Actual: 415,
		Body:   []byte("{\"error_code\": \"TESTERR.0002\", \"error_msg\": \"Operation completed\"}"),
	}

	// Step1: Check whether the function can normally identify the expected error code under follow input and convert
	// the 409 error into a 404 error.
	parseResult1 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "error_code", "TESTERR.0001")
	if _, ok := parseResult1.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 409 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult1).String())
	}
	// Step2: Check whether the function can normally recognize unexpected 409 error (strcture is undefined) under
	// follow input and terminate subsequent processing, and directly return error.
	parseResult2 := common.ConvertUndefinedErrInto404Err(input415Err, 409, "error_code", "TESTERR.0001")
	if !reflect.DeepEqual(parseResult2, input415Err) {
		t.Fatalf("The expected 415 error was not recognized and was incorrectly converted")
	}
	// Step3: Check whether the function can normally recognize unexpected 403 error (strcture is defined) under
	// follow input and terminate subsequent processing, and directly return error.
	parseResult3 := common.ConvertUndefinedErrInto404Err(input403Err, 409, "error_code", "TESTERR.0001")
	if _, ok := parseResult3.(golangsdk.ErrDefault403); !ok {
		t.Fatalf("The expected 403 error was not recognized and was incorrectly converted")
	}
	// Step4: Check whether the function can normally recognize unexpected error code key under follow input and
	// terminate subsequent processing, and return 400 error.
	parseResult4 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "err_code", "TESTERR.0001")
	if _, ok := parseResult4.(golangsdk.ErrDefault400); !ok {
		t.Fatalf("The expected error key 'err_code' was not recognized and the error was incorrectly converted")
	}
	// Step5: Check whether the function can normally recognize expected error code key (but the error code is
	// unexcepted) under follow input and terminate subsequent processing, and directly return original error.
	parseResult5 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "error_code", "TESTERR.0002")
	if !reflect.DeepEqual(parseResult5, input409Err) {
		t.Fatalf("Illegal recognition of unexpected error code key and convert the error to other type")
	}
	// Step6: Check whether the function can normally identify the expected error code (during a expected code list)
	// under follow input and convert the 409 error into a 404 error.
	parseResult6 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "error_code",
		[]string{"TESTERR.0001", "TESTERR.0002"}...)
	if _, ok := parseResult6.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("Unable to convert 409 error to 404 error, the result type of the convert function is: %s",
			reflect.TypeOf(parseResult6).String())
	}
	// Step7: Check whether the function can normally recognize expected error code under (during a unmatched code
	// list) follow input and terminate subsequent processing, and directly return error.
	parseResult7 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "error_code",
		[]string{"TESTERR.0002", "TESTERR.0003"}...)
	if !reflect.DeepEqual(parseResult7, input409Err) {
		t.Fatalf("error converting 409 error to 404 error via a non-exist error code")
	}
	// Step8: Check whether the function can normally recognize unexpected error status number follow input and
	// terminate subsequent processing, skipping the error code check and directly return 404 error.
	parseResult8 := common.ConvertUndefinedErrInto404Err(input409Err, 409, "")
	if _, ok := parseResult8.(golangsdk.ErrDefault404); !ok {
		t.Fatalf("error converting 409 error to 404 error with the omitted error code key")
	}
}
