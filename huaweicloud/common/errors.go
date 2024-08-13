package common

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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

	errCode, searchErr := jmespath.Search(errCodeKey, apiError)
	if searchErr != nil || errCode == nil {
		log.Printf("[WARN] Unable to find the expected error code key (%s)", errCodeKey)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 400 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if errCodeStr, ok := errCode.(string); ok && utils.StrSliceContains(specErrCodes, errCodeStr) {
		log.Printf("[INFO] Identified 400 error with code '%s' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%s), want %v", errCode, specErrCodes)
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

	errCode, searchErr := jmespath.Search(errCodeKey, apiError)
	if searchErr != nil || errCode == nil {
		log.Printf("[WARN] Unable to find the expected error code key (%s)", errCodeKey)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 401 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if errCodeStr, ok := errCode.(string); ok && utils.StrSliceContains(specErrCodes, errCodeStr) {
		log.Printf("[INFO] Identified 401 error with code '%s' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%s), want %v", errCode, specErrCodes)
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

	errCode, searchErr := jmespath.Search(errCodeKey, apiError)
	if searchErr != nil || errCode == nil {
		log.Printf("[WARN] Unable to find the expected error code key (%s)", errCodeKey)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 403 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if errCodeStr, ok := errCode.(string); ok && utils.StrSliceContains(specErrCodes, errCodeStr) {
		log.Printf("[INFO] Identified 403 error with code '%s' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%s), want %v", errCode, specErrCodes)
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

	errCode, searchErr := jmespath.Search(errCodeKey, apiError)
	if searchErr != nil || errCode == nil {
		log.Printf("[WARN] Unable to find the expected error code key (%s)", errCodeKey)
		return err
	}

	if len(specErrCodes) < 1 {
		log.Printf("[INFO] Identified 500 error parsed it as 404 error (without the error code control)")
		return golangsdk.ErrDefault404{}
	}
	if errCodeStr, ok := errCode.(string); ok && utils.StrSliceContains(specErrCodes, errCodeStr) {
		log.Printf("[INFO] Identified 500 error with code '%s' and parsed it as 404 error", errCode)
		return golangsdk.ErrDefault404{}
	}
	log.Printf("[WARN] Unable to recognize expected error code (%s), want %v", errCode, specErrCodes)
	return err
}
