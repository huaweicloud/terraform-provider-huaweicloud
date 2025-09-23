package codeartsbuild

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	buildTaskNotFoundErr = "DEVCB.00031006"
	templateNotFoundErr  = "DEV.CB.0520002"
)

// checkResponseError use to check whether the CodeArts Build API response with OkCode but body contains error code.
// An example of an error response body is as follows: {"error_code": "XXX", "error_msg": "XXX"}
func checkResponseError(respBody interface{}) error {
	errorCode := utils.PathSearch("error_code", respBody, "")
	if errorCode == "" {
		return nil
	}

	errorMsg := utils.PathSearch("error_msg", respBody, "")
	return fmt.Errorf("error code: %s, error message: %s", errorCode, errorMsg)
}

func parseJson(v string) interface{} {
	if v == "" {
		return nil
	}

	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return v
	}

	return data
}

func encodeJson(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rst, err := json.Marshal(v)
	if err != nil {
		log.Printf("[DEBUG] Unable to encode into json: %s", err)
		return nil
	}

	return string(rst)
}
