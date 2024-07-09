// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AS
// ---------------------------------------------------------------

package as

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_policy/{id}
// @API AS POST /autoscaling-api/v2/{project_id}/scaling_policy
// @API AS GET /autoscaling-api/v2/{project_id}/scaling_policy/{id}
// @API AS PUT /autoscaling-api/v2/{project_id}/scaling_policy/{id}
func ResourceASBandWidthPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASBandWidthPolicyCreate,
		UpdateContext: resourceASBandWidthPolicyUpdate,
		ReadContext:   resourceASBandWidthPolicyRead,
		DeleteContext: resourceASPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AS policy name.`,
			},
			"scaling_policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AS policy type.`,
			},
			"bandwidth_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the scaling bandwidth ID.`,
			},
			"alarm_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the alarm rule ID.`,
				ExactlyOneOf: []string{
					"scheduled_policy",
				},
			},
			"cool_down_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the cooldown period (in seconds).`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the AS policy.`,
			},
			"scaling_policy_action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     bandWidthPolicyActionSchema(),
				Optional: true,
				Computed: true,
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     bandWidthScheduledPolicySchema(),
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"alarm_id",
				},
			},
			"scaling_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `the scaling resource type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `the AS policy status.`,
			},
		},
	}
}

func bandWidthPolicyActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the operation to be performed. The default operation is ADD.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bandwidth (Mbit/s).`,
			},
			"limits": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the operation restrictions.`,
			},
		},
	}
	return &sc
}

func bandWidthScheduledPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"launch_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the time when the scaling action is triggered. The time format complies with UTC.`,
			},
			"recurrence_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the periodic triggering type.`,
			},
			"recurrence_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the day when a periodic scaling action is triggered.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the start time of the scaling action triggered periodically.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the end time of the scaling action triggered periodically.`,
			},
		},
	}
	return &sc
}

func resourceASBandWidthPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf                         = meta.(*config.Config)
		region                       = conf.GetRegion(d)
		createBandwidthPolicyHttpUrl = "autoscaling-api/v2/{project_id}/scaling_policy"
		createBandwidthPolicyProduct = "autoscaling"
	)

	client, err := conf.NewServiceClient(createBandwidthPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating AS bandwidth policy client: %s", err)
	}

	createBandwidthPolicyPath := client.Endpoint + createBandwidthPolicyHttpUrl
	createBandwidthPolicyPath = strings.ReplaceAll(createBandwidthPolicyPath, "{project_id}", client.ProjectID)
	createBandwidthPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateBandwidthPolicyBodyParams(d)),
	}

	createBandwidthPolicyResp, err := client.Request("POST", createBandwidthPolicyPath, &createBandwidthPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating AS bandwidth policy: %s", err)
	}

	createBandwidthPolicyRespBody, err := utils.FlattenResponse(createBandwidthPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("scaling_policy_id", createBandwidthPolicyRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating AS bandwidth policy: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceASBandWidthPolicyRead(ctx, d, meta)
}

func buildCreateBandwidthPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scaling_policy_name":   utils.ValueIgnoreEmpty(d.Get("scaling_policy_name")),
		"scaling_policy_type":   utils.ValueIgnoreEmpty(d.Get("scaling_policy_type")),
		"scaling_resource_id":   utils.ValueIgnoreEmpty(d.Get("bandwidth_id")),
		"scaling_resource_type": "BANDWIDTH",
		"alarm_id":              utils.ValueIgnoreEmpty(d.Get("alarm_id")),
		"cool_down_time":        utils.ValueIgnoreEmpty(d.Get("cool_down_time")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"scaling_policy_action": buildCreateBandwidthPolicyScalingPolicyActionChildBody(d),
		"scheduled_policy":      buildCreateBandwidthPolicyScheduledPolicyChildBody(d),
	}
	return bodyParams
}

func buildCreateBandwidthPolicyScalingPolicyActionChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("scaling_policy_action").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"operation": utils.ValueIgnoreEmpty(raw["operation"]),
		"size":      utils.ValueIgnoreEmpty(raw["size"]),
		"limits":    utils.ValueIgnoreEmpty(raw["limits"]),
	}

	return params
}

func buildCreateBandwidthPolicyScheduledPolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("scheduled_policy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"launch_time":      utils.ValueIgnoreEmpty(raw["launch_time"]),
		"recurrence_type":  utils.ValueIgnoreEmpty(raw["recurrence_type"]),
		"recurrence_value": utils.ValueIgnoreEmpty(raw["recurrence_value"]),
		"start_time":       utils.ValueIgnoreEmpty(raw["start_time"]),
		"end_time":         utils.ValueIgnoreEmpty(raw["end_time"]),
	}

	return params
}

func resourceASBandWidthPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf                      = meta.(*config.Config)
		region                    = conf.GetRegion(d)
		mErr                      *multierror.Error
		getBandwidthPolicyHttpUrl = "autoscaling-api/v2/{project_id}/scaling_policy/{id}"
		getBandwidthPolicyProduct = "autoscaling"
	)

	client, err := conf.NewServiceClient(getBandwidthPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating AS bandwidth policy client: %s", err)
	}

	getBandwidthPolicyPath := client.Endpoint + getBandwidthPolicyHttpUrl
	getBandwidthPolicyPath = strings.ReplaceAll(getBandwidthPolicyPath, "{project_id}", client.ProjectID)
	getBandwidthPolicyPath = strings.ReplaceAll(getBandwidthPolicyPath, "{id}", d.Id())
	getBandwidthPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBandwidthPolicyResp, err := client.Request("GET", getBandwidthPolicyPath, &getBandwidthPolicyOpt)
	if err != nil {
		// When the resource does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving AS bandwidth policy")
	}

	respBody, err := utils.FlattenResponse(getBandwidthPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("scaling_policy_name", utils.PathSearch("scaling_policy.scaling_policy_name", respBody, nil)),
		d.Set("scaling_policy_type", utils.PathSearch("scaling_policy.scaling_policy_type", respBody, nil)),
		d.Set("bandwidth_id", utils.PathSearch("scaling_policy.scaling_resource_id", respBody, nil)),
		d.Set("scaling_resource_type", utils.PathSearch("scaling_policy.scaling_resource_type", respBody, nil)),
		d.Set("alarm_id", utils.PathSearch("scaling_policy.alarm_id", respBody, nil)),
		d.Set("cool_down_time", utils.PathSearch("scaling_policy.cool_down_time", respBody, nil)),
		d.Set("description", utils.PathSearch("scaling_policy.description", respBody, nil)),
		d.Set("status", utils.PathSearch("scaling_policy.policy_status", respBody, nil)),
		d.Set("scaling_policy_action", flattenGetBandwidthPolicyResponseBodyScalingPolicyAction(respBody)),
		d.Set("scheduled_policy", flattenGetBandwidthPolicyResponseBodyScheduledPolicy(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetBandwidthPolicyResponseBodyScalingPolicyAction(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("scaling_policy.scaling_policy_action", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing scaling_policy_action from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"operation": utils.PathSearch("operation", curJson, nil),
			"size":      utils.PathSearch("size", curJson, nil),
			"limits":    utils.PathSearch("limits", curJson, nil),
		},
	}
	return rst
}

func flattenGetBandwidthPolicyResponseBodyScheduledPolicy(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("scaling_policy.scheduled_policy", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing scheduled_policy from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"launch_time":      utils.PathSearch("launch_time", curJson, nil),
			"recurrence_type":  utils.PathSearch("recurrence_type", curJson, nil),
			"recurrence_value": utils.PathSearch("recurrence_value", curJson, nil),
			"start_time":       utils.PathSearch("start_time", curJson, nil),
			"end_time":         utils.PathSearch("end_time", curJson, nil),
		},
	}
	return rst
}

func resourceASBandWidthPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	updateBandwidthPolicyChanges := []string{
		"scaling_policy_name",
		"scaling_policy_type",
		"bandwidth_id",
		"scaling_resource_type",
		"alarm_id",
		"cool_down_time",
		"description",
		"scaling_policy_action",
		"scheduled_policy",
	}

	if d.HasChanges(updateBandwidthPolicyChanges...) {
		// updateBandwidthPolicy: update the AS bandwidth scaling policy
		var (
			conf                         = meta.(*config.Config)
			region                       = conf.GetRegion(d)
			updateBandwidthPolicyHttpUrl = "autoscaling-api/v2/{project_id}/scaling_policy/{id}"
			updateBandwidthPolicyProduct = "autoscaling"
		)
		client, err := conf.NewServiceClient(updateBandwidthPolicyProduct, region)
		if err != nil {
			return diag.Errorf("error creating AS bandwidth policy client: %s", err)
		}

		updateBandwidthPolicyPath := client.Endpoint + updateBandwidthPolicyHttpUrl
		updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{project_id}", client.ProjectID)
		updateBandwidthPolicyPath = strings.ReplaceAll(updateBandwidthPolicyPath, "{id}", d.Id())
		updateBandwidthPolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateBandwidthPolicyBodyParams(d)),
		}

		_, err = client.Request("PUT", updateBandwidthPolicyPath, &updateBandwidthPolicyOpt)
		if err != nil {
			return diag.Errorf("error updating AS bandwidth policy: %s", err)
		}
	}
	return resourceASBandWidthPolicyRead(ctx, d, meta)
}

func buildUpdateBandwidthPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scaling_policy_name":   utils.ValueIgnoreEmpty(d.Get("scaling_policy_name")),
		"scaling_policy_type":   utils.ValueIgnoreEmpty(d.Get("scaling_policy_type")),
		"scaling_resource_id":   utils.ValueIgnoreEmpty(d.Get("bandwidth_id")),
		"scaling_resource_type": utils.ValueIgnoreEmpty(d.Get("scaling_resource_type")),
		"alarm_id":              utils.ValueIgnoreEmpty(d.Get("alarm_id")),
		"cool_down_time":        utils.ValueIgnoreEmpty(d.Get("cool_down_time")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"scaling_policy_action": buildUpdateBandwidthPolicyScalingPolicyActionChildBody(d),
		"scheduled_policy":      buildUpdateBandwidthPolicyScheduledPolicyChildBody(d),
	}
	return bodyParams
}

func buildUpdateBandwidthPolicyScalingPolicyActionChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("scaling_policy_action").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"operation": utils.ValueIgnoreEmpty(raw["operation"]),
		"size":      utils.ValueIgnoreEmpty(raw["size"]),
		"limits":    utils.ValueIgnoreEmpty(raw["limits"]),
	}

	return params
}

func buildUpdateBandwidthPolicyScheduledPolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("scheduled_policy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"launch_time":      utils.ValueIgnoreEmpty(raw["launch_time"]),
		"recurrence_type":  utils.ValueIgnoreEmpty(raw["recurrence_type"]),
		"recurrence_value": utils.ValueIgnoreEmpty(raw["recurrence_value"]),
		"start_time":       utils.ValueIgnoreEmpty(raw["start_time"]),
		"end_time":         utils.ValueIgnoreEmpty(raw["end_time"]),
	}

	return params
}
