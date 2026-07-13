package taurusdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var starrocksInstanceRestartNoneUpdatableParams = []string{
	"taurusdb_instance_id", "starrocks_instance_id",
}

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}
// @API TaurusDB PUT /v3/{project_id}/instances/{starrocks_instance_id}/starrocks/restart
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBHtapStarrocksInstanceRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksInstanceRestartCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksInstanceRestartRead,
		UpdateContext: resourceTaurusDBHtapStarrocksInstanceRestartUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksInstanceRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(starrocksInstanceRestartNoneUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"taurusdb_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"starrocks_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceTaurusDBHtapStarrocksInstanceRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		instanceId          = d.Get("taurusdb_instance_id").(string)
		starrocksInstanceId = d.Get("starrocks_instance_id").(string)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	err = restartStarrocksInstance(ctx, client, instanceId, starrocksInstanceId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(starrocksInstanceId)

	return nil
}

func restartStarrocksInstance(ctx context.Context, client *golangsdk.ServiceClient,
	instanceId, starrocksInstanceId string, timeout time.Duration) error {
	createPath := client.Endpoint + "v3/{project_id}/instances/{starrocks_instance_id}/starrocks/restart"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{starrocks_instance_id}", starrocksInstanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		shouldRetry, err := handleMultiOperationsError(err)
		return res, shouldRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, instanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting TaurusDB Htap StarRocks instance(%s): %s", starrocksInstanceId, err)
	}
	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error restarting StarRocks instance(%s), job_id is not found in the response",
			starrocksInstanceId)
	}

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for restarting StarRocks instance(%s) job to complete: %s",
			starrocksInstanceId, err)
	}
	return nil
}

func htapInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, htapInstanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Check TaurusDB HTAP instance status is normal and actions is empty list
		htapInstanceDetail, err := GetHtapInstanceDetail(client, instanceId, htapInstanceId)
		if err != nil {
			return nil, "", err
		}
		htapInstanceStatus := utils.PathSearch("status", htapInstanceDetail, "").(string)
		// If the status is abnormal or createfail, return error directly without retry
		if htapInstanceStatus == "abnormal" || htapInstanceStatus == "createfail" {
			return nil, "", fmt.Errorf("TaurusDB HTAP instance(%s) is in %s state, cannot proceed with restart",
				htapInstanceId, htapInstanceStatus)
		}
		if htapInstanceStatus != "normal" {
			return htapInstanceDetail, htapInstanceStatus, nil
		}
		htapActions := utils.PathSearch("actions", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
		if len(htapActions) == 0 {
			// Check status is normal and actions is empty list in all nodes
			groups := utils.PathSearch("groups", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
			if len(groups) == 0 {
				return htapInstanceDetail, "normal", nil
			}
			for _, group := range groups {
				nodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
				if len(nodes) == 0 {
					continue
				}
				for _, node := range nodes {
					nodeStatus := utils.PathSearch("status", node, "").(string)
					// If the node status is abnormal or createfail, return error directly without retry
					if nodeStatus == "abnormal" || nodeStatus == "createfail" {
						return nil, "", fmt.Errorf("StarRocks node is in %s state, cannot proceed with restart",
							nodeStatus)
					}
					if nodeStatus != "normal" {
						return node, nodeStatus, nil
					}
					nodeActions := utils.PathSearch("actions", node, make([]interface{}, 0)).([]interface{})
					if len(nodeActions) > 0 {
						// Return the action of the node as status
						nodeAction := utils.PathSearch("actions[0].action", node, "").(string)
						return node, nodeAction, nil
					}
				}
			}
			return htapInstanceDetail, "normal", nil
		}
		htapAction := utils.PathSearch("actions[0].action", htapInstanceDetail, "").(string)
		return htapInstanceDetail, htapAction, nil
	}
}

func GetHtapInstanceDetail(client *golangsdk.ServiceClient, instanceId, htapInstanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{starrocks_instance_id}", htapInstanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving TaurusDB instance(%s) StarRocks instance(%s): %s", instanceId, htapInstanceId, err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("instances[0]", getRespBody, nil), nil
}

func resourceTaurusDBHtapStarrocksInstanceRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksInstanceRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksInstanceRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the StarRocks instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
