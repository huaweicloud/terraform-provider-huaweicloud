package waf

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

// @API WAF POST /v2/{project_id}/waf/alert
// @API WAF PUT /v2/{project_id}/waf/alert/{alert_id}
// @API WAF GET /v2/{project_id}/waf/alerts
func ResourceWafAlarmNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafAlarmNotificationCreate,
		ReadContext:   resourceWafAlarmNotificationRead,
		UpdateContext: resourceWafAlarmNotificationUpdate,
		DeleteContext: resourceWafAlarmNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWafAlarmNotificationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			// This field is Optional in API document, but it is required actually.
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`Specifies the name of the alarm notification.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			// This field is Optional in API document, but it is required actually.
			"topic_urn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`Specifies the topic URN of the SMN.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			// This field is Optional in API document, but it is required actually.
			"notice_class": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`Specifies the type of the alarm notification.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			// This field `enterprise_project_id` only appears in the creation and reading API document,
			// not in the update API document. So this field does not participate in the editing operation.
			// The value of this field is always `0` in detail API response, so this field is not filled back when querying.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// For Optional fields, default values are not marked in the document, and the default value cannot be
			// subjectively judged.
			// Therefore, the default value of the Optional class field is not described in detail in the document.
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sendfreq": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"times": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// This field `threat` only appears in the edit API document, not in the create API document.
			// However, the actual test creation API also supports this field.
			"threat": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nearly_expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_all_enterprise_project": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// This field `description` only appears in the creation API document, not in the edit API document.
			// However, the actual test edit API also supports this field.
			// This field supporting to be updated to empty, so it is not marked as Computed.
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
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildWafAlarmNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                      utils.ValueIgnoreEmpty(d.Get("name")),
		"topic_urn":                 utils.ValueIgnoreEmpty(d.Get("topic_urn")),
		"notice_class":              utils.ValueIgnoreEmpty(d.Get("notice_class")),
		"enabled":                   d.Get("enabled"),
		"sendfreq":                  utils.ValueIgnoreEmpty(d.Get("sendfreq")),
		"locale":                    utils.ValueIgnoreEmpty(d.Get("locale")),
		"times":                     utils.ValueIgnoreEmpty(d.Get("times")),
		"threat":                    utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("threat").([]interface{}))),
		"nearly_expired_time":       utils.ValueIgnoreEmpty(d.Get("nearly_expired_time")),
		"is_all_enterprise_project": d.Get("is_all_enterprise_project"),
		"description":               utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceWafAlarmNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/waf/alert"
		epsID   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + fmt.Sprintf("%s?enterpriseProjectId=%s", httpUrl, epsID)
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
			"X-Language":   "en-us",
		},
		JSONBody: utils.RemoveNil(buildWafAlarmNotificationBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF alert notice configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the alert notice ID from the API response")
	}

	d.SetId(id)

	return resourceWafAlarmNotificationRead(ctx, d, meta)
}

func GetAlarmNotificationDetail(client *golangsdk.ServiceClient, alertID, epsID string) (interface{}, error) {
	requestPath := client.Endpoint + "v2/{project_id}/waf/alerts"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsID != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%s", requestPath, epsID)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("items[?id == '%s']|[0]", alertID)
	alarmNotificationDetail := utils.PathSearch(expression, respBody, nil)
	if alarmNotificationDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return alarmNotificationDetail, nil
}

func resourceWafAlarmNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsID  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	respDetail, err := GetAlarmNotificationDetail(client, d.Id(), epsID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF alert notice configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respDetail, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", respDetail, nil)),
		d.Set("notice_class", utils.PathSearch("notice_class", respDetail, nil)),
		d.Set("enabled", utils.PathSearch("enabled", respDetail, nil)),
		d.Set("sendfreq", utils.PathSearch("sendfreq", respDetail, nil)),
		d.Set("locale", utils.PathSearch("locale", respDetail, nil)),
		d.Set("times", utils.PathSearch("times", respDetail, nil)),
		d.Set("threat", utils.PathSearch("threat", respDetail, nil)),
		d.Set("nearly_expired_time", utils.PathSearch("nearly_expired_time", respDetail, nil)),
		d.Set("is_all_enterprise_project", utils.PathSearch("is_all_enterprise_project", respDetail, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respDetail, nil)),
		d.Set("description", utils.PathSearch("description", respDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateAlarmNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                      utils.ValueIgnoreEmpty(d.Get("name")),
		"topic_urn":                 utils.ValueIgnoreEmpty(d.Get("topic_urn")),
		"notice_class":              utils.ValueIgnoreEmpty(d.Get("notice_class")),
		"enabled":                   d.Get("enabled"),
		"sendfreq":                  utils.ValueIgnoreEmpty(d.Get("sendfreq")),
		"locale":                    utils.ValueIgnoreEmpty(d.Get("locale")),
		"times":                     utils.ValueIgnoreEmpty(d.Get("times")),
		"threat":                    utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("threat").([]interface{}))),
		"nearly_expired_time":       utils.ValueIgnoreEmpty(d.Get("nearly_expired_time")),
		"is_all_enterprise_project": d.Get("is_all_enterprise_project"),
		// Field `description` supports to be updated to empty.
		"description": d.Get("description"),
	}

	return bodyParams
}

func resourceWafAlarmNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/waf/alert/{alert_id}"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{alert_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
			"X-Language":   "en-us",
		},
		JSONBody: utils.RemoveNil(buildUpdateAlarmNotificationBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating WAF alert notice configuration: %s", err)
	}

	return resourceWafAlarmNotificationRead(ctx, d, meta)
}

func resourceWafAlarmNotificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not destroy the WAF alert notice configuration,
	but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceWafAlarmNotificationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf(`invalid format specified for import ID, must be <id>/<enterprise_project_id>,
		but got %s`, d.Id())
	}
	d.SetId(parts[0])
	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[1])
}
