package taurusdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		taurusdbInstanceId  = d.Get("taurusdb_instance_id").(string)
		starrocksInstanceId = d.Get("starrocks_instance_id").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	err = restartStarrocksInstance(ctx, client, taurusdbInstanceId, starrocksInstanceId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(starrocksInstanceId)
	return nil
}

func restartStarrocksInstance(ctx context.Context, client *golangsdk.ServiceClient,
	instanceId, starrocksInstanceId string, timeout time.Duration) error {
	restartPath := client.Endpoint + "v3/{project_id}/instances/{starrocks_instance_id}/starrocks/restart"
	restartPath = strings.ReplaceAll(restartPath, "{project_id}", client.ProjectID)
	restartPath = strings.ReplaceAll(restartPath, "{starrocks_instance_id}", starrocksInstanceId)

	restartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", restartPath, &restartOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
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
	restartRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", restartRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error restarting TaurusDB Htap StarRocks instance(%s), job_id is not found in the response",
			starrocksInstanceId)
	}

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for restarting TaurusDB Htap StarRocks instance(%s) job to complete: %s",
			starrocksInstanceId, err)
	}

	return nil
}

func htapInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, htapInstanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Check TaurusDB HTAP instance status is normal and actions is empty list
		htapInstanceDetail, err := GetHtapInstanceDetail(client, instanceId, htapInstanceId)
		if err != nil {
			return nil, "", err
		}
		htapInstanceStatus := utils.PathSearch("status", htapInstanceDetail, "").(string)
		htapInstanceActions := utils.PathSearch("actions", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
		if htapInstanceStatus == "normal" && len(htapInstanceActions) == 0 {
			return checkHtapInstanceNodesStatus(htapInstanceDetail)
		}
		if htapInstanceStatus == "abnormal" || htapInstanceStatus == "createfail" {
			return nil, "", fmt.Errorf("TaurusDB HTAP instance(%s) is in status(%s), cannot proceed with new job",
				htapInstanceId, htapInstanceStatus)
		}
		return htapInstanceDetail, "pending", nil
	}
}
func checkHtapInstanceNodesStatus(htapInstanceDetail interface{}) (interface{}, string, error) {
	// Check status is normal and actions is empty list in all nodes
	groups := utils.PathSearch("groups", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
	for _, group := range groups {
		nodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
		for _, node := range nodes {
			nodeStatus := utils.PathSearch("status", node, "").(string)
			nodeActions := utils.PathSearch("actions", node, make([]interface{}, 0)).([]interface{})
			if nodeStatus == "normal" && len(nodeActions) == 0 {
				// If the node status is normal and actions is empty list, continue to check next node
				continue
			}
			// If the node status is abnormal or createfail, return error directly without retry
			if nodeStatus == "abnormal" || nodeStatus == "createfail" {
				return nil, "", fmt.Errorf("StarRocks node is in status(%s), cannot proceed with new job", nodeStatus)
			}
			return node, "pending", nil
		}
	}
	// If all nodes are normal and actions is empty list, return normal
	return htapInstanceDetail, "normal", nil
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
