package rds

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventOperateNonUpdatableParams = []string{
	"event_instances", "event_instances.*.event_id", "event_instances.*.instance_id", "operation_type",
	"event_schedule_window", "event_schedule_window.*.planned_day", "event_schedule_window.*.start_time",
	"event_schedule_window.*.end_time",
}

// @API RDS POST /v3/{project_id}/schedule-events
func ResourceEventOperate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventOperateCreate,
		ReadContext:   resourceEventOperateRead,
		UpdateContext: resourceEventOperateUpdate,
		DeleteContext: resourceEventOperateDelete,

		CustomizeDiff: config.FlexibleForceNew(eventOperateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"event_instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     eventInstancesSchema(),
			},
			"operation_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_schedule_window": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     eventScheduleWindowSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     eventOperateResultSchema(),
			},
		},
	}
}

func eventInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func eventScheduleWindowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"planned_day": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func eventOperateResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"success": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceEventOperateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/schedule-events"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateEventOperateBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating RDS event: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("results", flattenEventOperateResults(createRespBody)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting RDS event operate results: %s", mErr)
	}

	return nil
}

func buildCreateEventOperateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"event_instances":       buildEventOperateEventInstancesBodyParams(d.Get("event_instances")),
		"operation_type":        d.Get("operation_type"),
		"event_schedule_window": buildEventOperateEventScheduleWindowBodyParams(d.Get("event_schedule_window")),
	}

	return bodyParams
}

func buildEventOperateEventInstancesBodyParams(eventInstancesRaw interface{}) []map[string]interface{} {
	rawParams := eventInstancesRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		if raw, ok := v.(map[string]interface{}); ok {
			params = append(params, map[string]interface{}{
				"event_id":    raw["event_id"],
				"instance_id": raw["instance_id"],
			})
		}
	}
	return params
}

func buildEventOperateEventScheduleWindowBodyParams(eventScheduleWindowRaw interface{}) map[string]interface{} {
	rawParams := eventScheduleWindowRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if raw, ok := rawParams[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"planned_day": raw["planned_day"],
			"start_time":  utils.ValueIgnoreEmpty(raw["start_time"]),
			"end_time":    utils.ValueIgnoreEmpty(raw["end_time"]),
		}
		return rst
	}
	return nil
}

func flattenEventOperateResults(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	results := utils.PathSearch("results", resp, make([]interface{}, 0)).([]interface{})
	if len(results) == 0 {
		return nil
	}

	resultList := make([]map[string]interface{}, len(results))
	for i, v := range results {
		resultList[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"instance_id": utils.PathSearch("instance_id", v, nil),
			"job_id":      utils.PathSearch("job_id", v, nil),
			"error_code":  utils.PathSearch("error_code", v, nil),
			"error_msg":   utils.PathSearch("error_msg", v, nil),
			"success":     utils.PathSearch("success", v, false),
		}
	}

	return resultList
}

func resourceEventOperateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventOperateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventOperateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS event operate resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
