package workspace

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func waitForJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshJobStatusFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "ERROR", err
	}

	status := "SUCCESS"
	// Check whether the latest sub job is success.
	subJobs := utils.PathSearch("sub_jobs", resp, make([]interface{}, 0)).([]interface{})
	for _, subJob := range subJobs {
		subJobStatus := utils.PathSearch("status", subJob, "NONE").(string)
		if subJobStatus != "SUCCESS" {
			status = "FAIL"
			log.Printf("[ERROR] error waiting for the job(ID: %s, type: %s) to become complete: %s",
				utils.PathSearch("id", subJob, "NOT_FOUND").(string),
				utils.PathSearch("job_type", subJob, "NOT_FOUND").(string),
				utils.PathSearch("fail_reason", subJob, "").(string))
		}
	}

	return status, nil
}

func refreshJobStatusFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/workspace-jobs/{job_id}"
			listOpts = golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			}
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{job_id}", jobId)
		resp, err := client.Request("GET", listPath, &listOpts)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		if status == "" {
			return respBody, "ERROR", errors.New("dispatch operation job failed")
		}

		if utils.StrSliceContains([]string{"COMPLETE", "SUCCESS", "FAIL"}, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}
