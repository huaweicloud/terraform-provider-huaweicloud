package as

import (
	"context"
	"fmt"
	"strconv"
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

// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_policy/{id}
// @API AS POST /autoscaling-api/v2/{project_id}/scaling_policy
// @API AS GET /autoscaling-api/v2/{project_id}/scaling_policy/{id}
// @API AS PUT /autoscaling-api/v2/{project_id}/scaling_policy/{id}
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_policy/{scaling_policy_id}/action
func ResourceASBandWidthPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASBandWidthPolicyCreate,
		UpdateContext: resourceASBandWidthPolicyUpdate,
		ReadContext:   resourceASBandWidthPolicyRead,
		DeleteContext: resourceASPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
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
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies identification of operation the AS bandwidth policy.`,
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
			"interval_alarm_actions": {
				Type:        schema.TypeSet,
				Elem:        bandWidthIntervalAlarmActions(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the alarm interval of the bandwidth policy.`,
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
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the bandwidth policy.`,
			},
			"meta_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The bandwidth policy additional information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_bandwidth_share_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bandwidth sharing type in the bandwidth policy.`,
						},
						"metadata_eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The EIP ID for the bandwidth in the bandwidth policy.`,
						},
						"metadata_eip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The EIP IP address for the bandwidth in the bandwidth policy.`,
						},
					},
				},
			},
		},
	}
}

func bandWidthIntervalAlarmActions() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"lower_bound": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the lower limit of the value range.`,
			},
			"upper_bound": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the upper limit of the value range.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the operation to be performed.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the operation size.`,
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

func buildIntervalAlarmActionsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	intervalActions := d.Get("interval_alarm_actions").(*schema.Set).List()
	if len(intervalActions) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(intervalActions))
	for _, v := range intervalActions {
		raw := v.(map[string]interface{})
		params := map[string]interface{}{
			"operation":   utils.ValueIgnoreEmpty(raw["operation"]),
			"size":        utils.ValueIgnoreEmpty(raw["size"]),
			"limits":      utils.ValueIgnoreEmpty(raw["limits"]),
			"lower_bound": bulidIntervalThreshold(raw["lower_bound"].(string)),
			"upper_bound": bulidIntervalThreshold(raw["upper_bound"].(string)),
		}

		rst = append(rst, params)
	}

	return rst
}

func bulidIntervalThreshold(str string) interface{} {
	resp, err := strconv.Atoi(str)
	if err != nil {
		return nil
	}

	return resp
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

	policyId := utils.PathSearch("scaling_policy_id", createBandwidthPolicyRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the AS bandwidth policy ID from the API response")
	}
	d.SetId(policyId)

	action := d.Get("action").(string)
	if action == "pause" {
		err := updateBandwidthPolicyStatus(ctx, client, d.Timeout(schema.TimeoutCreate), d.Id(), action)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceASBandWidthPolicyRead(ctx, d, meta)
}

func buildCreateBandwidthPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scaling_policy_name":    utils.ValueIgnoreEmpty(d.Get("scaling_policy_name")),
		"scaling_policy_type":    utils.ValueIgnoreEmpty(d.Get("scaling_policy_type")),
		"scaling_resource_id":    utils.ValueIgnoreEmpty(d.Get("bandwidth_id")),
		"scaling_resource_type":  "BANDWIDTH",
		"alarm_id":               utils.ValueIgnoreEmpty(d.Get("alarm_id")),
		"cool_down_time":         utils.ValueIgnoreEmpty(d.Get("cool_down_time")),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"scaling_policy_action":  buildCreateBandwidthPolicyScalingPolicyActionChildBody(d),
		"scheduled_policy":       buildCreateBandwidthPolicyScheduledPolicyChildBody(d),
		"interval_alarm_actions": buildIntervalAlarmActionsBodyParams(d),
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
		getBandwidthPolicyProduct = "autoscaling"
	)

	client, err := conf.NewServiceClient(getBandwidthPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating AS bandwidth policy client: %s", err)
	}

	respBody, err := GetBandwidthPolicy(client, d.Id())
	if err != nil {
		// When the resource does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving AS bandwidth policy")
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
		d.Set("interval_alarm_actions", flattenIntervalActions(utils.PathSearch("scaling_policy.interval_alarm_actions", respBody, nil))),
		d.Set("create_time", utils.PathSearch("scaling_policy.create_time", respBody, nil)),
		d.Set("meta_data", flattenBandwidthMetaData(utils.PathSearch("scaling_policy.meta_data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIntervalActions(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	rawArray := resp.([]interface{})
	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		params := map[string]interface{}{
			"lower_bound": flattenIntervalThreshold(utils.PathSearch("lower_bound", v, nil)),
			"upper_bound": flattenIntervalThreshold(utils.PathSearch("upper_bound", v, nil)),
			"operation":   utils.PathSearch("operation", v, nil),
			"size":        utils.PathSearch("size", v, nil),
			"limits":      utils.PathSearch("limits", v, nil),
		}
		rst[i] = params
	}

	return rst
}

func flattenIntervalThreshold(param interface{}) string {
	if param == nil {
		return ""
	}

	return fmt.Sprintf("%v", param)
}

func flattenBandwidthMetaData(metaData interface{}) []map[string]interface{} {
	if metaData == nil {
		return nil
	}

	bandwidthMetaData := map[string]interface{}{
		"metadata_bandwidth_share_type": utils.PathSearch("metadata_bandwidth_share_type", metaData, nil),
		"metadata_eip_id":               utils.PathSearch("metadata_eip_id", metaData, nil),
		"metadata_eip_address":          utils.PathSearch("metadata_eip_address", metaData, nil),
	}

	return []map[string]interface{}{bandwidthMetaData}
}

func flattenGetBandwidthPolicyResponseBodyScalingPolicyAction(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("scaling_policy.scaling_policy_action", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
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
	curJson := utils.PathSearch("scaling_policy.scheduled_policy", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
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
	var (
		cfg                          = meta.(*config.Config)
		region                       = cfg.GetRegion(d)
		updateBandwidthPolicyProduct = "autoscaling"
	)

	client, err := cfg.NewServiceClient(updateBandwidthPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating AS bandwidth policy client: %s", err)
	}

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
		"interval_alarm_actions",
	}

	if d.HasChanges(updateBandwidthPolicyChanges...) {
		// updateBandwidthPolicy: update the AS bandwidth scaling policy
		updateBandwidthPolicyHttpUrl := "autoscaling-api/v2/{project_id}/scaling_policy/{id}"
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

	if d.HasChange("action") {
		action := d.Get("action").(string)
		err := updateBandwidthPolicyStatus(ctx, client, d.Timeout(schema.TimeoutUpdate), d.Id(), action)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceASBandWidthPolicyRead(ctx, d, meta)
}

func buildUpdateBandwidthPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scaling_policy_name":    utils.ValueIgnoreEmpty(d.Get("scaling_policy_name")),
		"scaling_policy_type":    utils.ValueIgnoreEmpty(d.Get("scaling_policy_type")),
		"scaling_resource_id":    utils.ValueIgnoreEmpty(d.Get("bandwidth_id")),
		"scaling_resource_type":  utils.ValueIgnoreEmpty(d.Get("scaling_resource_type")),
		"alarm_id":               utils.ValueIgnoreEmpty(d.Get("alarm_id")),
		"cool_down_time":         utils.ValueIgnoreEmpty(d.Get("cool_down_time")),
		"description":            utils.ValueIgnoreEmpty(d.Get("description")),
		"scaling_policy_action":  buildUpdateBandwidthPolicyScalingPolicyActionChildBody(d),
		"scheduled_policy":       buildUpdateBandwidthPolicyScheduledPolicyChildBody(d),
		"interval_alarm_actions": buildIntervalAlarmActionsBodyParams(d),
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

func GetBandwidthPolicy(client *golangsdk.ServiceClient, policyId string) (interface{}, error) {
	httpUrl := "autoscaling-api/v2/{project_id}/scaling_policy/{scaling_policy_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_policy_id}", policyId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func updateBandwidthPolicyStatus(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration,
	policyId, action string) error {
	updateStatusHttpUrl := "autoscaling-api/v1/{project_id}/scaling_policy/{scaling_policy_id}/action"
	updateStatusPath := client.Endpoint + updateStatusHttpUrl
	updateStatusPath = strings.ReplaceAll(updateStatusPath, "{project_id}", client.ProjectID)
	updateStatusPath = strings.ReplaceAll(updateStatusPath, "{scaling_policy_id}", policyId)

	updatePolicyStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": action,
		},
		OkCodes: []int{
			200, 201, 204,
		},
	}

	_, err := client.Request("POST", updateStatusPath, &updatePolicyStatusOpt)
	if err != nil {
		return fmt.Errorf("error updating AS bandwidth policy status: %s", err)
	}

	err = waitingForBandWidthPolicyStatusCompleted(ctx, client, t, policyId, action)
	if err != nil {
		return fmt.Errorf("error waiting for the AS bandwidth policy status update to complete: %s", err)
	}

	return nil
}

func waitingForBandWidthPolicyStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration,
	policyId, action string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitBandwidthPolicyStatusRefreshFunc(client, policyId, action),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitBandwidthPolicyStatusRefreshFunc(client *golangsdk.ServiceClient, policyId, action string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetBandwidthPolicy(client, policyId)
		if err != nil {
			return nil, "ERROR", err
		}

		policyStatus := utils.PathSearch("scaling_policy.policy_status", respBody, "").(string)

		if action == "resume" {
			if policyStatus == "INSERVICE" {
				return respBody, "COMPLETED", nil
			}
		}

		if action == "pause" {
			if policyStatus == "PAUSED" {
				return respBody, "COMPLETED", nil
			}
		}

		return respBody, "PENDING", nil
	}
}
