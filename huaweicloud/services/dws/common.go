package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func waitClusterTaskStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, clusterId string) error {
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterTaskStateFun(client, clusterId),
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DWS cluster task to complete: %s", err)
	}

	return nil
}

func refreshClusterTaskStateFun(client *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetClusterInfoByClusterId(client, clusterId)
		if err != nil {
			return nil, "ERROR", err
		}

		taskStatus := utils.PathSearch("cluster.task_status", respBody, "").(string)
		status := utils.PathSearch("cluster.status", respBody, "").(string)
		if taskStatus == "" && utils.StrSliceContains([]string{"AVAILABLE", "ACTIVE"}, status) {
			return respBody, "COMPLETED", nil
		}

		// ACTIVE_STANDY_SWITCHOVER: The active-standby switchover is beibng performed.
		if strings.Contains(taskStatus, "ING") || taskStatus == "ACTIVE_STANDY_SWITCHOVER" {
			return respBody, "PENDING", nil
		}

		return respBody, taskStatus, nil
	}
}
