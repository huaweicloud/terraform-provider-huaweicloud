package rabbitmq

import (
	"encoding/json"
	"errors"
	"strings"

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

func getInstanceTaskById(client *golangsdk.ServiceClient, instanceId, taskId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/tasks/{task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}
