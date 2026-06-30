package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterActionNonUpdatableParams = []string{
	"cluster_id",
	"action",
}

// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/restart
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/stop
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/start
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceClusterAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterActionCreate,
		ReadContext:   resourceClusterActionRead,
		UpdateContext: resourceClusterActionUpdate,
		DeleteContext: resourceClusterActionDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the cluster is located.",
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the cluster to be operated.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The action type of the operation.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func doClusterAction(client *golangsdk.ServiceClient, clusterId, action string) error {
	var httpUrl string

	switch action {
	case "restart":
		httpUrl = "v1.0/{project_id}/clusters/{cluster_id}/restart"
	case "stop":
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/stop"
	case "start":
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/start"
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", actionPath, &opt)
	return err
}

func refreshClusterActionStateFunc(client *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetClusterInfoByClusterId(client, clusterId)
		if err != nil {
			return respBody, "ERROR", err
		}

		statusesInProgress := []string{"REBOOTING", "STARTING", "STOPPING"}
		taskStatus := utils.PathSearch("cluster.task_status", respBody, "").(string)
		if taskStatus == "" {
			return "NOT_IN_PROGRESS", "COMPLETED", nil
		}
		if utils.StrSliceContains(statusesInProgress, taskStatus) {
			return "IN_PROGRESS", "PENDING", nil
		}

		return respBody, "ERROR", fmt.Errorf("unexpected task status: %s", taskStatus)
	}
}

func waitForClusterActionCompleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterActionStateFunc(client, clusterId),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClusterActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		action    = d.Get("action").(string)
	)

	// For the same DWS cluster, it is not supported to run multiple tasks at the same time.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = doClusterAction(client, clusterId, action)
	if err != nil {
		return diag.Errorf("error operating cluster (%s) with action (%s): %s", clusterId, action, err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	err = waitForClusterActionCompleted(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for cluster (%s) to be available: %s", clusterId, err)
	}

	return resourceClusterActionRead(ctx, d, meta)
}

func resourceClusterActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the DWS cluster. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
