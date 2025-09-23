package ces

import (
	"context"
	"errors"
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

var notificationMaskNonUpdatableParams = []string{"relation_type"}

// @API CES PUT /v2/{project_id}/notification-masks
// @API CES PUT /v2/{project_id}/notification-masks/{notification_mask_id}
// @API CES POST /v2/{project_id}/notification-masks/batch-delete
// @API CES POST /v2/{project_id}/notification-masks/batch-query
// @API CES GET /v2/{project_id}/notification-masks/{notification_mask_id}/resources
func ResourceNotificationMask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationMaskCreate,
		ReadContext:   resourceNotificationMaskRead,
		UpdateContext: resourceNotificationMaskUpdate,
		DeleteContext: resourceNotificationMaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNotificationMaskImportState,
		},
		CustomizeDiff: config.FlexibleForceNew(notificationMaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"relation_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of a resource that is associated with an alarm notification masking rule.`,
			},
			"mask_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm notification masking type.`,
			},
			"relation_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   `Specifies the alarm policy ID.`,
				ConflictsWith: []string{"relation_id"},
			},
			"relation_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   `Specifies the alarm rule ID.`,
				ConflictsWith: []string{"relation_ids"},
			},
			"mask_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the masking rule name.`,
			},
			"start_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the masking start date.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the masking start time.`,
			},
			"end_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the masking end date.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the masking end time.`,
			},
			"resources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the resource for which alarm notifications will be masked.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the resource namespace in **service.item** format.`,
						},
						"dimensions": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: `Specifies the resource dimension information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the dimension of a resource.`,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the value of a resource dimension.`,
									},
								},
							},
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"mask_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the alarm notification masking status.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The alarm policy list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm policy ID.`,
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The metric name of a resource.`,
						},
						"extra_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The extended metric information.`,
							Elem:        notificationMaskPolicyExtendedInfo(),
						},
						"period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The period for determining whether to generate an alarm, in seconds.`,
						},
						"filter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data rollup method.`,
						},
						"comparison_operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operator.`,
						},
						"value": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The alarm threshold.`,
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data unit.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of consecutive times that alarm conditions are met.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm policy type.`,
						},
						"suppress_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The interval for triggering alarms.`,
						},
						"alarm_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The alarm severity.`,
						},
						"selected_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit you selected, which is used for subsequent metric data display and calculation.`,
						},
					},
				},
			},
		},
	}
}

func notificationMaskPolicyExtendedInfo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"origin_metric_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The original metric name.`,
			},
			"metric_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The metric name prefix.`,
			},
			"custom_proc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of a user process.`,
			},
			"metric_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The metric type.`,
			},
		},
	}
}

func resourceNotificationMaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createNotificationMaskHttpUrl = "v2/{project_id}/notification-masks"
		createNotificationMaskProduct = "ces"
	)
	createNotificationMaskClient, err := cfg.NewServiceClient(createNotificationMaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createNotificationMaskPath := createNotificationMaskClient.Endpoint + createNotificationMaskHttpUrl
	createNotificationMaskPath = strings.ReplaceAll(createNotificationMaskPath, "{project_id}", createNotificationMaskClient.ProjectID)

	createNotificationMaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createNotificationMaskOpt.JSONBody = utils.RemoveNil(buildNotificationMaskBodyParams(d))
	createNotificationMaskResp, err := createNotificationMaskClient.Request("PUT", createNotificationMaskPath, &createNotificationMaskOpt)
	if err != nil {
		return diag.Errorf("error creating CES notification mask: %s", err)
	}

	createNotificationMaskRespBody, err := utils.FlattenResponse(createNotificationMaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("notification_mask_id", createNotificationMaskRespBody, "").(string)
	// When relation_type is ALARM_RULE, no ID is returned
	if id != "" {
		d.SetId(id)
	}

	return resourceNotificationMaskRead(ctx, d, meta)
}

func buildNotificationMaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	relationType := d.Get("relation_type").(string)
	param := map[string]interface{}{
		"mask_name":     utils.ValueIgnoreEmpty(d.Get("mask_name")),
		"relation_type": relationType,
		"resources":     buildResourcesParams(d.Get("resources").(*schema.Set).List()),
		"mask_type":     d.Get("mask_type").(string),
		"start_date":    utils.ValueIgnoreEmpty(d.Get("start_date")),
		"start_time":    utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_date":      utils.ValueIgnoreEmpty(d.Get("end_date")),
		"end_time":      utils.ValueIgnoreEmpty(d.Get("end_time")),
	}
	if relationType == "ALARM_RULE" {
		param["relation_ids"] = []string{d.Get("relation_id").(string)}
	} else {
		param["relation_ids"] = d.Get("relation_ids").(*schema.Set).List()
	}

	return param
}

func buildResourcesParams(resources []interface{}) []map[string]interface{} {
	var resourcesParams []map[string]interface{}
	for _, resource := range resources {
		resourceParams := map[string]interface{}{
			"namespace":  resource.(map[string]interface{})["namespace"].(string),
			"dimensions": buildDimensionsParams(resource.(map[string]interface{})["dimensions"].(*schema.Set).List()),
		}
		resourcesParams = append(resourcesParams, resourceParams)
	}

	return resourcesParams
}

func buildDimensionsParams(dimensions []interface{}) []map[string]interface{} {
	var dimensionsParams []map[string]interface{}
	for _, dimension := range dimensions {
		dimensionParams := map[string]interface{}{
			"name":  dimension.(map[string]interface{})["name"].(string),
			"value": dimension.(map[string]interface{})["value"].(string),
		}
		dimensionsParams = append(dimensionsParams, dimensionParams)
	}

	return dimensionsParams
}

func resourceNotificationMaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	readNotificationMaskClient, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	relationType := d.Get("relation_type").(string)
	relationId := d.Get("relation_id").(string)
	id := d.Id()
	mask, err := GetNotificationMask(readNotificationMaskClient, relationType, relationId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES notification mask")
	}

	maskId := utils.PathSearch("notification_mask_id", mask, "").(string)
	if relationType == "ALARM_RULE" {
		d.SetId(maskId)
	}

	resources, err := getNotificationMaskResources(readNotificationMaskClient, maskId)
	if err != nil {
		return diag.Errorf("error retrieving CES notification mask resources: %s", err)
	}

	rawPolicies := utils.PathSearch("policies", mask, make([]interface{}, 0)).([]interface{})
	policies := flattenNotificationMaskPolicies(rawPolicies)
	relationIds := utils.PathSearch("policies[*].alarm_policy_id", mask, make([]interface{}, 0)).([]interface{})
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("relation_type", utils.PathSearch("relation_type", mask, nil)),
		d.Set("mask_type", utils.PathSearch("mask_type", mask, nil)),
		d.Set("relation_ids", relationIds),
		d.Set("relation_id", utils.PathSearch("relation_id", mask, nil)),
		d.Set("mask_name", utils.PathSearch("mask_name", mask, nil)),
		d.Set("start_date", utils.PathSearch("start_date", mask, nil)),
		d.Set("start_time", utils.PathSearch("start_time", mask, nil)),
		d.Set("end_date", utils.PathSearch("end_date", mask, nil)),
		d.Set("end_time", utils.PathSearch("end_time", mask, nil)),
		d.Set("resources", resources),
		d.Set("mask_status", utils.PathSearch("mask_status", mask, nil)),
		d.Set("policies", policies),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetNotificationMask(client *golangsdk.ServiceClient, relationType, relationId, id string) (interface{}, error) {
	getNotificationMaskHttpUrl := "v2/{project_id}/notification-masks/batch-query"
	getNotificationMaskPath := client.Endpoint + getNotificationMaskHttpUrl
	getNotificationMaskPath = strings.ReplaceAll(getNotificationMaskPath, "{project_id}", client.ProjectID)

	getNotificationMaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	param := map[string]interface{}{
		"relation_type": relationType,
	}

	if relationType == "ALARM_RULE" {
		param["relation_ids"] = []string{relationId}
	} else {
		param["relation_ids"] = []interface{}{}
		param["mask_id"] = id
	}

	getNotificationMaskOpt.JSONBody = param

	resp, err := client.Request("POST", getNotificationMaskPath, &getNotificationMaskOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	var mask interface{}
	if relationType == "ALARM_RULE" {
		mask = utils.PathSearch("[0]", respBody, nil)
	} else {
		mask = utils.PathSearch("notification_masks|[0]", respBody, nil)
	}

	if mask == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	maskId := utils.PathSearch("notification_mask_id", mask, "").(string)
	if maskId == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return mask, nil
}

func getNotificationMaskResources(client *golangsdk.ServiceClient, id string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/notification-masks/{notification_mask_id}/resources"
	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{notification_mask_id}", id)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	rst := make([]interface{}, 0)

	offset := 0
	for {
		path := fmt.Sprintf("%s?limit=100&offset=%d", basePath, offset)
		resp, err := client.Request("GET", path, &opt)

		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resources) == 0 {
			return rst, nil
		}

		rst = append(rst, resources...)
		offset += 100
		total := utils.PathSearch("count", respBody, float64(0))
		if int(total.(float64)) <= offset {
			break
		}
	}
	return rst, nil
}

func flattenNotificationMaskPolicies(policies []interface{}) []interface{} {
	rst := make([]interface{}, len(policies))
	for i, policy := range policies {
		rst[i] = map[string]interface{}{
			"alarm_policy_id":     utils.PathSearch("alarm_policy_id", policy, nil),
			"metric_name":         utils.PathSearch("metric_name", policy, nil),
			"extra_info":          flattenNotificationMaskExtInfo(utils.PathSearch("extra_info", policy, nil)),
			"period":              utils.PathSearch("period", policy, nil),
			"filter":              utils.PathSearch("filter", policy, nil),
			"comparison_operator": utils.PathSearch("comparison_operator", policy, nil),
			"value":               utils.PathSearch("value", policy, nil),
			"unit":                utils.PathSearch("unit", policy, nil),
			"count":               utils.PathSearch("count", policy, nil),
			"type":                utils.PathSearch("type", policy, nil),
			"suppress_duration":   utils.PathSearch("suppress_duration", policy, nil),
			"alarm_level":         utils.PathSearch("alarm_level", policy, nil),
			"selected_unit":       utils.PathSearch("selected_unit", policy, nil),
		}
	}
	return rst
}

func flattenNotificationMaskExtInfo(extInfo interface{}) []interface{} {
	if extInfo == nil {
		return nil
	}
	rst := map[string]interface{}{
		"origin_metric_name": utils.PathSearch("origin_metric_name", extInfo, nil),
		"metric_prefix":      utils.PathSearch("metric_prefix", extInfo, nil),
		"custom_proc_name":   utils.PathSearch("custom_proc_name", extInfo, nil),
		"metric_type":        utils.PathSearch("metric_type", extInfo, nil),
	}
	return []interface{}{rst}
}

func resourceNotificationMaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateNotificationMaskClient, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	relationType := d.Get("relation_type").(string)
	var updateNotificationMaskPath string

	updateNotificationMaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if relationType == "ALARM_RULE" {
		updateNotificationMaskHttpUrl := "v2/{project_id}/notification-masks"
		updateNotificationMaskPath = updateNotificationMaskClient.Endpoint + updateNotificationMaskHttpUrl
		updateNotificationMaskPath = strings.ReplaceAll(updateNotificationMaskPath, "{project_id}", updateNotificationMaskClient.ProjectID)
	} else {
		updateNotificationMaskHttpUrl := "v2/{project_id}/notification-masks/{notification_mask_id}"
		updateNotificationMaskPath = updateNotificationMaskClient.Endpoint + updateNotificationMaskHttpUrl
		updateNotificationMaskPath = strings.ReplaceAll(updateNotificationMaskPath, "{project_id}", updateNotificationMaskClient.ProjectID)
		updateNotificationMaskPath = strings.ReplaceAll(updateNotificationMaskPath, "{notification_mask_id}", d.Id())
		updateNotificationMaskOpt.OkCodes = []int{204}
	}

	updateNotificationMaskOpt.JSONBody = utils.RemoveNil(buildNotificationMaskBodyParams(d))

	_, err = updateNotificationMaskClient.Request("PUT", updateNotificationMaskPath, &updateNotificationMaskOpt)
	if err != nil {
		return diag.Errorf("error updating CES notification mask: %s", err)
	}

	return resourceNotificationMaskRead(ctx, d, meta)
}

func resourceNotificationMaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	deleteNotificationMaskClient, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	deleteNotificationMaskHttpUrl := "v2/{project_id}/notification-masks/batch-delete"
	deleteNotificationMaskPath := deleteNotificationMaskClient.Endpoint + deleteNotificationMaskHttpUrl
	deleteNotificationMaskPath = strings.ReplaceAll(deleteNotificationMaskPath, "{project_id}", deleteNotificationMaskClient.ProjectID)

	deleteNotificationMaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteNotificationMaskOpt.JSONBody = map[string]interface{}{
		"notification_mask_ids": []string{d.Id()},
	}

	_, err = deleteNotificationMaskClient.Request("POST", deleteNotificationMaskPath, &deleteNotificationMaskOpt)
	if err != nil {
		return diag.Errorf("error deleting CES notification mask: %s", err)
	}

	return nil
}

func resourceNotificationMaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import id, " +
			"must be <relation_type>/<relation_id> or <relation_type>/<notification_mask_id>")
	}
	relationType := parts[0]
	d.Set("relation_type", relationType)
	if relationType == "ALARM_RULE" {
		d.Set("relation_id", parts[1])
	} else {
		d.SetId(parts[1])
	}

	return []*schema.ResourceData{d}, nil
}
