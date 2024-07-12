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
