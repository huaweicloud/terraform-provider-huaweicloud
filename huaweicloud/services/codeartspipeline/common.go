package codeartspipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	pipelineNotFoundError = "DEVPIPE.00011401"
	templateNotFoundError = "DEVPIPE.00011203"
	projectNotFoundError  = "DEV_21_100169"
	projectNotFoundError2 = "DEVPIPE.00011412"
)

// checkResponseError use to check whether the CodeArts Pipeline API response body contains error code.
// Parameter 'respBody' is the response body and 'notFoundCode' is the error code when the resource is not found.
// An example of an error response body is as follows: {"error_code": "XXX", "error_msg": "XXX"}
// Another example of an error response body is as follows: {"error": {"code": "xxx","message": "xxx"}}
func checkResponseError(respBody interface{}, notFoundCode ...string) error {
	errorCode := utils.PathSearch("error_code", respBody, "")
	errorCodeInErrorStruct := utils.PathSearch("error.code", respBody, "")
	if errorCode == "" && errorCodeInErrorStruct == "" {
		return nil
	}

	var err error
	if errorCode != "" {
		errorMsg := utils.PathSearch("error_msg", respBody, "")
		err = fmt.Errorf("error code: %s, error message: %s", errorCode, errorMsg)
		if !utils.StrSliceContains(notFoundCode, errorCode.(string)) {
			return err
		}
	} else {
		errorMsg := utils.PathSearch("error.message", respBody, "")
		err = fmt.Errorf("error code: %s, error message: %s", errorCodeInErrorStruct, errorMsg)
		if !utils.StrSliceContains(notFoundCode, errorCodeInErrorStruct.(string)) {
			return err
		}
	}

	return golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(err.Error()),
		},
	}
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

func encodeIntoJson(v interface{}) interface{} {
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

func resourceImportStateFuncWithProjectIdAndId(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("project_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving project ID: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
