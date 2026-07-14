package taurusdb

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

var (
	// Some error codes that need to be retried coming from https://support.huaweicloud.com/api-gaussdbformysql/ErrorCode.html
	retryErrCodes = map[string]struct{}{
		"DBS.200019":   {},
		"DBS.201014":   {},
		"DBS.201015":   {},
		"DBS.200047":   {},
		"DBS.05000084": {},
	}
)

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
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
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}
