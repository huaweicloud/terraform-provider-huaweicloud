package rds

import (
	"encoding/json"
	"fmt"

	"github.com/chnsz/golangsdk"
)

type requestErr struct {
	// The error code.
	ErrCode string `json:"error_code"`
	// The error message.
	ErrMsg string `json:"error_message"`
}

// The RDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError requestErr
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}
		if apiError.ErrCode == "DBS.200019" {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}
