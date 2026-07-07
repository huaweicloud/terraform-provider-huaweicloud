package eg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventRouterClusterNonUpdatableParams = []string{
	"source_type",
	"sink_type",
	"vpc_id",
	"subnet_id",
	"availability_zones",
	"flavor",
}

// @API EG POST /v1/{project_id}/eventrouter/clusters
// @API EG GET /v1/{project_id}/eventrouter/clusters/{cluster_id}
// @API EG PUT /v1/{project_id}/eventrouter/clusters/{cluster_id}
// @API EG DELETE /v1/{project_id}/eventrouter/clusters/{cluster_id}
func ResourceEventRouterCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventRouterClusterCreate,
		ReadContext:   resourceEventRouterClusterRead,
		UpdateContext: resourceEventRouterClusterUpdate,
		DeleteContext: resourceEventRouterClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(eventRouterClusterNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the event router cluster is located.",
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the event router cluster.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The source type of the event router cluster.",
			},
			"sink_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The sink type of the event router cluster.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC ID to which the event router cluster belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet ID to which the event router cluster belongs.",
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the event router cluster.",
			},
			"availability_zones": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The availability zone names of the event router cluster.",
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The flavor of the event router cluster.",
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the event router cluster.",
			},
			"job_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of jobs running in the event router cluster.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the event router cluster, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the event router cluster, in RFC3339 format.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildEventRouterClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameter(s).
		"name":        d.Get("name"),
		"source_type": d.Get("source_type"),
		"sink_type":   d.Get("sink_type"),
		"vpc_id":      d.Get("vpc_id"),
		"subnet_id":   d.Get("subnet_id"),

		// Optional parameters(s).
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"zone_names":  utils.ValueIgnoreEmpty(d.Get("availability_zones")),
		"flavor":      utils.ValueIgnoreEmpty(d.Get("flavor")),
	}
}

func createEventRouterCluster(client *golangsdk.ServiceClient, bodyParams map[string]interface{}) (interface{}, error) {
	httpUrl := "v1/{project_id}/eventrouter/clusters"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         bodyParams,
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceEventRouterClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	respBody, err := createEventRouterCluster(client, buildEventRouterClusterBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating event router cluster: %s", err)
	}

	clusterId := utils.PathSearch("cluster_id", respBody, "").(string)
	if clusterId == "" {
		return diag.Errorf("unable to find the cluster ID from the API response")
	}
	d.SetId(clusterId)

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      eventRouterClusterStateRefreshFunc(client, clusterId, []string{"RUNNING", "RUNNING_PARTIALLY"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for event router cluster to become ready: %s", err)
	}

	return resourceEventRouterClusterRead(ctx, d, meta)
}

func eventRouterClusterStateRefreshFunc(client *golangsdk.ServiceClient, clusterId string, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetEventRouterClusterById(client, clusterId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "not_found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		unexpectStatus := []string{
			"ERROR",           // Cluster is down
			"FAILED",          // Create cluster failed
			"FREEZE_FAILED",   // Failed to freeze the cluster
			"COMPACT_FAILED",  // Failed to compact the cluster
			"UPGRADE_FAILED",  // Failed to upgrade the cluster
			"EXTENDED_FAILED", // Failed to extend the cluster
			"ROLLBACK_FAILED", // Failed to rollback the cluster
			"DELETED_FAILED",  // Failed to delete the cluster
		}
		if utils.StrSliceContains(unexpectStatus, status) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}

		// Some regions may return "DELETED" when the cluster is deleted.
		if status == "DELETED" && len(targets) < 1 {
			return "not_found", "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func GetEventRouterClusterById(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/eventrouter/clusters/{cluster_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	if status := utils.PathSearch("status", respBody, "").(string); status == "DELETED" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/eventrouter/clusters/{cluster_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The event router cluster (%s) not found", clusterId)),
			},
		}
	}
	return respBody, nil
}

func resourceEventRouterClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	cluster, err := GetEventRouterClusterById(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving event router cluster (%s)", clusterId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameter(s).
		d.Set("name", utils.PathSearch("name", cluster, nil)),
		d.Set("source_type", utils.PathSearch("source_type", cluster, nil)),
		d.Set("sink_type", utils.PathSearch("sink_type", cluster, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", cluster, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", cluster, nil)),
		// Optional parameter(s).
		d.Set("description", utils.PathSearch("description", cluster, nil)),
		d.Set("availability_zones", utils.PathSearch("zone_names", cluster, nil)),
		// Attributer(s).
		d.Set("status", utils.PathSearch("status", cluster, nil)),
		d.Set("job_count", utils.PathSearch("job_count", cluster, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_time",
			cluster, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_time",
			cluster, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEventRouterClusterUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Optional parameter(s).
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
	}
}

func updateEventRouterCluster(client *golangsdk.ServiceClient, clusterId string, bodyParams map[string]interface{}) error {
	httpUrl := "v1/{project_id}/eventrouter/clusters/{cluster_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         bodyParams,
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceEventRouterClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	err = updateEventRouterCluster(client, clusterId, buildEventRouterClusterUpdateBodyParams(d))
	if err != nil {
		return diag.Errorf("error updating event router cluster (%s): %s", clusterId, err)
	}

	return resourceEventRouterClusterRead(ctx, d, meta)
}

func deleteEventRouterCluster(client *golangsdk.ServiceClient, clusterId string) error {
	httpUrl := "v1/{project_id}/eventrouter/clusters/{cluster_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceEventRouterClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	err = deleteEventRouterCluster(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting event router cluster (%s)", clusterId))
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      eventRouterClusterStateRefreshFunc(client, clusterId, nil),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for event router cluster to be deleted: %s", err)
	}
	return nil
}
