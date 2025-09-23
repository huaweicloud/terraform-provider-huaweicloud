package coc

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var customEventReportNonUpdatableParams = []string{"integration_key", "alarm_id", "alarm_name", "alarm_level", "time",
	"namespace", "application_id", "alarm_desc", "alarm_source", "region_id", "resource_name", "resource_id", "url",
	"alarm_status", "additional"}

// @API COC POST /v1/event/huawei/custom/{integration_key}
func ResourceCustomEventReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomEventReportCreate,
		ReadContext:   resourceCustomEventReportRead,
		UpdateContext: resourceCustomEventReportUpdate,
		DeleteContext: resourceCustomEventReportDelete,

		CustomizeDiff: config.FlexibleForceNew(customEventReportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"integration_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_desc": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"additional": {
				Type:     schema.TypeString,
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

func buildCustomEventReportCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarmId":       d.Get("alarm_id"),
		"alarmName":     d.Get("alarm_name"),
		"alarmLevel":    d.Get("alarm_level"),
		"time":          d.Get("time"),
		"nameSpace":     d.Get("namespace"),
		"applicationId": d.Get("application_id"),
		"alarmDesc":     d.Get("alarm_desc"),
		"alarmSource":   d.Get("alarm_source"),
		"regionId":      utils.ValueIgnoreEmpty(d.Get("region_id")),
		"resourceName":  utils.ValueIgnoreEmpty(d.Get("resource_name")),
		"resourceId":    utils.ValueIgnoreEmpty(d.Get("resource_id")),
		"URL":           utils.ValueIgnoreEmpty(d.Get("url")),
		"alarmStatus":   utils.ValueIgnoreEmpty(d.Get("alarm_status")),
		"additional":    parseJson(d.Get("additional").(string)),
	}

	return bodyParams
}

func parseJson(v string) interface{} {
	if v == "" {
		return nil
	}

	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return v
	}

	return data
}

func resourceCustomEventReportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/event/huawei/custom/{integration_key}"
	integrationKey := d.Get("integration_key").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{integration_key}", integrationKey)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCustomEventReportCreateOpts(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error reporting the COC custom event (%s): %s", integrationKey, err)
	}

	d.SetId(integrationKey)

	return resourceCustomEventReportRead(ctx, d, meta)
}

func resourceCustomEventReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCustomEventReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCustomEventReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting custom event report resource is not supported. The custom event report resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
