package css

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
)

func ConvertExpectedHwSdkErrInto404Err(err error, httpStatusCode int, expectCode, errMsg string) error {
	if sdkErr, ok := err.(*sdkerr.ServiceResponseError); ok {
		if sdkErr.StatusCode != httpStatusCode {
			log.Printf("[WARN] Unable to recognize expected error type, want '%d', but got '%d'", httpStatusCode, sdkErr.StatusCode)
			return err
		}

		var apiError ResponseError
		pErr := json.Unmarshal([]byte(sdkErr.ErrorMessage), &apiError)
		if pErr != nil {
			log.Printf("[WARN] failed to parse response error message: %s", pErr)
			return err
		}

		if apiError.ErrorCode == expectCode {
			if errMsg != "" && !strings.Contains(apiError.ErrorMsg, errMsg) {
				return err
			}
			return golangsdk.ErrDefault404{}
		}
	}

	return err
}

func ConvertExpectOtherErrInto404Err(err error, expectStatusCode int, errCodeKey, expectErrCode string) error {
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == expectStatusCode {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search(errCodeKey, response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse %s from response body: %s", errCodeKey, parseErr)
			}

			if errorCode == expectErrCode {
				return golangsdk.ErrDefault404{}
			}
		}
	}

	return err
}
