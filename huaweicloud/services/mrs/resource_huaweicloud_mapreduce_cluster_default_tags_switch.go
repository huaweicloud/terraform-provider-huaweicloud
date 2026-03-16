package mrs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	clusterDefaultTagsActionCreate = "create"
	clusterDefaultTagsActionDelete = "delete"
)

var clusterDefaultTagsSwitchNonUpdatable = []string{"cluster_id"}

// @API MRS POST /v2/{project_id}/clusters/{cluster_id}/tags/switch
// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/tags/status
func ResourceClusterDefaultTagsSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterDefaultTagsSwitchCreate,
		ReadContext:   resourceClusterDefaultTagsSwitchRead,
		UpdateContext: resourceClusterDefaultTagsSwitchUpdate,
		DeleteContext: resourceClusterDefaultTagsSwitchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(clusterDefaultTagsSwitchNonUpdatable),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the cluster default tags switch is located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the action.`,
			},
			"default_tags_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the default tags switch is enabled.`,
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

func switchClusterDefaultTags(client *golangsdk.ServiceClient, clusterId, action string) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/tags/switch"
	switchPath := client.Endpoint + httpUrl
	switchPath = strings.ReplaceAll(switchPath, "{project_id}", client.ProjectID)
	switchPath = strings.ReplaceAll(switchPath, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": action,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("POST", switchPath, &opt)
	return err
}

func GetClusterDefaultTagsSwitchStatus(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/tags/status"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitForClusterDefaultTagsStatus(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterDefaultTagsStatusFunc(client, clusterId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshClusterDefaultTagsStatusFunc(client *golangsdk.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetClusterDefaultTagsSwitchStatus(client, clusterID)
		if err != nil {
			return nil, "ERROR", err
		}
		// status options: "failed", "succeed", "processing"
		status := utils.PathSearch("status", respBody, "").(string)
		if status == "failed" {
			return nil, "ERROR", fmt.Errorf("unexpected status: %s", status)
		}

		if status == "succeed" {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceClusterDefaultTagsSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	action := d.Get("action").(string)
	if err = switchClusterDefaultTags(client, clusterId, action); err != nil {
		return diag.Errorf("unable to %s default tags of cluster (%s): %s", action, clusterId, err)
	}

	if err = waitForClusterDefaultTagsStatus(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for default tags to be enabled of cluster (%s): %s", clusterId, err)
	}

	d.SetId(clusterId)

	return resourceClusterDefaultTagsSwitchRead(ctx, d, meta)
}

func resourceClusterDefaultTagsSwitchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)
	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	respBody, err := GetClusterDefaultTagsSwitchStatus(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", clusterNotFoundCodes...),
			"error querying cluster default tags status",
		)
	}

	enabled := utils.PathSearch("default_tags_enable", respBody, false).(bool)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("action", parseClusterDefaultTagsSwitchAction(enabled)),
		d.Set("default_tags_enable", enabled),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseClusterDefaultTagsSwitchAction(enabled bool) string {
	if enabled {
		return clusterDefaultTagsActionCreate
	}

	return clusterDefaultTagsActionDelete
}

func resourceClusterDefaultTagsSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterDefaultTagsSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	// Currently, repeating the operation to disable the default tags switch will not report an error, but it may
	// report an error in the future, so we need to check if it has already been disabled.
	if !d.Get("default_tags_enable").(bool) {
		return nil
	}

	if err = switchClusterDefaultTags(client, clusterId, "delete"); err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", clusterNotFoundCodes...),
			"error disabling default tags for cluster",
		)
	}

	if err = waitForClusterDefaultTagsStatus(ctx, client, clusterId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for default tags to be disabled of cluster (%s): %s", clusterId, err)
	}

	return nil
}
