package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var metricDataAddNonUpdatableFields = []string{
	"metric", "metric.*.namespace", "metric.*.metric_name", "metric.*.dimensions",
	"metric.*.dimensions.*.name", "metric.*.dimensions.*.value", "ttl",
	"collect_time", "value", "unit", "type",
}

// @API CES POST /V1.0/{project_id}/metric-data
func ResourceMetricDataAdd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricDataAddCreate,
		UpdateContext: resourceMetricDataAddUpdate,
		ReadContext:   resourceMetricDataAddRead,
		DeleteContext: resourceMetricDataAddDelete,

		CustomizeDiff: config.FlexibleForceNew(metricDataAddNonUpdatableFields),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"metric": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        metricDataAddMetricSchema(),
				Description: `Specifies the CES monitoring metric data.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the monitoring metric data retention period.`,
			},
			"collect_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the collect time.`,
			},
			"value": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: `Specifies the value of the monitoring metric data.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the unit of the monitoring metric data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the monitoring metric data.`,
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

func metricDataAddMetricSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the customized namespace.`,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the metric ID.`,
			},
			"dimensions": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        metricDataAdddimensionsSchema(),
				Description: `Specifies the metric dimension.`,
			},
		},
	}
	return &sc
}

func metricDataAdddimensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dimension.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dimension value.`,
			},
		},
	}
}

func resourceMetricDataAddCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createMetricDataAddHttpUrl = "V1.0/{project_id}/metric-data"
		createMetricDataAddProduct = "ces"
	)
	createMetricDataAddClient, err := cfg.NewServiceClient(createMetricDataAddProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createMetricDataAddPath := createMetricDataAddClient.Endpoint + createMetricDataAddHttpUrl
	createMetricDataAddPath = strings.ReplaceAll(createMetricDataAddPath, "{project_id}",
		createMetricDataAddClient.ProjectID)

	createMetricDataAddOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	rawParams, err := buildMetricDataAddBodyParams(d)
	if err != nil {
		return diag.Errorf("error building CES metric data add body params: %s", err)
	}

	params := utils.RemoveNil(rawParams)
	createMetricDataAddOpt.JSONBody = []map[string]interface{}{params}
	_, err = createMetricDataAddClient.Request("POST", createMetricDataAddPath, &createMetricDataAddOpt)
	if err != nil {
		return diag.Errorf("error creating CES metric data add: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceMetricDataAddRead(ctx, d, meta)
}

func buildMetricDataAddBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	rawCollectTime := d.Get("collect_time").(string)
	collectTime, err := parseTimeToTimestamp(rawCollectTime)
	if err != nil {
		return nil, err
	}

	param := map[string]interface{}{
		"metric":       buildMetricBodyParams(d.Get("metric")),
		"ttl":          d.Get("ttl").(int),
		"collect_time": collectTime * 1000,
		"value":        d.Get("value").(float64),
		"unit":         utils.ValueIgnoreEmpty(d.Get("unit")),
		"type":         utils.ValueIgnoreEmpty(d.Get("type")),
	}

	return param, nil
}

func buildMetricBodyParams(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		metricMap := rawArray[0].(map[string]interface{})
		metricParams := map[string]interface{}{
			"namespace":   metricMap["namespace"].(string),
			"metric_name": metricMap["metric_name"].(string),
			"dimensions":  buildMetricDataAddDimensionsBodyParams(metricMap["dimensions"]),
		}
		return metricParams
	}
	return nil
}

func buildMetricDataAddDimensionsBodyParams(rawParam interface{}) []map[string]interface{} {
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

func resourceMetricDataAddRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMetricDataAddUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMetricDataAddDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
