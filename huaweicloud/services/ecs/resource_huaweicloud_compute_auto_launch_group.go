package ecs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS POST /v2/{domain_id}/auto-launch-groups
// @API ECS GET /v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}
// @API ECS PUT /v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}
// @API ECS DELETE /v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}
func ResourceComputeAutoLaunchGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutoLaunchGroupCreate,
		ReadContext:   resourceAutoLaunchGroupRead,
		UpdateContext: resourceAutoLaunchGroupUpdate,
		DeleteContext: resourceAutoLaunchGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"overrides": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"spot_price": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"weighted_capacity": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return hashcode.String(m["availability_zone"].(string) + m["flavor_id"].(string))
				},
			},
			"delete_instances": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stable_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"excess_fulfilled_capacity_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instances_behavior_with_expiration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"valid_since": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"valid_until": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"allocation_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"supply_option": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spot_price": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_capacity": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"current_stable_capacity": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutoLaunchGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cms"
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CMS client: %s", err)
	}

	createAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups"
	createAutoLaunchGroupPath := client.Endpoint + createAutoLaunchGroupHttpUrl
	createAutoLaunchGroupPath = strings.ReplaceAll(createAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)

	createAutoLaunchGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAutoLaunchGroupBodyParams(d, cfg.GetRegion(d))),
	}
	createAutoLaunchGroupResp, err := client.Request("POST", createAutoLaunchGroupPath, &createAutoLaunchGroupOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	createAutoLaunchGroupRespBody, err := utils.FlattenResponse(createAutoLaunchGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("auto_launch_group_id", createAutoLaunchGroupRespBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable to find the auto launch group ID from the API response")
	}
	d.SetId(groupId)

	return resourceAutoLaunchGroupRead(ctx, d, meta)
}

func buildCreateAutoLaunchGroupBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                               d.Get("name"),
		"target_capacity":                    d.Get("target_capacity"),
		"stable_capacity":                    utils.ValueIgnoreEmpty(d.Get("stable_capacity")),
		"excess_fulfilled_capacity_behavior": utils.ValueIgnoreEmpty(d.Get("excess_fulfilled_capacity_behavior")),
		"type":                               utils.ValueIgnoreEmpty(d.Get("type")),
		"instances_behavior_with_expiration": utils.ValueIgnoreEmpty(d.Get("instances_behavior_with_expiration")),
		"valid_since":                        utils.ValueIgnoreEmpty(d.Get("valid_since")),
		"valid_until":                        utils.ValueIgnoreEmpty(d.Get("valid_until")),
		"allocation_strategy":                utils.ValueIgnoreEmpty(d.Get("allocation_strategy")),
		"supply_option":                      utils.ValueIgnoreEmpty(d.Get("supply_option")),
		"spot_price":                         utils.ValueIgnoreEmpty(d.Get("spot_price")),
		"region_specs":                       buildAutoLaunchGroupRequestBodyRegionSpecs(d, region),
	}
	return bodyParams
}

func buildAutoLaunchGroupRequestBodyRegionSpecs(d *schema.ResourceData, region string) []map[string]interface{} {
	regionSpecs := make([]map[string]interface{}, 1)
	params := map[string]interface{}{
		"region_id":              region,
		"expect_target_capacity": d.Get("target_capacity"),
		"expect_stable_capacity": d.Get("stable_capacity"),
		"launch_template_config": map[string]interface{}{
			"launch_template": map[string]interface{}{
				"launch_template_id": d.Get("launch_template_id"),
				"version":            d.Get("launch_template_version"),
			},
			"overrides": buildAutoLaunchGroupRequestBodyRegionSpecsOverrides(d.Get("overrides").(*schema.Set)),
		},
	}
	regionSpecs[0] = params
	return regionSpecs
}

func buildAutoLaunchGroupRequestBodyRegionSpecsOverrides(rawParams *schema.Set) []map[string]interface{} {
	if rawParams.Len() == 0 {
		return nil
	}
	overrides := make([]map[string]interface{}, rawParams.Len())
	for i, val := range rawParams.List() {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"availability_zone_id": raw["availability_zone"],
			"flavor_id":            raw["flavor_id"],
			"spot_price":           utils.ValueIgnoreEmpty(raw["spot_price"]),
			"priority":             utils.ValueIgnoreEmpty(raw["priority"]),
			"weighted_capacity":    utils.ValueIgnoreEmpty(raw["weighted_capacity"]),
		}
		overrides[i] = params
	}
	return overrides
}

func resourceAutoLaunchGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cms"
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CMS client: %s", err)
	}

	getAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}"
	getAutoLaunchGroupPath := client.Endpoint + getAutoLaunchGroupHttpUrl
	getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)
	getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{auto_launch_group_id}", d.Id())
	getAutoLaunchGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAutoLaunchGroupResp, err := client.Request("GET", getAutoLaunchGroupPath, &getAutoLaunchGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving auto launch group")
	}

	getAutoLaunchGroupRespBody, err := utils.FlattenResponse(getAutoLaunchGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// `supply_option` does not in return.
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("name", getAutoLaunchGroupRespBody, nil)),
		d.Set("target_capacity", utils.PathSearch("target_capacity", getAutoLaunchGroupRespBody, 0)),
		d.Set("stable_capacity", utils.PathSearch("stable_capacity", getAutoLaunchGroupRespBody, 0)),
		d.Set("excess_fulfilled_capacity_behavior", utils.PathSearch(
			"excess_fulfilled_capacity_behavior", getAutoLaunchGroupRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getAutoLaunchGroupRespBody, nil)),
		d.Set("instances_behavior_with_expiration", utils.PathSearch(
			"instances_behavior_with_expiration", getAutoLaunchGroupRespBody, nil)),
		d.Set("valid_since", utils.PathSearch("valid_since", getAutoLaunchGroupRespBody, nil)),
		d.Set("valid_until", utils.PathSearch("valid_until", getAutoLaunchGroupRespBody, nil)),
		d.Set("allocation_strategy", utils.PathSearch("allocation_strategy", getAutoLaunchGroupRespBody, nil)),
		d.Set("spot_price", utils.PathSearch("spot_price", getAutoLaunchGroupRespBody, float64(0))),
		d.Set("created_at", utils.PathSearch("created_at", getAutoLaunchGroupRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getAutoLaunchGroupRespBody, nil)),
		d.Set("current_capacity", utils.PathSearch("current_capacity", getAutoLaunchGroupRespBody, float64(0))),
		d.Set("current_stable_capacity", utils.PathSearch("current_stable_capacity", getAutoLaunchGroupRespBody, float64(0))),
		d.Set("task_state", utils.PathSearch("task_state", getAutoLaunchGroupRespBody, nil)),
		d.Set("launch_template_id", utils.PathSearch(
			"region_specs|[0].launch_template_config.launch_template.launch_template_id", getAutoLaunchGroupRespBody, nil)),
		d.Set("launch_template_version", utils.PathSearch(
			"region_specs|[0].launch_template_config.launch_template.version", getAutoLaunchGroupRespBody, nil)),
		d.Set("overrides", flattenRegionSpecsOverrides(utils.PathSearch(
			"region_specs|[0].launch_template_config.overrides", getAutoLaunchGroupRespBody, make([]interface{}, 0)))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting auto launch group fields: %s", err)
	}

	return nil
}

func flattenRegionSpecsOverrides(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"availability_zone": utils.PathSearch("availability_zone_id", val, nil),
			"flavor_id":         utils.PathSearch("flavor_id", val, nil),
			"spot_price":        utils.PathSearch("spot_price", val, float64(0)),
			"priority":          utils.PathSearch("priority", val, 0),
			"weighted_capacity": utils.PathSearch("weighted_capacity", val, float64(0)),
		}
		rst[i] = params
	}
	return rst
}

func resourceAutoLaunchGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cms"
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CMS client: %s", err)
	}

	updateChanges := []string{
		"name",
		"target_capacity",
		"stable_capacity",
		"excess_fulfilled_capacity_behavior",
		"instances_behavior_with_expiration",
		"spot_price",
	}

	if d.HasChanges(updateChanges...) {
		updateAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}"
		updateAutoLaunchGroupPath := client.Endpoint + updateAutoLaunchGroupHttpUrl
		updateAutoLaunchGroupPath = strings.ReplaceAll(updateAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)
		updateAutoLaunchGroupPath = strings.ReplaceAll(updateAutoLaunchGroupPath, "{auto_launch_group_id}", d.Id())
		updateAutoLaunchGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateAutoLaunchGroupBodyParams(d)),
		}

		_, err = client.Request("PUT", updateAutoLaunchGroupPath, &updateAutoLaunchGroupOpt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If only update `delete_instances`, just return READ.
	return resourceAutoLaunchGroupRead(ctx, d, meta)
}

func buildUpdateAutoLaunchGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                               d.Get("name"),
		"target_capacity":                    utils.ValueIgnoreEmpty(d.Get("target_capacity")),
		"stable_capacity":                    utils.ValueIgnoreEmpty(d.Get("stable_capacity")),
		"excess_fulfilled_capacity_behavior": utils.ValueIgnoreEmpty(d.Get("excess_fulfilled_capacity_behavior")),
		"instances_behavior_with_expiration": utils.ValueIgnoreEmpty(d.Get("instances_behavior_with_expiration")),
		"spot_price":                         utils.ValueIgnoreEmpty(d.Get("spot_price")),
	}
	return bodyParams
}

func resourceAutoLaunchGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cms"
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CMS client: %s", err)
	}

	deleteAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}"
	deleteAutoLaunchGroupPath := client.Endpoint + deleteAutoLaunchGroupHttpUrl
	deleteAutoLaunchGroupPath = strings.ReplaceAll(deleteAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)
	deleteAutoLaunchGroupPath = strings.ReplaceAll(deleteAutoLaunchGroupPath, "{auto_launch_group_id}", d.Id())

	deleteAutoLaunchGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"delete_instances": utils.ValueIgnoreEmpty(d.Get("delete_instances")),
		}),
	}

	_, err = client.Request("DELETE", deleteAutoLaunchGroupPath, &deleteAutoLaunchGroupOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING"},
		Target:       []string{"DELETED"},
		Refresh:      autoLaunchGroupStatusRefreshFunc(d.Id(), cfg, client),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting auto launch groups: %s", err)
	}

	return nil
}

func autoLaunchGroupStatusRefreshFunc(id string, cfg *config.Config,
	client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}"
		getAutoLaunchGroupPath := client.Endpoint + getAutoLaunchGroupHttpUrl
		getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)
		getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{auto_launch_group_id}", id)
		getPrivateCAOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getPrivateCAResp, err := client.Request("GET", getAutoLaunchGroupPath, &getPrivateCAOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}
		getPrivateCARespBody, err := utils.FlattenResponse(getPrivateCAResp)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status", getPrivateCARespBody, "")
		return getPrivateCARespBody, status.(string), nil
	}
}
