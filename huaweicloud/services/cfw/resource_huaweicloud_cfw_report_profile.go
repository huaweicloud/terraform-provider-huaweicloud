package cfw

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
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

// @API CFW POST /v1/{project_id}/report-profile
// @API CFW PUT /v1/{project_id}/report-profile/{report_profile_id}
// @API CFW GET /v1/{project_id}/report-profile/{report_profile_id}
// @API CFW DELETE /v1/{project_id}/report-profile/{report_profile_id}
func ResourceReportProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReportProfileCreate,
		ReadContext:   resourceReportProfileRead,
		UpdateContext: resourceReportProfileUpdate,
		DeleteContext: resourceReportProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceReportProfileImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"fw_instance_id",
			"category",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			// Unable to retrieve the value of field `fw_instance_id` from the details API.
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			// API requires the fields `send_period`, `send_week_day`, `status` and `subscription_type` to be numeric
			// types, but the default value of `0` has practical significance, so their types have been changed to strings.
			"send_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"send_week_day": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"statistic_period": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     statisticPeriodSchema(),
			},
			// Unable to retrieve the value of field `description` from the details API. So no Computed attribute was added.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func statisticPeriodSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildCreateReportProfileStatisticPeriodBodyParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"start_time": rawMap["start_time"],
		"end_time":   rawMap["end_time"],
	}
}

func convertStringToInt(rawValue string) interface{} {
	if rawValue == "" {
		return nil
	}

	r, err := strconv.Atoi(rawValue)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", rawValue)
		return nil
	}

	return r
}

func buildCreateReportProfileBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"category":          d.Get("category"),
		"name":              d.Get("name"),
		"topic_urn":         d.Get("topic_urn"),
		"send_period":       convertStringToInt(d.Get("send_period").(string)),
		"send_week_day":     convertStringToInt(d.Get("send_week_day").(string)),
		"status":            convertStringToInt(d.Get("status").(string)),
		"subscription_type": convertStringToInt(d.Get("subscription_type").(string)),
		"statistic_period":  buildCreateReportProfileStatisticPeriodBodyParam(d.Get("statistic_period").([]interface{})),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceReportProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/report-profile"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateReportProfileBodyParam(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CFW report profile: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CFW report profile: ID is not found in API response")
	}
	d.SetId(id)

	return resourceReportProfileRead(ctx, d, meta)
}

func flattenStatisticPeriodAttribute(respMap interface{}) []map[string]interface{} {
	if respMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"start_time": utils.PathSearch("start_time", respMap, nil),
			"end_time":   utils.PathSearch("end_time", respMap, nil),
		},
	}
}

func GetReportProfileDetail(client *golangsdk.ServiceClient, id, fwInstanceId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/report-profile/{report_profile_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{report_profile_id}", id)
	requestPath += fmt.Sprintf("?fw_instance_id=%s", fwInstanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenReportProfileNumberValue(respValue interface{}) string {
	if respValue == nil {
		return ""
	}

	// The response value type for this field should be float64.
	float64Value, ok := respValue.(float64)
	if !ok {
		log.Printf("[WARN] failed to convert %v to float64", respValue)
		return ""
	}

	return strconv.Itoa(int(float64Value))
}

func resourceReportProfileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	respBody, err := GetReportProfileDetail(client, d.Id(), d.Get("fw_instance_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving CFW report profile",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("category", utils.PathSearch("data.category", respBody, nil)),
		d.Set("name", utils.PathSearch("data.name", respBody, nil)),
		d.Set("topic_urn", utils.PathSearch("data.topic_urn", respBody, nil)),
		d.Set("send_period", flattenReportProfileNumberValue(utils.PathSearch("data.send_period", respBody, nil))),
		d.Set("send_week_day", flattenReportProfileNumberValue(utils.PathSearch("data.send_week_day", respBody, nil))),
		d.Set("status", flattenReportProfileNumberValue(utils.PathSearch("data.status", respBody, nil))),
		d.Set("subscription_type", flattenReportProfileNumberValue(utils.PathSearch("data.subscription_type", respBody, nil))),
		d.Set("statistic_period", flattenStatisticPeriodAttribute(utils.PathSearch("data.statistic_period", respBody, nil))),
		d.Set("topic_name", utils.PathSearch("data.topic_name", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateReportProfileStatisticPeriodBodyParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"start_time": rawMap["start_time"],
		"end_time":   rawMap["end_time"],
	}
}

func buildUpdateReportProfileBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":              d.Get("name"),
		"topic_urn":         d.Get("topic_urn"),
		"send_period":       convertStringToInt(d.Get("send_period").(string)),
		"send_week_day":     convertStringToInt(d.Get("send_week_day").(string)),
		"status":            convertStringToInt(d.Get("status").(string)),
		"subscription_type": convertStringToInt(d.Get("subscription_type").(string)),
		"statistic_period":  buildUpdateReportProfileStatisticPeriodBodyParam(d.Get("statistic_period").([]interface{})),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceReportProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/report-profile/{report_profile_id}"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{report_profile_id}", d.Id())
	requestPath += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateReportProfileBodyParam(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating CFW report profile: %s", err)
	}

	return resourceReportProfileRead(ctx, d, meta)
}

func resourceReportProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/report-profile/{report_profile_id}"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{report_profile_id}", d.Id())
	requestPath += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CFW report profile: %s", err)
	}

	return nil
}

func resourceReportProfileImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <fw_instance_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("fw_instance_id", parts[0])
}
