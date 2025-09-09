package coc

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

var diagnosisTaskNonUpdatableParams = []string{"resource_id", "type", "extra_properties"}

// @API COC POST /v1/diagnosis/tasks
// @API COC GET /v1/diagnosis/tasks/{task_id}
func ResourceDiagnosisTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiagnosisTaskCreate,
		ReadContext:   resourceDiagnosisTaskRead,
		UpdateContext: resourceDiagnosisTaskUpdate,
		DeleteContext: resourceDiagnosisTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDiagnosisTaskImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(diagnosisTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extra_properties": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"work_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     diagnosisTaskNodeListSchema(),
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func diagnosisTaskNodeListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_zh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"diagnosis_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDiagnosisTaskCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_ids": []interface{}{d.Get("resource_id")},
		"type":         d.Get("type"),
	}

	if v, ok := d.GetOk("extra_properties"); ok {
		bodyParams["extra_properties"] = []interface{}{v}
	}

	return bodyParams
}

func resourceDiagnosisTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	httpUrl := "v1/diagnosis/tasks"
	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDiagnosisTaskCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating the COC diagnosis task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening COC diagnosis response: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC diagnosis task ID from the API response")
	}

	d.SetId(id)

	return resourceDiagnosisTaskRead(ctx, d, meta)
}

func resourceDiagnosisTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	getRespBody, err := GetDiagnosisTask(client, d.Id(), d.Get("resource_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00010002"),
			"error retrieving diagnosis task")
	}
	if utils.PathSearch("data", getRespBody, nil) == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving diagnosis task")
	}

	mErr := multierror.Append(
		nil,
		d.Set("code", utils.PathSearch("data.code", getRespBody, nil)),
		d.Set("project_id", utils.PathSearch("data.project_id", getRespBody, nil)),
		d.Set("user_id", utils.PathSearch("data.user_id", getRespBody, nil)),
		d.Set("user_name", utils.PathSearch("data.user_name", getRespBody, nil)),
		d.Set("progress", utils.PathSearch("data.progress", getRespBody, nil)),
		d.Set("work_order_id", utils.PathSearch("data.work_order_id", getRespBody, nil)),
		d.Set("instance_name", utils.PathSearch("data.instance_name", getRespBody, nil)),
		d.Set("type", utils.PathSearch("data.type", getRespBody, nil)),
		d.Set("status", utils.PathSearch("data.status", getRespBody, nil)),
		d.Set("start_time", utils.PathSearch("data.start_time", getRespBody, nil)),
		d.Set("end_time", utils.PathSearch("data.end_time", getRespBody, nil)),
		d.Set("instance_num", utils.PathSearch("data.instance_num", getRespBody, nil)),
		d.Set("os_type", utils.PathSearch("data.os_type", getRespBody, nil)),
		d.Set("region", utils.PathSearch("data.region", getRespBody, nil)),
		d.Set("node_list", flattenCocDiagnosisTaskNodeList(utils.PathSearch("data.node_list", getRespBody, nil))),
		d.Set("message", utils.PathSearch("data.message", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDiagnosisTask(client *golangsdk.ServiceClient, diagnosisTaskID string, resourceID string) (interface{}, error) {
	httpUrl := "v1/diagnosis/tasks/{task_id}?instance_id={instance_id}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{task_id}", diagnosisTaskID)
	readPath = strings.ReplaceAll(readPath, "{instance_id}", resourceID)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	readResp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}
	readRespBody, err := utils.FlattenResponse(readResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening diagnosis response: %s", err)
	}
	return readRespBody, nil
}

func flattenCocDiagnosisTaskNodeList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"id":                utils.PathSearch("id", raw, nil),
				"code":              utils.PathSearch("code", raw, nil),
				"name":              utils.PathSearch("name", raw, nil),
				"name_zh":           utils.PathSearch("name_zh", raw, nil),
				"diagnosis_task_id": utils.PathSearch("diagnosis_task_id", raw, nil),
				"status":            utils.PathSearch("status", raw, nil),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func resourceDiagnosisTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDiagnosisTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting diagnosis task resource is not supported. The diagnosis task resource is only removed from" +
		" the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDiagnosisTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <resource_id>/<id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("resource_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
