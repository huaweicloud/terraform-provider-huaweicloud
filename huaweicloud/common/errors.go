package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var definedStatusNumbers = []interface{}{
	400, 401, 402, 403, 404, 405, 408, 429, 500, 503,
}

// ConvertExpected400ErrInto404Err is a method used to parsing 400 error and try to convert it to 404 error according
// to the right error code.
// Arguments:
// + err: The error response obtained through HTTP/HTTPS request.
// + errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
// + specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'APIGW.0001'.
// Notes: If you missing specErrCodes input, this function will convert all 400 errors into 404 errors.
// How to use it:
// + For the general cases, their error code key is 'error_code', and we should call as follow:
//   - utils.ConvertExpected400ErrInto404Err(err, "error_code")
//   - utils.ConvertExpected400ErrInto404Err(err, "error_code", "DLM.3027")
//   - utils.ConvertExpected400ErrInto404Err(err, "error_code", []string{"DLM.3027", "DLM.3028"}...)
func ConvertExpected400ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err400 golangsdk.ErrDefault400
	if !errors.As(err, &err400) {
		log.Printf("[WARN] Unable to recognize expected error type, want 'golangsdk.ErrDefault400', but got '%s'",
			reflect.TypeOf(err).String())
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err400.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := utils.PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		log.Printf("[WARN] Unable to find the error code from the error body using given error code key (%s), the error is: %#v",
			errCodeKey, apiError)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 400 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if utils.StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 400 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// ConvertExpected401ErrInto404Err is a method used to parsing 401 error and try to convert it to 404 error according
// to the right error code.
// Arguments:
// + err: The error response obtained through HTTP/HTTPS request.
// + errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
// + specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'FSS.0401'.
// Notes: If you missing specErrCodes input, this function will convert all 401 errors into 404 errors.
// How to use it:
// + For the general cases, their error code key is 'error_code', and we should call as follow:
//   - utils.ConvertExpected401ErrInto404Err(err, "error_code")
//   - utils.ConvertExpected401ErrInto404Err(err, "error_code", "FSS.0401")
//   - utils.ConvertExpected401ErrInto404Err(err, "error_code", []string{"FSS.0401", "FSS.0402"}...)
func ConvertExpected401ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err401 golangsdk.ErrDefault401
	if !errors.As(err, &err401) {
		log.Printf("[WARN] Unable to recognize expected error type, want 'golangsdk.ErrDefault401', but got '%s'",
			reflect.TypeOf(err).String())
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err401.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := utils.PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		log.Printf("[WARN] Unable to find the error code from the error body using given error code key (%s), the error is: %#v",
			errCodeKey, apiError)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 401 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if utils.StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 401 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// ConvertExpected403ErrInto404Err is a method used to parsing 403 error and try to convert it to 404 error according
// to the right error code.
// Arguments:
// + err: The error response obtained through HTTP/HTTPS request.
// + errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
// + specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'APIGW.0001'.
// Notes: If you missing specErrCodes input, this function will convert all 403 errors into 404 errors.
// How to use it:
// + For the general cases, their error code key is 'error_code', and we should call as follow:
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code")
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code", "DWS.0001")
//   - utils.ConvertExpected403ErrInto404Err(err, "error_code", []string{"DWS.0001", "DLM.3028"}...)
func ConvertExpected403ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err403 golangsdk.ErrDefault403
	if !errors.As(err, &err403) {
		log.Printf("[WARN] Unable to recognize expected error type, want 'golangsdk.ErrDefault403', but got '%s'",
			reflect.TypeOf(err).String())
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err403.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := utils.PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		log.Printf("[WARN] Unable to find the error code from the error body using given error code key (%s), the error is: %#v",
			errCodeKey, apiError)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 403 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if utils.StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 403 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// ConvertExpected500ErrInto404Err is a method used to parsing 500 error and try to convert it to 404 error according
// to the right error code.
// Arguments:
// + err: The error response obtained through HTTP/HTTPS request.
// + errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
// + specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'APIGW.0001'.
// Notes: If you missing specErrCodes input, this function will convert all 500 errors into 404 errors.
// How to use it:
// + For the general cases, their error code key is 'error_code', and we should call as follow:
//   - utils.ConvertExpected500ErrInto404Err(err, "error_code")
//   - utils.ConvertExpected500ErrInto404Err(err, "error_code", "SCM.0016")
//   - utils.ConvertExpected500ErrInto404Err(err, "error_code", []string{"SCM.0016", "SCM.0017"}...)
func ConvertExpected500ErrInto404Err(err error, errCodeKey string, specErrCodes ...string) error {
	var err500 golangsdk.ErrDefault500
	if !errors.As(err, &err500) {
		log.Printf("[WARN] Unable to recognize expected error type, want 'golangsdk.ErrDefault500', but got '%s'",
			reflect.TypeOf(err).String())
		return err
	}
	var apiError interface{}
	if jsonErr := json.Unmarshal(err500.Body, &apiError); jsonErr != nil {
		return err
	}

	errCode := utils.PathSearch(errCodeKey, apiError, nil)
	if errCode == nil {
		log.Printf("[WARN] Unable to find the error code from the error body using given error code key (%s), the error is: %#v",
			errCodeKey, apiError)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 500 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if utils.StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified 500 error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// ConvertUndefinedErrInto404Err is a method used to parsing errors which related structures undefined in the golangsdk
// package (so a general method for parsing such errors is needed) and try to convert it to 404 error according to the
// right error number and code (omitted means all passed).
// Unsupported status codes (may not enter the convert logic and returned directly):
//   - 400, 401, 402, 403, 404, 405, 408, 429, 500, 503
//
// Arguments:
//   - err: The error response obtained through HTTP/HTTPS request.
//   - errStatusNum: The status number of the error, e.g. 409, 415 (The common status numbers are defined in the
//     golangsdk package, e.g. 400, 403, 500)
//   - errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
//   - specErrCodes: One or more error codes that you wish to match against the current error, e.g. 'CSS.0004'.
//
// Notes: If you missing specErrCodes input, this function will convert all errors which matched your specifies status
// into 404 errors.
//
// How to use it:
//   - utils.ConvertUndefinedErrInto404Err(err, 415, "error_code")
//   - utils.ConvertUndefinedErrInto404Err(err, 415, "error_code", "CSS.0004")
//   - utils.ConvertUndefinedErrInto404Err(err, 415, "error_code", []string{"CSS.0004", "CSS.0005"}...)
//   - utils.ConvertUndefinedErrInto404Err(err, 415, "")
//
// The corresponding processing of this method are as follows:
//   - The status code has a corresponding structure definition in golangsdk: the error is recorded in the log and the
//     original error is returned directly.
//   - The status code (key) not found: return error (404 error: if the errCodeKey is omitted).
//   - The input of the expected error code(s) is omitted: return 404 error.
//   - The expected error code(s) matched: return 404 error.
//   - Other situations: return original error.
func ConvertUndefinedErrInto404Err(err error, errStatusNum int, errCodeKey string, specErrCodes ...string) error {
	if utils.SliceContains(definedStatusNumbers, errStatusNum) {
		log.Printf("[INFO] The error with status code %d already has a corresponding structure definition in the "+
			"golangsdk package, please use the corresponding conversion method to process", errStatusNum)
		return err
	}
	errCode, optErr := getExpectedErrCode(err, errStatusNum, errCodeKey)
	if optErr != nil {
		return optErr
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if utils.StrSliceContains(specErrCodes, fmt.Sprint(errCode)) {
		log.Printf("[INFO] Identified error with code '%v' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%v), want %v", errCode, specErrCodes)
	return err
}

// getExpectedErrCode is a method used to identify whether the status code of the target error meets expectations, and
// attempts to obtain the error code corresponding to the target error according to the error code key of the target
// error.
//
// Arguments:
//   - err: The error response obtained through HTTP/HTTPS request.
//   - errStatusNum: The status number of the error, e.g. 409, 415 (The common status numbers are defined in the
//     golangsdk package, e.g. 400, 403, 500)
//   - errCodeKey: The key name of the error code in the error response body, e.g. 'error_code', 'err_code'.
//
// How to use it:
//   - getExpectedErrCode(err, 409, "error_code")
//   - getExpectedErrCode(err, 409, "")
//
// The corresponding processing of this method are as follows:
//   - The status code has a corresponding structure definition in golangsdk: the error is recorded in the log and the
//     original error is returned directly, and the error code is returned as null.
//   - The status code does not match: the error is recorded in the log and the original error is returned directly,
//     and the error code is returned as null.
//   - The status code matches and the error code key is null: a 404 error is returned directly, skipping the error code
//     check.
//   - The status code matches but the error cannot be parsed normally: a 400 error is returned, and the content is
//     error parsing failure.
//   - The status code matches but the error code is not found: a 400 error is returned, and the content is that the
//     error code does not match.
//   - The status code matches and the target error code is found: the corresponding error code is returned.
func getExpectedErrCode(err error, errStatusNum int, errCodeKey string) (string, error) {
	var apiError interface{}
	parsedErr, ok := err.(golangsdk.ErrUnexpectedResponseCode)
	if !ok {
		log.Printf("[WARN] Failed to recognize error type, want 'golangsdk.ErrUnexpectedResponseCode', but got '%s'",
			reflect.TypeOf(err).String())
		return "", err
	} else if parsedErr.Actual != errStatusNum {
		log.Printf("[WARN] Unable to recognize expected error status number, want '%d', but got '%d'",
			errStatusNum, parsedErr.Actual)
		return "", err
	}

	if errCodeKey == "" {
		return "", golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("The string input of errCodeKey is not detected, so the error is directly " +
					"converted to a 404 error and returned"),
			},
		}
	}
	jsonErr := json.Unmarshal(parsedErr.Body, &apiError)
	if jsonErr != nil {
		return "", golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("error parsing the request error: %s", jsonErr)),
			},
		}
	}
	errCode := utils.PathSearch(errCodeKey, apiError, "").(string)
	if errCode == "" {
		// 4xx means the client parsing was failed.
		return errCode, golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("Unable to find the error code from the error body using given status "+
					"number (%d) and the error code key (%s), the error is: '%v'", errStatusNum, errCodeKey, err)),
			},
		}
	}
	return errCode, nil
}
