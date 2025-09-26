package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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

func instanceTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getInstanceTask(client, instanceID, taskId)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		// DELETED, SUCCESS, EXECUTING, FAILED
		status := utils.PathSearch("tasks[0].status", resp, "").(string)
		if status == "FAILED" || status == "DELETED" {
			return resp, "FAILED", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func waitForInstanceTaskStatusComplete(ctx context.Context, client *golangsdk.ServiceClient, instanceId, taskId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
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
