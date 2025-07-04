package sms

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

var consistencyResultNonUpdatableParams = []string{"task_id", "consistency_result", "consistency_result.*.dir_check",
	"consistency_result.*.num_total_files", "consistency_result.*.num_different_files",
	"consistency_result.*.num_target_miss_files", "consistency_result.*.num_target_more_files"}

// @API SMS POST /v3/tasks/{task_id}/consistency-result
func ResourceTaskConsistencyResultReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskConsistencyResultReportCreate,
		ReadContext:   resourceTaskConsistencyResultReportRead,
		UpdateContext: resourceTaskConsistencyResultReportUpdate,
		DeleteContext: resourceTaskConsistencyResultReportDelete,

		CustomizeDiff: config.FlexibleForceNew(consistencyResultNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consistency_result": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dir_check": {
							Type:     schema.TypeString,
							Required: true,
						},
						"num_total_files": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"num_different_files": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"num_target_miss_files": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"num_target_more_files": {
							Type:     schema.TypeInt,
							Required: true,
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
		},
	}
}

func buildTaskConsistencyResultReportCreateOpts(d *schema.ResourceData) map[string]interface{} {
	consistencyResultParams := make([]map[string]interface{}, 0)
	if rawArray, ok := d.Get("consistency_result").([]interface{}); ok {
		if len(rawArray) != 0 {
			consistencyResultParams = make([]map[string]interface{}, len(rawArray))
			for i, v := range rawArray {
				raw := v.(map[string]interface{})
				consistencyResultParams[i] = map[string]interface{}{
					"dir_check":             raw["dir_check"],
					"num_total_files":       raw["num_total_files"],
					"num_different_files":   raw["num_different_files"],
					"num_target_miss_files": raw["num_target_miss_files"],
					"num_target_more_files": raw["num_target_more_files"],
				}
			}
		}
	}

	bodyParams := map[string]interface{}{
		"consistency_result": consistencyResultParams,
	}

	return bodyParams
}

func resourceTaskConsistencyResultReportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createHttpUrl := "v3/tasks/{task_id}/consistency-result"
	taskID := d.Get("task_id").(string)
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildTaskConsistencyResultReportCreateOpts(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SMS task consistency result report: %s", err)
	}

	d.SetId(taskID)

	return nil
}

func resourceTaskConsistencyResultReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskConsistencyResultReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskConsistencyResultReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting task consistency result report resource is not supported. The task consistency result report resource is" +
		" only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
