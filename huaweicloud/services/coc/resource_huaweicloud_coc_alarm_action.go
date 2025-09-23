package coc

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

var alarmActionNonUpdatableParams = []string{"alarm_id", "task_type", "associated_task_id", "associated_task_type",
	"associated_task_name", "associated_task_enterprise_project_id", "runbook_instance_mode", "input_param",
	"target_instances", "target_instances.*.target_selection", "target_instances.*.order_no",
	"target_instances.*.batch_strategy", "target_instances.*.sub_target_instances", "target_instances.*.target_instances",
	"target_instances.*.sub_target_instances.*.target_selection", "target_instances.*.sub_target_instances.*.order_no",
	"target_instances.*.sub_target_instances.*.batch_strategy",
	"target_instances.*.sub_target_instances.*.target_instances", "region_id"}

// @API COC POST /v1/alarm-mgmt/alarm/{alarm_id}/auto-process
func ResourceAlarmAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmActionCreate,
		ReadContext:   resourceAlarmActionRead,
		UpdateContext: resourceAlarmActionUpdate,
		DeleteContext: resourceAlarmActionDelete,

		CustomizeDiff: config.FlexibleForceNew(alarmActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"alarm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runbook_instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_param": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_selection": {
							Type:     schema.TypeString,
							Required: true,
						},
						"order_no": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"batch_strategy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_instances": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sub_target_instances": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_selection": {
										Type:     schema.TypeString,
										Required: true,
									},
									"order_no": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"batch_strategy": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_instances": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"region_id": {
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

func buildAlarmActionCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_type":                             d.Get("task_type"),
		"associated_task_id":                    d.Get("associated_task_id"),
		"associated_task_type":                  d.Get("associated_task_type"),
		"associated_task_name":                  d.Get("associated_task_name"),
		"associated_task_enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("associated_task_enterprise_project_id")),
		"runbook_instance_mode":                 utils.ValueIgnoreEmpty(d.Get("runbook_instance_mode")),
		"input_param":                           utils.ValueIgnoreEmpty(d.Get("input_param")),
		"target_instances":                      buildAlarmActionTargetInstancesCreateOpts(d.Get("target_instances")),
		"region_id":                             utils.ValueIgnoreEmpty(d.Get("region_id")),
	}

	return bodyParams
}

func buildAlarmActionTargetInstancesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"target_selection":     raw["target_selection"],
				"order_no":             raw["order_no"],
				"target_instances":     utils.ValueIgnoreEmpty(raw["target_instances"]),
				"batch_strategy":       utils.ValueIgnoreEmpty(raw["batch_strategy"]),
				"sub_target_instances": buildAlarmActionTargetInstancesSubTargetInstancesCreateOpts(raw["sub_target_instances"]),
			}
		}
		return params
	}

	return nil
}

func buildAlarmActionTargetInstancesSubTargetInstancesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"target_selection": raw["target_selection"],
				"order_no":         raw["order_no"],
				"target_instances": utils.ValueIgnoreEmpty(raw["target_instances"]),
				"batch_strategy":   utils.ValueIgnoreEmpty(raw["batch_strategy"]),
			}
		}
		return params
	}

	return nil
}

func resourceAlarmActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/alarm-mgmt/alarm/{alarm_id}/auto-process"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{alarm_id}", d.Get("alarm_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAlarmActionCreateOpts(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC alarm action: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening alarm action: %s", err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC alarm action ID from the API response")
	}

	d.SetId(id)

	return resourceAlarmActionRead(ctx, d, meta)
}

func resourceAlarmActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting alarm action resource is not supported. The alarm action resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
