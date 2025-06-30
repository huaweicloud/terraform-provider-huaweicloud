package rds

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

var (
	// Some error codes that need to be retried coming from https://console-intl.huaweicloud.com/apiexplorer/#/errorcenter/RDS.
	retryErrCodes = map[string]struct{}{
		"DBS.201202":   {},
		"DBS.200011":   {},
		"DBS.200018":   {},
		"DBS.200019":   {},
		"DBS.200047":   {},
		"DBS.200076":   {},
		"DBS.200611":   {},
		"DBS.200080":   {},
		"DBS.200463":   {}, // create replica instance
		"DBS.201015":   {},
		"DBS.201206":   {},
		"DBS.212033":   {}, // http response code is 403
		"DBS.280011":   {},
		"DBS.280343":   {},
		"DBS.280816":   {},
		"DBS.01010337": {},
		"DBS.01280030": {}, // instance status is illegal
	}
)

// The RDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code||errCode", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// if the error code is RDS.0005, it indicates that the SSL is changed, and the db is rebooted
		if errorCode.(string) == "RDS.0005" {
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}

// The RDS instance can not be deleted or unsubscribe if another operation is being performed.
func handleDeletionError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	// unsubscribe fail
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// CBC.99003651: Another operation is being performed.
		if errorCode == "CBC.99003651" {
			return true, err
		}
	}
	// delete fail
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	return false, err
}

// The RDS cross region backup strategy can not be updated if another operation is being performed.
func handleCrossRegionBackupStrategyError(err error) (bool, error) {
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// DBS.280228: Another operation is being performed.
		if errorCode == "DBS.280228" {
			return true, err
		}
	}
	return false, err
}

func handleApiNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	if errCode, ok := err.(golangsdk.ErrDefault404); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false
		}
		if errorCode.(string) == "APIGW.0101" {
			return true
		}
	}
	return false
}

func handleTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	// if the http status code is 500 and the error code is DBS.111205, it indicates timeout for the service
	// error should be ignored, just wait for success
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false
		}
		if errorCode.(string) == "DBS.111205" {
			return true
		}
	}
	return false
}
