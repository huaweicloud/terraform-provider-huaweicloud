package rocketmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
	getTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks/{task_id}"
	getTaskPath := client.Endpoint + getTaskHttpUrl
	getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", client.ProjectID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{instance_id}", instanceId)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", taskId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getTaskPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func refreshInstanceTaskStatus(client *golangsdk.ServiceClient, instanceID, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getInstanceTaskById(client, instanceID, taskId)
		if err != nil {
			return respBody, "ERROR", err
		}

		// status options: DELETED, SUCCESS, EXECUTING, FAILED, CREATED
		// CREATED means the task is in progress.
		status := utils.PathSearch("tasks[0].status", respBody, "").(string)
		if status == "FAILED" {
			return respBody, "FAILED", fmt.Errorf("unexpect status (%s)", status)
		}

		if status == "SUCCESS" {
			return respBody, "COMPLETED", nil
		}

		return "continue", "PENDING", nil
	}
}

func waitForInstanceTaskStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceId, taskId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshInstanceTaskStatus(client, instanceId, taskId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
