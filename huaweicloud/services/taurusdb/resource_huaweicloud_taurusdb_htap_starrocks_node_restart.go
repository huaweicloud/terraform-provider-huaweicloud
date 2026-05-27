package taurusdb

import (
	"context"
	"errors"
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

var starrocksNodeRestartNoneUpdatableParams = []string{
	"taurusdb_instance_id", "starrocks_instance_id", "starrocks_node_id",
}

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}
// @API TaurusDB PUT /v3/{project_id}/instances/{starrocks_instance_id}/starrocks/{starrocks_node_id}/restart
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBHtapStarrocksNodeRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksNodeRestartCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksNodeRestartRead,
		UpdateContext: resourceTaurusDBHtapStarrocksNodeRestartUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksNodeRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(starrocksNodeRestartNoneUpdatableParams),

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
			"starrocks_node_id": {
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

func resourceTaurusDBHtapStarrocksNodeRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	starrocksInstanceId := d.Get("starrocks_instance_id").(string)
	starrocksNodeId := d.Get("starrocks_node_id").(string)

	jobId, err := restartStarrocksNode(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(starrocksNodeId)

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for restarting StarRocks instance(%s) node(%s) job to complete: %s",
			starrocksInstanceId, starrocksNodeId, err)
	}

	return nil
}

func restartStarrocksNode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{starrocks_instance_id}/starrocks/{starrocks_node_id}/restart"
	)

	instanceId := d.Get("taurusdb_instance_id").(string)
	starrocksInstanceId := d.Get("starrocks_instance_id").(string)
	starrocksNodeId := d.Get("starrocks_node_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{starrocks_instance_id}", starrocksInstanceId)
	createPath = strings.ReplaceAll(createPath, "{starrocks_node_id}", starrocksNodeId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceNodeStateRefreshFunc(client, instanceId, starrocksInstanceId, starrocksNodeId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return "", fmt.Errorf("error restarting TaurusDB instance(%s) node(%s): %s", starrocksInstanceId, starrocksNodeId, err)
	}
	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return "", err
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return "", fmt.Errorf("error restarting StarRocks instance(%s) node(%s), job_id is not found in the response",
			starrocksInstanceId, starrocksNodeId)
	}

	return jobId, nil
}

func htapInstanceNodeStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, htapInstanceId, htapNodeId string) resource.StateRefreshFunc {
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
			// Check TaurusDB instance node status is normal and actions is empty list
			nodeDetail, err := GetHtapInstanceNodeDetail(htapInstanceDetail, htapNodeId)
			if err != nil {
				return nil, "", err
			}
			nodeStatus := utils.PathSearch("status", nodeDetail, "").(string)
			// If the node status is abnormal or createfail, return error directly without retry
			if nodeStatus == "abnormal" || nodeStatus == "createfail" {
				return nil, "", fmt.Errorf("StarRocks node(%s) is in %s state, cannot proceed with restart",
					htapNodeId, nodeStatus)
			}
			if nodeStatus != "normal" {
				return nodeDetail, nodeStatus, nil
			}
			nodeActions := utils.PathSearch("actions", nodeDetail, make([]interface{}, 0)).([]interface{})
			if len(nodeActions) == 0 {
				return nodeDetail, "normal", nil
			}
			nodeAction := utils.PathSearch("actions[0].action", nodeDetail, "").(string)
			if nodeAction != "" {
				return nodeDetail, nodeAction, nil
			}
			return nodeDetail, "normal", nil
		}
		htapAction := utils.PathSearch("actions[0].action", htapInstanceDetail, "").(string)
		return htapInstanceDetail, htapAction, nil
	}
}

func GetHtapInstanceNodeDetail(instanceDetail interface{}, htapNodeId string) (interface{}, error) {
	groups := utils.PathSearch("groups", instanceDetail, make([]interface{}, 0)).([]interface{})
	if len(groups) == 0 {
		return nil, errors.New("no groups found in TaurusDB StarRocks instance")
	}
	for _, group := range groups {
		nodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
		if len(nodes) == 0 {
			continue
		}

		for _, node := range nodes {
			nodeId := utils.PathSearch("id", node, "").(string)
			if nodeId == htapNodeId {
				return node, nil
			}
		}
	}
	return nil, fmt.Errorf("node (%s) not found in HTAP instance", htapNodeId)
}

func resourceTaurusDBHtapStarrocksNodeRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksNodeRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksNodeRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the StarRocks instance node remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
