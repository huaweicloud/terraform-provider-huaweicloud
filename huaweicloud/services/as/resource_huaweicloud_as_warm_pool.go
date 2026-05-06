package as

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var warmPoolNonUpdatableParams = []string{"scaling_group_id"}

// @API AS PUT /v2/{project_id}/scaling_groups/{scaling_group_id}/warm-pool
// @API AS GET /v2/{project_id}/scaling_groups/{scaling_group_id}/warm-pool
// @API AS DELETE /v2/{project_id}/scaling_groups/{scaling_group_id}/warm-pool
func ResourceAsWarmPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsWarmPoolCreateOrUpdate,
		UpdateContext: resourceAsWarmPoolCreateOrUpdate,
		ReadContext:   resourceAsWarmPoolRead,
		DeleteContext: resourceAsWarmPoolDelete,
		CustomizeDiff: config.FlexibleForceNew(warmPoolNonUpdatableParams),
		Importer: &schema.ResourceImporter{
			StateContext: resourceAsWarmPoolImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"min_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_init_wait_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAsWarmPoolCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool"
		product        = "autoscaling"
		scalingGroupID = d.Get("scaling_group_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{scaling_group_id}", scalingGroupID)
	putOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildCreateOrUpdateASWarmPoolBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &putOpt)
	if err != nil {
		return diag.Errorf("error creating or updating AS warm pool: %s", err)
	}

	d.SetId(d.Get("scaling_group_id").(string))
	return resourceAsWarmPoolRead(ctx, d, meta)
}

func buildCreateOrUpdateASWarmPoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"min_capacity":            utils.ValueIgnoreEmpty(d.Get("min_capacity")),
		"max_capacity":            utils.ValueIgnoreEmpty(d.Get("max_capacity")),
		"instance_init_wait_time": utils.ValueIgnoreEmpty(d.Get("instance_init_wait_time")),
	}
	return bodyParams
}

func resourceAsWarmPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "autoscaling"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	getRespBody, err := getWarmPool(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	status := utils.PathSearch("warm_pool.status", getRespBody, "").(string)
	if status == "CLOSED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("min_capacity", utils.PathSearch("warm_pool.min_capacity", getRespBody, nil)),
		d.Set("max_capacity", utils.PathSearch("warm_pool.max_capacity", getRespBody, nil)),
		d.Set("instance_init_wait_time", utils.PathSearch("warm_pool.instance_init_wait_time", getRespBody, nil)),
		d.Set("status", status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAsWarmPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		deleteUrl = "v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool"
		product   = "autoscaling"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	deletePath := client.Endpoint + deleteUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{scaling_group_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007"), "error deleting AS warm pool")
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      waitForWarmPoolDelete(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for AS warm pool (%s) to deleted: %s", d.Id(), err)
	}
	return nil
}

func waitForWarmPoolDelete(client *golangsdk.ServiceClient, serverId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getWarmPool(client, serverId)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("warm_pool.status", getRespBody, "").(string)
		if status == "CLOSED" {
			return getRespBody, "SUCCESS", nil
		}
		return getRespBody, "PENDING", nil
	}
}

func getWarmPool(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_group_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceAsWarmPoolImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {

	mErr := multierror.Append(
		nil,
		d.Set("scaling_group_id", d.Id()),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
