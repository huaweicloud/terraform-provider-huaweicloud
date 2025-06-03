package codeartspipeline

import (
	"fmt"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	templateNotFoundError = "DEVPIPE.00011203"
)

// checkResponseError use to check whether the CodeArts Pipeline API response body contains error code.
// Parameter 'respBody' is the response body and 'notFoundCode' is the error code when the resource is not found.
// An example of an error response body is as follows: {"error_code": "XXX", "error_msg": "XXX"}
func checkResponseError(respBody interface{}, notFoundCode string) error {
	errorCode := utils.PathSearch("error_code", respBody, "")
	if errorCode == "" {
		return nil
	}

	errorMsg := utils.PathSearch("error_msg", respBody, "")
	err := fmt.Errorf("error code: %s, error message: %s", errorCode, errorMsg)
	if errorCode != notFoundCode {
		return err
	}

	return golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(err.Error()),
		},
	}
}
