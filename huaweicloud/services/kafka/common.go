package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceTask(client *golangsdk.ServiceClient, instanceId, taskId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/tasks/{task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func instanceTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskId string, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getInstanceTask(client, instanceID, taskId)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		// The status enumeration values are as follows: `DELETED`, `SUCCESS`, `EXECUTING` and `FAILED`.
		status := utils.PathSearch("tasks[0].status", resp, "").(string)
		if status == "FAILED" || status == "DELETED" {
			return resp, "FAILED", fmt.Errorf("unexpected status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

// For Job ID task, use waitForInstanceTaskStatusComplete method.
func waitForInstanceTaskStatusComplete(ctx context.Context, client *golangsdk.ServiceClient, instanceId, taskId string,
	timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceTaskStatusRefreshFunc(client, instanceId, taskId, []string{"SUCCESS"}),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getInstanceTasks(client *golangsdk.ServiceClient, instanceId string) ([]interface{}, error) {
	var (
		listHttpUrl = "v2/{project_id}/instances/{instance_id}/tasks"
		// limit range is `1` to `100`.
		limit = 100
		// `start` must be greater than or equal to `1`.
		start  = 1
		result = make([]interface{}, 0)
	)
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath += fmt.Sprintf("?limit=%v", limit)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithStart := listPath + fmt.Sprintf("&start=%d", start)
		resp, err := client.Request("GET", listPathWithStart, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		tasks := utils.PathSearch("tasks", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)
		if len(tasks) < limit {
			break
		}

		start += len(tasks)
	}

	return result, nil
}

func refreshTaskStatusByName(client *golangsdk.ServiceClient, instanceId, taskName string, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tasks, err := getInstanceTasks(client, instanceId)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		task := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", taskName), tasks, nil)
		// The status enumeration values are as follows: `DELETED`, `SUCCESS`, `EXECUTING` and `FAILED`.
		status := utils.PathSearch("status", task, "").(string)
		if status == "FAILED" || status == "DELETED" {
			return nil, "FAILED", fmt.Errorf("unexpected status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return task, "COMPLETED", nil
		}

		return task, "PENDING", nil
	}
}

// For task without job ID, use task name to wait for the task to complete.
func waitForInstanceTaskStatusCompleteByName(ctx context.Context, client *golangsdk.ServiceClient, instanceId, taskName string,
	timeout time.Duration, targets ...[]string) error {
	completedStatuses := []string{"SUCCESS"}
	if len(targets) > 0 {
		completedStatuses = targets[0]
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshTaskStatusByName(client, instanceId, taskName, completedStatuses),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
