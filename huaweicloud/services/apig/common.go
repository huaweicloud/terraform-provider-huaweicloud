package apig

import (
	"encoding/json"
	"fmt"

	"github.com/chnsz/golangsdk"
)

type requestErr struct {
	// The error code.
	ErrCode string `json:"error_code"`
	// The error message.
	ErrMsg string `json:"error_msg"`
}

// The APIG API is limited to only one attach operation at a time for standrad policy and plugin policy.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}

	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError requestErr
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		// Some error codes that need to be retried coming from https://console-intl.huaweicloud.com/apiexplorer/#/errorcenter/APIG.
		retryErrCodes := map[string]struct{}{
			"APIG.3500": {},
		}
		if _, ok := retryErrCodes[apiError.ErrCode]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}
