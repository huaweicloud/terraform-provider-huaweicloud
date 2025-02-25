package rabbitmq

import (
	"encoding/json"
	"errors"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, errors.New("unmarshal the response body failed: " + jsonErr.Error())
		}

		errorCode := utils.PathSearch("error_code", apiError, "").(string)
		if errorCode == "" {
			return false, errors.New("unable to find error code from the API response")
		}

		// CBC.99003651: unsubscribe fail, another operation is being performed
		if errorCode == "DMS.00400026" || errorCode == "CBC.99003651" {
			return true, err
		}
	}
	return false, err
}
