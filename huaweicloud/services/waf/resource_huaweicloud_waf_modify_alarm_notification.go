package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var modifyAlarmNotificationNonUpdatableParams = []string{"alert_id", "name", "topic_urn", "notice_class", "enabled",
	"sendfreq", "locale", "times", "threat", "nearly_expired_time", "is_all_enterprise_project"}

// @API WAF PUT /v2/{project_id}/waf/alert/{alert_id}
func ResourceModifyAlarmNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModifyAlarmNotificationCreate,
		ReadContext:   resourceModifyAlarmNotificationRead,
		UpdateContext: resourceModifyAlarmNotificationUpdate,
		DeleteContext: resourceModifyAlarmNotificationDelete,

		CustomizeDiff: config.FlexibleForceNew(modifyAlarmNotificationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alert_id": {
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
			"notice_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sendfreq": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// This field does not take effect due to the API issue.
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"times": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"threat": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nearly_expired_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_all_enterprise_project": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildModifyAlarmNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                      d.Get("name"),
		"topic_urn":                 d.Get("topic_urn"),
		"notice_class":              d.Get("notice_class"),
		"enabled":                   d.Get("enabled"),
		"sendfreq":                  utils.ValueIgnoreEmpty(d.Get("sendfreq")),
		"locale":                    utils.ValueIgnoreEmpty(d.Get("locale")),
		"times":                     utils.ValueIgnoreEmpty(d.Get("times")),
		"threat":                    utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("threat").([]interface{}))),
		"nearly_expired_time":       utils.ValueIgnoreEmpty(d.Get("nearly_expired_time")),
		"is_all_enterprise_project": d.Get("is_all_enterprise_project"),
	}

	return bodyParams
}

func resourceModifyAlarmNotificationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		alertId = d.Get("alert_id").(string)
		httpUrl = "v2/{project_id}/waf/alert/{alert_id}"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{alert_id}", alertId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
			"X-Language":   "en-us",
		},
		JSONBody: utils.RemoveNil(buildModifyAlarmNotificationBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating WAF alert notice configuration: %s", err)
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

	return nil
}

func resourceModifyAlarmNotificationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceModifyAlarmNotificationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceModifyAlarmNotificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
