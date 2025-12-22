package mrs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getClusterById(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1.1/{project_id}/cluster_infos/{cluster_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("cluster", respBody, nil), nil
}

func waitForClusterStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId string, timeout time.Duration,
	targets ...string) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"COMPLETED"},
		Pending:      []string{"PENDING"},
		Refresh:      refreshClusterStatusFunc(client, clusterId, targets...),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshClusterStatusFunc(client *golangsdk.ServiceClient, clusterId string, targets ...string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getClusterById(client, clusterId)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("clusterState", resp, "").(string)

		if status == "failed" {
			return resp, "ERROR", fmt.Errorf("unexpected status: %s", status)
		}

		if len(targets) == 0 {
			targets = []string{"running", "abnormal"}
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}
