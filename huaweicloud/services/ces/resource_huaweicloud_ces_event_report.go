package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventReportNonUpdatableFields = []string{
	"name", "source", "time", "detail", "detail.*.type", "detail.*.state", "detail.*.level",
	"detail.*.content", "detail.*.group_id", "detail.*.resource_id", "detail.*.resource_name",
	"detail.*.user", "detail.*.dimensions", "detail.*.dimensions.*.name", "detail.*.dimensions.*.value",
}

// @API CES POST /V1.0/{project_id}/events
func ResourceCesEventReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCesEventReportCreate,
		UpdateContext: resourceCesEventReportUpdate,
		ReadContext:   resourceCesEventReportRead,
		DeleteContext: resourceCesEventReportDelete,

		CustomizeDiff: config.FlexibleForceNew(eventReportNonUpdatableFields),

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
				Description: `Specifies the CES event name.`,
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event source.`,
			},
			"time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the occurrence time of the event.`,
			},
			"detail": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        eventDetailSchema(),
				MaxItems:    1,
				Description: `Specifies the detail of the CES event.`,
			},
		},
	}
}

func eventDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"state": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event status.`,
			},
			"level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event level.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event type.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event content.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group that the event belongs to.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource ID.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource name.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event user.`,
			},
			"dimensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        eventReportDimensionsSchema(),
				Description: `Specifies the resource dimensions.`,
			},
		},
	}
	return &sc
}

func eventReportDimensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource dimension name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource dimension value.`,
			},
		},
	}
}

func resourceCesEventReportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createEventReportHttpUrl = "V1.0/{project_id}/events"
		createEventReportProduct = "ces"
	)
	createEventReportClient, err := cfg.NewServiceClient(createEventReportProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createEventReportPath := createEventReportClient.Endpoint + createEventReportHttpUrl
	createEventReportPath = strings.ReplaceAll(createEventReportPath, "{project_id}",
		createEventReportClient.ProjectID)

	createEventReportOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	rawParams, err := buildEventReportBodyParams(d)
	if err != nil {
		return diag.Errorf("error building CES event report body params: %s", err)
	}

	params := utils.RemoveNil(rawParams)
	createEventReportOpt.JSONBody = []map[string]interface{}{params}
	createEventReportResp, err := createEventReportClient.Request("POST",
		createEventReportPath, &createEventReportOpt)
	if err != nil {
		return diag.Errorf("error creating CES event report: %s", err)
	}

	createEventReportRespBody, err := utils.FlattenResponse(createEventReportResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("[0].event_id", createEventReportRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES event report: ID is not found in API response")
	}
	d.SetId(id)

	return resourceCesEventReportRead(ctx, d, meta)
}

func buildEventReportBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	rawEventTime := d.Get("time").(string)
	eventTime, err := parseTimeToTimestamp(rawEventTime)
	if err != nil {
		return nil, err
	}

	param := map[string]interface{}{
		"event_name":   d.Get("name").(string),
		"event_source": d.Get("source").(string),
		"time":         eventTime * 1000,
		"detail":       buildEventDetailBodyParams(d.Get("detail")),
	}

	return param, nil
}

func buildEventDetailBodyParams(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		detailMap := rawArray[0].(map[string]interface{})
		detailParams := map[string]interface{}{
			"event_state":   detailMap["state"],
			"event_level":   detailMap["level"],
			"event_type":    utils.ValueIgnoreEmpty(detailMap["type"]),
			"content":       utils.ValueIgnoreEmpty(detailMap["content"]),
			"group_id":      utils.ValueIgnoreEmpty(detailMap["group_id"]),
			"resource_id":   utils.ValueIgnoreEmpty(detailMap["resource_id"]),
			"resource_name": utils.ValueIgnoreEmpty(detailMap["resource_name"]),
			"event_user":    utils.ValueIgnoreEmpty(detailMap["user"]),
			"dimensions":    buildEventDimensionsBodyParams(detailMap["dimensions"]),
		}
		return detailParams
	}
	return nil
}

func buildEventDimensionsBodyParams(rawParam interface{}) []map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok && len(rawArray) > 0 {
		dimensionsParams := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			dimensionsMap := v.(map[string]interface{})
			dimensionsParams = append(dimensionsParams, map[string]interface{}{
				"name":  dimensionsMap["name"].(string),
				"value": dimensionsMap["value"].(string),
			})
		}
		return dimensionsParams
	}
	return nil
}

func resourceCesEventReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCesEventReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCesEventReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
