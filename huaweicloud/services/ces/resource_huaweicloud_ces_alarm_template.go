// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CES
// ---------------------------------------------------------------

package ces

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var alarmTemplateNonUpdatableParams = []string{"type", "is_overwrite"}

// @API CES POST /v2/{project_id}/alarm-templates
// @API CES GET /v2/{project_id}/alarm-templates/{template_id}
// @API CES PUT /v2/{project_id}/alarm-templates/{template_id}
// @API CES POST /v2/{project_id}/alarm-templates/batch-delete
func ResourceCesAlarmTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCesAlarmTemplateCreate,
		UpdateContext: resourceCesAlarmTemplateUpdate,
		ReadContext:   resourceCesAlarmTemplateRead,
		DeleteContext: resourceCesAlarmTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(alarmTemplateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CES alarm template.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Elem:        AlarmTemplatePolicySchema(),
				Required:    true,
				Description: `Specifies the policy list of the CES alarm template.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the type of the CES alarm template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the CES alarm template.`,
			},
			"is_overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to overwrite an existing alarm template with the same template name.`,
			},
			"delete_associate_alarm": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether delete the alarm rule which the alarm template associated with.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"association_alarm_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total num of the alarm that associated with the alarm template.`,
			},
		},
	}
}

func AlarmTemplatePolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace of the service.`,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm metric name.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the judgment period of alarm condition.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data rollup methods.`,
			},
			"comparison_operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the comparison conditions for alarm threshold.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of consecutive triggering of alarms.`,
			},
			"suppress_duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the alarm suppression cycle.`,
			},
			"value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the alarm threshold.`,
			},
			"hierarchical_value": {
				Type:        schema.TypeList,
				Elem:        AlarmTemplatePolicyHierarchicalValueSchema(),
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the multiple levels of alarm thresholds.`,
			},
			"alarm_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the alarm level.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit string of the alarm threshold.`,
			},
			"dimension_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource dimension.`,
			},
		},
	}
	return &sc
}

func AlarmTemplatePolicyHierarchicalValueSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"critical": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the threshold for the critical level.`,
			},
			"major": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the threshold for the major level.`,
			},
			"minor": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the threshold for the minor level.`,
			},
			"info": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the threshold for the info level.`,
			},
		},
	}
	return &sc
}

func resourceCesAlarmTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAlarmTemplate: create CES alarm template
	var (
		createAlarmTemplateHttpUrl = "v2/{project_id}/alarm-templates"
		createAlarmTemplateProduct = "ces"
	)
	createAlarmTemplateClient, err := cfg.NewServiceClient(createAlarmTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createAlarmTemplatePath := createAlarmTemplateClient.Endpoint + createAlarmTemplateHttpUrl
	createAlarmTemplatePath = strings.ReplaceAll(createAlarmTemplatePath, "{project_id}",
		createAlarmTemplateClient.ProjectID)

	createAlarmTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAlarmTemplateOpt.JSONBody = utils.RemoveNil(buildAlarmTemplateCreateBodyParams(d))
	createAlarmTemplateResp, err := createAlarmTemplateClient.Request("POST",
		createAlarmTemplatePath, &createAlarmTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating CES alarm template: %s", err)
	}

	createAlarmTemplateRespBody, err := utils.FlattenResponse(createAlarmTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("template_id", createAlarmTemplateRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES alarm template: ID is not found in API response")
	}
	d.SetId(id)

	return resourceCesAlarmTemplateRead(ctx, d, meta)
}

func resourceCesAlarmTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAlarmTemplateHasChanges := []string{
		"name",
		"description",
		"policies",
	}

	if d.HasChanges(updateAlarmTemplateHasChanges...) {
		// updateAlarmTemplate: update CES alarm template
		var (
			updateAlarmTemplateHttpUrl = "v2/{project_id}/alarm-templates/{template_id}"
			updateAlarmTemplateProduct = "ces"
		)
		updateAlarmTemplateClient, err := cfg.NewServiceClient(updateAlarmTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating CES client: %s", err)
		}

		updateAlarmTemplatePath := updateAlarmTemplateClient.Endpoint + updateAlarmTemplateHttpUrl
		updateAlarmTemplatePath = strings.ReplaceAll(updateAlarmTemplatePath, "{project_id}",
			updateAlarmTemplateClient.ProjectID)
		updateAlarmTemplatePath = strings.ReplaceAll(updateAlarmTemplatePath, "{template_id}", d.Id())

		updateAlarmTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateAlarmTemplateOpt.JSONBody = utils.RemoveNil(buildAlarmTemplateUpdateBodyParams(d))
		_, err = updateAlarmTemplateClient.Request("PUT", updateAlarmTemplatePath, &updateAlarmTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating CES alarm template: %s", err)
		}
	}
	return resourceCesAlarmTemplateRead(ctx, d, meta)
}

func buildAlarmTemplateCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_name":        d.Get("name"),
		"template_type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"template_description": utils.ValueIgnoreEmpty(d.Get("description")),
		"is_overwrite":         utils.ValueIgnoreEmpty(d.Get("is_overwrite")),
		"policies":             buildAlarmTemplatePoliciesChildBody(d),
	}
	return bodyParams
}

func buildAlarmTemplateUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_name":        d.Get("name"),
		"template_type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"template_description": utils.ValueIgnoreEmpty(d.Get("description")),
		"policies":             buildAlarmTemplatePoliciesChildBody(d),
	}
	return bodyParams
}

func buildAlarmTemplatePoliciesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("policies").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(rawParams))
	for _, rawParam := range rawParams {
		raw := rawParam.(map[string]interface{})
		param := map[string]interface{}{
			"namespace":           utils.ValueIgnoreEmpty(raw["namespace"]),
			"dimension_name":      utils.ValueIgnoreEmpty(raw["dimension_name"]),
			"metric_name":         utils.ValueIgnoreEmpty(raw["metric_name"]),
			"period":              raw["period"],
			"filter":              utils.ValueIgnoreEmpty(raw["filter"]),
			"comparison_operator": utils.ValueIgnoreEmpty(raw["comparison_operator"]),
			"value":               raw["value"],
			"hierarchical_value":  buildAlarmTemplatePoliciesHierarchicalValueChildBody(raw["hierarchical_value"]),
			"unit":                utils.ValueIgnoreEmpty(raw["unit"]),
			"count":               utils.ValueIgnoreEmpty(raw["count"]),
			"alarm_level":         utils.ValueIgnoreEmpty(raw["alarm_level"]),
			"suppress_duration":   raw["suppress_duration"],
		}
		params = append(params, param)
	}
	return params
}

func buildAlarmTemplatePoliciesHierarchicalValueChildBody(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"critical": utils.ValueIgnoreEmpty(raw["critical"]),
			"major":    utils.ValueIgnoreEmpty(raw["major"]),
			"minor":    utils.ValueIgnoreEmpty(raw["minor"]),
			"info":     utils.ValueIgnoreEmpty(raw["info"]),
		}

		return param
	}

	return nil
}

func resourceCesAlarmTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAlarmTemplate: Query CES alarm template
	var (
		getAlarmTemplateHttpUrl = "v2/{project_id}/alarm-templates/{template_id}"
		getAlarmTemplateProduct = "ces"
	)
	getAlarmTemplateClient, err := cfg.NewServiceClient(getAlarmTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	getAlarmTemplatePath := getAlarmTemplateClient.Endpoint + getAlarmTemplateHttpUrl
	getAlarmTemplatePath = strings.ReplaceAll(getAlarmTemplatePath, "{project_id}",
		getAlarmTemplateClient.ProjectID)
	getAlarmTemplatePath = strings.ReplaceAll(getAlarmTemplatePath, "{template_id}", d.Id())

	getAlarmTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAlarmTemplateResp, err := getAlarmTemplateClient.Request("GET", getAlarmTemplatePath, &getAlarmTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm template")
	}

	getAlarmTemplateRespBody, err := utils.FlattenResponse(getAlarmTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rawTemplateType := utils.PathSearch("template_type", getAlarmTemplateRespBody, "")
	templateType := -1
	if rawTemplateType == "custom" {
		templateType = 0
	} else if rawTemplateType == "custom_event" {
		templateType = 2
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("template_name", getAlarmTemplateRespBody, nil)),
		d.Set("type", templateType),
		d.Set("description", utils.PathSearch("template_description",
			getAlarmTemplateRespBody, nil)),
		d.Set("association_alarm_total", utils.PathSearch("association_alarm_total",
			getAlarmTemplateRespBody, nil)),
		d.Set("policies", flattenGetAlarmTemplateResponseBodyPolicy(d, getAlarmTemplateRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAlarmTemplateResponseBodyPolicy(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for i, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"namespace":           utils.PathSearch("namespace", v, nil),
			"dimension_name":      utils.PathSearch("dimension_name", v, nil),
			"metric_name":         utils.PathSearch("metric_name", v, nil),
			"period":              utils.PathSearch("period", v, nil),
			"filter":              utils.PathSearch("filter", v, nil),
			"comparison_operator": utils.PathSearch("comparison_operator", v, nil),
			"value":               utils.PathSearch("value", v, nil),
			"hierarchical_value": flattenGetAlarmTemplateResponseBodyPolicyHierarchicalValue(d, i,
				utils.PathSearch("hierarchical_value", v, nil)),
			"unit":              utils.PathSearch("unit", v, nil),
			"count":             utils.PathSearch("count", v, nil),
			"alarm_level":       utils.PathSearch("alarm_level", v, nil),
			"suppress_duration": utils.PathSearch("suppress_duration", v, nil),
		})
	}
	return rst
}

func flattenGetAlarmTemplateResponseBodyPolicyHierarchicalValue(d *schema.ResourceData, i int, param interface{}) interface{} {
	// The hierarchical_value has a default value and takes precedence over alarm_level and value. When updated, it
	// causes alarm_value to be unable to update.
	hierarchicalValueIndex := fmt.Sprintf("policies.%d.hierarchical_value", i)
	if _, ok := d.GetOk(hierarchicalValueIndex); !ok || param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"critical": utils.PathSearch("critical", param, nil),
			"major":    utils.PathSearch("major", param, nil),
			"minor":    utils.PathSearch("minor", param, nil),
			"info":     utils.PathSearch("info", param, nil),
		},
	}

	return rst
}

func resourceCesAlarmTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAlarmTemplate: Delete CES alarm template
	var (
		deleteAlarmTemplateHttpUrl = "v2/{project_id}/alarm-templates/batch-delete"
		deleteAlarmTemplateProduct = "ces"
	)
	deleteAlarmTemplateClient, err := cfg.NewServiceClient(deleteAlarmTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	deleteAlarmTemplatePath := deleteAlarmTemplateClient.Endpoint + deleteAlarmTemplateHttpUrl
	deleteAlarmTemplatePath = strings.ReplaceAll(deleteAlarmTemplatePath, "{project_id}",
		deleteAlarmTemplateClient.ProjectID)

	deleteAlarmTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteAlarmTemplateOpt.JSONBody = utils.RemoveNil(buildDeleteAlarmTemplateBodyParams(d))
	fmt.Println("")
	_, err = deleteAlarmTemplateClient.Request("POST", deleteAlarmTemplatePath, &deleteAlarmTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting CES alarm template: %s", err)
	}

	return nil
}

func buildDeleteAlarmTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_ids":           []interface{}{d.Id()},
		"delete_associate_alarm": utils.ValueIgnoreEmpty(d.Get("delete_associate_alarm")),
	}
	return bodyParams
}
